package handlers

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserEventHandler,
	NewDataPermissionEventHandler,
	NewDepartmentEventHandler,
	NewPermissionEventHandler,
	NewRoleEventHandler,
	NewTenantEventHandler,
	NewHandlerEvent,
)
