package repository

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/repository"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/persistence/entity"
)

// ICategoryRepository 分类仓储接口
type ICategoryRepository interface {
	baserepo.IBaseRepo[entity.Category, string]
	FindByCode(ctx context.Context, code string) (*entity.Category, error)
	FindAll(ctx context.Context) ([]*entity.Category, error)
}

// CategoryRepository 分类仓储实现
type CategoryRepository struct {
	repo ICategoryRepository
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(repo ICategoryRepository) repository.CategoryRepository {
	return &CategoryRepository{
		repo: repo,
	}
}

// Create 创建分类
func (d CategoryRepository) Create(ctx context.Context, category *model.Category) error {
	// 转换为数据库实体
	entity := &entity.Category{
		Name:        category.Name,
		Code:        category.Code,
		Description: category.Description,
		Sort:        category.Sort,
		Status:      category.Status,
	}
	_, err := d.repo.Add(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

// Update 更新分类
func (d CategoryRepository) Update(ctx context.Context, category *model.Category) error {
	// 转换为数据库实体
	entity := &entity.Category{
		ID:          category.ID,
		Name:        category.Name,
		Code:        category.Code,
		Description: category.Description,
		Sort:        category.Sort,
		Status:      category.Status,
	}
	return d.repo.EditById(ctx, category.ID, entity)
}

// Delete 删除分类
func (d CategoryRepository) Delete(ctx context.Context, id string) error {
	return d.repo.DelByIdUnScoped(ctx, id)
}

// FindByID 根据ID查询分类
func (d CategoryRepository) FindByID(ctx context.Context, id string) (*model.Category, error) {
	entity, err := d.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return d.toModel(entity), nil
}

// FindByCode 根据编码查询分类
func (d CategoryRepository) FindByCode(ctx context.Context, code string) (*model.Category, error) {
	entity, err := d.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return d.toModel(entity), nil
}

// FindAll 查询所有分类
func (d CategoryRepository) FindAll(ctx context.Context) ([]*model.Category, error) {
	entities, err := d.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return d.toModels(entities), nil
}

// toModel 将实体转换为领域模型
func (d CategoryRepository) toModel(entity *entity.Category) *model.Category {
	if entity == nil {
		return nil
	}
	return &model.Category{
		ID:          entity.ID,
		Name:        entity.Name,
		Code:        entity.Code,
		Description: entity.Description,
		Sort:        entity.Sort,
		Status:      entity.Status,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// toModels 将实体切片转换为领域模型切片
func (d CategoryRepository) toModels(entities []*entity.Category) []*model.Category {
	if entities == nil {
		return nil
	}
	models := make([]*model.Category, len(entities))
	for i, entity := range entities {
		models[i] = d.toModel(entity)
	}
	return models
}
