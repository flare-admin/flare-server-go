package handlers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	pkgEvents "github.com/flare-admin/flare-server-go/framework/pkg/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query/cache/handlers"
)

// TenantEventHandler 租户事件处理器
type TenantEventHandler struct {
	ch *handlers.CacheHandler
}

func NewTenantEventHandler(ch *handlers.CacheHandler) *TenantEventHandler {
	return &TenantEventHandler{ch: ch}
}

func (h *TenantEventHandler) Handle(ctx context.Context, event pkgEvents.Event) error {
	switch e := event.(type) {
	case *events.TenantEvent:
		return h.handleTenantEvent(ctx, e)
	case *events.TenantPermissionEvent:
		return h.handleTenantPermissionEvent(ctx, e)
	default:
		return nil
	}
}

func (h *TenantEventHandler) handleTenantEvent(ctx context.Context, event *events.TenantEvent) error {
	switch event.EventName() {
	case events.TenantCreated:
		hlog.CtxDebugf(ctx, "租户创建事件: 租户ID=%s", event.TenantID)
	case events.TenantUpdated:
		hlog.CtxDebugf(ctx, "租户更新事件: 租户ID=%s", event.TenantID)
	case events.TenantDeleted:
		hlog.CtxDebugf(ctx, "租户删除事件: 租户ID=%s", event.TenantID)
	case events.TenantLocked, events.TenantUnlocked:
		hlog.CtxDebugf(ctx, "租户状态变更事件: 租户ID=%s", event.TenantID)
	}
	return h.ch.InvalidateTenantAllCache(ctx, event.TenantID)
}

// 添加处理租户权限事件的方法
func (h *TenantEventHandler) handleTenantPermissionEvent(ctx context.Context, event *events.TenantPermissionEvent) error {
	hlog.CtxDebugf(ctx, "租户权限变更事件: 租户ID=%s, 权限IDs=%v", event.TenantID, event.PermissionIDs)
	// 清除租户权限相关的缓存
	return h.ch.InvalidateTenantAllCache(ctx, event.TenantID)
}
