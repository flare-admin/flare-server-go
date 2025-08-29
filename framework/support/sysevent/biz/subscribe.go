package biz

import (
	"context"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/mqevent/manager"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/dtm-labs/rockscache"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/dto"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/event_err"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/model"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/service"
)

type ISubscribeRepo interface {
	baserepo.IBaseRepo[model.Subscribe, string]
	GetByStatus(ctx context.Context, status int8) ([]*model.Subscribe, error)
	GetAll(ctx context.Context) ([]*model.Subscribe, error)
	GetByTopicAndGroup(ctx context.Context, topic, group string) (*model.Subscribe, error)
}
type ISubscribeParameterRepo interface {
	baserepo.IBaseRepo[model.SubscribeParameter, string]
	GetBySubscribe(ctx context.Context, subscribe string) ([]*model.SubscribeParameter, error)
	DeleteBySubscribe(ctx context.Context, subscribe string) error
	DeleteByIds(ctx context.Context, ids ...string) error
}

type SubscribeUseCase struct {
	repo  ISubscribeRepo
	par   ISubscribeParameterRepo
	drepo IDeadLetterSubscribeRepo
	db    database.IDataBase
	sm    manager.EventManager
	rc    *rockscache.Client
}

func NewSubscribeUseCase(repo ISubscribeRepo, db database.IDataBase, par ISubscribeParameterRepo, sm manager.EventManager,
	drepo IDeadLetterSubscribeRepo, rc *rockscache.Client) service.ISubscribeServerApi {
	return &SubscribeUseCase{repo: repo, db: db, par: par, sm: sm, drepo: drepo, rc: rc}
}

func (s SubscribeUseCase) Add(ctx context.Context, req *dto.AddSubscribeReq) herrors.Herr {
	vo := &model.Subscribe{
		Id:    s.db.GenStringId(),
		Name:  req.Name,
		Topic: req.Topic,
		Group: req.Group,
		Dis:   req.Dis,
	}
	subscribe, err1 := s.repo.GetByTopicAndGroup(ctx, req.Topic, req.Topic)
	if err1 != nil && !database.IfErrorNotFound(err1) {
		hlog.CtxErrorf(ctx, "add subscribe GetByTopicAndGroup err:%v", err1)
		return event_err.SubscriptionEventFail(err1)
	}
	if subscribe != nil && subscribe.Id != "" {
		return event_err.TheSameSubscriptionAlreadyExists
	}

	if err := s.db.InTx(ctx, func(ctx1 context.Context) error {
		if _, err := s.repo.Add(ctx1, vo); err != nil {
			hlog.CtxErrorf(ctx, "add subscribe template tx err: %v", err)
			return err
		}
		return s.BathAddParameter(ctx1, req.Parameter, vo.Id)
	}); err != nil {
		hlog.CtxErrorf(ctx, "add subscribe failed: %v", err)
		return event_err.AddSubscribeFail(err)
	}
	return nil
}

func (s SubscribeUseCase) GetDetails(ctx context.Context, id string) (*dto.SubscribeModel, herrors.Herr) {
	en, err := s.repo.FindById(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "find subscribe failed: %v", err)
		return nil, event_err.GetSubscribeFail(err)
	}
	// 获取属性
	parameters, err := s.par.GetBySubscribe(ctx, en.Id)
	if err != nil {
		hlog.CtxErrorf(ctx, "get subscribe parameter err: %v", err)
		return nil, event_err.GetSubscribeFail(err)
	}
	return dto.SubscribeToDtoWithParams(en, parameters), nil
}

func (s SubscribeUseCase) UpdateStatus(ctx context.Context, id string, status int32) herrors.Herr {
	su, err2 := s.repo.FindById(ctx, id)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "find subscribe failed: %v", err2)
		return event_err.EditSubscribeFail(err2)
	}
	if su.Status == int8(status) {
		return nil
	}
	su.Status = int8(status)
	err := s.repo.EditById(ctx, su)
	if err != nil {
		hlog.CtxErrorf(ctx, "edit subscribe failed: %v", err)
		return event_err.EditSubscribeFail(err)
	}
	return nil
}

