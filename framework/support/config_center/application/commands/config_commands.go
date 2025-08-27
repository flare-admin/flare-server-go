package commands

import (
	"github.com/flare-admin/flare-server-go/framework/support/config_center/domain/model"
)

// CreateConfigCommand 创建配置命令
type CreateConfigCommand struct {
	Name        string           `json:"name"`        //  配置名称
	Key         string           `json:"key"`         // 配置键
	Value       string           `json:"value"`       // 配置值
	Type        model.ConfigType `json:"type"`        // 配置类型
	Group       string           `json:"group"`       // 配置分组
	Description string           `json:"description"` // 配置描述
	I18nKey     string           `json:"i18n_key"`    // 国际化键
	IsSystem    bool             `json:"is_system"`   // 是否系统配置
	IsEnabled   bool             `json:"is_enabled"`  // 是否启用
	Sort        int              `json:"sort"`        // 排序
}

// UpdateConfigCommand 更新配置命令
type UpdateConfigCommand struct {
	ID          string           `json:"id"`          // 配置ID
	Name        string           `json:"name"`        //  配置名称
	Key         string           `json:"key"`         // 配置键
	Value       string           `json:"value"`       // 配置值
	Type        model.ConfigType `json:"type"`        // 配置类型
	Group       string           `json:"group"`       // 配置分组
	Description string           `json:"description"` // 配置描述
	I18nKey     string           `json:"i18n_key"`    // 国际化键
	IsSystem    bool             `json:"is_system"`   // 是否系统配置
	IsEnabled   bool             `json:"is_enabled"`  // 是否启用
	Sort        int              `json:"sort"`        // 排序
}

// DeleteConfigCommand 删除配置命令
type DeleteConfigCommand struct {
	ID string `json:"id"` // 配置ID
}

// UpdateConfigStatusCommand 更新配置状态命令
type UpdateConfigStatusCommand struct {
	ID        string `json:"id"`         // 配置ID
	IsEnabled bool   `json:"is_enabled"` // 是否启用
}

// CreateConfigGroupCommand 创建配置分组命令
type CreateConfigGroupCommand struct {
	Name        string `json:"name"`        // 分组名称
	Code        string `json:"code"`        // 分组编码
	Description string `json:"description"` // 分组描述
	I18nKey     string `json:"i18n_key"`    // 国际化键
	IsSystem    bool   `json:"is_system"`   // 是否系统分组
	IsEnabled   bool   `json:"is_enabled"`  // 是否启用
	Sort        int    `json:"sort"`        // 排序
}

// UpdateConfigGroupCommand 更新配置分组命令
type UpdateConfigGroupCommand struct {
	ID          string `json:"id"`          // 分组ID
	Name        string `json:"name"`        // 分组名称
	Code        string `json:"code"`        // 分组编码
	Description string `json:"description"` // 分组描述
	I18nKey     string `json:"i18n_key"`    // 国际化键
	IsSystem    bool   `json:"is_system"`   // 是否系统分组
	IsEnabled   bool   `json:"is_enabled"`  // 是否启用
	Sort        int    `json:"sort"`        // 排序
}

// DeleteConfigGroupCommand 删除配置分组命令
type DeleteConfigGroupCommand struct {
	ID string `json:"id"` // 分组ID
}

// UpdateConfigGroupStatusCommand 更新配置分组状态命令
type UpdateConfigGroupStatusCommand struct {
	ID        string `json:"id"`         // 分组ID
	IsEnabled bool   `json:"is_enabled"` // 是否启用
}
