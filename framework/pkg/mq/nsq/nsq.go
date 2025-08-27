package nsq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/models"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/nsqio/go-nsq"
	"sync"
	"time"
)

// Producer NSQ生产者实现
type Producer struct {
	producer *nsq.Producer
	config   *Config
}

// NewNSQProducer 创建NSQ生产者
func NewNSQProducer(config *models.Config) (*Producer, error) {
	// 转换为NSQ配置
	nsqConfig, err := FromMQConfig(config)
	if err != nil {
		return nil, fmt.Errorf("转换配置失败: %v", err)
	}

	// 创建NSQ配置
	cfg := nsq.NewConfig()
	cfg.MaxInFlight = nsqConfig.MaxInFlight

	// 设置压缩选项
	if nsqConfig.Deflate {
		cfg.Deflate = true
		cfg.DeflateLevel = nsqConfig.DeflateLevel
	}
	if nsqConfig.Snappy {
		cfg.Snappy = true
	}

	// 创建生产者
	producer, err := nsq.NewProducer(nsqConfig.Address, cfg)
	if err != nil {
		return nil, fmt.Errorf("创建NSQ生产者失败: %v", err)
	}

	return &Producer{
		producer: producer,
		config:   nsqConfig,
	}, nil
}

// Publish 发布消息
func (p *Producer) Publish(ctx context.Context, topic string, payload []byte, headers map[string]string) error {
	msg := models.NewBaseMessage(topic, payload, headers)
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	return p.producer.Publish(topic, data)
}

// PublishDelay 发布延迟消息
func (p *Producer) PublishDelay(ctx context.Context, topic string, payload []byte, headers map[string]string, delay time.Duration) error {
	msg := models.NewBaseMessage(topic, payload, headers)
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	return p.producer.DeferredPublish(topic, delay, data)
}

// Close 关闭生产者
func (p *Producer) Close() error {
	p.producer.Stop()
	return nil
}

// Consumer NSQ消费者实现
type Consumer struct {
	config            *Config
	handlers          map[string][]func(msg *models.BaseMessage) error
	deadLetterHandler func(msg models.DeadLetterMessage) error
	consumers         map[string]*nsq.Consumer // 保存每个 topic:channel 的消费者
	mu                sync.RWMutex
}

// NewNSQConsumer 创建NSQ消费者
func NewNSQConsumer(config *models.Config) (*Consumer, error) {
	// 转换为NSQ配置
	nsqConfig, err := FromMQConfig(config)
	if err != nil {
		return nil, fmt.Errorf("转换配置失败: %v", err)
	}

	return &Consumer{
		config:    nsqConfig,
		handlers:  make(map[string][]func(msg *models.BaseMessage) error),
		consumers: make(map[string]*nsq.Consumer),
	}, nil
}

