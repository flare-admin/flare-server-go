package handlers

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	pkgEvents "github.com/flare-admin/flare-server-go/framework/pkg/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/events"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query/cache/handlers"
)

// UserEventHandler 用户事件处理器
type UserEventHandler struct {
	ch *handlers.CacheHandler
}

func NewUserEventHandler(ch *handlers.CacheHandler) *UserEventHandler {
	return &UserEventHandler{ch: ch}
}

// Handle 处理事件
func (h *UserEventHandler) Handle(ctx context.Context, event pkgEvents.Event) error {
	switch e := event.(type) {
	case *events.UserEvent:
		return h.handleUserEvent(ctx, e)
	default:
		return nil
	}
}

// handleUserEvent 处理用户基础事件
func (h *UserEventHandler) handleUserEvent(ctx context.Context, event *events.UserEvent) error {
	switch event.EventName() {
	case events.UserCreated:
		hlog.CtxDebugf(ctx, "用户创建事件: 租户ID=%s, 用户ID=%s", event.TenantID, event.UserID)
		return nil
	case events.UserUpdated:
		hlog.CtxDebugf(ctx, "用户更新事件: 租户ID=%s, 用户ID=%s", event.TenantID, event.UserID)
	case events.UserDeleted:
		hlog.CtxDebugf(ctx, "用户删除事件: 租户ID=%s, 用户ID=%s", event.TenantID, event.UserID)
	}
	return h.ch.InvalidateUserAllCache(ctx, event.UserID)
}

// handleUserDeleted 处理用户删除事件
func (h *UserEventHandler) handleUserDeleted(ctx context.Context, event *events.UserEvent) error {
	return nil
}

// cleanUserCache 清除用户缓存
func (h *UserEventHandler) cleanUserCache(ctx context.Context, userID string) error {
	return nil
}

// cleanUserRelatedCache 清除用户相关的其他数据缓存
func (h *UserEventHandler) cleanUserRelatedCache(ctx context.Context, userID string) error {

	return nil
}