func (s SubscribeUseCase) Update(ctx context.Context, req *dto.UpdateSubscribeReq) herrors.Herr {
	sub, err := s.repo.FindById(ctx, req.Id)
	if err != nil {
		hlog.CtxErrorf(ctx, "find subscribe failed: %v", err)
		return event_err.EditSubscribeFail(err)
	}
	if req.Name != "" {
		sub.Name = req.Name
	}
	if req.Topic != "" {
		sub.Topic = req.Topic
	}
	if req.Group != "" {
		sub.Group = req.Group
	}
	if req.Dis != "" {
		sub.Dis = req.Dis
	}
	if err = s.db.InTx(ctx, func(ctx1 context.Context) error {
		if err = s.repo.EditById(ctx, sub); err != nil {
			hlog.CtxErrorf(ctx, "edit subscribe fiels tx err: %v", err)
			return err
		}
		err = s.BathUpdateParameter(ctx1, req.Parameter, req.Id)
		if err != nil {
			hlog.CtxErrorf(ctx, "bath update parameter err: %v", err)
			return err
		}
		return s.DelSubordinatesParameterCache(ctx, sub.Topic, sub.Group)
	}); err != nil {
		hlog.CtxErrorf(ctx, "edit subscribe failed: %v", err)
		return event_err.EditSubscribeFail(err)
	}
	return nil
}

// DelSubordinatesParameterCache 根据群删除群成员缓存
func (s SubscribeUseCase) DelSubordinatesParameterCache(ctx context.Context, topic, group string) error {
	key := model.GetSubordinatesParameterCacheKey(topic, group, actx.GetTenantId(ctx))
	return s.rc.TagAsDeleted2(ctx, key)
}

func (s SubscribeUseCase) GetList(ctx context.Context, req *dto.GetSubscribeListReq) (models.PageRes[dto.SubscribeModel], herrors.Herr) {
	qb := db_query.NewQueryBuilder()
	if req.Topic != "" {
		qb.Where("topic", db_query.Eq, req.Topic)
	}
	if req.Status != 0 {
		qb.Where("status", db_query.Eq, req.Status)
	}
	if req.Group != "" {
		qb.Where("group_name", db_query.Eq, req.Group)
	}
	if req.Name != "" {
		qb.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	qb.OrderBy("status", true)
	qb.OrderBy("created_at", false)
	qb.WithPage(&req.Page)

	res := models.PageRes[dto.SubscribeModel]{}
	// 查询总数
	total, err := s.repo.Count(ctx, qb)
	if err != nil {
		return res, herrors.QueryFail(err)
	}

	// 查询数据
	subscribes, err := s.repo.Find(ctx, qb)
	if err != nil {
		return res, herrors.QueryFail(err)
	}

	dtos := make([]*dto.SubscribeModel, 0, len(subscribes))
	for _, subscribe := range subscribes {
		parameters, _ := s.par.GetBySubscribe(ctx, subscribe.Id)
		if dto := dto.SubscribeToDtoWithParams(subscribe, parameters); dto != nil {
			dtos = append(dtos, dto)
		}
	}

	res.Total = total
	res.List = dtos
	return res, nil
}

func (s SubscribeUseCase) Enable(ctx context.Context, id string, ignoringHistory int32) herrors.Herr {
	sub, err := s.repo.FindById(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "find subscribe list failed: %v", err)
		return event_err.GetSubscribeFail(err)
	}
	sub.Status = int8(manager.StatusEnable)
	sub.Start = utils.GetDateUnix()
	if err = s.db.InTx(ctx, func(ctx context.Context) error {
		if err := s.repo.EditById(ctx, sub); err != nil {
			hlog.CtxErrorf(ctx, "update subscribe status err: %v", err)
			return fmt.Errorf("update subscribe status err: %v", err)
		}
		if hr := s.sm.ActivateSubscription(sub.Topic, sub.Group); herrors.HaveError(hr) {
			hlog.CtxErrorf(ctx, "start subscribe  err: %v", hr)
			return fmt.Errorf("start subscribe  err: %v", hr)
		}
		return nil
	}); err != nil {
		hlog.CtxErrorf(ctx, "Enable subscribe err: %v", err)
		return event_err.EditSubscribeFail(err)
	}
	return nil
}

