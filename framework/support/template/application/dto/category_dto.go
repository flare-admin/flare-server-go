package dto

import (
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/persistence/entity"
)

// CategoryDTO 分类数据传输对象
type CategoryDTO struct {
	ID          string `json:"id"`          // 分类ID
	Name        string `json:"name"`        // 分类名称
	Code        string `json:"code"`        // 分类编码
	Description string `json:"description"` // 分类描述
	Sort        int    `json:"sort"`        // 排序
	Status      int    `json:"status"`      // 状态
	CreatedAt   int64  `json:"created_at"`  // 创建时间
	UpdatedAt   int64  `json:"updated_at"`  // 更新时间
}

// ToCategoryDTO 将领域模型转换为DTO
func ToCategoryDTO(category *model.Category) *CategoryDTO {
	if category == nil {
		return nil
	}

	return &CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Code:        category.Code,
		Description: category.Description,
		Sort:        category.Sort,
		Status:      category.Status,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// ToCategoryDTOs 将领域模型列表转换为DTO列表
func ToCategoryDTOs(categories []*model.Category) []*CategoryDTO {
	if categories == nil {
		return nil
	}

	dtos := make([]*CategoryDTO, 0, len(categories))
	for _, category := range categories {
		dtos = append(dtos, ToCategoryDTO(category))
	}
	return dtos
}

// FromEntity 从实体转换为DTO
func FromEntity(entity *entity.Category) *CategoryDTO {
	if entity == nil {
		return nil
	}

	return &CategoryDTO{
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

// FromEntities 从实体列表转换为DTO列表
func FromEntities(entities []*entity.Category) []*CategoryDTO {
	if entities == nil {
		return nil
	}

	dtos := make([]*CategoryDTO, 0, len(entities))
	for _, entity := range entities {
		dtos = append(dtos, FromEntity(entity))
	}
	return dtos
}
