package queries

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/domain/model"
)

// GetConfigQuery 获取配置查询
type GetConfigQuery struct {
	ID string `json:"id" query:"id"` // 配置ID
}

// ListConfigsQuery 获取配置列表查询
type ListConfigsQuery struct {
	db_query.Page
	Key       string           `json:"key" query:"key"`               // 配置键
	Type      model.ConfigType `json:"type" query:"type"`             // 配置类型
	Group     string           `json:"group" query:"group"`           // 配置分组
	IsSystem  *bool            `json:"is_system" query:"is_system"`   // 是否系统配置
	IsEnabled *bool            `json:"is_enabled" query:"is_enabled"` // 是否启用
}

// GetConfigValueQuery 获取配置值查询
type GetConfigValueQuery struct {
	Key          string      `json:"key" query:"key"`                     // 配置键
	DefaultValue interface{} `json:"default_value" query:"default_value"` // 默认值
}

// GetConfigValueMapQuery 获取配置值映射查询
type GetConfigValueMapQuery struct {
	Keys         []string    `json:"keys" query:"keys"`                   // 配置键列表
	DefaultValue interface{} `json:"default_value" query:"default_value"` // 默认值
}

// GetConfigGroupQuery 获取配置分组查询
type GetConfigGroupQuery struct {
	ID string `json:"id" query:"id"` // 分组ID
}

// ListConfigGroupsQuery 获取配置分组列表查询
type ListConfigGroupsQuery struct {
	db_query.Page
	Name      string `json:"name" query:"name"`             // 分组名称
	Code      string `json:"code" query:"code"`             // 分组编码
	IsSystem  *bool  `json:"is_system" query:"is_system"`   // 是否系统分组
	IsEnabled *bool  `json:"is_enabled" query:"is_enabled"` // 是否启用
}
