package data

import (
	"context"
	"gorm.io/gorm"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/infrastructure/persistence/entity"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/infrastructure/repository"
)

// ruleCategoryRepository 规则分类数据访问层
type ruleCategoryRepository struct {
	*baserepo.BaseRepo[entity.RuleCategory, string]
}

// NewRuleCategoryRepository 创建规则分类数据访问层
func NewRuleCategoryRepository(data database.IDataBase) repository.IRuleCategoryRepository {
	// 同步表
	tables := []interface{}{
		&entity.RuleCategory{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync rule category tables error: %v", err)
	}
	return &ruleCategoryRepository{
		BaseRepo: baserepo.NewBaseRepo[entity.RuleCategory, string](data, entity.RuleCategory{}),
	}
}

// FindByCode 根据编码查询分类
func (r *ruleCategoryRepository) FindByCode(ctx context.Context, code string) (*entity.RuleCategory, error) {
	var category entity.RuleCategory
	err := r.Db(ctx).Model(&entity.RuleCategory{}).Where("code = ?", code).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindByParentID 根据父分类ID查询子分类列表
func (r *ruleCategoryRepository) FindByParentID(ctx context.Context, parentID string) ([]*entity.RuleCategory, error) {
	categories := make([]*entity.RuleCategory, 0)
	err := r.Db(ctx).Model(&entity.RuleCategory{}).Where("parent_id = ?", parentID).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByType 根据类型查询分类列表
func (r *ruleCategoryRepository) FindByType(ctx context.Context, categoryType string) ([]*entity.RuleCategory, error) {
	categories := make([]*entity.RuleCategory, 0)
	err := r.Db(ctx).Model(&entity.RuleCategory{}).Scopes(func(db *gorm.DB) *gorm.DB {
		if categoryType != "" {
			db.Where("type = ?", categoryType)
		}
		return db
	}).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByBusinessType 根据业务类型查询分类列表
func (r *ruleCategoryRepository) FindByBusinessType(ctx context.Context, businessType string) ([]*entity.RuleCategory, error) {
	categories := make([]*entity.RuleCategory, 0)
	err := r.Db(ctx).Model(&entity.RuleCategory{}).Where("business_type = ?", businessType).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByScope 根据作用域查询分类列表
func (r *ruleCategoryRepository) FindByScope(ctx context.Context, scope string) ([]*entity.RuleCategory, error) {
	categories := make([]*entity.RuleCategory, 0)
	query := r.Db(ctx).Model(&entity.RuleCategory{})

	if scope == "global" {
		query = query.Where("scope = ?", scope)
	} else {
		query = query.Where("(scope = ? OR scope = ?)", "global", scope)
	}

	err := query.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// ExistsByCode 检查编码是否存在
func (r *ruleCategoryRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.Db(ctx).Model(&entity.RuleCategory{}).Where("code = ?", code).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll 查询所有分类
func (r *ruleCategoryRepository) FindAll(ctx context.Context) ([]*entity.RuleCategory, error) {
	categories := make([]*entity.RuleCategory, 0)
	err := r.Db(ctx).Model(&entity.RuleCategory{}).Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Find 根据查询条件查找分类列表
func (r *ruleCategoryRepository) Find(ctx context.Context, query *db_query.QueryBuilder) ([]*entity.RuleCategory, error) {
	categories := make([]*entity.RuleCategory, 0)
	dbQuery := r.Db(ctx).Model(&entity.RuleCategory{})

	// 应用查询条件
	if query != nil {
		if where, values := query.BuildWhere(); where != "" {
			dbQuery = dbQuery.Where(where, values...)
		}
		if orderBy := query.BuildOrderBy(); orderBy != "" {
			dbQuery = dbQuery.Order(orderBy)
		}
		if limit, values := query.BuildLimit(); limit != "" {
			dbQuery = dbQuery.Offset(values[0]).Limit(values[1])
		}
	}

	err := dbQuery.Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// Count 根据查询条件统计分类数量
func (r *ruleCategoryRepository) Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error) {
	var count int64
	dbQuery := r.Db(ctx).Model(&entity.RuleCategory{})

	// 应用查询条件
	if query != nil {
		if where, values := query.BuildWhere(); where != "" {
			dbQuery = dbQuery.Where(where, values...)
		}
	}

	err := dbQuery.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FindRootCategories 查询根分类列表
func (r *ruleCategoryRepository) FindRootCategories(ctx context.Context) ([]*entity.RuleCategory, error) {
	categories := make([]*entity.RuleCategory, 0)
	err := r.Db(ctx).Model(&entity.RuleCategory{}).Where("parent_id = ?", "").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// FindByPath 根据路径查询分类
func (r *ruleCategoryRepository) FindByPath(ctx context.Context, path string) (*entity.RuleCategory, error) {
	var category entity.RuleCategory
	err := r.Db(ctx).Model(&entity.RuleCategory{}).Where("path = ?", path).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindDescendants 查询后代分类列表
func (r *ruleCategoryRepository) FindDescendants(ctx context.Context, path string) ([]*entity.RuleCategory, error) {
	categories := make([]*entity.RuleCategory, 0)
	err := r.Db(ctx).Model(&entity.RuleCategory{}).Where("path LIKE ?", path+"/%").Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}
