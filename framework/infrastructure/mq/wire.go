package mq

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/nats"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/server"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewMqServer,
)

func NewMqServer(cof *configs.Bootstrap) (mq.Server, func(), error) {
	//1. 创建 NATS 配置
	natsConfig := nats.NewConfigBuilder(cof.NATSConfig.Address).
		MaxRetries(cof.NATSConfig.MaxRetries).
		ReconnectWait(cof.NATSConfig.ReconnectWait). // 5秒
		AckWait(cof.NATSConfig.AckWait).             // 60秒
		QueueGroup(cof.NATSConfig.QueueGroup).
		DurableName(cof.NATSConfig.DurableName).
		Build()

	// 2. 转换为通用 MQ 配置
	mqConfig := natsConfig.ToMQConfig()

	//config := nsq.NewConfigBuilder(cof.NSQConfig.Address).
	//	MaxRetries(cof.NSQConfig.MaxRetries).
	//	MaxInFlight(cof.NSQConfig.MaxInFlight).
	//	MaxBackoffDuration(cof.NSQConfig.MaxBackoffDuration).
	//	LookupdPollInterval(cof.NSQConfig.LookupdPollInterval).
	//	Build()
	//mqConfig := config.ToMQConfig()

	// 3. 创建统一服务器
	server, err := server.NewServer(mqConfig)
	if err != nil {
		panic(fmt.Sprintf("创建mq服务器失败: %v", err))
	}
	cleanup := func() {
		server.Close()
	}
	return server, cleanup, nil
}
