package errors

import "github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"

var (
	// 配置相关错误
	AddConfigFail         = herrors.NewServerError("AddConfigFail")                 // 新增配置失败
	EditConfigFail        = herrors.NewServerError("EditConfigFail")                // 修改配置失败
	DeleteConfigFail      = herrors.NewServerError("DeleteConfigFail")              // 删除配置失败
	GetConfigFail         = herrors.NewServerError("GetConfigFail")                 // 获取配置失败
	ConfigNotExistFail    = herrors.NewBusinessServerError("ConfigNotExistFail")    // 配置不存在
	ConfigKeyExistFail    = herrors.NewBusinessServerError("ConfigKeyExistFail")    // 配置键已存在
	ConfigNotEnableFail   = herrors.NewBusinessServerError("ConfigNotEnableFail")   // 配置未启用
	ConfigTypeInvalidFail = herrors.NewBusinessServerError("ConfigTypeInvalidFail") // 配置类型无效

	// 配置分组相关错误
	AddConfigGroupFail       = herrors.NewServerError("AddConfigGroupFail")               // 新增配置分组失败
	EditConfigGroupFail      = herrors.NewServerError("EditConfigGroupFail")              // 修改配置分组失败
	DeleteConfigGroupFail    = herrors.NewServerError("DeleteConfigGroupFail")            // 删除配置分组失败
	GetConfigGroupFail       = herrors.NewServerError("GetConfigGroupFail")               // 获取配置分组失败
	ConfigGroupNotExistFail  = herrors.NewBusinessServerError("ConfigGroupNotExistFail")  // 配置分组不存在
	ConfigGroupCodeExistFail = herrors.NewBusinessServerError("ConfigGroupCodeExistFail") // 配置分组编码已存在
	ConfigGroupNotEnableFail = herrors.NewBusinessServerError("ConfigGroupNotEnableFail") // 配置分组未启用

	// 缓存相关错误
	SetConfigCacheFail    = herrors.NewServerError("SetConfigCacheFail")    // 设置配置缓存失败
	GetConfigCacheFail    = herrors.NewServerError("GetConfigCacheFail")    // 获取配置缓存失败
	DeleteConfigCacheFail = herrors.NewServerError("DeleteConfigCacheFail") // 删除配置缓存失败
	SetGroupCacheFail     = herrors.NewServerError("SetGroupCacheFail")     // 设置分组缓存失败
	GetGroupCacheFail     = herrors.NewServerError("GetGroupCacheFail")     // 获取分组缓存失败
	DeleteGroupCacheFail  = herrors.NewServerError("DeleteGroupCacheFail")  // 删除分组缓存失败
)
