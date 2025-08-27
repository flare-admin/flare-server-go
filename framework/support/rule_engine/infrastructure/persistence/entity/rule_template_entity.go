package entity

import "github.com/flare-admin/flare-server-go/framework/pkg/database"

// RuleTemplate 规则模板实体
type RuleTemplate struct {
	database.BaseModel
	ID          string `gorm:"primarykey"`
	Code        string `gorm:"size:100;not null;uniqueIndex;comment:模板编码"`
	Name        string `gorm:"size:100;not null;comment:模板名称"`
	Description string `gorm:"size:500;comment:模板描述"`
	CategoryID  string `gorm:"not null;index;comment:分类ID"`
	Type        string `gorm:"size:50;not null;comment:模板类型：condition(条件模板) lua(lua脚本模板) formula(公式模板)"`
	Version     string `gorm:"size:20;not null;default:'1.0.0';comment:模板版本"`
	Status      int    `gorm:"not null;default:1;comment:状态：1-启用 2-禁用"`
	Conditions  string `gorm:"type:json;comment:条件表达式(JSON格式)"`
	LuaScript   string `gorm:"type:text;comment:Lua脚本代码"`
	Formula     string `gorm:"type:text;comment:计算公式"`
	FormulaVars string `gorm:"type:json;comment:公式变量映射(JSON格式)"`
	Parameters  string `gorm:"type:json;comment:模板参数定义(JSON格式)"`
	Priority    int32  `gorm:"not null;default:0;comment:优先级"`
	Sorting     int32  `gorm:"not null;default:0;comment:排序权重"`
	TenantID    string `gorm:"size:50;comment:租户ID"`
}

func (RuleTemplate) TableName() string {
	return "rule_templates"
}

// GetPrimaryKey 获取主键
func (RuleTemplate) GetPrimaryKey() string {
	return "id"
}
