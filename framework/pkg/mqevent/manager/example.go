package manager

//
//import (
//	"encoding/json"
//	"fmt"
//	"time"
//
//	"github.com/cloudwego/hertz/pkg/common/hlog"
//	"go.uber.org/zap"
//)
//
//// ExampleEvent 示例事件
//type ExampleEvent struct {
//	ID      string    `json:"id"`
//	Message string    `json:"message"`
//	Time    time.Time `json:"time"`
//}
//
//// ExampleParameterProvider 示例参数提供者
//type ExampleParameterProvider struct {
//	params map[string]map[string]interface{}
//}
//
//// GetParameters 实现参数提供者接口
//func (p *ExampleParameterProvider) GetParameters(topic, channel string) (map[string]interface{}, error) {
//	key := fmt.Sprintf("%s:%s", topic, channel)
//	if params, exists := p.params[key]; exists {
//		return params, nil
//	}
//	return make(map[string]interface{}), nil
//}
//
//// ExampleUsage 示例用法
//func ExampleUsage() {
//	// 创建配置
//	config := &Config{
//		Type:                 EventTypeNats,
//		URL:                  "nats://localhost:4222",
//		StreamName:           "example-stream",
//		StreamSubjects:       []string{"example.>"},
//		MaxPendingMessages:   100,
//		MaxRetries:           3,
//		RetryDelay:           time.Second * 5,
//		DeadLetterEnabled:    true,
//		DeadLetterTopic:      "dead.letter",
//		DeadLetterMaxRetries: 3,
//		IdempotencyEnabled:   true,
//		IdempotencyTTL:       time.Hour * 24,
//		ParameterEnabled:     true,
//		ParameterTTL:         time.Hour,
//	}
//
//	// 创建事件管理器
//	manager, err := NewEventManager(config)
//	if err != nil {
//		hlog.Fatalf("Failed to create event manager: %v", err)
//	}
//	defer manager.Close()
//
//	// 设置参数提供者
//	parameterProvider := &ExampleParameterProvider{
//		params: map[string]map[string]interface{}{
//			"example.topic:example-group": {
//				"maxRetries": 3,
//				"timeout":    "5s",
//			},
//		},
//	}
//	manager.SetParameterProvider(parameterProvider)
//
//	// 注册事件处理器
//	err = manager.RegisterHandler("example.topic", func(ctx *EventContext) error {
//		// 解析事件数据
//		var event ExampleEvent
//		err := json.Unmarshal(ctx.event.Payload, &event)
//		if err != nil {
//			return err
//		}
//
//		// 获取参数
//		params, err := manager.GetParameters(ctx.event.Topic, ctx.event.Group)
//		if err != nil {
//			hlog.Error("Failed to get parameters", zap.Error(err))
//		}
//
//		// 处理事件
//		hlog.Info("Processing event",
//			zap.String("id", event.ID),
//			zap.String("message", event.Message),
//			zap.Time("time", event.Time),
//			zap.Any("parameters", params))
//
//		// 确认消息
//		return ctx.ack()
//	})
//	if err != nil {
//		hlog.Fatalf("Failed to register handler: %v", err)
//	}
//
//	// 订阅主题
//	err = manager.Subscribe("example.topic", "example-group")
//	if err != nil {
//		hlog.Fatalf("Failed to subscribe: %v", err)
//	}
//
//	// 处理死信消息
//	err = manager.ProcessDeadLetter(func(ctx *EventContext) error {
//		hlog.Info("Processing dead letter",
//			zap.String("topic", ctx.event.Topic),
//			zap.String("group", ctx.event.Group),
//			zap.Int("retryCount", ctx.retryCount))
//
//		// 重试消息
//		return manager.RetryDeadLetter(ctx)
//	})
//	if err != nil {
//		hlog.Fatalf("Failed to process dead letter: %v", err)
//	}
//
//	// 发布事件
//	event := ExampleEvent{
//		ID:      "123",
//		Message: "Hello, World!",
//		Time:    time.Now(),
//	}
//	payload, err := json.Marshal(event)
//	if err != nil {
//		hlog.Fatalf("Failed to marshal event: %v", err)
//	}
//
//	headers := map[string]string{
//		"source": "example",
//	}
//	err = manager.Publish("example.topic", payload, headers)
//	if err != nil {
//		hlog.Fatalf("Failed to publish event: %v", err)
//	}
//
//	// 等待一段时间以处理消息
//	time.Sleep(time.Second * 5)
//}
