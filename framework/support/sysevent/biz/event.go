package biz

import (
	"context"
	"github.com/dtm-labs/rockscache"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/dto"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/event_err"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/model"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/service"
)

type IEventRepo interface {
	baserepo.IBaseRepo[model.Event, string]
	GetByTopic(ctx context.Context, topic string) (*model.Event, error)
	GetAllByStatusList(ctx context.Context, status int32) ([]*model.Event, error)
}

type EventUseCase struct {
	repo IEventRepo
	db   database.IDataBase
	rc   *rockscache.Client
}

func NewEventUseCase(repo IEventRepo, db database.IDataBase, rc *rockscache.Client) service.IEventServerApi {
	return &EventUseCase{repo: repo, db: db, rc: rc}
}

func (e EventUseCase) Add(ctx context.Context, req *dto.AddEventReq) herrors.Herr {
	topic, err2 := e.repo.GetByTopic(ctx, req.Topic)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "get topic error: %v", err2)
		return event_err.GetEventFail(err2)
	}
	if topic != nil && topic.Topic != "" {
		return event_err.TopicIsExistFail
	}
	vo := e.toVo(req)
	_, err := e.repo.Add(ctx, vo)
	if err != nil {
		hlog.CtxErrorf(ctx, "add event failed: %v", err)
		return event_err.AddEventFail(err)
	}
	return nil
}

func (e EventUseCase) GetDetails(ctx context.Context, id string) (*dto.EventModel, herrors.Herr) {
	en, err := e.repo.FindById(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "find event failed: %v", err)
		return nil, event_err.GetEventFail(err)
	}
	return dto.EventToDto(en), nil
}

func (e EventUseCase) UpdateStatus(ctx context.Context, id string, status int32) herrors.Herr {
	ev, err2 := e.repo.FindById(ctx, id)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "find event failed: %v", err2)
		return event_err.GetEventFail(err2)
	}
	ev.Status = int8(status)
	err := e.repo.EditById(ctx, ev)
	if err != nil {
		hlog.CtxErrorf(ctx, "edit event failed: %v", err)
		return event_err.EditEventFail(err)
	}
	// 删除缓存
	e.DelEventCache(ctx, ev.Topic)
	return nil
}

func (e EventUseCase) Update(ctx context.Context, req *dto.UpdateEventReq) herrors.Herr {
	ev, err2 := e.repo.FindById(ctx, req.Id)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "find event failed: %v", err2)
		return event_err.GetEventFail(err2)
	}
	if req.Name != "" {
		ev.Name = req.Name
	}
	if req.Topic != "" {
		ev.Topic = req.Topic
	}
	if req.Dis != "" {
		ev.Dis = req.Dis
	}
	err := e.repo.EditById(ctx, ev)
	if err != nil {
		hlog.CtxErrorf(ctx, "edit event failed: %v", err)
		return event_err.EditEventFail(err)
	}
	// 删除缓存
	e.DelEventCache(ctx, ev.Topic)
	return nil
}

// DelEventCache 根据群删除群成员缓存
func (e EventUseCase) DelEventCache(ctx context.Context, topic string) error {
	key := model.GetEventCacheKey(topic)
	return e.rc.TagAsDeleted2(ctx, key)
}
func (e EventUseCase) GetList(ctx context.Context, req *dto.GetEventListReq) (models.PageRes[dto.EventModel], herrors.Herr) {
	qb := db_query.NewQueryBuilder()
	if req.Topic != "" {
		qb.Where("topic", db_query.Like, "%"+req.Topic+"%")
	}
	if req.Status != 0 {
		qb.Where("status", db_query.Eq, req.Status)
	}
	qb.OrderBy("created_at", false)
	qb.WithPage(&req.Page)
	res := models.PageRes[dto.EventModel]{}
	// 查询总数
	total, err := e.repo.Count(ctx, qb)
	if err != nil {
		return res, herrors.QueryFail(err)
	}

	// 查询数据
	tasks, err := e.repo.Find(ctx, qb)
	if err != nil {
		return res, herrors.QueryFail(err)
	}
	tds := make([]*dto.EventModel, 0)
	for _, task := range tasks {
		tds = append(tds, dto.EventToDto(task))
	}
	res.Total = total
	res.List = tds
	return res, nil
}

func (e EventUseCase) GetAllByStatusList(ctx context.Context, status int32) ([]*dto.EventModel, herrors.Herr) {
	list, err := e.repo.GetAllByStatusList(ctx, status)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to list news from db error: %v", err)
		return nil, event_err.GetEventFail(err)
	}
	dts := make([]*dto.EventModel, 0)
	for _, task := range list {
		dts = append(dts, dto.EventToDto(task))
	}
	return dts, nil
}
func (e EventUseCase) GetByTopic(ctx context.Context, topic string) (*dto.EventModel, herrors.Herr) {
	event, err := e.repo.GetByTopic(ctx, topic)
	if err != nil {
		if database.IfErrorNotFound(err) {
			return nil, event_err.EventNotExistFail
		}
		return nil, herrors.QueryFail(err)
	}
	return dto.EventToDto(event), nil
}
func (e EventUseCase) toVo(vo *dto.AddEventReq) *model.Event {
	return &model.Event{
		Id:    e.db.GenStringId(),
		Name:  vo.Name,
		Topic: vo.Topic,
		Dis:   vo.Dis,
	}
}
