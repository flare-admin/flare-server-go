package mqevent

import (
	"context"
	"time"
)

// Event 定义事件接口
type Event interface {
	// GetID 获取事件ID
	GetID() string
	// GetType 获取事件类型
	GetType() string
	// GetData 获取事件数据
	GetData() interface{}
	// GetTimestamp 获取事件发生时间
	GetTimestamp() time.Time
	// GetMetadata 获取事件元数据
	GetMetadata() map[string]string
	// SetMetadata 设置事件元数据
	SetMetadata(key, value string)
	// GetTenantID 获取租户ID
	GetTenantID() string
	// SetTenantID 设置租户ID
	SetTenantID(tenantID string)
}

// EventHandlerFunc 定义事件处理函数类型
type EventHandlerFunc func(ctx context.Context, event Event) error

// Handle 实现 EventHandler 接口
func (f EventHandlerFunc) Handle(ctx context.Context, event Event) error {
	return f(ctx, event)
}

// EventHandler 定义事件处理器接口
type EventHandler interface {
	// Handle 处理事件
	Handle(ctx context.Context, event Event) error
}

// DeadLetterEventHandlerFunc 定义死信队列处理函数类型
type DeadLetterEventHandlerFunc func(ctx context.Context, event *DeadLetterEvent) error

// IMQEventBus 定义事件总线接口
type IMQEventBus interface {
	// Publish 发布事件
	Publish(ctx context.Context, event Event) error
	// Subscribe 订阅事件
	// channel 用于区分不同的业务场景，例如：
	// - "user-init" 用于用户注册后的初始化
	// - "user-notify" 用于用户注册后的通知
	// - "user-sync" 用于用户数据同步
	Subscribe(eventType string, channel string, handler EventHandler) (string, error)
	// Unsubscribe 取消订阅
	Unsubscribe(subscriptionID string) error
	// Close 关闭事件总线
	Close() error
	// SubscribeDeadLetter 订阅死信队列
	SubscribeDeadLetter(handler DeadLetterEventHandlerFunc) (string, error)
}
