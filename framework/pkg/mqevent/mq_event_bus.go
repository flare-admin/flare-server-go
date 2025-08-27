package mqevent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq/models"
	"sync"
	"time"
)

// Bus 基于 MQ Server 的事件总线实现
type Bus struct {
	server          mq.Server
	subscriptions   map[string]string // subscriptionID -> topic:channel
	deadLetterSubID string
	mu              sync.RWMutex
}

// NewMQEventBus 创建基于 MQ Server 的事件总线
func NewMQEventBus(server mq.Server) IMQEventBus {
	return &Bus{
		server:        server,
		subscriptions: make(map[string]string),
	}
}

// Publish 发布事件
func (b *Bus) Publish(ctx context.Context, event Event) error {
	// 构建消息头
	headers := make(map[string]string)
	headers["event_id"] = event.GetID()
	headers["event_type"] = event.GetType()
	headers["timestamp"] = event.GetTimestamp().Format(time.RFC3339)
	tenantId := event.GetTenantID()
	if tenantId == "" {
		tenantId = actx.GetTenantId(ctx)
	}
	headers["tenant_id"] = tenantId
	event.SetTenantID(tenantId)

	// 合并事件元数据
	for k, v := range event.GetMetadata() {
		headers[k] = v
	}

	// 序列化事件数据
	data, err := json.Marshal(event.GetData())
	if err != nil {
		return fmt.Errorf("序列化事件数据失败: %v", err)
	}

	// 创建基础消息
	msg := models.NewBaseMessage(event.GetType(), data, headers)

	// 发布消息
	return b.server.Publish(ctx, msg)
}

// Subscribe 订阅事件
func (b *Bus) Subscribe(eventType string, channel string, handler EventHandler) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 生成订阅ID
	subscriptionID := fmt.Sprintf("%s:%s:%d", eventType, channel, time.Now().UnixNano())

	// 订阅消息
	err := b.server.Subscribe(context.Background(), eventType, channel, func(msg *models.BaseMessage) error {
		// 解析消息头
		headers := msg.GetHeaders()
		eventID := headers["event_id"]
		eventType := headers["event_type"]
		tenantID := headers["tenant_id"]
		timestamp, _ := time.Parse(time.RFC3339, headers["timestamp"])

		// 创建基础事件
		event := NewBaseEvent(eventType, nil,
			WithID(eventID),
			WithTimestamp(timestamp),
			WithTenantID(tenantID),
			WithMetadata(headers),
		)

		// 解析事件数据
		var data interface{}
		if err := json.Unmarshal(msg.GetPayload(), &data); err != nil {
			return fmt.Errorf("反序列化事件数据失败: %v", err)
		}
		event.Data = data

		// 处理事件
		return handler.Handle(context.Background(), event)
	})

	if err != nil {
		return "", fmt.Errorf("订阅事件失败: %v", err)
	}

	// 保存订阅信息
	b.subscriptions[subscriptionID] = fmt.Sprintf("%s:%s", eventType, channel)

	return subscriptionID, nil
}

// Unsubscribe 取消订阅
func (b *Bus) Unsubscribe(subscriptionID string) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 获取订阅信息
	topicChannel, ok := b.subscriptions[subscriptionID]
	if !ok {
		return fmt.Errorf("订阅不存在: %s", subscriptionID)
	}

	// 解析主题和通道
	var topic, channel string
	fmt.Sscanf(topicChannel, "%s:%s", &topic, &channel)

	// 取消订阅
	if err := b.server.Unsubscribe(topic, channel); err != nil {
		return fmt.Errorf("取消订阅失败: %v", err)
	}

	// 删除订阅信息
	delete(b.subscriptions, subscriptionID)

	return nil
}

// SubscribeDeadLetter 订阅死信队列
func (b *Bus) SubscribeDeadLetter(handler DeadLetterEventHandlerFunc) (string, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 生成订阅ID
	subscriptionID := fmt.Sprintf("dead_letter:%d", time.Now().UnixNano())

	// 订阅死信队列
	err := b.server.SubscribeDeadLetter(context.Background(), func(msg models.DeadLetterMessage) error {
		// 解析原始事件
		var originalEvent *BaseEvent
		if err := json.Unmarshal(msg.Payload, &originalEvent); err != nil {
			return fmt.Errorf("反序列化原始事件失败: %v", err)
		}

		// 创建死信事件
		deadLetterEvent := &DeadLetterEvent{
			OriginalEvent: originalEvent,
			Error:         msg.Error,
			RetryCount:    msg.RetryCount,
			LastAttempt:   msg.Timestamp,
			NextRetry:     msg.Timestamp.Add(time.Second * time.Duration(msg.RetryCount+1)),
			EventType:     msg.Topic,
			OriginalTopic: msg.Topic,
			Channel:       msg.Channel,
		}

		// 处理死信事件
		return handler(context.Background(), deadLetterEvent)
	})

	if err != nil {
		return "", fmt.Errorf("订阅死信队列失败: %v", err)
	}

	// 保存订阅ID
	b.deadLetterSubID = subscriptionID

	return subscriptionID, nil
}

// Close 关闭事件总线
func (b *Bus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	// 取消所有订阅
	for subscriptionID := range b.subscriptions {
		if err := b.Unsubscribe(subscriptionID); err != nil {
			return fmt.Errorf("取消订阅失败: %v", err)
		}
	}

	// 取消死信队列订阅
	if b.deadLetterSubID != "" {
		if err := b.server.UnsubscribeDeadLetter(); err != nil {
			return fmt.Errorf("取消死信队列订阅失败: %v", err)
		}
	}

	// 关闭服务器
	return b.server.Close()
}