// Subscribe 订阅主题
func (c *Consumer) Subscribe(ctx context.Context, topic string, channel string, handler func(msg *models.BaseMessage) error) error {
	if handler == nil {
		return fmt.Errorf("handler 不能为空")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	handlerKey := topic + ":" + channel
	c.handlers[handlerKey] = append(c.handlers[handlerKey], handler)

	// 如果已经有消费者，直接返回
	if _, ok := c.consumers[handlerKey]; ok {
		return nil
	}

	// 创建NSQ配置
	cfg := nsq.NewConfig()
	if c.config.MaxInFlight > 0 {
		cfg.MaxInFlight = c.config.MaxInFlight
	}
	if c.config.MaxRetries > 0 {
		cfg.MaxAttempts = uint16(c.config.MaxRetries)
	}
	if c.config.LookupdPollInterval > 0 {
		cfg.LookupdPollInterval = time.Duration(c.config.LookupdPollInterval) * time.Millisecond
	}
	if c.config.MaxBackoffDuration > 0 {
		cfg.MaxBackoffDuration = time.Duration(c.config.MaxBackoffDuration) * time.Millisecond
	}

	if c.config.LookupdPollJitter > 0 {
		cfg.LookupdPollJitter = c.config.LookupdPollJitter
	}
	if c.config.MaxBackoffDuration > 0 {
		cfg.BackoffMultiplier = time.Duration(c.config.BackoffMultiplier) * time.Millisecond
	}
	// 设置压缩选项
	if c.config.Deflate {
		cfg.Deflate = true
		cfg.DeflateLevel = c.config.DeflateLevel
	}
	if c.config.Snappy {
		cfg.Snappy = true
	}

	// 创建消费者
	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		hlog.CtxInfof(ctx, "创建NSQ消费者失败: %v", err)
		return fmt.Errorf("创建NSQ消费者失败: %v", err)
	}

	// 设置消息处理器
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		defer func() {
			if r := recover(); r != nil {
				hlog.CtxErrorf(ctx, "panic recovered: %v", r)
				err = fmt.Errorf("panic: %v", r)
			}
		}()
		msg := models.BaseMessage{}
		if err1 := json.Unmarshal(message.Body, &msg); err1 != nil {
			hlog.CtxErrorf(ctx, "反序列化消息失败: %v, body: %s", err1, string(message.Body))
			return fmt.Errorf("反序列化消息失败: %v", err1)
		}

		if handlers, ok := c.handlers[handlerKey]; ok {
			for _, handler := range handlers {
				if handler == nil {
					hlog.Errorf("handler 为空，handlerKey: %s", handlerKey)
					continue
				}
				if err := handler(&msg); err != nil {
					// 如果超过最大重试次数，发送到死信队列
					if int(message.Attempts) >= c.config.MaxRetries {
						deadLetterMsg := models.DeadLetterMessage{
							Topic:      topic,
							Channel:    channel,
							Payload:    msg.GetPayload(),
							Headers:    msg.GetHeaders(),
							Error:      err.Error(),
							RetryCount: int(message.Attempts),
							Timestamp:  time.Unix(0, message.Timestamp),
							DeadTime:   utils.GetTimeNow(),
						}
						// 调用死信处理器
						if c.deadLetterHandler != nil {
							c.deadLetterHandler(deadLetterMsg)
						}
						// 超过重试次数，不再重试
						return nil
					}
					// 未超过重试次数，返回错误以触发重试
					return err
				}
			}
		}
		return nil
	}))

	// 连接到NSQ服务器
	if err := consumer.ConnectToNSQD(c.config.Address); err != nil {
		return fmt.Errorf("连接NSQ失败: %v", err)
	}

	// 保存消费者
	c.consumers[handlerKey] = consumer
	return nil
}

// Unsubscribe 取消订阅
func (c *Consumer) Unsubscribe(topic string, channel string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	handlerKey := topic + ":" + channel

	// 停止并删除消费者
	if consumer, ok := c.consumers[handlerKey]; ok {
		consumer.Stop()
		delete(c.consumers, handlerKey)
	}

	// 删除处理器
	delete(c.handlers, handlerKey)
	return nil
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 停止所有消费者
	for _, consumer := range c.consumers {
		consumer.Stop()
	}

	// 清理资源
	c.consumers = make(map[string]*nsq.Consumer)
	c.handlers = make(map[string][]func(msg *models.BaseMessage) error)
	return nil
}

// SetDeadLetterHandler 设置死信处理器
func (c *Consumer) SetDeadLetterHandler(handler func(msg models.DeadLetterMessage) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deadLetterHandler = handler
}

// NSQ NSQ实现
type NSQ struct {
	config *Config
}

// NewNSQ 创建NSQ实例
func NewNSQ(config *models.Config) *NSQ {
	nsqConfig, _ := FromMQConfig(config)
	return &NSQ{
		config: nsqConfig,
	}
}

// NewProducer 创建生产者
func (n *NSQ) NewProducer(config *models.Config) (mq.Producer, error) {
	return NewNSQProducer(config)
}

// NewConsumer 创建消费者
func (n *NSQ) NewConsumer(config *models.Config) (mq.Consumer, error) {
	return NewNSQConsumer(config)
}

// Close 关闭连接
func (n *NSQ) Close() error {
	return nil
}
