package handlers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	pkgEvents "github.com/flare-admin/flare-server-go/framework/pkg/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query/cache/handlers"
)

// DepartmentEventHandler 部门事件处理器
type DepartmentEventHandler struct {
	ch *handlers.CacheHandler
}

func NewDepartmentEventHandler(ch *handlers.CacheHandler) *DepartmentEventHandler {
	return &DepartmentEventHandler{ch: ch}
}

func (h *DepartmentEventHandler) Handle(ctx context.Context, event pkgEvents.Event) error {
	switch e := event.(type) {
	case *events.DepartmentEvent:
		return h.handleDepartmentEvent(ctx, e)
	case *events.DepartmentMovedEvent:
		return h.handleDepartmentMovedEvent(ctx, e)
	case *events.UserAssignedEvent:
		return h.handleUserAssignedEvent(ctx, e)
	case *events.UserRemovedEvent:
		return h.handleUserRemovedEvent(ctx, e)
	case *events.UserTransferredEvent:
		return h.handleUserTransferredEvent(ctx, e)
	default:
		return nil
	}
}

func (h *DepartmentEventHandler) handleDepartmentEvent(ctx context.Context, event *events.DepartmentEvent) error {
	switch event.EventName() {
	case events.DepartmentCreated:
		hlog.CtxDebugf(ctx, "部门创建事件: 租户ID=%s, 部门ID=%s", event.TenantID, event.DeptID)
		return h.ch.InvalidateDepartmentAllCache(ctx, event.DeptID)
	case events.DepartmentUpdated:
		hlog.CtxDebugf(ctx, "部门更新事件: 租户ID=%s, 部门ID=%s", event.TenantID, event.DeptID)
		return h.ch.InvalidateDepartmentAllCache(ctx, event.DeptID)
	case events.DepartmentDeleted:
		hlog.CtxDebugf(ctx, "部门删除事件: 租户ID=%s, 部门ID=%s", event.TenantID, event.DeptID)
		return h.ch.InvalidateDepartmentAllCache(ctx, event.DeptID)
	case events.DepartmentMoved:
		hlog.CtxDebugf(ctx, "部门移动事件: 租户ID=%s, 部门ID=%s", event.TenantID, event.DeptID)
		return h.ch.InvalidateDepartmentAllCache(ctx, event.DeptID)
	}
	return nil
}

func (h *DepartmentEventHandler) handleDepartmentMovedEvent(ctx context.Context, event *events.DepartmentMovedEvent) error {

	return h.ch.InvalidateDepartmentAllCache(ctx, event.DeptID)
}

func (h *DepartmentEventHandler) handleUserAssignedEvent(ctx context.Context, event *events.UserAssignedEvent) error {
	for _, v := range event.UserIDs {
		err := h.ch.InvalidateUserAllCache(ctx, v)
		if err != nil {
			hlog.CtxErrorf(ctx, "部门删除事件: 租户ID=%s, 部门ID=%s", event.TenantID, event.DeptID)
			return err
		}
	}
	return h.ch.InvalidateDepartmentAllCache(ctx, event.DeptID)
}

func (h *DepartmentEventHandler) handleUserRemovedEvent(ctx context.Context, event *events.UserRemovedEvent) error {
	for _, v := range event.UserIDs {
		err := h.ch.InvalidateUserAllCache(ctx, v)
		if err != nil {
			hlog.CtxErrorf(ctx, "部门删除事件: 租户ID=%s, 部门ID=%s", event.TenantID, event.DeptID)
		}
	}
	return h.ch.InvalidateDepartmentAllCache(ctx, event.DeptID)
}

func (h *DepartmentEventHandler) handleUserTransferredEvent(ctx context.Context, event *events.UserTransferredEvent) error {
	err := h.ch.InvalidateUserAllCache(ctx, event.UserID)
	if err != nil {
		hlog.CtxErrorf(ctx, "部门删除事件: 租户ID=%s, 部门ID=%s", event.TenantID, event.DeptID)
		return err
	}
	err = h.ch.InvalidateDepartmentAllCache(ctx, event.FromDeptID)
	if err != nil {
		hlog.CtxErrorf(ctx, "部门删除事件: 租户ID=%s, 部门ID=%s", event.TenantID, event.DeptID)
		return err
	}
	return h.ch.InvalidateDepartmentAllCache(ctx, event.DeptID)
}
