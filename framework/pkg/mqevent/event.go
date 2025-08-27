package mqevent

import (
	"time"
)

// EventOptions 事件选项
type EventOptions struct {
	// RetryCount 重试次数
	RetryCount int
	// RetryDelay 重试延迟
	RetryDelay time.Duration
	// Timeout 超时时间
	Timeout time.Duration
}

// EventOption 事件选项函数
type EventOption func(*EventOptions)

// WithRetry 设置重试选项
func WithRetry(count int, delay time.Duration) EventOption {
	return func(o *EventOptions) {
		o.RetryCount = count
		o.RetryDelay = delay
	}
}

// WithTimeout 设置超时选项
func WithTimeout(timeout time.Duration) EventOption {
	return func(o *EventOptions) {
		o.Timeout = timeout
	}
}

// NewEvent 创建新事件
func NewEvent(eventType string, data interface{}, opts ...EventOption) Event {
	options := &EventOptions{
		RetryCount: 3,
		RetryDelay: time.Second,
		Timeout:    time.Second * 30,
	}

	for _, opt := range opts {
		opt(options)
	}

	return NewBaseEvent(eventType, data)
}

// 生成事件ID
func generateEventID() string {
	return time.Now().Format("20060102150405.000000000")
}
