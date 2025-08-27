package model

import (
	"encoding/json"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
)

// Template 模板领域模型
type Template struct {
	ID          string
	Code        string
	Name        string
	Description string
	CategoryID  string
	Attributes  []Attribute
	Status      int
	CreatedAt   int64
	UpdatedAt   int64
}

// Attribute 模板属性
type Attribute struct {
	Key         string      `json:"key"`         // 属性键
	Name        string      `json:"name"`        // 属性名称
	Type        string      `json:"type"`        // 属性类型：string、number、date、select、radio、checkbox
	Required    bool        `json:"required"`    // 是否必填
	IsQuery     bool        `json:"is_query"`    // 是否作为查询条件
	I18nKey     string      `json:"i18n_key"`    // 国际化标识
	Options     []Option    `json:"options"`     // 选项列表（用于select、radio、checkbox类型）
	Default     interface{} `json:"default"`     // 默认值
	Validation  Validation  `json:"validation"`  // 验证规则
	Description string      `json:"description"` // 属性描述
}

// Option 选项
type Option struct {
	Label string      `json:"label"` // 选项标签
	Value interface{} `json:"value"` // 选项值
	Sort  int         `json:"sort"`  // 排序值
}

// Validation 验证规则
type Validation struct {
	Min     *float64 `json:"min,omitempty"`     // 最小值（数字类型）
	Max     *float64 `json:"max,omitempty"`     // 最大值（数字类型）
	Pattern string   `json:"pattern,omitempty"` // 正则表达式（字符串类型）
	Length  *int     `json:"length,omitempty"`  // 长度限制（字符串类型）
}

// NewTemplate 创建模板
func NewTemplate(code, name, description, categoryID string) *Template {
	return &Template{
		Code:        code,
		Name:        name,
		Description: description,
		CategoryID:  categoryID,
		Status:      1,
		Attributes:  make([]Attribute, 0),
		CreatedAt:   utils.GetDateUnix(),
		UpdatedAt:   utils.GetDateUnix(),
	}
}

// AddAttribute 添加属性
func (t *Template) AddAttribute(attr Attribute) {
	t.Attributes = append(t.Attributes, attr)
	t.UpdatedAt = utils.GetDateUnix()
}

// ToJSON 转换为JSON字符串
func (t *Template) ToJSON() (string, error) {
	bytes, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON 从JSON字符串解析
func (t *Template) FromJSON(data string) error {
	return json.Unmarshal([]byte(data), t)
}
