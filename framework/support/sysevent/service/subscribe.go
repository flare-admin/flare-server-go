package service

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/dto"
)

type ISubscribeServerApi interface {
	//Add  新增
	Add(ctx context.Context, req *dto.AddSubscribeReq) herrors.Herr
	//GetDetails 获取事件详情
	GetDetails(ctx context.Context, id string) (*dto.SubscribeModel, herrors.Herr)
	//UpdateStatus 更新事件状态
	UpdateStatus(ctx context.Context, id string, status int32) herrors.Herr
	//Update 更新事件
	Update(ctx context.Context, req *dto.UpdateSubscribeReq) herrors.Herr
	//GetList 获取事件
	GetList(ctx context.Context, req *dto.GetSubscribeListReq) (models.PageRes[dto.SubscribeModel], herrors.Herr)
	//Enable 开启订阅
	Enable(ctx context.Context, id string, ignoringHistory int32) herrors.Herr
	// Disable 关闭订阅
	Disable(ctx context.Context, id string) herrors.Herr
	// GetSubordinatesPar 获取订阅的参数
	GetSubordinatesPar(ctx context.Context, topic, group string) (map[string]interface{}, error)
	//GetByStatus 根据状态获取订阅
	GetByStatus(ctx context.Context, status int32) ([]*dto.SubscribeModel, herrors.Herr)
	// GetSubordinatesParameter 获取订阅的参数
	GetSubordinatesParameter(ctx context.Context, topic, group string) ([]*dto.SubscribeParameterModel, error)
}
