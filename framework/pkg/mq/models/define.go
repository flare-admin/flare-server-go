package models

import "time"

// Config 配置
type Config struct {
	// Type 类型 (nsq, nats等)
	Type string
	// Address 地址
	Address string
	// MaxRetries 最大重试次数
	MaxRetries int
	// Options 其他选项
	Options map[string]interface{}
}

// DeadLetterMessage 死信消息模型
type DeadLetterMessage struct {
	Topic      string            `json:"topic"`       // 原始主题
	Channel    string            `json:"channel"`     // 原始通道
	Payload    []byte            `json:"payload"`     // 原始消息内容
	Headers    map[string]string `json:"headers"`     // 原始消息头
	Error      string            `json:"error"`       // 处理失败原因
	RetryCount int               `json:"retry_count"` // 重试次数
	Timestamp  time.Time         `json:"timestamp"`   // 消息时间戳
	DeadTime   time.Time         `json:"dead_time"`   // 进入死信队列时间
}
