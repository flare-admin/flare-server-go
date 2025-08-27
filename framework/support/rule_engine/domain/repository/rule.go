package repository

import (
	"context"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
)

// IRuleRepository 规则仓储接口
type IRuleRepository interface {
	// Create 创建规则
	Create(ctx context.Context, rule *model.Rule) error

	// Update 更新规则
	Update(ctx context.Context, rule *model.Rule) error

	// Delete 删除规则
	Delete(ctx context.Context, id string) error

	// FindByID 根据ID查找规则
	FindByID(ctx context.Context, id string) (*model.Rule, error)

	// FindByCode 根据编码查找规则
	FindByCode(ctx context.Context, code string) (*model.Rule, error)

	// FindByTemplateID 根据模板ID查找规则列表
	FindByTemplateID(ctx context.Context, templateID string) ([]*model.Rule, error)

	// FindByCategoryID 根据分类ID查找规则列表
	FindByCategoryID(ctx context.Context, categoryID string) ([]*model.Rule, error)

	// FindByType 根据类型查找规则列表
	FindByType(ctx context.Context, ruleType string) ([]*model.Rule, error)

	// FindByTrigger 根据触发条件查找规则列表
	FindByTrigger(ctx context.Context, trigger string) ([]*model.Rule, error)

	// FindByScope 根据作用域查找规则列表
	FindByScope(ctx context.Context, scope string) ([]*model.Rule, error)

	// FindByBusinessType 根据业务类型查找规则列表
	FindByBusinessType(ctx context.Context, businessType string) ([]*model.Rule, error)

	// ExistsByCode 检查编码是否存在
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// Find 根据查询条件查找规则列表
	Find(ctx context.Context, query *db_query.QueryBuilder) ([]*model.Rule, error)

	// Count 根据查询条件统计规则数量
	Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error)

	// FindAll 查找所有规则
	FindAll(ctx context.Context) ([]*model.Rule, error)
}
