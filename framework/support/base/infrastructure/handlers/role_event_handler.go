package handlers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	pkgEvents "github.com/flare-admin/flare-server-go/framework/pkg/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query/cache/handlers"
)

// RoleEventHandler 角色事件处理器
type RoleEventHandler struct {
	ch *handlers.CacheHandler
}

func NewRoleEventHandler(ch *handlers.CacheHandler) *RoleEventHandler {
	return &RoleEventHandler{ch: ch}
}

func (h *RoleEventHandler) Handle(ctx context.Context, event pkgEvents.Event) error {
	switch e := event.(type) {
	case *events.RoleEvent:
		return h.handleRoleEvent(ctx, e)
	case *events.RolePermissionsAssignedEvent:
		return h.handleRolePermissionsAssignedEvent(ctx, e)
	default:
		return nil
	}
}

func (h *RoleEventHandler) handleRoleEvent(ctx context.Context, event *events.RoleEvent) error {
	switch event.EventName() {
	case events.RoleCreated:
		hlog.CtxDebugf(ctx, "角色创建事件: 租户ID=%s, 角色ID=%d", event.TenantID, event.RoleID)
		return h.ch.InvalidateRoleAllCache(ctx, event.RoleID)
	case events.RoleUpdated:
		hlog.CtxDebugf(ctx, "角色更新事件: 租户ID=%s, 角色ID=%d", event.TenantID, event.RoleID)
		return h.ch.InvalidateRoleAllCache(ctx, event.RoleID)
	case events.RoleDeleted:
		hlog.CtxDebugf(ctx, "角色删除事件: 租户ID=%s, 角色ID=%d", event.TenantID, event.RoleID)
		return h.ch.InvalidateRoleAllCache(ctx, event.RoleID)
	case events.RolePermissionsChanged:
		hlog.CtxDebugf(ctx, "角色权限变更事件: 租户ID=%s, 角色ID=%d", event.TenantID, event.RoleID)
		return h.ch.InvalidateRoleAllCache(ctx, event.RoleID)
	}
	return nil
}

func (h *RoleEventHandler) handleRolePermissionsAssignedEvent(ctx context.Context, event *events.RolePermissionsAssignedEvent) error {
	return h.ch.InvalidateRoleAllCache(ctx, event.RoleID)
}
