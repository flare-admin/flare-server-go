package repository

import (
	"context"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
)

// ITemplateRepository 规则模板仓储接口
type ITemplateRepository interface {
	// Create 创建模板
	Create(ctx context.Context, template *model.RuleTemplate) error

	// Update 更新模板
	Update(ctx context.Context, template *model.RuleTemplate) error

	// Delete 删除模板
	Delete(ctx context.Context, id string) error

	// FindByID 根据ID查找模板
	FindByID(ctx context.Context, id string) (*model.RuleTemplate, error)

	// FindByCode 根据编码查找模板
	FindByCode(ctx context.Context, code string) (*model.RuleTemplate, error)

	// FindByCategoryID 根据分类ID查找模板列表
	FindByCategoryID(ctx context.Context, categoryID string) ([]*model.RuleTemplate, error)

	// FindByType 根据类型查找模板列表
	FindByType(ctx context.Context, templateType string) ([]*model.RuleTemplate, error)

	// FindEnabledByCategoryID 根据分类ID查找启用的模板列表
	FindEnabledByCategoryID(ctx context.Context, categoryID string) ([]*model.RuleTemplate, error)

	// FindEnabledByType 根据类型查找启用的模板列表
	FindEnabledByType(ctx context.Context, templateType string) ([]*model.RuleTemplate, error)

	// FindByBusinessType 根据业务类型查找模板列表
	FindByBusinessType(ctx context.Context, businessType string) ([]*model.RuleTemplate, error)

	// FindEnabledByBusinessType 根据业务类型查找启用的模板列表
	FindEnabledByBusinessType(ctx context.Context, businessType string) ([]*model.RuleTemplate, error)

	// FindByScope 根据作用域查找模板列表
	FindByScope(ctx context.Context, scope string) ([]*model.RuleTemplate, error)

	// FindEnabledByScope 根据作用域查找启用的模板列表
	FindEnabledByScope(ctx context.Context, scope string) ([]*model.RuleTemplate, error)

	// ExistsByCode 检查编码是否存在
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// Find 根据查询条件查找模板列表
	Find(ctx context.Context, query *db_query.QueryBuilder) ([]*model.RuleTemplate, error)

	// Count 根据查询条件统计模板数量
	Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error)

	// FindAll 查找所有模板
	FindAll(ctx context.Context) ([]*model.RuleTemplate, error)
}
