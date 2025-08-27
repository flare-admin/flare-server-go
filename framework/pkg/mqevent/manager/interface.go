package manager

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/mqevent"
)

// EventType 事件类型
type EventType string

const (
	EventTypeNats  EventType = "nats"
	EventTypeNSQ   EventType = "nsq"
	EventTypeKafka EventType = "kafka"
)

const (
	StatusAdd     = 1
	StatusEnable  = 2
	StatusDisable = 3
)

// EventHandlerFunc 定义事件处理函数类型
type EventHandlerFunc func(*mqevent.EventContext) error

// Handle 实现 EventHandler 接口
func (f EventHandlerFunc) Handle(ctx *mqevent.EventContext) error {
	return f(ctx)
}

// EventHandler 定义事件处理器接口
type EventHandler interface {
	// Handle 处理事件
	Handle(ctx *mqevent.EventContext) error
}

// ISubscribeSmServerApi 订阅管理接口提供者
type ISubscribeSmServerApi interface {
	//GetByStatus 根据状态获取订阅
	GetByStatus(ctx context.Context, status int32) ([]*Subscribe, error)
	// GetParameters 获取订阅参数
	GetParameters(ctx context.Context, topic, channel string) (map[string]interface{}, error)
	// DadEventSave 保存死信事件
	DadEventSave(ctx context.Context, event *mqevent.DeadLetterEvent) error
}

// EventManager 事件管理器接口
type EventManager interface {
	// RegisterSubscribeHandel 注册处理器
	RegisterSubscribeHandel(topic, channel string, handler EventHandler) error
	// RegisterIdempotenceSubscribeHandel 注册幂等性订阅
	RegisterIdempotenceSubscribeHandel(topic, channel string, handler EventHandler) error
	// ActivateSubscription 激活订阅
	ActivateSubscription(topic, channel string) error
	//DeactivateSubscription 停止订阅
	DeactivateSubscription(topic, channel string) error

	// 死信处理
	RetryDeadLetter(ctx context.Context, dead *mqevent.DeadLetterEvent) error

	// 生命周期管理
	Start() error
	Stop() error
	Close() error
}
