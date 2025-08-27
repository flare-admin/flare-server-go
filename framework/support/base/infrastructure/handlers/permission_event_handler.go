package handlers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	pkgEvents "github.com/flare-admin/flare-server-go/framework/pkg/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query/cache/handlers"
)

// PermissionEventHandler 权限事件处理器
type PermissionEventHandler struct {
	ch *handlers.CacheHandler
}

func NewPermissionEventHandler(ch *handlers.CacheHandler) *PermissionEventHandler {
	return &PermissionEventHandler{ch: ch}
}

func (h *PermissionEventHandler) Handle(ctx context.Context, event pkgEvents.Event) error {
	switch e := event.(type) {
	case *events.PermissionEvent:
		return h.handlePermissionEvent(ctx, e)
	default:
		return nil
	}
}

func (h *PermissionEventHandler) handlePermissionEvent(ctx context.Context, event *events.PermissionEvent) error {
	switch event.EventName() {
	case events.PermissionCreated:
		hlog.CtxDebugf(ctx, "权限创建事件: 租户ID=%s, 权限ID=%d", event.TenantID, event.PermID)
		return h.ch.InvalidatePermissionAllCache(ctx, event.PermID)
	case events.PermissionUpdated:
		hlog.CtxDebugf(ctx, "权限更新事件: 租户ID=%s, 权限ID=%d", event.TenantID, event.PermID)
		return h.ch.InvalidatePermissionAllCache(ctx, event.PermID)
	case events.PermissionDeleted:
		hlog.CtxDebugf(ctx, "权限删除事件: 租户ID=%s, 权限ID=%d", event.TenantID, event.PermID)
		return h.ch.InvalidatePermissionAllCache(ctx, event.PermID)
	case events.PermissionStatusChange:
		hlog.CtxDebugf(ctx, "权限状态变更事件: 租户ID=%s, 权限ID=%d", event.TenantID, event.PermID)
		return h.ch.InvalidatePermissionAllCache(ctx, event.PermID)
	}
	return nil
}
