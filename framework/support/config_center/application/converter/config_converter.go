package converter

import (
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/dto"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/infrastructure/entity"
)

// ToConfigDTO 将配置实体转换为 DTO
func ToConfigDTO(e *entity.Config) *dto.ConfigDTO {
	if e == nil {
		return nil
	}
	return &dto.ConfigDTO{
		ID:          e.ID,
		Name:        e.Name,
		Key:         e.Key,
		Value:       e.Value,
		Type:        e.Type,
		Group:       e.Group,
		Description: e.Description,
		I18nKey:     e.I18nKey,
		IsSystem:    e.IsSystem,
		IsEnabled:   e.IsEnabled,
		Sort:        e.Sort,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// ToConfigEntity 将配置 DTO 转换为实体
func ToConfigEntity(d *dto.ConfigDTO) *entity.Config {
	if d == nil {
		return nil
	}
	return &entity.Config{
		ID:          d.ID,
		Key:         d.Key,
		Name:        d.Name,
		Value:       d.Value,
		Type:        d.Type,
		Group:       d.Group,
		Description: d.Description,
		I18nKey:     d.I18nKey,
		IsSystem:    d.IsSystem,
		IsEnabled:   d.IsEnabled,
		Sort:        d.Sort,
	}
}

// ToConfigGroupDTO 将配置分组实体转换为 DTO
func ToConfigGroupDTO(e *entity.ConfigGroup) *dto.ConfigGroupDTO {
	if e == nil {
		return nil
	}
	return &dto.ConfigGroupDTO{
		ID:          e.ID,
		Name:        e.Name,
		Code:        e.Code,
		Description: e.Description,
		I18nKey:     e.I18nKey,
		IsSystem:    e.IsSystem,
		IsEnabled:   e.IsEnabled,
		Sort:        e.Sort,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

// ToConfigGroupEntity 将配置分组 DTO 转换为实体
func ToConfigGroupEntity(d *dto.ConfigGroupDTO) *entity.ConfigGroup {
	if d == nil {
		return nil
	}
	return &entity.ConfigGroup{
		ID:          d.ID,
		Name:        d.Name,
		Code:        d.Code,
		Description: d.Description,
		I18nKey:     d.I18nKey,
		IsSystem:    d.IsSystem,
		IsEnabled:   d.IsEnabled,
		Sort:        d.Sort,
	}
}

// ToConfigDTOList 将配置实体列表转换为 DTO 列表
func ToConfigDTOList(entities []*entity.Config) []*dto.ConfigDTO {
	if entities == nil {
		return nil
	}
	dtos := make([]*dto.ConfigDTO, len(entities))
	for i, e := range entities {
		dtos[i] = ToConfigDTO(e)
	}
	return dtos
}

// ToConfigGroupDTOList 将配置分组实体列表转换为 DTO 列表
func ToConfigGroupDTOList(entities []*entity.ConfigGroup) []*dto.ConfigGroupDTO {
	if entities == nil {
		return nil
	}
	dtos := make([]*dto.ConfigGroupDTO, len(entities))
	for i, e := range entities {
		dtos[i] = ToConfigGroupDTO(e)
	}
	return dtos
}
