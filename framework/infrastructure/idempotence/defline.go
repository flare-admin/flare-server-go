package idempotence

import (
	"time"
)

type IdempotencyRecord struct {
	ID        uint      `gorm:"primaryKey"`                                                     // 自增主键
	Topic     string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_topic_group_message"` // 主题
	Channel   string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_topic_group_message"` // 消费者组
	MessageID string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_topic_group_message"` // 消息ID
	Status    string    `gorm:"type:varchar(50)"`                                               // 消息处理状态
	CreatedAt time.Time `gorm:"type:timestamp;default:now();not null"`                          // 消息处理时间
}

// TableName  为了确保 (Topic, ConsumerGroup, MessageID) 的唯一性，可以定义如下索引
func (IdempotencyRecord) TableName() string {
	return "idempotency_records"
}
