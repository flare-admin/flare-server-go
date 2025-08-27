package domain

import (
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	service.NewAuthService,
	service.NewRoleCommandService,
	service.NewPermissionService,
	service.NewTenantCommandService,
	service.NewDepartmentService,
	service.NewUserCommandService,
	service.NewDataPermissionService,
)
