package repository

import (
	"context"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/repository"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/infrastructure/persistence/entity"
)

// IRuleTemplateRepository 规则模板数据访问接口
type IRuleTemplateRepository interface {
	baserepo.IBaseRepo[entity.RuleTemplate, string]
	FindByCode(ctx context.Context, code string) (*entity.RuleTemplate, error)
	FindByCategoryID(ctx context.Context, categoryID string) ([]*entity.RuleTemplate, error)
	FindByType(ctx context.Context, templateType string) ([]*entity.RuleTemplate, error)
	FindByBusinessType(ctx context.Context, businessType string) ([]*entity.RuleTemplate, error)
	FindByScope(ctx context.Context, scope string) ([]*entity.RuleTemplate, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	FindAll(ctx context.Context) ([]*entity.RuleTemplate, error)
	Find(ctx context.Context, query *db_query.QueryBuilder) ([]*entity.RuleTemplate, error)
	Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error)
}

// RuleTemplateRepository 规则模板仓储实现
type RuleTemplateRepository struct {
	repo IRuleTemplateRepository
}

// NewRuleTemplateRepository 创建规则模板仓储
func NewRuleTemplateRepository(repo IRuleTemplateRepository) repository.ITemplateRepository {
	return &RuleTemplateRepository{
		repo: repo,
	}
}

