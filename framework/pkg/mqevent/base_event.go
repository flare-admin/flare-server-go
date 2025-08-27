package mqevent

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"time"
)

// BaseEvent 基础事件实现
type BaseEvent struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Data      interface{}       `json:"data"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
	TenantID  string            `json:"tenant_id"` // 租户ID
}

// BaseEventOption 事件选项函数类型
type BaseEventOption func(*BaseEvent)

// WithID 设置事件ID
func WithID(id string) BaseEventOption {
	return func(e *BaseEvent) {
		e.ID = id
	}
}

// WithTimestamp 设置事件时间戳
func WithTimestamp(timestamp time.Time) BaseEventOption {
	return func(e *BaseEvent) {
		e.Timestamp = timestamp
	}
}

// WithMetadata 设置事件元数据
func WithMetadata(metadata map[string]string) BaseEventOption {
	return func(e *BaseEvent) {
		e.Metadata = metadata
	}
}

// WithMetadataItem 设置单个元数据项
func WithMetadataItem(key, value string) BaseEventOption {
	return func(e *BaseEvent) {
		if e.Metadata == nil {
			e.Metadata = make(map[string]string)
		}
		e.Metadata[key] = value
	}
}

// WithTenantID 设置租户ID
func WithTenantID(tenantID string) BaseEventOption {
	return func(e *BaseEvent) {
		e.TenantID = tenantID
	}
}

// GetID 获取事件ID
func (b *BaseEvent) GetID() string {
	return b.ID
}

// GetType 获取事件类型
func (b *BaseEvent) GetType() string {
	return b.Type
}

// GetData 获取事件数据
func (b *BaseEvent) GetData() interface{} {
	return b.Data
}

// GetTimestamp 获取事件发生时间
func (b *BaseEvent) GetTimestamp() time.Time {
	return b.Timestamp
}

// GetMetadata 获取事件元数据
func (b *BaseEvent) GetMetadata() map[string]string {
	return b.Metadata
}

// SetMetadata 设置元数据
func (b *BaseEvent) SetMetadata(key, value string) {
	if b.Metadata == nil {
		b.Metadata = make(map[string]string)
	}
	b.Metadata[key] = value
}

// GetTenantID 获取租户ID
func (b *BaseEvent) GetTenantID() string {
	return b.TenantID
}

// SetTenantID 设置租户ID
func (b *BaseEvent) SetTenantID(tenantID string) {
	b.TenantID = tenantID
}

// NewBaseEvent 创建基础事件
func NewBaseEvent(eventType string, data interface{}, opts ...BaseEventOption) *BaseEvent {
	event := &BaseEvent{
		ID:        generateEventID(),
		Type:      eventType,
		Data:      data,
		Timestamp: utils.GetTimeNow(),
		Metadata:  make(map[string]string),
	}

	// 应用所有选项
	for _, opt := range opts {
		opt(event)
	}

	return event
}
