package repository

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/model"
)

// CategoryRepository 分类仓储接口
type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Category, error)
	FindByCode(ctx context.Context, code string) (*model.Category, error)
	FindAll(ctx context.Context) ([]*model.Category, error)
}
