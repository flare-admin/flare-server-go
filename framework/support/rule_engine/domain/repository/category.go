package repository

import (
	"context"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
)

// ICategoryRepository 规则分类仓储接口
type ICategoryRepository interface {
	// Create 创建分类
	Create(ctx context.Context, category *model.RuleCategory) error

	// Update 更新分类
	Update(ctx context.Context, category *model.RuleCategory) error

	// Delete 删除分类
	Delete(ctx context.Context, id string) error

	// FindByID 根据ID查找分类
	FindByID(ctx context.Context, id string) (*model.RuleCategory, error)

	// FindByCode 根据编码查找分类
	FindByCode(ctx context.Context, code string) (*model.RuleCategory, error)

	// FindByParentID 根据父分类ID查找子分类列表
	FindByParentID(ctx context.Context, parentID string) ([]*model.RuleCategory, error)

	// FindRootCategories 查找根分类列表
	FindRootCategories(ctx context.Context) ([]*model.RuleCategory, error)

	// FindByType 根据类型查找分类列表
	FindByType(ctx context.Context, categoryType string) ([]*model.RuleCategory, error)

	// FindByBusinessType 根据业务类型查找分类列表
	FindByBusinessType(ctx context.Context, businessType string) ([]*model.RuleCategory, error)

	// FindEnabledByParentID 根据父分类ID查找启用的子分类列表
	FindEnabledByParentID(ctx context.Context, parentID string) ([]*model.RuleCategory, error)

	// FindEnabledRootCategories 查找启用的根分类列表
	FindEnabledRootCategories(ctx context.Context) ([]*model.RuleCategory, error)

	// FindEnabledByType 根据类型查找启用的分类列表
	FindEnabledByType(ctx context.Context, categoryType string) ([]*model.RuleCategory, error)

	// FindEnabledByBusinessType 根据业务类型查找启用的分类列表
	FindEnabledByBusinessType(ctx context.Context, businessType string) ([]*model.RuleCategory, error)

	// FindDescendants 查找指定分类的所有后代分类
	FindDescendants(ctx context.Context, path string) ([]*model.RuleCategory, error)

	// ExistsByCode 检查编码是否存在
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// Find 根据查询条件查找分类列表
	Find(ctx context.Context, query *db_query.QueryBuilder) ([]*model.RuleCategory, error)

	// Count 根据查询条件统计分类数量
	Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error)

	// FindAll 查找所有分类
	FindAll(ctx context.Context) ([]*model.RuleCategory, error)
}
