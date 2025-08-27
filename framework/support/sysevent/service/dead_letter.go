package service

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/dto"
)

type IDeadLetterServiceApi interface {
	// GetList 获取列表
	GetList(ctx context.Context, req *dto.GetDeadLetterSubscribeListReq) (models.PageRes[dto.DeadLetterSubscribeModel], error)
	// Retry 重试
	Retry(ctx context.Context, id string) error
}
