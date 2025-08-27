package model

import (
	"encoding/json"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/mqevent"
	"time"
)

// DeadLetterSubscribe 死信订阅模型
type DeadLetterSubscribe struct {
	database.BaseIntTime
	Id          string            `gorm:"column:id;primary_key" json:"id"`                                         // 主键ID
	Name        string            `json:"name" gorm:"column:name;size:255;comment:订阅名称"`                           // 订阅名称
	MsgId       string            `json:"msgId" gorm:"column:msg_id;size:255;comment:消息ID"`                        // 消息ID
	Topic       string            `json:"topic" gorm:"column:topic;size:150;not null;comment:事件主题"`                // 事件主题
	Channel     string            `json:"channel" gorm:"column:channel;size:150;not null;comment:消费通道"`            // 消费通道
	EventType   string            `json:"eventType" gorm:"column:event_type;size:150;comment:事件类型"`                // 事件类型
	Data        []byte            `json:"data" gorm:"column:data;type:bytea;comment:事件数据"`                         // 事件数据
	Error       string            `json:"error" gorm:"column:error;type:text;comment:错误信息"`                        // 错误信息
	RetryCount  int               `json:"retryCount" gorm:"column:retry_count;default:0;comment:重试次数"`             // 重试次数
	LastAttempt time.Time         `json:"lastAttempt" gorm:"column:last_attempt;comment:最后尝试时间"`                   // 最后尝试时间
	NextRetry   time.Time         `json:"nextRetry" gorm:"column:next_retry;comment:下次重试时间"`                       // 下次重试时间
	Status      int8              `gorm:"column:status;default:1;comment:状态 1->待处理, 2->已处理，3->处理失败" json:"status"` // 处理状态
	EventId     string            `json:"eventId" gorm:"column:event_id;size:255;comment:事件ID"`                    // 事件ID
	Timestamp   time.Time         `json:"timestamp" gorm:"column:timestamp;comment:事件发生时间"`                        // 事件发生时间
	Metadata    map[string]string `json:"metadata" gorm:"column:metadata;type:jsonb;comment:事件元数据"`                // 事件元数据
	TenantID    string            `json:"tenantId" gorm:"column:tenant_id;size:255;comment:租户ID"`                  // 租户ID
}

// TableName 指定表名
func (DeadLetterSubscribe) TableName() string {
	return "sys_event_dead_letter_subscribe"
}

// GetPrimaryKey 获取主键
func (DeadLetterSubscribe) GetPrimaryKey() string {
	return "id"
}

// FromDeadLetterEvent 从 DeadLetterEvent 转换为 DeadLetterSubscribe
func (d *DeadLetterSubscribe) FromDeadLetterEvent(event *mqevent.DeadLetterEvent) error {
	// 设置基本信息
	d.MsgId = event.GetID()
	d.Topic = event.OriginalTopic
	d.Channel = event.Channel
	d.EventType = event.OriginalEvent.GetType()
	d.Error = event.Error
	d.RetryCount = event.RetryCount
	d.LastAttempt = event.LastAttempt
	d.NextRetry = event.NextRetry
	d.Status = 1 // 默认状态为待处理

	// 存储原始事件的所有信息
	d.EventId = event.OriginalEvent.GetID()
	d.Timestamp = event.OriginalEvent.GetTimestamp()
	d.Metadata = event.OriginalEvent.GetMetadata()
	d.TenantID = event.OriginalEvent.GetTenantID()

	// 序列化事件数据
	data, err := json.Marshal(event.OriginalEvent.GetData())
	if err != nil {
		return err
	}
	d.Data = data

	return nil
}

// ToDeadLetterEvent 转换为 DeadLetterEvent
func (d *DeadLetterSubscribe) ToDeadLetterEvent() (*mqevent.DeadLetterEvent, error) {
	// 反序列化事件数据
	var eventData interface{}
	if len(d.Data) > 0 {
		if err := json.Unmarshal(d.Data, &eventData); err != nil {
			return nil, err
		}
	}

	// 创建基础事件
	baseEvent := mqevent.NewBaseEvent(
		d.EventType,
		eventData,
		mqevent.WithID(d.EventId),
		mqevent.WithTimestamp(d.Timestamp),
		mqevent.WithMetadata(d.Metadata),
		mqevent.WithTenantID(d.TenantID),
	)

	// 创建死信事件
	event := &mqevent.DeadLetterEvent{
		Error:         d.Error,
		RetryCount:    d.RetryCount,
		LastAttempt:   d.LastAttempt,
		NextRetry:     d.NextRetry,
		EventType:     d.EventType,
		OriginalTopic: d.Topic,
		Channel:       d.Channel,
		OriginalEvent: baseEvent,
	}

	return event, nil
}
