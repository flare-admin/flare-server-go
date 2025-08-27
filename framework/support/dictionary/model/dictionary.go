package model

import "github.com/flare-admin/flare-server-go/framework/pkg/database"

// Category 表示字典分类
type Category struct {
	database.BaseModel
	ID          string `json:"id" gorm:"column:id;type:varchar(64);primary_key;comment:分类ID"`
	Name        string `json:"name" gorm:"column:name;type:varchar(100);comment:分类名称"`
	I18nKey     string `json:"i18n_key" gorm:"column:i18n_key;type:varchar(100);comment:国际化key"`
	Description string `json:"description" gorm:"column:description;type:varchar(500);comment:描述"`
	TenantID    string `json:"tenantId" gorm:"index:tenantIndex,column:tenant_id;default:'';comment:租户ID"`
}

// TableName 定义表名称
func (Category) TableName() string {
	return "dict_categories"
}

// GetPrimaryKey 定义表主键
func (Category) GetPrimaryKey() string {
	return "id"
}

// Option 表示字典选项
type Option struct {
	database.BaseModel
	ID         string `json:"id" gorm:"column:id;type:varchar(64);primary_key;comment:选项ID"`
	CategoryID string `json:"category_id" gorm:"column:category_id;type:varchar(64);index:idx_category;comment:分类ID"`
	Label      string `json:"label" gorm:"column:label;type:varchar(50);comment:默认名称"`
	Value      string `json:"value" gorm:"column:value;type:varchar(100);comment:选项值"`
	I18nKey    string `json:"i18n_key" gorm:"column:i18n_key;type:varchar(100);comment:国际化key"`
	Sort       int    `json:"sort" gorm:"column:sort;type:int;default:0;comment:排序号"`
	Status     int    `json:"status" gorm:"column:status;default:1;comment:状态:1-启用,0-禁用"`
	Remark     string `json:"remark" gorm:"column:remark;type:varchar(500);comment:备注"`
	TenantID   string `json:"tenantId" gorm:"index:tenantIndex,column:tenant_id;default:'';comment:租户ID"`
}

// TableName 定义表名称
func (Option) TableName() string {
	return "dict_options"
}

// GetPrimaryKey 定义表主键
func (Option) GetPrimaryKey() string {
	return "id"
}
