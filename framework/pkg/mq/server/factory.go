package server

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/models"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/nats"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/nsq"
	"sync"
)

var (
	factory *Factory
	once    sync.Once
)

// Factory MQ工厂
type Factory struct {
	mqs map[string]mq.MQ
	mu  sync.RWMutex
}

// GetFactory 获取工厂实例
func GetFactory() *Factory {
	once.Do(func() {
		factory = &Factory{
			mqs: make(map[string]mq.MQ),
		}
	})
	return factory
}

// Register 注册MQ实现
func (f *Factory) Register(name string, mq mq.MQ) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.mqs[name] = mq
}

// GetMQ 获取MQ实现
func (f *Factory) GetMQ(name string) (mq.MQ, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	mq, ok := f.mqs[name]
	if !ok {
		return nil, fmt.Errorf("未找到MQ实现: %s", name)
	}
	return mq, nil
}

// NewMQ 创建MQ实例
func NewMQ(config *models.Config) (mq.MQ, error) {
	switch config.Type {
	case "nsq":
		return nsq.NewNSQ(config), nil
	case "nats":
		return nats.NewNATS(config), nil
	default:
		return nil, fmt.Errorf("不支持的MQ类型: %s", config.Type)
	}
}

// Init 初始化MQ工厂
func Init() {
	factory := GetFactory()
	factory.Register("nsq", &nsq.NSQ{})
	factory.Register("nats", &nats.NATS{})
}
