package data

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/infrastructure/persistence/entity"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/infrastructure/repository"
)

// ruleRepository 规则数据访问层
type ruleRepository struct {
	*baserepo.BaseRepo[entity.Rule, string]
}

// NewRuleRepository 创建规则数据访问层
func NewRuleRepository(data database.IDataBase) repository.IRuleRepository {
	// 同步表
	tables := []interface{}{
		&entity.Rule{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync rule tables error: %v", err)
	}
	return &ruleRepository{
		BaseRepo: baserepo.NewBaseRepo[entity.Rule, string](data),
	}
}

// FindByCode 根据编码查询规则
func (r *ruleRepository) FindByCode(ctx context.Context, code string) (*entity.Rule, error) {
	var rule entity.Rule
	err := r.Db(ctx).Model(&entity.Rule{}).Where("code = ?", code).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// FindByCategoryID 根据分类ID查询规则列表
func (r *ruleRepository) FindByCategoryID(ctx context.Context, categoryID string) ([]*entity.Rule, error) {
	rules := make([]*entity.Rule, 0)
	err := r.Db(ctx).Model(&entity.Rule{}).Where("category_id = ?", categoryID).Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// FindByTemplateID 根据模板ID查询规则列表
func (r *ruleRepository) FindByTemplateID(ctx context.Context, templateID string) ([]*entity.Rule, error) {
	rules := make([]*entity.Rule, 0)
	err := r.Db(ctx).Model(&entity.Rule{}).Where("template_id = ?", templateID).Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// FindByTrigger 根据触发动作查询规则列表
func (r *ruleRepository) FindByTrigger(ctx context.Context, trigger string) ([]*entity.Rule, error) {
	rules := make([]*entity.Rule, 0)
	err := r.Db(ctx).Model(&entity.Rule{}).
		Where("JSON_CONTAINS(triggers, ?)", trigger).
		Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// FindByScope 根据作用域查询规则列表
func (r *ruleRepository) FindByScope(ctx context.Context, scope string) ([]*entity.Rule, error) {
	rules := make([]*entity.Rule, 0)
	query := r.Db(ctx).Model(&entity.Rule{})

	if scope == "global" {
		query = query.Where("scope = ?", scope)
	} else {
		query = query.Where("(scope = ? OR scope = ?)", "global", scope)
	}

	err := query.Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// FindByType 根据类型查询规则列表
func (r *ruleRepository) FindByType(ctx context.Context, ruleType string) ([]*entity.Rule, error) {
	rules := make([]*entity.Rule, 0)
	err := r.Db(ctx).Model(&entity.Rule{}).Where("type = ?", ruleType).Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// ExistsByCode 检查编码是否存在
func (r *ruleRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.Db(ctx).Model(&entity.Rule{}).Where("code = ?", code).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindAll 查询所有规则
func (r *ruleRepository) FindAll(ctx context.Context) ([]*entity.Rule, error) {
	rules := make([]*entity.Rule, 0)
	err := r.Db(ctx).Model(&entity.Rule{}).Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// Find 根据查询条件查找规则列表
func (r *ruleRepository) Find(ctx context.Context, query *db_query.QueryBuilder) ([]*entity.Rule, error) {
	rules := make([]*entity.Rule, 0)
	dbQuery := r.Db(ctx).Model(&entity.Rule{})

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

	err := dbQuery.Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// Count 根据查询条件统计规则数量
func (r *ruleRepository) Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error) {
	var count int64
	dbQuery := r.Db(ctx).Model(&entity.Rule{})

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

// UpdateExecuteStats 更新执行统计
func (r *ruleRepository) UpdateExecuteStats(ctx context.Context, ruleID string, success bool) error {
	rule, err := r.FindById(ctx, ruleID)
	if err != nil {
		return err
	}

	rule.ExecuteCount++
	if success {
		rule.SuccessCount++
	}
	rule.LastExecuteAt = utils.GetDateUnix()

	return r.EditById(ctx, rule)
}

// FindByBusinessType 根据业务类型查询规则列表
func (r *ruleRepository) FindByBusinessType(ctx context.Context, businessType string) ([]*entity.Rule, error) {
	rules := make([]*entity.Rule, 0)
	err := r.Db(ctx).Model(&entity.Rule{}).Where("business_type = ?", businessType).Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}
