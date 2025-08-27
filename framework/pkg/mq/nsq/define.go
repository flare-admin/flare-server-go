package nsq

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/models"
)

const (
	defaultMaxRetries          = 3
	defaultMaxInFlight         = 100
	defaultLookupdPollInterval = 2
	defaultLookupdPollJitter   = 0.3
	defaultMaxBackoffDuration  = 60
	defaultBackoffMultiplier   = 1.5
	defaultDeflateLevel        = 6
)

// Config NSQ配置
type Config struct {
	// Address 地址
	Address string
	// MaxRetries 最大重试次数
	MaxRetries int
	// MaxInFlight 最大飞行消息数
	MaxInFlight int
	// LookupdPollInterval 查询间隔(秒)
	LookupdPollInterval int
	// LookupdPollJitter 查询抖动
	LookupdPollJitter float64
	// MaxBackoffDuration 最大退避时间(秒)
	MaxBackoffDuration int
	// BackoffMultiplier 退避乘数
	BackoffMultiplier float64
	// Deflate 是否启用压缩
	Deflate bool
	// DeflateLevel 压缩级别
	DeflateLevel int
	// Snappy 是否启用Snappy压缩
	Snappy bool
	// TLSConfig TLS配置
	TLSConfig map[string]interface{}
}

// ConfigBuilder NSQ配置构建器
type ConfigBuilder struct {
	cfg *Config
}

// NewConfigBuilder 创建NSQ配置构建器
func NewConfigBuilder(address string) *ConfigBuilder {
	return &ConfigBuilder{
		cfg: &Config{
			Address:             address,
			MaxRetries:          defaultMaxRetries,
			MaxInFlight:         defaultMaxInFlight,
			LookupdPollInterval: defaultLookupdPollInterval,
			LookupdPollJitter:   defaultLookupdPollJitter,
			MaxBackoffDuration:  defaultMaxBackoffDuration,
			BackoffMultiplier:   defaultBackoffMultiplier,
			Deflate:             false,
			DeflateLevel:        defaultDeflateLevel,
			Snappy:              false,
			TLSConfig:           nil,
		},
	}
}

// MaxRetries 设置最大重试次数
func (b *ConfigBuilder) MaxRetries(v int) *ConfigBuilder {
	b.cfg.MaxRetries = v
	return b
}

// MaxInFlight 设置最大飞行消息数
func (b *ConfigBuilder) MaxInFlight(v int) *ConfigBuilder {
	b.cfg.MaxInFlight = v
	return b
}

// LookupdPollInterval 设置查询间隔(秒)
func (b *ConfigBuilder) LookupdPollInterval(v int) *ConfigBuilder {
	b.cfg.LookupdPollInterval = v
	return b
}

// LookupdPollJitter 设置查询抖动
func (b *ConfigBuilder) LookupdPollJitter(v float64) *ConfigBuilder {
	b.cfg.LookupdPollJitter = v
	return b
}

// MaxBackoffDuration 设置最大退避时间(秒)
func (b *ConfigBuilder) MaxBackoffDuration(v int) *ConfigBuilder {
	b.cfg.MaxBackoffDuration = v
	return b
}

// BackoffMultiplier 设置退避乘数
func (b *ConfigBuilder) BackoffMultiplier(v float64) *ConfigBuilder {
	b.cfg.BackoffMultiplier = v
	return b
}

// Deflate 设置是否启用压缩
func (b *ConfigBuilder) Deflate(v bool) *ConfigBuilder {
	b.cfg.Deflate = v
	return b
}

// DeflateLevel 设置压缩级别
func (b *ConfigBuilder) DeflateLevel(v int) *ConfigBuilder {
	b.cfg.DeflateLevel = v
	return b
}

// Snappy 设置是否启用Snappy压缩
func (b *ConfigBuilder) Snappy(v bool) *ConfigBuilder {
	b.cfg.Snappy = v
	return b
}

// TLSConfig 设置TLS配置
func (b *ConfigBuilder) TLSConfig(v map[string]interface{}) *ConfigBuilder {
	b.cfg.TLSConfig = v
	return b
}

// Build 构建配置
func (b *ConfigBuilder) Build() *Config {
	return b.cfg
}

// ToMQConfig 转换为MQ配置
func (c *Config) ToMQConfig() *models.Config {
	return &models.Config{
		Type:       "nsq",
		Address:    c.Address,
		MaxRetries: c.MaxRetries,
		Options: map[string]interface{}{
			"max_in_flight":         c.MaxInFlight,
			"lookupd_poll_interval": c.LookupdPollInterval,
			"lookupd_poll_jitter":   c.LookupdPollJitter,
			"max_backoff_duration":  c.MaxBackoffDuration,
			"backoff_multiplier":    c.BackoffMultiplier,
			"deflate":               c.Deflate,
			"deflate_level":         c.DeflateLevel,
			"snappy":                c.Snappy,
			"tls_config":            c.TLSConfig,
		},
	}
}

// FromMQConfig 从MQ配置创建NSQ配置
func FromMQConfig(config *models.Config) (*Config, error) {
	if config == nil {
		return nil, fmt.Errorf("配置不能为空")
	}

	builder := NewConfigBuilder(config.Address)
	builder.MaxRetries(config.MaxRetries)

	if config.Options != nil {
		if v, ok := config.Options["max_in_flight"].(int); ok {
			builder.MaxInFlight(v)
		}
		if v, ok := config.Options["lookupd_poll_interval"].(int); ok {
			builder.LookupdPollInterval(v)
		}
		if v, ok := config.Options["lookupd_poll_jitter"].(float64); ok {
			builder.LookupdPollJitter(v)
		}
		if v, ok := config.Options["max_backoff_duration"].(int); ok {
			builder.MaxBackoffDuration(v)
		}
		if v, ok := config.Options["backoff_multiplier"].(float64); ok {
			builder.BackoffMultiplier(v)
		}
		if v, ok := config.Options["deflate"].(bool); ok {
			builder.Deflate(v)
		}
		if v, ok := config.Options["deflate_level"].(int); ok {
			builder.DeflateLevel(v)
		}
		if v, ok := config.Options["snappy"].(bool); ok {
			builder.Snappy(v)
		}
		if v, ok := config.Options["tls_config"].(map[string]interface{}); ok {
			builder.TLSConfig(v)
		}
	}

	return builder.Build(), nil
}
