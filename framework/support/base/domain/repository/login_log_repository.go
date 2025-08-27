package repository

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/model"
)

type ILoginLogRepository interface {
	Create(ctx context.Context, log *model.LoginLog) error
}
