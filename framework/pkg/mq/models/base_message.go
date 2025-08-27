package models

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
)

// BaseMessage 基础消息实现
type BaseMessage struct {
	ID        string            `json:"id"`
	Topic     string            `json:"topic"`
	Payload   string            `json:"payload"`
	Timestamp int64             `json:"timestamp"` // 纳秒时间戳
	Headers   map[string]string `json:"headers"`
}

// todo nsqd的消息结构
//
//	type T struct {
//		Id        string    `json:"id"`
//		Topic     string    `json:"topic"`
//		Payload   string    `json:"payload"`
//		Timestamp time.Time `json:"timestamp"`
//		Headers   struct {
//			EventId   string    `json:"event_id"`
//			EventType string    `json:"event_type"`
//			TenantId  string    `json:"tenant_id"`
//			Timestamp time.Time `json:"timestamp"`
//		} `json:"headers"`
//	}
//
// NewBaseMessage 创建基础消息
func NewBaseMessage(topic string, payload []byte, headers map[string]string) *BaseMessage {
	return &BaseMessage{
		ID:        GenerateID(),
		Topic:     topic,
		Payload:   string(payload),
		Timestamp: utils.GetTimeNow().UnixNano(), // 纳秒时间戳
		Headers:   headers,
	}
}

// GetID 获取消息ID
func (m *BaseMessage) GetID() string {
	return m.ID
}

// GetTopic 获取主题
func (m *BaseMessage) GetTopic() string {
	return m.Topic
}

// GetPayload 获取消息内容
func (m *BaseMessage) GetPayload() []byte {
	return []byte(m.Payload)
}

// GetTimestamp 获取时间戳
func (m *BaseMessage) GetTimestamp() int64 {
	return m.Timestamp
}

// GetHeaders 获取消息头
func (m *BaseMessage) GetHeaders() map[string]string {
	return m.Headers
}

func (m *BaseMessage) GetHeadersItem(key string) (string, bool) {
	v, ok := m.Headers[key]
	return v, ok
}

// GenerateID 生成消息ID
func GenerateID() string {
	return utils.GetTimeNow().Format("20060102150405.000000000")
}
