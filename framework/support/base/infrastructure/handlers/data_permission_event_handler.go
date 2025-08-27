package handlers

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query/cache/handlers"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	pkgEvents "github.com/flare-admin/flare-server-go/framework/pkg/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/events"
)

// DataPermissionEventHandler 数据权限事件处理器
type DataPermissionEventHandler struct {
	ch *handlers.CacheHandler
}

func NewDataPermissionEventHandler(ch *handlers.CacheHandler) *DataPermissionEventHandler {
	return &DataPermissionEventHandler{ch: ch}
}

func (h *DataPermissionEventHandler) Handle(ctx context.Context, event pkgEvents.Event) error {
	switch e := event.(type) {
	case *events.DataPermissionEvent:
		return h.handleDataPermissionEvent(ctx, e)
	default:
		return nil
	}
}

func (h *DataPermissionEventHandler) handleDataPermissionEvent(ctx context.Context, event *events.DataPermissionEvent) error {
	switch event.EventName() {
	case events.DataPermissionAssigned:
		hlog.CtxDebugf(ctx, "数据权限分配事件: 租户ID=%s, 权限ID=%s", event.TenantID, event.Permission)
	case events.DataPermissionRemoved:
		hlog.CtxDebugf(ctx, "数据权限移除事件: 租户ID=%s, 权限ID=%s", event.TenantID, event.Permission)
	}
	return h.ch.InvalidatePermissionAllCache(ctx, event.Permission.ID)
}
