package biz

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/pkg/mqevent/manager"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/dto"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/model"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/service"
)

type IDeadLetterSubscribeRepo interface {
	baserepo.IBaseRepo[model.DeadLetterSubscribe, string]
	GetBy(ctx context.Context, topic, channel, messageId string) (*model.DeadLetterSubscribe, error)
	UpdateStatus(ctx context.Context, id string, status int8) error
}

type DeadLetterSubscribeUseCase struct {
	repo IDeadLetterSubscribeRepo
	par  ISubscribeRepo
	em   manager.EventManager
	db   database.IDataBase
}

func NewDeadLetterSubscribeUseCase(repo IDeadLetterSubscribeRepo, par ISubscribeRepo, em manager.EventManager, db database.IDataBase) service.IDeadLetterServiceApi {
	return &DeadLetterSubscribeUseCase{repo: repo, par: par, em: em, db: db}
}
func (d *DeadLetterSubscribeUseCase) Retry(ctx context.Context, id string) error {
	dead, err := d.repo.FindById(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "find subscribe list failed: %v", err)
		return err
	}
	event, err := dead.ToDeadLetterEvent()
	if err != nil {
		hlog.CtxErrorf(ctx, "convert dead letter event failed: %v", err)
		return err
	}
	return d.em.RetryDeadLetter(ctx, event)
}

func (d *DeadLetterSubscribeUseCase) GetList(ctx context.Context, req *dto.GetDeadLetterSubscribeListReq) (models.PageRes[dto.DeadLetterSubscribeModel], error) {
	qb := db_query.NewQueryBuilder()
	if req.Topic != "" {
		qb.Where("topic", db_query.Eq, req.Topic)
	}
	if req.Group != "" {
		qb.Where("group", db_query.Eq, req.Group)
	}
	if req.Name != "" {
		qb.Where("name", db_query.Eq, req.Name)
	}
	qb.OrderBy("status", true)
	qb.OrderBy("created_at", false)
	res := models.PageRes[dto.DeadLetterSubscribeModel]{}
	// 查询总数
	total, err := d.repo.Count(ctx, qb)
	if err != nil {
		return res, herrors.QueryFail(err)
	}

	// 查询数据
	tasks, err := d.repo.Find(ctx, qb)
	if err != nil {
		return res, herrors.QueryFail(err)
	}
	tds := make([]*dto.DeadLetterSubscribeModel, 0)
	for _, task := range tasks {
		tds = append(tds, dto.DeadLetterToDto(task))
	}
	res.Total = total
	res.List = tds
	return res, nil
}