func (s SubscribeUseCase) Disable(ctx context.Context, id string) herrors.Herr {
	sub, err := s.repo.FindById(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "find subscribe list failed: %v", err)
		return event_err.GetSubscribeFail(err)
	}
	sub.Status = int8(manager.StatusDisable)
	sub.Start = utils.GetDateUnix()
	if err = s.db.InTx(ctx, func(ctx context.Context) error {
		if err = s.repo.EditById(ctx, sub); err != nil {
			hlog.CtxErrorf(ctx, "update subscribe status err: %v", err)
			return fmt.Errorf("update subscribe status err: %v", err)
		}
		if err = s.sm.DeactivateSubscription(sub.Topic, sub.Group); err != nil {
			hlog.CtxErrorf(ctx, "stop subscribe  err: %v", err)
			return fmt.Errorf("stop subscribe  err: %v", err)
		}
		return nil
	}); err != nil {
		hlog.CtxErrorf(ctx, "Enable subscribe err: %v", err)
		return event_err.EditSubscribeFail(err)
	}
	return nil
}
func (s SubscribeUseCase) GetSubordinatesPar(ctx context.Context, topic, group string) (map[string]interface{}, error) {
	subscribe, err := s.repo.GetByTopicAndGroup(ctx, topic, group)
	if err != nil {
		hlog.CtxErrorf(ctx, "get subscribe list failed: %v", err)
		return nil, err
	}
	if subscribe == nil {
		return nil, fmt.Errorf("subscribe is nil")
	}
	parameters, err := s.par.GetBySubscribe(ctx, subscribe.Id)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return make(map[string]interface{}), nil
		}
		hlog.CtxErrorf(ctx, "get subscribe list failed: %v", err)
		return nil, err
	}
	if parameters == nil || len(parameters) == 0 {
		return make(map[string]interface{}), nil
	}
	parametersToMap, err := model.ConvertParametersToMap(parameters)
	if err != nil {
		hlog.CtxErrorf(ctx, "convert parameters to map failed: %v", err)
		return nil, err
	}
	return parametersToMap, nil
}

func (s SubscribeUseCase) GetByStatus(ctx context.Context, status int32) ([]*dto.SubscribeModel, herrors.Herr) {
	subscribes, err := s.repo.GetByStatus(ctx, int8(status))
	if err != nil {
		if database.IfErrorNotFound(err) {
			return make([]*dto.SubscribeModel, 0), nil
		}
		hlog.CtxErrorf(ctx, "get subscribe list failed: %v", err)
		return nil, event_err.GetSubscribeFail(err)
	}

	dtos := make([]*dto.SubscribeModel, 0, len(subscribes))
	for _, subscribe := range subscribes {
		parameters, _ := s.par.GetBySubscribe(ctx, subscribe.Id)
		if dto := dto.SubscribeToDtoWithParams(subscribe, parameters); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos, nil
}

func (s SubscribeUseCase) GetSubordinatesParameter(ctx context.Context, topic, group string) ([]*dto.SubscribeParameterModel, error) {
	subscribe, err := s.repo.GetByTopicAndGroup(ctx, topic, group)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return make([]*dto.SubscribeParameterModel, 0), nil
		}
		hlog.CtxErrorf(ctx, "get subscribe list failed: %v", err)
		return nil, err
	}
	parameters, err := s.par.GetBySubscribe(ctx, subscribe.Id)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return make([]*dto.SubscribeParameterModel, 0), nil
		}
		hlog.CtxErrorf(ctx, "get subscribe list failed: %v", err)
		return nil, err
	}
	return dto.SubscribeParametersToDtos(parameters), nil
}

func (s SubscribeUseCase) BathAddParameter(ctx context.Context, data []*dto.SubscribeParameterModel, subscribeId string) error {
	if len(data) == 0 {
		return nil
	}
	pars := make([]*model.SubscribeParameter, 0, len(data))
	for _, item := range data {
		id := item.Id
		if id == "" {
			id = s.db.GenStringId()
		}
		pars = append(pars, &model.SubscribeParameter{
			Id:          id,
			Key:         item.Key,
			DataType:    item.DataType,
			Value:       item.Value,
			Dis:         item.Dis,
			SubscribeId: subscribeId,
		})
	}
	if len(pars) > 0 {
		if err := s.par.BathAdd(ctx, pars...); err != nil {
			hlog.CtxErrorf(ctx, "BathAddParameter bath add parameter err: %v", err)
			return err
		}
	}
	return nil
}

