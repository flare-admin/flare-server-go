package repository

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/repository"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/infrastructure/persistence/entity"
)

// IRuleRepository 规则数据访问接口
type IRuleRepository interface {
	baserepo.IBaseRepo[entity.Rule, string]
	FindByCode(ctx context.Context, code string) (*entity.Rule, error)
	FindByCategoryID(ctx context.Context, categoryID string) ([]*entity.Rule, error)
	FindByTemplateID(ctx context.Context, templateID string) ([]*entity.Rule, error)
	FindByTrigger(ctx context.Context, trigger string) ([]*entity.Rule, error)
	FindByScope(ctx context.Context, scope string) ([]*entity.Rule, error)
	FindByType(ctx context.Context, ruleType string) ([]*entity.Rule, error)
	FindByBusinessType(ctx context.Context, businessType string) ([]*entity.Rule, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	FindAll(ctx context.Context) ([]*entity.Rule, error)
	UpdateExecuteStats(ctx context.Context, ruleID string, success bool) error
	Find(ctx context.Context, query *db_query.QueryBuilder) ([]*entity.Rule, error)
	Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error)
}

// RuleRepository 规则仓储实现
type RuleRepository struct {
	repo IRuleRepository
}

// NewRuleRepository 创建规则仓储
func NewRuleRepository(repo IRuleRepository) repository.IRuleRepository {
	return &RuleRepository{
		repo: repo,
	}
}

// Create 创建规则
func (r *RuleRepository) Create(ctx context.Context, rule *model.Rule) error {
	// 转换为数据库实体
	entity := r.toEntity(rule)
	_, err := r.repo.Add(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

// Update 更新规则
func (r *RuleRepository) Update(ctx context.Context, rule *model.Rule) error {
	// 转换为数据库实体
	entity := r.toEntity(rule)
	return r.repo.EditById(ctx, entity)
}

// Delete 删除规则
func (r *RuleRepository) Delete(ctx context.Context, id string) error {
	return r.repo.DelByIdUnScoped(ctx, id)
}

// FindByID 根据ID查询规则
func (r *RuleRepository) FindByID(ctx context.Context, id string) (*model.Rule, error) {
	entity, err := r.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toModel(entity), nil
}

// FindByCode 根据编码查询规则
func (r *RuleRepository) FindByCode(ctx context.Context, code string) (*model.Rule, error) {
	entity, err := r.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return r.toModel(entity), nil
}

// FindByCategoryID 根据分类ID查询规则列表
func (r *RuleRepository) FindByCategoryID(ctx context.Context, categoryID string) ([]*model.Rule, error) {
	entities, err := r.repo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindByTemplateID 根据模板ID查询规则列表
func (r *RuleRepository) FindByTemplateID(ctx context.Context, templateID string) ([]*model.Rule, error) {
	entities, err := r.repo.FindByTemplateID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindByTrigger 根据触发动作查询规则列表
func (r *RuleRepository) FindByTrigger(ctx context.Context, trigger string) ([]*model.Rule, error) {
	entities, err := r.repo.FindByTrigger(ctx, trigger)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindByScope 根据作用域查询规则列表
func (r *RuleRepository) FindByScope(ctx context.Context, scope string) ([]*model.Rule, error) {
	entities, err := r.repo.FindByScope(ctx, scope)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return make([]*model.Rule, 0), nil
		}
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindByType 根据类型查询规则列表
func (r *RuleRepository) FindByType(ctx context.Context, ruleType string) ([]*model.Rule, error) {
	entities, err := r.repo.FindByType(ctx, ruleType)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindByBusinessType 根据业务类型查询规则列表
func (r *RuleRepository) FindByBusinessType(ctx context.Context, businessType string) ([]*model.Rule, error) {
	entities, err := r.repo.FindByBusinessType(ctx, businessType)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return make([]*model.Rule, 0), nil
		}
		return nil, err
	}
	return r.toModels(entities), nil
}

// ExistsByCode 检查编码是否存在
func (r *RuleRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return r.repo.ExistsByCode(ctx, code)
}

// FindAll 查询所有规则
func (r *RuleRepository) FindAll(ctx context.Context) ([]*model.Rule, error) {
	entities, err := r.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// Find 根据查询条件查找规则列表
func (r *RuleRepository) Find(ctx context.Context, query *db_query.QueryBuilder) ([]*model.Rule, error) {
	entities, err := r.repo.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// Count 根据查询条件统计规则数量
func (r *RuleRepository) Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error) {
	return r.repo.Count(ctx, query)
}

// RecordExecution 记录执行统计
func (r *RuleRepository) RecordExecution(ctx context.Context, ruleID string, success bool) error {
	return r.repo.UpdateExecuteStats(ctx, ruleID, success)
}

// toEntity 将领域模型转换为实体
func (r *RuleRepository) toEntity(rule *model.Rule) *entity.Rule {
	return &entity.Rule{
		ID:            rule.ID,
		Code:          rule.Code,
		Name:          rule.Name,
		Description:   rule.Description,
		CategoryID:    rule.CategoryID,
		TemplateID:    rule.TemplateID,
		Type:          rule.Type,
		Version:       rule.Version,
		Status:        int(rule.Status),
		Triggers:      strings.Join(rule.Triggers, ","),
		Scope:         rule.Scope,
		ScopeID:       rule.ScopeID,
		Conditions:    rule.Conditions,
		LuaScript:     rule.LuaScript,
		Formula:       rule.Formula,
		FormulaVars:   rule.FormulaVars,
		Action:        rule.Action,
		Priority:      rule.Priority,
		Sorting:       rule.Sorting,
		ExecuteCount:  rule.ExecuteCount,
		SuccessCount:  rule.SuccessCount,
		LastExecuteAt: rule.LastExecuteAt,
		TenantID:      rule.TenantID,
	}
}

// toModel 将实体转换为领域模型
func (r *RuleRepository) toModel(entity *entity.Rule) *model.Rule {
	if entity == nil {
		return nil
	}
	triggers := make([]string, 0)
	if entity.Triggers != "" {
		triggers = strings.Split(entity.Triggers, ",")
	}
	return &model.Rule{
		ID:              entity.ID,
		Code:            entity.Code,
		Name:            entity.Name,
		Description:     entity.Description,
		CategoryID:      entity.CategoryID,
		TemplateID:      entity.TemplateID,
		Type:            entity.Type,
		Version:         entity.Version,
		Status:          int32(entity.Status),
		Triggers:        triggers,
		Scope:           entity.Scope,
		ScopeID:         entity.ScopeID,
		ExecutionTiming: entity.ExecutionTiming,
		Conditions:      entity.Conditions,
		LuaScript:       entity.LuaScript,
		Formula:         entity.Formula,
		FormulaVars:     entity.FormulaVars,
		Action:          entity.Action,
		Priority:        entity.Priority,
		Sorting:         entity.Sorting,
		ExecuteCount:    entity.ExecuteCount,
		SuccessCount:    entity.SuccessCount,
		LastExecuteAt:   entity.LastExecuteAt,
		CreatedAt:       entity.CreatedAt,
		UpdatedAt:       entity.UpdatedAt,
		TenantID:        entity.TenantID,
	}
}

// toModels 将实体切片转换为领域模型切片
func (r *RuleRepository) toModels(entities []*entity.Rule) []*model.Rule {
	if entities == nil {
		return nil
	}
	models := make([]*model.Rule, len(entities))
	for i, entity := range entities {
		models[i] = r.toModel(entity)
	}
	return models
}
