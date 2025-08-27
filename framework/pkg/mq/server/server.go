package server

import (
	"context"
	"encoding/json"
	"fmt"
	mq2 "github.com/flare-admin/flare-server-go/framework/pkg/mq"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/models"
	"sync"
	"time"
)

// UnifiedServer 统一服务实现
type UnifiedServer struct {
	mq              mq2.MQ
	producer        mq2.Producer
	consumer        mq2.Consumer
	config          *models.Config
	mu              sync.RWMutex
	deadLetterTopic string
	deadLetterSub   *struct {
		handler func(msg models.DeadLetterMessage) error
		cancel  context.CancelFunc
	}
}

// NewServer 创建统一服务
func NewServer(config *models.Config) (mq2.Server, error) {
	// 验证配置
	if config.Type != "nsq" && config.Type != "nats" {
		return nil, fmt.Errorf("不支持的MQ类型: %s", config.Type)
	}

	// 创建MQ实例
	mq, err := NewMQ(config)
	if err != nil {
		return nil, fmt.Errorf("创建MQ失败: %v", err)
	}

	// 创建生产者
	producer, err := mq.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("创建生产者失败: %v", err)
	}

	// 创建消费者
	consumer, err := mq.NewConsumer(config)
	if err != nil {
		producer.Close()
		return nil, fmt.Errorf("创建消费者失败: %v", err)
	}

	// 设置死信队列主题
	deadLetterTopic := "dead_letter"
	if config.Options != nil {
		if topic, ok := config.Options["dead_letter_topic"].(string); ok {
			deadLetterTopic = topic
		}
	}

	server := &UnifiedServer{
		mq:              mq,
		producer:        producer,
		consumer:        consumer,
		config:          config,
		deadLetterTopic: deadLetterTopic,
	}

	// 设置死信处理器
	consumer.SetDeadLetterHandler(func(msg models.DeadLetterMessage) error {
		// 将死信消息发布到死信队列
		data, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("序列化死信消息失败: %v", err)
		}

		// 创建死信消息
		deadLetterMsg := &models.BaseMessage{
			Topic:   server.deadLetterTopic,
			Payload: string(data),
			Headers: map[string]string{
				"original_topic":   msg.Topic,
				"original_channel": msg.Channel,
				"error":            msg.Error,
				"retry_count":      fmt.Sprintf("%d", msg.RetryCount),
			},
		}

		// 发布死信消息
		return server.producer.Publish(context.Background(), deadLetterMsg.GetTopic(), deadLetterMsg.GetPayload(), deadLetterMsg.GetHeaders())
	})

	return server, nil
}

// Publish 发布消息
func (s *UnifiedServer) Publish(ctx context.Context, message *models.BaseMessage) error {
	if s.producer == nil {
		return fmt.Errorf("生产者未初始化")
	}
	return s.producer.Publish(ctx, message.GetTopic(), message.GetPayload(), message.GetHeaders())
}

// PublishDelay 发布延迟消息
func (s *UnifiedServer) PublishDelay(ctx context.Context, message *models.BaseMessage, delay time.Duration) error {
	if s.producer == nil {
		return fmt.Errorf("生产者未初始化")
	}
	return s.producer.PublishDelay(ctx, message.GetTopic(), message.GetPayload(), message.GetHeaders(), delay)
}

// Subscribe 订阅主题
func (s *UnifiedServer) Subscribe(ctx context.Context, topic string, channel string, handler func(msg *models.BaseMessage) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.consumer.Subscribe(ctx, topic, channel, handler)
}

// Unsubscribe 取消订阅
func (s *UnifiedServer) Unsubscribe(topic string, channel string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.consumer.Unsubscribe(topic, channel)
}

// SubscribeDeadLetter 订阅死信队列
func (s *UnifiedServer) SubscribeDeadLetter(ctx context.Context, handler func(msg models.DeadLetterMessage) error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果已经有订阅，先取消
	if s.deadLetterSub != nil {
		s.deadLetterSub.cancel()
	}

	// 创建新的上下文
	ctx, cancel := context.WithCancel(ctx)

	// 订阅死信队列
	err := s.consumer.Subscribe(ctx, s.deadLetterTopic, "dead_letter", func(msg *models.BaseMessage) error {
		var deadLetterMsg models.DeadLetterMessage
		if err := json.Unmarshal(msg.GetPayload(), &deadLetterMsg); err != nil {
			return fmt.Errorf("反序列化死信消息失败: %v", err)
		}
		return handler(deadLetterMsg)
	})

	if err != nil {
		cancel()
		return fmt.Errorf("订阅死信队列失败: %v", err)
	}

	// 保存订阅信息
	s.deadLetterSub = &struct {
		handler func(msg models.DeadLetterMessage) error
		cancel  context.CancelFunc
	}{
		handler: handler,
		cancel:  cancel,
	}

	return nil
}

// UnsubscribeDeadLetter 取消订阅死信队列
func (s *UnifiedServer) UnsubscribeDeadLetter() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.deadLetterSub != nil {
		// 取消订阅
		if err := s.consumer.Unsubscribe(s.deadLetterTopic, "dead_letter"); err != nil {
			return fmt.Errorf("取消订阅死信队列失败: %v", err)
		}

		// 取消上下文
		s.deadLetterSub.cancel()
		s.deadLetterSub = nil
	}

	return nil
}

// Close 关闭服务
func (s *UnifiedServer) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 取消死信队列订阅
	if s.deadLetterSub != nil {
		s.deadLetterSub.cancel()
		s.deadLetterSub = nil
	}

	// 关闭消费者
	if err := s.consumer.Close(); err != nil {
		return fmt.Errorf("关闭消费者失败: %v", err)
	}

	// 关闭生产者
	if err := s.producer.Close(); err != nil {
		return fmt.Errorf("关闭生产者失败: %v", err)
	}

	// 关闭MQ
	if err := s.mq.Close(); err != nil {
		return fmt.Errorf("关闭MQ失败: %v", err)
	}

	return nil
}
