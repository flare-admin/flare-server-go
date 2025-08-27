package command

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
	Status int    `json:"status" binding:"required" comment:"状态"`
}

// DeleteTemplateCommand 删除模板命令
type DeleteTemplateCommand struct {
	ID string `json:"id" binding:"required" comment:"模板ID"`
}

// TemplateAttribute 模板属性
type TemplateAttribute struct {
	Key         string      `json:"key" binding:"required" comment:"属性键"`
	Name        string      `json:"name" binding:"required" comment:"属性名称"`
	Type        string      `json:"type" binding:"required" comment:"属性类型"`
	Required    bool        `json:"required" comment:"是否必填"`
	IsQuery     bool        `json:"is_query" comment:"是否作为查询条件"`
	I18nKey     string      `json:"i18n_key" comment:"国际化标识"`
	Options     []Option    `json:"options" comment:"选项列表"`
	Default     interface{} `json:"default" comment:"默认值"`
	Validation  Validation  `json:"validation" comment:"验证规则"`
	Description string      `json:"description" comment:"属性描述"`
}

// Option 选项
type Option struct {
	Label string      `json:"label" binding:"required" comment:"选项标签"`
	Value interface{} `json:"value" binding:"required" comment:"选项值"`
	Sort  int         `json:"sort" comment:"排序值"`
}

// Validation 验证规则
type Validation struct {
	Min     *float64 `json:"min,omitempty" comment:"最小值"`
	Max     *float64 `json:"max,omitempty" comment:"最大值"`
	Pattern string   `json:"pattern,omitempty" comment:"正则表达式"`
	Length  *int     `json:"length,omitempty" comment:"长度限制"`
}
