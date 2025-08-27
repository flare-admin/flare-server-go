package valueobject

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	template_err "github.com/flare-admin/flare-server-go/framework/support/template/domain/err"
)

// TemplateAttribute 模板属性值对象
type TemplateAttribute struct {
	Key         string      `json:"key" is_query:"true"`         // 属性键
	Name        string      `json:"name" is_query:"true"`        // 属性名称
	Type        string      `json:"type" is_query:"true"`        // 属性类型
	Required    bool        `json:"required"`                    // 是否必填
	IsQuery     bool        `json:"is_query"`                    // 是否作为查询条件
	I18nKey     string      `json:"i18n_key"`                    // 国际化标识
	Options     []Option    `json:"options"`                     // 选项列表
	Default     interface{} `json:"default"`                     // 默认值
	Validation  Validation  `json:"validation"`                  // 验证规则
	Description string      `json:"description" is_query:"true"` // 属性描述
}

// Option 选项值对象
type Option struct {
	Label string      `json:"label"` // 选项标签
	Value interface{} `json:"value"` // 选项值
	Sort  int         `json:"sort"`  // 排序值
}

// Validation 验证规则值对象
type Validation struct {
	Min     *float64 `json:"min,omitempty"`     // 最小值
	Max     *float64 `json:"max,omitempty"`     // 最大值
	Pattern string   `json:"pattern,omitempty"` // 正则表达式
	Length  *int     `json:"length,omitempty"`  // 长度限制
}

// Template 模板值对象
type Template struct {
	ID          string              `json:"id" is_query:"true"`          // 模板ID
	Code        string              `json:"code" is_query:"true"`        // 模板编码
	Name        string              `json:"name" is_query:"true"`        // 模板名称
	Description string              `json:"description" is_query:"true"` // 模板描述
	CategoryID  string              `json:"category_id" is_query:"true"` // 分类ID
	Attributes  []TemplateAttribute `json:"attributes"`                  // 模板属性
	Status      int                 `json:"status" is_query:"true"`      // 模板状态
	CreatedAt   int64               `json:"created_at" is_query:"true"`  // 创建时间
	UpdatedAt   int64               `json:"updated_at" is_query:"true"`  // 更新时间
}

// CreateTemplateCommand 创建模板命令
type CreateTemplateCommand struct {
	Code        string              `json:"code" binding:"required" comment:"模板编码"`
	Name        string              `json:"name" binding:"required" comment:"模板名称"`
	Description string              `json:"description" comment:"模板描述"`
	CategoryID  string              `json:"category_id" binding:"required" comment:"分类ID"`
	Attributes  []TemplateAttribute `json:"attributes" binding:"required" comment:"模板属性"`
}

// UpdateTemplateCommand 更新模板命令
type UpdateTemplateCommand struct {
	ID          string              `json:"id" binding:"required" comment:"模板ID"`
	Code        string              `json:"code" binding:"required" comment:"模板编码"`
	Name        string              `json:"name" binding:"required" comment:"模板名称"`
	Description string              `json:"description" comment:"模板描述"`
	CategoryID  string              `json:"category_id" binding:"required" comment:"分类ID"`
	Attributes  []TemplateAttribute `json:"attributes" binding:"required" comment:"模板属性"`
}

// UpdateTemplateStatusCommand 更新模板状态命令
type UpdateTemplateStatusCommand struct {
	ID     string `json:"id" binding:"required" comment:"模板ID"`
	Status int    `json:"status" binding:"required" comment:"模板状态"`
}

// DeleteTemplateCommand 删除模板命令
type DeleteTemplateCommand struct {
	ID string `json:"id" binding:"required" comment:"模板ID"`
}

// Validate 验证模板属性
func (a *TemplateAttribute) Validate() *herrors.HError {
	if a.Key == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("属性键不能为空"))
	}
	if a.Name == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("属性名称不能为空"))
	}
	if a.Type == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("属性类型不能为空"))
	}
	return nil
}

// Validate 验证创建模板命令
func (c *CreateTemplateCommand) Validate() *herrors.HError {
	if c.Code == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("模板编码不能为空"))
	}
	if c.Name == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("模板名称不能为空"))
	}
	if c.CategoryID == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("分类ID不能为空"))
	}
	if len(c.Attributes) == 0 {
		return template_err.TemplateCreateFailed(fmt.Errorf("模板属性不能为空"))
	}
	for _, attr := range c.Attributes {
		if err := attr.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate 验证更新模板命令
func (c *UpdateTemplateCommand) Validate() *herrors.HError {
	if c.ID == "" {
		return template_err.TemplateUpdateFailed(fmt.Errorf("模板ID不能为空"))
	}
	if c.Code == "" {
		return template_err.TemplateUpdateFailed(fmt.Errorf("模板编码不能为空"))
	}
	if c.Name == "" {
		return template_err.TemplateUpdateFailed(fmt.Errorf("模板名称不能为空"))
	}
	if c.CategoryID == "" {
		return template_err.TemplateUpdateFailed(fmt.Errorf("分类ID不能为空"))
	}
	if len(c.Attributes) == 0 {
		return template_err.TemplateUpdateFailed(fmt.Errorf("模板属性不能为空"))
	}
	for _, attr := range c.Attributes {
		if err := attr.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate 验证更新模板状态命令
func (c *UpdateTemplateStatusCommand) Validate() *herrors.HError {
	if c.ID == "" {
		return template_err.TemplateUpdateFailed(fmt.Errorf("模板ID不能为空"))
	}
	if c.Status != 1 && c.Status != 2 {
		return template_err.TemplateUpdateFailed(fmt.Errorf("无效的模板状态"))
	}
	return nil
}

// Validate 验证删除模板命令
func (c *DeleteTemplateCommand) Validate() *herrors.HError {
	if c.ID == "" {
		return template_err.TemplateDeleteFailed(fmt.Errorf("模板ID不能为空"))
	}
	return nil
}
