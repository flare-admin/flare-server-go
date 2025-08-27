package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/models"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"strings"
	"time"

	"github.com/nats-io/nats.go"
)

// Producer NATS生产者实现
type Producer struct {
	conn   *nats.Conn
	config *Config
}

// NewNATSProducer 创建NATS生产者
func NewNATSProducer(config *models.Config) (*Producer, error) {
	// 转换为NATS配置
	natsConfig, err := FromMQConfig(config)
	if err != nil {
		return nil, fmt.Errorf("转换配置失败: %v", err)
	}

	// 创建NATS选项
	opts := []nats.Option{
		nats.ReconnectWait(time.Duration(natsConfig.ReconnectWait) * time.Millisecond),
		nats.MaxReconnects(natsConfig.MaxReconnects),
	}

	// 连接NATS服务器
	conn, err := nats.Connect(natsConfig.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("连接NATS失败: %v", err)
	}

	return &Producer{
		conn:   conn,
		config: natsConfig,
	}, nil
}

// Publish 发布消息
func (p *Producer) Publish(ctx context.Context, topic string, payload []byte, headers map[string]string) error {
	msg := models.NewBaseMessage(topic, payload, headers)
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	return p.conn.Publish(topic, data)
}

// PublishDelay 发布延迟消息
func (p *Producer) PublishDelay(ctx context.Context, topic string, payload []byte, headers map[string]string, delay time.Duration) error {
	msg := models.NewBaseMessage(topic, payload, headers)
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	// NATS不支持原生延迟消息，使用定时器实现
	time.AfterFunc(delay, func() {
		p.conn.Publish(topic, data)
	})
	return nil
}

// Close 关闭生产者
func (p *Producer) Close() error {
	p.conn.Close()
	return nil
}

// Consumer NATS消费者实现
type Consumer struct {
	conn              *nats.Conn
	config            *Config
	handlers          map[string][]func(msg *models.BaseMessage) error
	deadLetterHandler func(msg models.DeadLetterMessage) error
	retryCounts       map[string]int                // 记录每个消息的重试次数
	subscriptions     map[string]*nats.Subscription // 保存订阅信息
}

// NewNATSConsumer 创建NATS消费者
func NewNATSConsumer(config *models.Config) (*Consumer, error) {
	// 转换为NATS配置
	natsConfig, err := FromMQConfig(config)
	if err != nil {
		return nil, fmt.Errorf("转换配置失败: %v", err)
	}

	// 创建NATS选项
	opts := []nats.Option{
		nats.ReconnectWait(time.Duration(natsConfig.ReconnectWait) * time.Second),
		nats.MaxReconnects(natsConfig.MaxReconnects),
	}

	// 连接NATS服务器
	conn, err := nats.Connect(natsConfig.Address, opts...)
	if err != nil {
		return nil, fmt.Errorf("连接NATS失败: %v", err)
	}

	return &Consumer{
		conn:          conn,
		config:        natsConfig,
		handlers:      make(map[string][]func(msg *models.BaseMessage) error),
		retryCounts:   make(map[string]int),
		subscriptions: make(map[string]*nats.Subscription),
	}, nil
}

// SetDeadLetterHandler 设置死信处理器
func (c *Consumer) SetDeadLetterHandler(handler func(msg models.DeadLetterMessage) error) {
	c.deadLetterHandler = handler
}

// Subscribe 订阅主题
func (c *Consumer) Subscribe(ctx context.Context, topic string, channel string, handler func(msg *models.BaseMessage) error) error {
	handlerKey := topic + ":" + channel
	c.handlers[handlerKey] = append(c.handlers[handlerKey], handler)

	queueGroup := fmt.Sprintf("%s-%s", topic, channel)
	sub, err := c.conn.QueueSubscribe(topic, queueGroup, func(msg *nats.Msg) {

		defer func() {
			if r := recover(); r != nil {
				hlog.CtxErrorf(ctx, "panic recovered: %v", r)
			}
		}()
		var baseMsg models.BaseMessage
		if err := json.Unmarshal(msg.Data, &baseMsg); err != nil {
			hlog.CtxErrorf(ctx, "反序列化消息失败: %v, body: %s", err, string(msg.Data))
			return
		}
		// 生成消息唯一标识
		msgID := fmt.Sprintf("%s:%s:%s", topic, channel, baseMsg.GetID())

		if handlers, ok := c.handlers[handlerKey]; ok {
			for _, h := range handlers {
				if err := h(&baseMsg); err != nil {

					// 增加重试次数
					c.retryCounts[msgID]++
					retryCount := c.retryCounts[msgID]

					// 如果超过最大重试次数，发送到死信队列
					if retryCount >= c.config.MaxRetries {
						deadLetterMsg := models.DeadLetterMessage{
							Topic:      topic,
							Channel:    channel,
							Payload:    baseMsg.GetPayload(),
							Headers:    baseMsg.GetHeaders(),
							Error:      err.Error(),
							RetryCount: retryCount,
							Timestamp:  utils.GetTimeNow(),
							DeadTime:   utils.GetTimeNow(),
						}

						// 调用死信处理器
						if c.deadLetterHandler != nil {
							c.deadLetterHandler(deadLetterMsg)
						}
						// 清理重试计数
						delete(c.retryCounts, msgID)
						return
					}

					// 未超过重试次数，重新发布消息
					time.Sleep(time.Second * time.Duration(retryCount)) // 指数退避
					c.conn.Publish(topic, msg.Data)
					return
				}
			}
		}
		// 处理成功，清理重试计数
		delete(c.retryCounts, msgID)
	})

	if err != nil {
		return fmt.Errorf("订阅主题失败: %v", err)
	}

	// 保存订阅信息
	c.subscriptions[handlerKey] = sub
	return nil
}

// Unsubscribe 取消订阅
func (c *Consumer) Unsubscribe(topic string, channel string) error {
	handlerKey := topic + ":" + channel

	// 取消订阅
	if sub, ok := c.subscriptions[handlerKey]; ok {
		if err := sub.Unsubscribe(); err != nil {
			return fmt.Errorf("取消订阅失败: %v", err)
		}
		delete(c.subscriptions, handlerKey)
	}

	// 清理处理器和重试计数
	delete(c.handlers, handlerKey)

	// 清理相关的重试计数
	for msgID := range c.retryCounts {
		if strings.HasPrefix(msgID, handlerKey) {
			delete(c.retryCounts, msgID)
		}
	}

	return nil
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	// 取消所有订阅
	for _, sub := range c.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			return fmt.Errorf("取消订阅失败: %v", err)
		}
	}

	// 清理所有资源
	c.subscriptions = make(map[string]*nats.Subscription)
	c.handlers = make(map[string][]func(msg *models.BaseMessage) error)
	c.retryCounts = make(map[string]int)

	// 关闭连接
	c.conn.Close()
	return nil
}

// NATS NATS实现
type NATS struct {
	config *models.Config
}

// NewNATS 创建NATS实例
func NewNATS(config *models.Config) *NATS {
	return &NATS{
		config: config,
	}
}

// NewProducer 创建生产者
func (n *NATS) NewProducer(config *models.Config) (mq.Producer, error) {
	return NewNATSProducer(config)
}

// NewConsumer 创建消费者
func (n *NATS) NewConsumer(config *models.Config) (mq.Consumer, error) {
	return NewNATSConsumer(config)
}

// Close 关闭连接
func (n *NATS) Close() error {
	return nil
}
