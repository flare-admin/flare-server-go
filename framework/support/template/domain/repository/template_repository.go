package repository

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/model"
)

// TemplateRepository 模板仓储接口
type TemplateRepository interface {
	Create(ctx context.Context, template *model.Template) error
	Update(ctx context.Context, template *model.Template) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Template, error)
	FindByCategoryID(ctx context.Context, categoryID string) ([]*model.Template, error)
	FindByCode(ctx context.Context, code string) (*model.Template, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	FindAll(ctx context.Context) ([]*model.Template, error)
}
