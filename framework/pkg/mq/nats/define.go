package nats

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/models"
)

const (
	defaultMaxRetries    = 3
	defaultReconnectWait = 2
	defaultMaxReconnects = 10
	defaultMaxInFlight   = 100
	defaultQueueGroup    = "default"
	defaultDurableName   = "default"
	defaultAckWait       = 30
	defaultMaxDeliver    = 3
)

// Config NATS配置
type Config struct {
	// Address 地址
	Address string
	// MaxRetries 最大重试次数
	MaxRetries int
	// ReconnectWait 重连等待时间(秒)
	ReconnectWait int
	// MaxReconnects 最大重连次数
	MaxReconnects int
	// QueueGroup 队列组名称
	QueueGroup string
	// DurableName 持久化名称
	DurableName string
	// AckWait 确认等待时间(秒)
	AckWait int
	// MaxDeliver 最大投递次数
	MaxDeliver int
}

// ConfigBuilder NATS配置构建器
type ConfigBuilder struct {
	cfg *Config
}

// NewConfigBuilder 创建NATS配置构建器
func NewConfigBuilder(address string) *ConfigBuilder {
	return &ConfigBuilder{
		cfg: &Config{
			Address:       address,
			MaxRetries:    defaultMaxRetries,
			ReconnectWait: defaultReconnectWait,
			MaxReconnects: defaultMaxReconnects,
			QueueGroup:    defaultQueueGroup,
			DurableName:   defaultDurableName,
			AckWait:       defaultAckWait,
			MaxDeliver:    defaultMaxDeliver,
		},
	}
}

// MaxRetries 设置最大重试次数
func (b *ConfigBuilder) MaxRetries(v int) *ConfigBuilder {
	b.cfg.MaxRetries = v
	return b
}

// ReconnectWait 设置重连等待时间(秒)
func (b *ConfigBuilder) ReconnectWait(v int) *ConfigBuilder {
	b.cfg.ReconnectWait = v
	return b
}

// MaxReconnects 设置最大重连次数
func (b *ConfigBuilder) MaxReconnects(v int) *ConfigBuilder {
	b.cfg.MaxReconnects = v
	return b
}

// QueueGroup 设置队列组名称
func (b *ConfigBuilder) QueueGroup(v string) *ConfigBuilder {
	b.cfg.QueueGroup = v
	return b
}

// DurableName 设置持久化名称
func (b *ConfigBuilder) DurableName(v string) *ConfigBuilder {
	b.cfg.DurableName = v
	return b
}

// AckWait 设置确认等待时间(秒)
func (b *ConfigBuilder) AckWait(v int) *ConfigBuilder {
	b.cfg.AckWait = v
	return b
}

// MaxDeliver 设置最大投递次数
func (b *ConfigBuilder) MaxDeliver(v int) *ConfigBuilder {
	b.cfg.MaxDeliver = v
	return b
}

// Build 构建配置
func (b *ConfigBuilder) Build() *Config {
	return b.cfg
}

// ToMQConfig 转换为MQ配置
func (c *Config) ToMQConfig() *models.Config {
	return &models.Config{
		Type:       "nats",
		Address:    c.Address,
		MaxRetries: c.MaxRetries,
		Options: map[string]interface{}{
			"reconnect_wait": c.ReconnectWait,
			"max_reconnects": c.MaxReconnects,
			"queue_group":    c.QueueGroup,
			"durable_name":   c.DurableName,
			"ack_wait":       c.AckWait,
			"max_deliver":    c.MaxDeliver,
		},
	}
}

// FromMQConfig 从MQ配置创建NATS配置
func FromMQConfig(config *models.Config) (*Config, error) {
	if config == nil {
		return nil, fmt.Errorf("配置不能为空")
	}

	builder := NewConfigBuilder(config.Address)
	builder.MaxRetries(config.MaxRetries)

	if config.Options != nil {
		if v, ok := config.Options["reconnect_wait"].(int); ok {
			builder.ReconnectWait(v)
		}
		if v, ok := config.Options["max_reconnects"].(int); ok {
			builder.MaxReconnects(v)
		}
		if v, ok := config.Options["queue_group"].(string); ok {
			builder.QueueGroup(v)
		}
		if v, ok := config.Options["durable_name"].(string); ok {
			builder.DurableName(v)
		}
		if v, ok := config.Options["ack_wait"].(int); ok {
			builder.AckWait(v)
		}
		if v, ok := config.Options["max_deliver"].(int); ok {
			builder.MaxDeliver(v)
		}
	}

	return builder.Build(), nil
}
