package repository

import (
	"context"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/repository"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/infrastructure/persistence/entity"
)

// IRuleCategoryRepository 规则分类数据访问接口
type IRuleCategoryRepository interface {
	baserepo.IBaseRepo[entity.RuleCategory, string]
	FindByCode(ctx context.Context, code string) (*entity.RuleCategory, error)
	FindByParentID(ctx context.Context, parentID string) ([]*entity.RuleCategory, error)
	FindByBusinessType(ctx context.Context, businessType string) ([]*entity.RuleCategory, error)
	FindByType(ctx context.Context, categoryType string) ([]*entity.RuleCategory, error)
	FindRootCategories(ctx context.Context) ([]*entity.RuleCategory, error)
	FindByPath(ctx context.Context, path string) (*entity.RuleCategory, error)
	FindDescendants(ctx context.Context, path string) ([]*entity.RuleCategory, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	FindAll(ctx context.Context) ([]*entity.RuleCategory, error)
	Find(ctx context.Context, query *db_query.QueryBuilder) ([]*entity.RuleCategory, error)
	Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error)
}

// RuleCategoryRepository 规则分类仓储实现
type RuleCategoryRepository struct {
	repo IRuleCategoryRepository
}

// NewRuleCategoryRepository 创建规则分类仓储
func NewRuleCategoryRepository(repo IRuleCategoryRepository) repository.ICategoryRepository {
	return &RuleCategoryRepository{
		repo: repo,
	}
}