func (s SubscribeUseCase) BathUpdateParameter(ctx context.Context, data []*dto.SubscribeParameterModel, subscribeId string) error {
	if len(data) == 0 {
		return nil
	}
	parameters := s.buildModelParameter(ctx, data, subscribeId)
	//获取现有的
	oldParameters, err := s.par.GetBySubscribe(ctx, subscribeId)
	if err != nil {
		hlog.CtxErrorf(ctx, "get subscribe parameter template err: %v", err)
		return err
	}
	// 处理新增修改和删除
	add, update, toDelete := s.processParameter(oldParameters, parameters)
	if len(add) > 0 {
		if err = s.par.BathAdd(ctx, add...); err != nil {
			hlog.CtxErrorf(ctx, "BathUpdateParameter bath add parameter err: %v", err)
			return err
		}
	}
	if len(update) > 0 {
		//处理更新
		for _, attribute := range update {
			if err = s.par.EditById(ctx, attribute); err != nil {
				hlog.CtxErrorf(ctx, "BathUpdateParameter update parameter err: %v", err)
				return err
			}
		}
	}
	if len(toDelete) > 0 {
		//处理删除
		attributesIds := make([]string, 0, len(toDelete))
		for _, item := range toDelete {
			attributesIds = append(attributesIds, item.Id)
		}
		//删除属性
		err = s.par.DeleteByIds(ctx, attributesIds...)
		if err != nil {
			hlog.CtxErrorf(ctx, "BathUpdateParameter delete subscribe parameter err: %v", err)
			return err
		}
	}
	return nil
}

// processParameter 处理参数
func (s SubscribeUseCase) processParameter(oldAttributes, newAttributes []*model.SubscribeParameter) (toAdd, toUpdate, toDelete []*model.SubscribeParameter) {
	// 当 oldAttributes 为空时，全部新增
	if len(oldAttributes) == 0 {
		toAdd = newAttributes
		return toAdd, toUpdate, toDelete
	}

	// 当 newOptions 为空时，全部删除
	if len(newAttributes) == 0 {
		toDelete = oldAttributes
		return toAdd, toUpdate, toDelete
	}
	// 将旧的数据放入map，方便查找
	oldMap := make(map[string]*model.SubscribeParameter)
	for _, attr := range oldAttributes {
		oldMap[attr.Id] = attr
	}
	// 将新数据逐个检查，找到新增或修改的数据
	newMap := make(map[string]*model.SubscribeParameter)
	for _, attr := range newAttributes {
		newMap[attr.Id] = attr

		if oldAttr, exists := oldMap[attr.Id]; exists {
			// 存在于旧的数据中，检查是否需要更新
			if oldAttr != attr { // 如果两个数据不一样，则表示需要更新
				toUpdate = append(toUpdate, attr)
			}
		} else {
			// 不存在于旧的数据中，表示需要新增
			toAdd = append(toAdd, attr)
		}
	}
	// 查找旧数据中存在但新数据中不存在的数据，表示需要删除
	for _, attr := range oldAttributes {
		if _, exists := newMap[attr.Id]; !exists {
			toDelete = append(toDelete, attr)
		}
	}
	return toAdd, toUpdate, toDelete
}

// buildModelParameter 将请求数据构建成模型参数
func (s SubscribeUseCase) buildModelParameter(_ context.Context, data []*dto.SubscribeParameterModel, subscribeId string) []*model.SubscribeParameter {
	if len(data) == 0 {
		return nil
	}
	attributes := make([]*model.SubscribeParameter, 0)
	for _, item := range data {
		id := item.Id
		if id == "" {
			id = s.db.GenStringId()
		}
		attributes = append(attributes, &model.SubscribeParameter{
			Id:       id,
			Key:      item.Key,
			DataType: item.DataType,
			Value:    item.Value,
			Dis:      item.Dis,
			SubscribeId: func() string {
				if subscribeId != "" {
					return subscribeId
				}
				return item.SubscribeId
			}(),
		})
	}
	return attributes
}
