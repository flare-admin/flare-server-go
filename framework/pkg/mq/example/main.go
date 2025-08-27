package main

import (
	"context"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/models"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/nats"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/server"
	"time"
)

func main() {
	// 1. 初始化 MQ 工厂
	server.Init()

	// 2. 创建 NATS 配置
	natsConfig := nats.NewConfigBuilder("nats://127.0.0.1:4222").
		MaxRetries(5).
		ReconnectWait(5). // 5秒
		AckWait(60).      // 60秒
		QueueGroup("group1").
		Build()

	// 3. 转换为通用 MQ 配置
	mqConfig := natsConfig.ToMQConfig()

	// 4. 创建 MQ 实例
	mqInstance, err := server.NewMQ(mqConfig)
	if err != nil {
		panic(fmt.Sprintf("创建 MQ 失败: %v", err))
	}

	// 5. 创建生产者和消费者
	producer, err := mqInstance.NewProducer(mqConfig)
	if err != nil {
		panic(fmt.Sprintf("创建生产者失败: %v", err))
	}

	consumer, err := mqInstance.NewConsumer(mqConfig)
	if err != nil {
		panic(fmt.Sprintf("创建消费者失败: %v", err))
	}

	// 6. 设置死信处理器
	consumer.SetDeadLetterHandler(func(msg models.DeadLetterMessage) error {
		fmt.Printf("死信消息: %+v\n", msg)
		return nil
	})

	// 7. 订阅消息
	ctx := context.Background()
	err = consumer.Subscribe(ctx, "test.topic", "test.channel", func(msg mq.Message) error {
		fmt.Printf("收到消息: %s\n", string(msg.GetPayload()))
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("订阅失败: %v", err))
	}

	// 8. 发布消息
	headers := map[string]string{
		"tenant_id": "123",
		"user_id":   "456",
	}

	// 发布普通消息
	err = producer.Publish(ctx, "test.topic", []byte("Hello MQ!"), headers)
	if err != nil {
		panic(fmt.Sprintf("发布失败: %v", err))
	}

	// 发布延迟消息
	err = producer.PublishDelay(ctx, "test.topic", []byte("Delayed Message!"), headers, 5*time.Second)
	if err != nil {
		panic(fmt.Sprintf("发布延迟消息失败: %v", err))
	}

	// 9. 等待消息处理
	time.Sleep(10 * time.Second)

	// 10. 清理资源
	producer.Close()
	consumer.Close()
}
