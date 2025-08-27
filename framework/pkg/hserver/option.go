package hserver

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
)

type ServerConfig struct {
	Port               int    `json:"port"`
	RateQPS            int    `json:"rate_qps"`
	TracerPort         int    `json:"tracer_port"`
	Name               string `json:"name"`
	MaxRequestBodySize int    `json:"max_request_body_size"`
}

// Option 定义一个函数类型，用于修改Server配置
type Option func(*Serve)

// WithTokenizer 设置token
func WithTokenizer(token token.IToken) Option {
	return func(a *Serve) {
		a.Tokenizer = token
	}
}

//// WithConfigs 设置config
//func WithConfigs(config *configs.Bootstrap) Option {
//	return func(a *Serve) {
//		a.config = config
//	}
//}

// WithEnv 设置环境变量
func WithEnv(env string) Option {
	return func(a *Serve) {
		a.Env = env
	}
}

// WithBaseUrl 设置baseUrl
func WithBaseUrl(baseUrl string) Option {
	return func(a *Serve) {
		a.baseUrl = baseUrl
	}
}
