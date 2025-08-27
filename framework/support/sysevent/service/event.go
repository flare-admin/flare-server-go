package service

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/dto"
)

type IEventServerApi interface {
	// Add 新增
	Add(ctx context.Context, req *dto.AddEventReq) herrors.Herr
	// GetDetails 获取事件详情
	GetDetails(ctx context.Context, id string) (*dto.EventModel, herrors.Herr)
	// UpdateStatus 更新事件状态
	UpdateStatus(ctx context.Context, id string, status int32) herrors.Herr
	// Update 更新事件
	Update(ctx context.Context, req *dto.UpdateEventReq) herrors.Herr
	// GetList 获取事件
	GetList(ctx context.Context, req *dto.GetEventListReq) (models.PageRes[dto.EventModel], herrors.Herr)
	// GetAllByStatusList 根据状态获取所有
	GetAllByStatusList(ctx context.Context, status int32) ([]*dto.EventModel, herrors.Herr)
	// GetByTopic 根据主题获取事件
	GetByTopic(ctx context.Context, topic string) (*dto.EventModel, herrors.Herr)
}
