package entity

import "github.com/flare-admin/flare-server-go/framework/pkg/database"

// Template 模板实体
type Template struct {
	database.BaseModel
	ID          string `gorm:"primarykey"`
	Code        string `gorm:"size:100;not null;uniqueIndex;comment:模板编码"`
	Name        string `gorm:"size:100;not null;comment:模板名称"`
	Description string `gorm:"size:500;comment:模板描述"`
	CategoryID  string `gorm:"not null;index;comment:分类ID"`
	Attributes  string `gorm:"type:json;comment:模板属性"`
	Status      int    `gorm:"not null;default:1;comment:状态：1-启用 2-禁用"`
}

func (Template) TableName() string {
	return "templates"
}

// GetPrimaryKey 获取主键
func (Template) GetPrimaryKey() string {
	return "id"
}

// Category 分类实体
type Category struct {
	database.BaseModel
	ID          string `gorm:"primarykey"`
	Name        string `gorm:"size:100;not null;comment:分类名称"`
	Code        string `gorm:"size:50;not null;uniqueIndex;comment:分类编码"`
	Description string `gorm:"size:500;comment:分类描述"`
	Sort        int    `gorm:"not null;default:0;comment:排序"`
	Status      int    `gorm:"not null;default:1;comment:状态：1-启用 2-禁用"`
}

func (Category) TableName() string {
	return "template_categories"
}

// GetPrimaryKey 获取主键
func (Category) GetPrimaryKey() string {
	return "id"
}
