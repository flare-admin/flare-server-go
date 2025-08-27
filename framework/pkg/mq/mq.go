package mq

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/models"
	"time"
)

// Producer 生产者接口
type Producer interface {
	// Publish 发布消息
	Publish(ctx context.Context, topic string, payload []byte, headers map[string]string) error
	// PublishDelay 发布延迟消息
	PublishDelay(ctx context.Context, topic string, payload []byte, headers map[string]string, delay time.Duration) error
	// Close 关闭生产者
	Close() error
}

// Consumer 消费者接口
type Consumer interface {
	// Subscribe 订阅主题
	Subscribe(ctx context.Context, topic string, channel string, handler func(msg *models.BaseMessage) error) error
	// Unsubscribe 取消订阅
	Unsubscribe(topic string, channel string) error
	// SetDeadLetterHandler 设置死信处理器
	SetDeadLetterHandler(handler func(msg models.DeadLetterMessage) error)
	// Close 关闭消费者
	Close() error
}

// MQ 消息队列接口
type MQ interface {
	// NewProducer 创建生产者
	NewProducer(config *models.Config) (Producer, error)
	// NewConsumer 创建消费者
	NewConsumer(config *models.Config) (Consumer, error)
	// Close 关闭连接
	Close() error
}

// Server 消息队列服务接口
type Server interface {
	// Publish 发布消息
	Publish(ctx context.Context, message *models.BaseMessage) error
	// PublishDelay 发布延迟消息
	PublishDelay(ctx context.Context, message *models.BaseMessage, delay time.Duration) error
	// Subscribe 订阅主题
	Subscribe(ctx context.Context, topic string, channel string, handler func(msg *models.BaseMessage) error) error
	// Unsubscribe 取消订阅
	Unsubscribe(topic string, channel string) error
	// SubscribeDeadLetter 订阅死信队列
	SubscribeDeadLetter(ctx context.Context, handler func(msg models.DeadLetterMessage) error) error
	// UnsubscribeDeadLetter 取消订阅死信队列
	UnsubscribeDeadLetter() error
	// Close 关闭服务
	Close() error
}