// Create 创建模板
func (r *RuleTemplateRepository) Create(ctx context.Context, template *model.RuleTemplate) error {
	// 转换为数据库实体
	entity := &entity.RuleTemplate{
		Code:        template.Code,
		Name:        template.Name,
		Description: template.Description,
		CategoryID:  template.CategoryID,
		Type:        template.Type,
		Version:     template.Version,
		Status:      int(template.Status),
		Conditions:  template.Conditions,
		LuaScript:   template.LuaScript,
		Formula:     template.Formula,
		FormulaVars: template.FormulaVars,
		Parameters:  template.Parameters,
		Priority:    template.Priority,
		Sorting:     template.Sorting,
		TenantID:    template.TenantID,
	}
	_, err := r.repo.Add(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

// Update 更新模板
func (r *RuleTemplateRepository) Update(ctx context.Context, template *model.RuleTemplate) error {
	// 转换为数据库实体
	entity := &entity.RuleTemplate{
		ID:          template.ID,
		Code:        template.Code,
		Name:        template.Name,
		Description: template.Description,
		CategoryID:  template.CategoryID,
		Type:        template.Type,
		Version:     template.Version,
		Status:      int(template.Status),
		Conditions:  template.Conditions,
		LuaScript:   template.LuaScript,
		Formula:     template.Formula,
		FormulaVars: template.FormulaVars,
		Parameters:  template.Parameters,
		Priority:    template.Priority,
		Sorting:     template.Sorting,
		TenantID:    template.TenantID,
	}
	return r.repo.EditById(ctx, template.ID, entity)
}

// Delete 删除模板
func (r *RuleTemplateRepository) Delete(ctx context.Context, id string) error {
	return r.repo.DelByIdUnScoped(ctx, id)
}

// FindByID 根据ID查询模板
func (r *RuleTemplateRepository) FindByID(ctx context.Context, id string) (*model.RuleTemplate, error) {
	entity, err := r.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toModel(entity), nil
}

// FindByCode 根据编码查询模板
func (r *RuleTemplateRepository) FindByCode(ctx context.Context, code string) (*model.RuleTemplate, error) {
	entity, err := r.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return r.toModel(entity), nil
}

// FindByCategoryID 根据分类ID查询模板列表
func (r *RuleTemplateRepository) FindByCategoryID(ctx context.Context, categoryID string) ([]*model.RuleTemplate, error) {
	entities, err := r.repo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindByType 根据类型查询模板列表
func (r *RuleTemplateRepository) FindByType(ctx context.Context, templateType string) ([]*model.RuleTemplate, error) {
	entities, err := r.repo.FindByType(ctx, templateType)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindByBusinessType 根据业务类型查询模板列表
func (r *RuleTemplateRepository) FindByBusinessType(ctx context.Context, businessType string) ([]*model.RuleTemplate, error) {
	entities, err := r.repo.FindByBusinessType(ctx, businessType)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindByScope 根据作用域查询模板列表
func (r *RuleTemplateRepository) FindByScope(ctx context.Context, scope string) ([]*model.RuleTemplate, error) {
	entities, err := r.repo.FindByScope(ctx, scope)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// ExistsByCode 检查编码是否存在
func (r *RuleTemplateRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return r.repo.ExistsByCode(ctx, code)
}

// FindAll 查询所有模板
func (r *RuleTemplateRepository) FindAll(ctx context.Context) ([]*model.RuleTemplate, error) {
	entities, err := r.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// Find 根据查询条件查找模板列表
func (r *RuleTemplateRepository) Find(ctx context.Context, query *db_query.QueryBuilder) ([]*model.RuleTemplate, error) {
	entities, err := r.repo.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// Count 根据查询条件统计模板数量
func (r *RuleTemplateRepository) Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error) {
	return r.repo.Count(ctx, query)
}

// FindEnabledByCategoryID 根据分类ID查找启用的模板列表
func (r *RuleTemplateRepository) FindEnabledByCategoryID(ctx context.Context, categoryID string) ([]*model.RuleTemplate, error) {
	// 实现启用状态的查询逻辑
	entities, err := r.repo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// 过滤启用的模板
	var enabledTemplates []*model.RuleTemplate
	for _, entity := range entities {
		if entity.Status == 1 { // 1表示启用状态
			enabledTemplates = append(enabledTemplates, r.toModel(entity))
		}
	}
	return enabledTemplates, nil
}

// FindEnabledByType 根据类型查找启用的模板列表
func (r *RuleTemplateRepository) FindEnabledByType(ctx context.Context, templateType string) ([]*model.RuleTemplate, error) {
	// 实现启用状态的查询逻辑
	entities, err := r.repo.FindByType(ctx, templateType)
	if err != nil {
		return nil, err
	}

	// 过滤启用的模板
	var enabledTemplates []*model.RuleTemplate
	for _, entity := range entities {
		if entity.Status == 1 { // 1表示启用状态
			enabledTemplates = append(enabledTemplates, r.toModel(entity))
		}
	}
	return enabledTemplates, nil
}

// FindEnabledByBusinessType 根据业务类型查找启用的模板列表
func (r *RuleTemplateRepository) FindEnabledByBusinessType(ctx context.Context, businessType string) ([]*model.RuleTemplate, error) {
	// 实现启用状态的查询逻辑
	entities, err := r.repo.FindByBusinessType(ctx, businessType)
	if err != nil {
		return nil, err
	}

	// 过滤启用的模板
	var enabledTemplates []*model.RuleTemplate
	for _, entity := range entities {
		if entity.Status == 1 { // 1表示启用状态
			enabledTemplates = append(enabledTemplates, r.toModel(entity))
		}
	}
	return enabledTemplates, nil
}

// FindEnabledByScope 根据作用域查找启用的模板列表
func (r *RuleTemplateRepository) FindEnabledByScope(ctx context.Context, scope string) ([]*model.RuleTemplate, error) {
	// 实现启用状态的查询逻辑
	entities, err := r.repo.FindByScope(ctx, scope)
	if err != nil {
		return nil, err
	}

	// 过滤启用的模板
	var enabledTemplates []*model.RuleTemplate
	for _, entity := range entities {
		if entity.Status == 1 { // 1表示启用状态
			enabledTemplates = append(enabledTemplates, r.toModel(entity))
		}
	}
	return enabledTemplates, nil
}

// toModel 将实体转换为领域模型
func (r *RuleTemplateRepository) toModel(entity *entity.RuleTemplate) *model.RuleTemplate {
	if entity == nil {
		return nil
	}

	return &model.RuleTemplate{
		ID:          entity.ID,
		Code:        entity.Code,
		Name:        entity.Name,
		Description: entity.Description,
		CategoryID:  entity.CategoryID,
		Type:        entity.Type,
		Version:     entity.Version,
		Status:      int32(entity.Status),
		Conditions:  entity.Conditions,
		LuaScript:   entity.LuaScript,
		Formula:     entity.Formula,
		FormulaVars: entity.FormulaVars,
		Parameters:  entity.Parameters,
		Priority:    entity.Priority,
		Sorting:     entity.Sorting,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		TenantID:    entity.TenantID,
	}
}

// toModels 将实体切片转换为领域模型切片
func (r *RuleTemplateRepository) toModels(entities []*entity.RuleTemplate) []*model.RuleTemplate {
	if entities == nil {
		return nil
	}
	models := make([]*model.RuleTemplate, len(entities))
	for i, entity := range entities {
		models[i] = r.toModel(entity)
	}
	return models
}
