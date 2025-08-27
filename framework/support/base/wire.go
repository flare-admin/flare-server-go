package base

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/lua_engine"
	"github.com/flare-admin/flare-server-go/framework/support/base/application/handlers"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure"
	base_api "github.com/flare-admin/flare-server-go/framework/support/base/interfaces/api"
	"github.com/flare-admin/flare-server-go/framework/support/base/interfaces/rest"
	"github.com/google/wire"
)

// ProviderSet is monitoring providers.
var ProviderSet = wire.NewSet(
	domain.ProviderSet,
	handlers.ProviderSet,
	infrastructure.ProviderSet,
	rest.NewSysRoleController,
	rest.NewSysUserController,
	rest.NewSysTenantController,
	rest.NewSysPermissionsController,
	rest.NewAuthController,
	rest.NewLoginLogController,
	rest.NewOperationLogController,
	rest.NewDepartmentController,
	rest.NewDataPermissionController,
	lua_engine.NewRuleExecutorWithDB,
	base_api.NewTenantApi,
	base_api.NewSysUserApi,
	NewBaseServer,
)

// BaseProviderSet 基础依赖
var BaseProviderSet = wire.NewSet(
	base_api.NewTenantApi,
	base_api.NewSysUserApi,
	infrastructure.ProviderSet,
	lua_engine.NewRuleExecutorWithDB,
)
