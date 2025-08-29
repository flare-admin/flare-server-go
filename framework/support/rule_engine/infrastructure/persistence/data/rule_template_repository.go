package data

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/infrastructure/persistence/entity"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/infrastructure/repository"
)

// ruleTemplateRepository 规则模板数据访问层
type ruleTemplateRepository struct {
	*baserepo.BaseRepo[entity.RuleTemplate, string]
}

// NewRuleTemplateRepository 创建规则模板数据访问层
func NewRuleTemplateRepository(data database.IDataBase) repository.IRuleTemplateRepository {
	// 同步表
	tables := []interface{}{
		&entity.RuleTemplate{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync rule template tables error: %v", err)
	}
	return &ruleTemplateRepository{
		BaseRepo: baserepo.NewBaseRepo[entity.RuleTemplate, string](data),
	}
}

// FindByCode 根据编码查询模板
func (r *ruleTemplateRepository) FindByCode(ctx context.Context, code string) (*entity.RuleTemplate, error) {
	var template entity.RuleTemplate
	err := r.Db(ctx).Model(&entity.RuleTemplate{}).Where("code = ?", code).First(&template).Error
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// FindByCategoryID 根据分类ID查询模板列表
func (r *ruleTemplateRepository) FindByCategoryID(ctx context.Context, categoryID string) ([]*entity.RuleTemplate, error) {
	templates := make([]*entity.RuleTemplate, 0)
	err := r.Db(ctx).Model(&entity.RuleTemplate{}).Where("category_id = ?", categoryID).Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// FindByType 根据类型查询模板列表
func (r *ruleTemplateRepository) FindByType(ctx context.Context, templateType string) ([]*entity.RuleTemplate, error) {
	templates := make([]*entity.RuleTemplate, 0)
	err := r.Db(ctx).Model(&entity.RuleTemplate{}).Where("type = ?", templateType).Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// FindByBusinessType 根据业务类型查询模板列表
func (r *ruleTemplateRepository) FindByBusinessType(ctx context.Context, businessType string) ([]*entity.RuleTemplate, error) {
	templates := make([]*entity.RuleTemplate, 0)
	err := r.Db(ctx).Model(&entity.RuleTemplate{}).Where("business_type = ?", businessType).Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// FindByScope 根据作用域查询模板列表
func (r *ruleTemplateRepository) FindByScope(ctx context.Context, scope string) ([]*entity.RuleTemplate, error) {
	templates := make([]*entity.RuleTemplate, 0)
	query := r.Db(ctx).Model(&entity.RuleTemplate{})

	if scope == "global" {
		query = query.Where("scope = ?", scope)
	} else {
		query = query.Where("(scope = ? OR scope = ?)", "global", scope)
	}

	err := query.Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// ExistsByCode 检查编码是否存在
func (r *ruleTemplateRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.Db(ctx).Model(&entity.RuleTemplate{}).Where("code = ?", code).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll 查询所有模板
func (r *ruleTemplateRepository) FindAll(ctx context.Context) ([]*entity.RuleTemplate, error) {
	templates := make([]*entity.RuleTemplate, 0)
	err := r.Db(ctx).Model(&entity.RuleTemplate{}).Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// Find 根据查询条件查找模板列表
func (r *ruleTemplateRepository) Find(ctx context.Context, query *db_query.QueryBuilder) ([]*entity.RuleTemplate, error) {
	templates := make([]*entity.RuleTemplate, 0)
	dbQuery := r.Db(ctx).Model(&entity.RuleTemplate{})

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

	err := dbQuery.Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// Count 根据查询条件统计模板数量
func (r *ruleTemplateRepository) Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error) {
	var count int64
	dbQuery := r.Db(ctx).Model(&entity.RuleTemplate{})

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
