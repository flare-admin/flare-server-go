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
	// 1. 创建 NATS 配置
	natsConfig := nats.NewConfigBuilder("nats://127.0.0.1:4222").
		MaxRetries(5).
		ReconnectWait(5). // 5秒
		AckWait(60).      // 60秒
		QueueGroup("group1").
		Build()

	// 2. 转换为通用 MQ 配置
	mqConfig := natsConfig.ToMQConfig()

	// 3. 创建统一服务器
	server, err := server.NewServer(mqConfig)
	if err != nil {
		panic(fmt.Sprintf("创建服务器失败: %v", err))
	}
	defer server.Close()

	// 4. 订阅死信队列
	ctx := context.Background()
	err = server.SubscribeDeadLetter(ctx, func(msg models.DeadLetterMessage) error {
		fmt.Printf("收到死信消息:\n")
		fmt.Printf("  主题: %s\n", msg.Topic)
		fmt.Printf("  通道: %s\n", msg.Channel)
		fmt.Printf("  错误: %s\n", msg.Error)
		fmt.Printf("  重试次数: %d\n", msg.RetryCount)
		fmt.Printf("  时间戳: %s\n", msg.Timestamp)
		fmt.Printf("  死信时间: %s\n", msg.DeadTime)
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("订阅死信队列失败: %v", err))
	}

	// 5. 订阅主题
	err = server.Subscribe(ctx, "test.topic", "test.channel", func(msg mq.Message) error {
		fmt.Printf("收到消息:\n")
		fmt.Printf("  主题: %s\n", msg.GetTopic())
		fmt.Printf("  负载: %s\n", string(msg.GetPayload()))
		fmt.Printf("  头信息: %v\n", msg.GetHeaders())

		// 模拟处理失败，触发重试和死信
		if string(msg.GetPayload()) == "trigger_error" {
			return fmt.Errorf("模拟处理失败")
		}
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("订阅主题失败: %v", err))
	}

	// 6. 发布消息
	// 创建基础消息
	normalMsg := models.NewBaseMessage("test.topic", []byte("Hello MQ!"), map[string]string{
		"tenant_id": "123",
		"user_id":   "456",
	})

	// 发布普通消息
	err = server.Publish(ctx, normalMsg)
	if err != nil {
		panic(fmt.Sprintf("发布消息失败: %v", err))
	}

	// 创建延迟消息
	delayMsg := models.NewBaseMessage("test.topic", []byte("Delayed Message!"), map[string]string{
		"tenant_id": "123",
		"user_id":   "456",
	})

	// 发布延迟消息
	err = server.PublishDelay(ctx, delayMsg, 5*time.Second)
	if err != nil {
		panic(fmt.Sprintf("发布延迟消息失败: %v", err))
	}

	// 创建会触发错误的消息
	errorMsg := models.NewBaseMessage("test.topic", []byte("trigger_error"), map[string]string{
		"tenant_id": "123",
		"user_id":   "456",
	})

	// 发布会触发错误的消息
	err = server.Publish(ctx, errorMsg)
	if err != nil {
		panic(fmt.Sprintf("发布错误消息失败: %v", err))
	}

	// 7. 等待消息处理
	time.Sleep(10 * time.Second)
}
