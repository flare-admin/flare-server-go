package entity

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
)

// Config 配置实体
type Config struct {
	database.BaseModel
	ID          string `gorm:"column:id;primaryKey;pk:true"`
	Name        string `gorm:"column:name;type:varchar(50)"`
	Key         string `gorm:"column:key;type:varchar(100);not null;uniqueIndex"`
	Value       string `gorm:"column:value;type:text;not null"`
	Type        string `gorm:"column:type;type:varchar(20);not null"`
	Group       string `gorm:"column:group_id;type:varchar(50);index"`
	Description string `gorm:"column:description;type:varchar(200)"`
	I18nKey     string `gorm:"column:i18n_key;type:varchar(100)"`
	IsSystem    bool   `gorm:"column:is_system;type:boolean;not null;default:false"`
	IsEnabled   bool   `gorm:"column:is_enabled;type:boolean;not null;default:true"`
	Sort        int    `gorm:"column:sort;type:integer;not null;default:0"`
}

func (g Config) GetPrimaryKey() string {
	return "id"
}

// TableName 表名
func (Config) TableName() string {
	return "configs"
}

// ConfigGroup 配置分组实体
type ConfigGroup struct {
	database.BaseModel
	ID          string `gorm:"column:id;primaryKey;pk:true"`
	Name        string `gorm:"column:name;type:varchar(50);not null"`
	Code        string `gorm:"column:code;type:varchar(50);not null;uniqueIndex"`
	Description string `gorm:"column:description;type:varchar(200)"`
	I18nKey     string `gorm:"column:i18n_key;type:varchar(100)"`
	IsSystem    bool   `gorm:"column:is_system;type:boolean;not null;default:false"`
	IsEnabled   bool   `gorm:"column:is_enabled;type:boolean;not null;default:true"`
	Sort        int    `gorm:"column:sort;type:integer;not null;default:0"`
}

func (g ConfigGroup) GetPrimaryKey() string {
	return "id"
}

// TableName 表名
func (ConfigGroup) TableName() string {
	return "config_groups"
}