// Create 创建分类
func (r *RuleCategoryRepository) Create(ctx context.Context, category *model.RuleCategory) error {
	// 转换为数据库实体
	entity := &entity.RuleCategory{
		Code:         category.Code,
		Name:         category.Name,
		Description:  category.Description,
		Type:         category.Type,
		ParentID:     category.ParentID,
		Level:        category.Level,
		Path:         category.Path,
		Sorting:      category.Sorting,
		Status:       int(category.Status),
		IsLeaf:       category.IsLeaf,
		BusinessType: category.BusinessType,
		TenantID:     category.TenantID,
	}
	_, err := r.repo.Add(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

// Update 更新分类
func (r *RuleCategoryRepository) Update(ctx context.Context, category *model.RuleCategory) error {
	// 转换为数据库实体
	entity := &entity.RuleCategory{
		ID:           category.ID,
		Code:         category.Code,
		Name:         category.Name,
		Description:  category.Description,
		Type:         category.Type,
		ParentID:     category.ParentID,
		Level:        category.Level,
		Path:         category.Path,
		Sorting:      category.Sorting,
		Status:       int(category.Status),
		IsLeaf:       category.IsLeaf,
		BusinessType: category.BusinessType,
		TenantID:     category.TenantID,
	}
	return r.repo.EditById(ctx, entity)
}

// Delete 删除分类
func (r *RuleCategoryRepository) Delete(ctx context.Context, id string) error {
	return r.repo.DelByIdUnScoped(ctx, id)
}

// FindByID 根据ID查询分类
func (r *RuleCategoryRepository) FindByID(ctx context.Context, id string) (*model.RuleCategory, error) {
	entity, err := r.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toModel(entity), nil
}

// FindByCode 根据编码查询分类
func (r *RuleCategoryRepository) FindByCode(ctx context.Context, code string) (*model.RuleCategory, error) {
	entity, err := r.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return r.toModel(entity), nil
}

// FindByParentID 根据父分类ID查询子分类列表
func (r *RuleCategoryRepository) FindByParentID(ctx context.Context, parentID string) ([]*model.RuleCategory, error) {
	entities, err := r.repo.FindByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindByBusinessType 根据业务类型查询分类列表
func (r *RuleCategoryRepository) FindByBusinessType(ctx context.Context, businessType string) ([]*model.RuleCategory, error) {
	entities, err := r.repo.FindByBusinessType(ctx, businessType)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindByType 根据分类类型查询分类列表
func (r *RuleCategoryRepository) FindByType(ctx context.Context, categoryType string) ([]*model.RuleCategory, error) {
	entities, err := r.repo.FindByType(ctx, categoryType)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindRootCategories 查询根分类列表
func (r *RuleCategoryRepository) FindRootCategories(ctx context.Context) ([]*model.RuleCategory, error) {
	entities, err := r.repo.FindRootCategories(ctx)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// FindByPath 根据路径查询分类
func (r *RuleCategoryRepository) FindByPath(ctx context.Context, path string) (*model.RuleCategory, error) {
	entity, err := r.repo.FindByPath(ctx, path)
	if err != nil {
		return nil, err
	}
	return r.toModel(entity), nil
}

// FindDescendants 查询后代分类列表
func (r *RuleCategoryRepository) FindDescendants(ctx context.Context, path string) ([]*model.RuleCategory, error) {
	entities, err := r.repo.FindDescendants(ctx, path)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// ExistsByCode 检查编码是否存在
func (r *RuleCategoryRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return r.repo.ExistsByCode(ctx, code)
}

// FindAll 查询所有分类
func (r *RuleCategoryRepository) FindAll(ctx context.Context) ([]*model.RuleCategory, error) {
	entities, err := r.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// Find 根据查询条件查找分类列表
func (r *RuleCategoryRepository) Find(ctx context.Context, query *db_query.QueryBuilder) ([]*model.RuleCategory, error) {
	entities, err := r.repo.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	return r.toModels(entities), nil
}

// Count 根据查询条件统计分类数量
func (r *RuleCategoryRepository) Count(ctx context.Context, query *db_query.QueryBuilder) (int64, error) {
	return r.repo.Count(ctx, query)
}

// FindEnabledByParentID 根据父分类ID查找启用的子分类列表
func (r *RuleCategoryRepository) FindEnabledByParentID(ctx context.Context, parentID string) ([]*model.RuleCategory, error) {
	// 实现启用状态的查询逻辑
	entities, err := r.repo.FindByParentID(ctx, parentID)
	if err != nil {
		return nil, err
	}

	// 过滤启用的分类
	var enabledCategories []*model.RuleCategory
	for _, entity := range entities {
		if entity.Status == 1 { // 1表示启用状态
			enabledCategories = append(enabledCategories, r.toModel(entity))
		}
	}
	return enabledCategories, nil
}

// FindEnabledRootCategories 查找启用的根分类列表
func (r *RuleCategoryRepository) FindEnabledRootCategories(ctx context.Context) ([]*model.RuleCategory, error) {
	// 实现启用状态的查询逻辑
	entities, err := r.repo.FindRootCategories(ctx)
	if err != nil {
		return nil, err
	}

	// 过滤启用的分类
	var enabledCategories []*model.RuleCategory
	for _, entity := range entities {
		if entity.Status == 1 { // 1表示启用状态
			enabledCategories = append(enabledCategories, r.toModel(entity))
		}
	}
	return enabledCategories, nil
}

// FindEnabledByType 根据类型查找启用的分类列表
func (r *RuleCategoryRepository) FindEnabledByType(ctx context.Context, categoryType string) ([]*model.RuleCategory, error) {
	// 实现启用状态的查询逻辑
	entities, err := r.repo.FindByType(ctx, categoryType)
	if err != nil {
		return nil, err
	}

	// 过滤启用的分类
	var enabledCategories []*model.RuleCategory
	for _, entity := range entities {
		if entity.Status == 1 { // 1表示启用状态
			enabledCategories = append(enabledCategories, r.toModel(entity))
		}
	}
	return enabledCategories, nil
}

// FindEnabledByBusinessType 根据业务类型查找启用的分类列表
func (r *RuleCategoryRepository) FindEnabledByBusinessType(ctx context.Context, businessType string) ([]*model.RuleCategory, error) {
	// 实现启用状态的查询逻辑
	entities, err := r.repo.FindByBusinessType(ctx, businessType)
	if err != nil {
		return nil, err
	}

	// 过滤启用的分类
	var enabledCategories []*model.RuleCategory
	for _, entity := range entities {
		if entity.Status == 1 { // 1表示启用状态
			enabledCategories = append(enabledCategories, r.toModel(entity))
		}
	}
	return enabledCategories, nil
}

// toModel 将实体转换为领域模型
func (r *RuleCategoryRepository) toModel(entity *entity.RuleCategory) *model.RuleCategory {
	if entity == nil {
		return nil
	}

	return &model.RuleCategory{
		ID:           entity.ID,
		Code:         entity.Code,
		Name:         entity.Name,
		Description:  entity.Description,
		Type:         entity.Type,
		ParentID:     entity.ParentID,
		Level:        entity.Level,
		Path:         entity.Path,
		Sorting:      entity.Sorting,
		Status:       int32(entity.Status),
		IsLeaf:       entity.IsLeaf,
		BusinessType: entity.BusinessType,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
		TenantID:     entity.TenantID,
	}
}

// toModels 将实体切片转换为领域模型切片
func (r *RuleCategoryRepository) toModels(entities []*entity.RuleCategory) []*model.RuleCategory {
	if entities == nil {
		return nil
	}
	models := make([]*model.RuleCategory, len(entities))
	for i, entity := range entities {
		models[i] = r.toModel(entity)
	}
	return models
}
