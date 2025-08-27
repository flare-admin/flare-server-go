package repository

import (
	"context"
	"encoding/json"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/repository"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/persistence/entity"
)

// ITemplateRepository 模板仓储接口
type ITemplateRepository interface {
	baserepo.IBaseRepo[entity.Template, string]
	FindByCode(ctx context.Context, code string) (*entity.Template, error)
	FindByCategoryID(ctx context.Context, categoryID string) ([]*entity.Template, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
	FindAll(ctx context.Context) ([]*entity.Template, error)
}

// TemplateRepository 模板仓储实现
type TemplateRepository struct {
	repo ITemplateRepository
}

// NewTemplateRepository 创建模板仓储
func NewTemplateRepository(repo ITemplateRepository) repository.TemplateRepository {
	return &TemplateRepository{
		repo: repo,
	}
}

// Create 创建模板
func (d TemplateRepository) Create(ctx context.Context, template *model.Template) error {
	// 将Attributes转换为JSON字符串
	attributesJSON, err := json.Marshal(template.Attributes)
	if err != nil {
		return err
	}

	// 转换为数据库实体
	entity := &entity.Template{
		Code:        template.Code,
		Name:        template.Name,
		Description: template.Description,
		CategoryID:  template.CategoryID,
		Attributes:  string(attributesJSON),
		Status:      template.Status,
	}
	_, err = d.repo.Add(ctx, entity)
	if err != nil {
		return err
	}
	return nil
}

// Update 更新模板
func (d TemplateRepository) Update(ctx context.Context, template *model.Template) error {
	// 将Attributes转换为JSON字符串
	attributesJSON, err := json.Marshal(template.Attributes)
	if err != nil {
		return err
	}

	// 转换为数据库实体
	entity := &entity.Template{
		ID:          template.ID,
		Code:        template.Code,
		Name:        template.Name,
		Description: template.Description,
		CategoryID:  template.CategoryID,
		Attributes:  string(attributesJSON),
		Status:      template.Status,
	}
	return d.repo.EditById(ctx, template.ID, entity)
}

// Delete 删除模板
func (d TemplateRepository) Delete(ctx context.Context, id string) error {
	return d.repo.DelByIdUnScoped(ctx, id)
}

// FindByID 根据ID查询模板
func (d TemplateRepository) FindByID(ctx context.Context, id string) (*model.Template, error) {
	entity, err := d.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return d.toModel(entity), nil
}

// FindByCategoryID 根据分类ID查询模板列表
func (d TemplateRepository) FindByCategoryID(ctx context.Context, categoryID string) ([]*model.Template, error) {
	entities, err := d.repo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	return d.toModels(entities), nil
}

// FindByCode 根据编码查询模板
func (d TemplateRepository) FindByCode(ctx context.Context, code string) (*model.Template, error) {
	entity, err := d.repo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return d.toModel(entity), nil
}

// ExistsByCode 检查编码是否存在
func (d TemplateRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return d.repo.ExistsByCode(ctx, code)
}

// FindAll 查询所有模板
func (d TemplateRepository) FindAll(ctx context.Context) ([]*model.Template, error) {
	entities, err := d.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return d.toModels(entities), nil
}

// toModel 将实体转换为领域模型
func (d TemplateRepository) toModel(entity *entity.Template) *model.Template {
	if entity == nil {
		return nil
	}

	// 解析Attributes JSON字符串
	var attributes []model.Attribute
	if entity.Attributes != "" {
		if err := json.Unmarshal([]byte(entity.Attributes), &attributes); err != nil {
			// 如果解析失败，返回空属性列表
			attributes = make([]model.Attribute, 0)
		}
	}

	return &model.Template{
		ID:          entity.ID,
		Code:        entity.Code,
		Name:        entity.Name,
		Description: entity.Description,
		CategoryID:  entity.CategoryID,
		Attributes:  attributes,
		Status:      entity.Status,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}

// toModels 将实体切片转换为领域模型切片
func (d TemplateRepository) toModels(entities []*entity.Template) []*model.Template {
	if entities == nil {
		return nil
	}
	models := make([]*model.Template, len(entities))
	for i, entity := range entities {
		models[i] = d.toModel(entity)
	}
	return models
}
