package mqevent

import (
	"time"
)

// DeadLetterEvent 死信事件结构
type DeadLetterEvent struct {
	// OriginalEvent 原始事件
	OriginalEvent Event `json:"original_event"`
	// Error 错误信息
	Error string `json:"error"`
	// RetryCount 重试次数
	RetryCount int `json:"retry_count"`
	// LastAttempt 最后尝试时间
	LastAttempt time.Time `json:"last_attempt"`
	// NextRetry 下次重试时间
	NextRetry time.Time `json:"next_retry"`
	// EventType 原始事件类型
	EventType string `json:"event_type"`
	// OriginalTopic 原始主题
	OriginalTopic string `json:"original_topic"`
	// Channel 消费通道
	Channel string `json:"channel"`
}

// GetID 获取事件ID
func (d *DeadLetterEvent) GetID() string {
	return d.OriginalEvent.GetID()
}

// GetType 获取事件类型
func (d *DeadLetterEvent) GetType() string {
	return d.EventType
}

// GetData 获取事件数据
func (d *DeadLetterEvent) GetData() interface{} {
	return d.OriginalEvent.GetData()
}

// GetTimestamp 获取事件发生时间
func (d *DeadLetterEvent) GetTimestamp() time.Time {
	return d.OriginalEvent.GetTimestamp()
}

// GetMetadata 获取事件元数据
func (d *DeadLetterEvent) GetMetadata() map[string]string {
	return d.OriginalEvent.GetMetadata()
}

// SetMetadata 设置元数据
func (d *DeadLetterEvent) SetMetadata(key, value string) {
	d.OriginalEvent.SetMetadata(key, value)
}
