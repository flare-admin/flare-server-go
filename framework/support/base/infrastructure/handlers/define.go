package handlers

import (
	pkgEvent "github.com/flare-admin/flare-server-go/framework/pkg/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/events"
)

type HandlerEvent struct {
	uh       *UserEventHandler
	rh       *RoleEventHandler
	dh       *DepartmentEventHandler
	ph       *PermissionEventHandler
	dph      *DataPermissionEventHandler
	th       *TenantEventHandler
	eventBus pkgEvent.IEventBus
}

func NewHandlerEvent(
	eventBus pkgEvent.IEventBus,
	uh *UserEventHandler,
	rh *RoleEventHandler,
	dh *DepartmentEventHandler,
	ph *PermissionEventHandler,
	dph *DataPermissionEventHandler,
	th *TenantEventHandler,
) *HandlerEvent {
	return &HandlerEvent{
		uh:       uh,
		rh:       rh,
		dh:       dh,
		ph:       ph,
		dph:      dph,
		th:       th,
		eventBus: eventBus,
	}
}

func (h *HandlerEvent) Register() {
	// 注册用户相关事件
	h.eventBus.Subscribe(events.UserLoggedIn, h.uh)
	h.eventBus.Subscribe(events.UserCreated, h.uh)
	h.eventBus.Subscribe(events.UserUpdated, h.uh)
	h.eventBus.Subscribe(events.UserDeleted, h.uh)
	h.eventBus.Subscribe(events.UserRoleChanged, h.uh)

	// 角色事件
	h.eventBus.Subscribe(events.RoleCreated, h.rh)
	h.eventBus.Subscribe(events.RoleUpdated, h.rh)
	h.eventBus.Subscribe(events.RoleDeleted, h.rh)
	h.eventBus.Subscribe(events.RolePermissionsChanged, h.rh)

	// 部门事件
	h.eventBus.Subscribe(events.DepartmentCreated, h.dh)
	h.eventBus.Subscribe(events.DepartmentUpdated, h.dh)
	h.eventBus.Subscribe(events.DepartmentDeleted, h.dh)
	h.eventBus.Subscribe(events.DepartmentMoved, h.dh)
	h.eventBus.Subscribe(events.UserAssigned, h.dh)
	h.eventBus.Subscribe(events.UserRemoved, h.dh)
	h.eventBus.Subscribe(events.UserTransferred, h.dh)

	// 权限事件
	h.eventBus.Subscribe(events.PermissionCreated, h.ph)
	h.eventBus.Subscribe(events.PermissionUpdated, h.ph)
	h.eventBus.Subscribe(events.PermissionDeleted, h.ph)
	h.eventBus.Subscribe(events.PermissionStatusChange, h.ph)

	// 数据权限事件
	h.eventBus.Subscribe(events.DataPermissionAssigned, h.dph)
	h.eventBus.Subscribe(events.DataPermissionRemoved, h.dph)

	// 租户事件
	h.eventBus.Subscribe(events.TenantCreated, h.th)
	h.eventBus.Subscribe(events.TenantUpdated, h.th)
	h.eventBus.Subscribe(events.TenantDeleted, h.th)
	h.eventBus.Subscribe(events.TenantLocked, h.th)
	h.eventBus.Subscribe(events.TenantUnlocked, h.th)
}
