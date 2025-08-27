package dto

import (
	"encoding/json"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/persistence/entity"
)

// TemplateDTO 模板数据传输对象
type TemplateDTO struct {
	ID          string         `json:"id"`          // 模板ID
	Code        string         `json:"code"`        // 模板编码
	Name        string         `json:"name"`        // 模板名称
	Description string         `json:"description"` // 模板描述
	CategoryID  string         `json:"category_id"` // 分类ID
	Attributes  []AttributeDTO `json:"attributes"`  // 模板属性
	Status      int            `json:"status"`      // 状态
	CreatedAt   int64          `json:"created_at"`  // 创建时间
	UpdatedAt   int64          `json:"updated_at"`  // 更新时间
}

// AttributeDTO 模板属性数据传输对象
type AttributeDTO struct {
	Key         string        `json:"key"`         // 属性键
	Name        string        `json:"name"`        // 属性名称
	Type        string        `json:"type"`        // 属性类型
	Required    bool          `json:"required"`    // 是否必填
	IsQuery     bool          `json:"is_query"`    // 是否作为查询条件
	I18nKey     string        `json:"i18n_key"`    // 国际化标识
	Options     []OptionDTO   `json:"options"`     // 选项列表
	Default     interface{}   `json:"default"`     // 默认值
	Validation  ValidationDTO `json:"validation"`  // 验证规则
	Description string        `json:"description"` // 属性描述
}

// OptionDTO 选项数据传输对象
type OptionDTO struct {
	Label string      `json:"label"` // 选项标签
	Value interface{} `json:"value"` // 选项值
	Sort  int         `json:"sort"`  // 排序值
}

// ValidationDTO 验证规则数据传输对象
type ValidationDTO struct {
	Min     *float64 `json:"min,omitempty"`     // 最小值
	Max     *float64 `json:"max,omitempty"`     // 最大值
	Pattern string   `json:"pattern,omitempty"` // 正则表达式
	Length  *int     `json:"length,omitempty"`  // 长度限制
}

// TemplateFromEntity 从实体转换为DTO
func TemplateFromEntity(entity *entity.Template) *TemplateDTO {
	if entity == nil {
		return nil
	}

	dto := &TemplateDTO{
		ID:          entity.ID,
		Code:        entity.Code,
		Name:        entity.Name,
		Description: entity.Description,
		CategoryID:  entity.CategoryID,
		Status:      entity.Status,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}

	// 解析属性JSON
	if entity.Attributes != "" {
		var attributes []AttributeDTO
		if err := json.Unmarshal([]byte(entity.Attributes), &attributes); err == nil {
			dto.Attributes = attributes
		}
	}

	return dto
}

// TemplateFromEntities 从实体列表转换为DTO列表
func TemplateFromEntities(entities []*entity.Template) []*TemplateDTO {
	if entities == nil {
		return nil
	}

	dtos := make([]*TemplateDTO, 0, len(entities))
	for _, entity := range entities {
		dtos = append(dtos, TemplateFromEntity(entity))
	}
	return dtos
}
