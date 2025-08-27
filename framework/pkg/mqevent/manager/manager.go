package manager

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/idempotence"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/mqevent"
	"strings"
	"sync"
)

// EventBusManager 事件总线管理器
type EventBusManager struct {
	subsampling     ISubscribeSmServerApi
	ebs             mqevent.IMQEventBus
	idempotenceTool idempotence.IdempotencyTool
	// 存储订阅ID的映射关系
	subscriptions sync.Map
	// 存储处理器映射关系
	handlers sync.Map
	// 存储处理器是否需要幂等性处理
	handlerIdempotence sync.Map
	// 存储订阅状态
	subscriptionStatus sync.Map
}

// NewEventBusManager 创建事件总线管理器
func NewEventBusManager(subsampling ISubscribeSmServerApi, ebs mqevent.IMQEventBus, idempotenceTool idempotence.IdempotencyTool) EventManager {
	return &EventBusManager{
		subsampling:     subsampling,
		idempotenceTool: idempotenceTool,
		ebs:             ebs,
	}
}

// RegisterSubscribeHandel 注册普通事件处理器
func (e *EventBusManager) RegisterSubscribeHandel(topic, channel string, handler EventHandler) error {
	key := getHandlerKey(topic, channel)

	// 检查处理器是否已注册
	if _, exists := e.handlers.Load(key); exists {
		return fmt.Errorf("handler already registered for topic: %s, channel: %s", topic, channel)
	}

	e.handlers.Store(key, handler)
	// 标记不需要幂等性处理
	e.handlerIdempotence.Store(key, false)
	return nil
}

// RegisterIdempotenceSubscribeHandel 注册幂等性事件处理器
func (e *EventBusManager) RegisterIdempotenceSubscribeHandel(topic, channel string, handler EventHandler) error {
	key := getHandlerKey(topic, channel)

	// 检查处理器是否已注册
	if _, exists := e.handlers.Load(key); exists {
		return fmt.Errorf("idempotence handler already registered for topic: %s, channel: %s", topic, channel)
	}

	e.handlers.Store(key, handler)
	// 标记需要幂等性处理
	e.handlerIdempotence.Store(key, true)
	return nil
}

// ActivateSubscription 激活订阅
func (e *EventBusManager) ActivateSubscription(topic, channel string) error {
	key := getHandlerKey(topic, channel)

	// 检查订阅是否存在
	if _, ok := e.subscriptions.Load(key); !ok {
		return fmt.Errorf("subscription not found: %s", key)
	}

	// 设置订阅状态为激活
	e.subscriptionStatus.Store(key, true)
	hlog.Infof("Subscription activated: %s", key)

	return nil
}

// DeactivateSubscription 停用订阅
func (e *EventBusManager) DeactivateSubscription(topic, channel string) error {
	key := getHandlerKey(topic, channel)

	// 检查订阅是否存在
	if _, ok := e.subscriptions.Load(key); !ok {
		return fmt.Errorf("subscription not found: %s", key)
	}

	// 设置订阅状态为停用
	e.subscriptionStatus.Store(key, false)
	hlog.Infof("Subscription deactivated: %s", key)

	return nil
}

// handel 处理事件
func (e *EventBusManager) handel(ctx context.Context, event mqevent.Event, channel string) error {
	key := getHandlerKey(event.GetType(), channel)

	// 检查订阅状态
	if status, ok := e.subscriptionStatus.Load(key); ok {
		if !status.(bool) {
			// 如果订阅已停用，直接返回
			hlog.Debugf("Subscription is deactivated, skipping event: %s", key)
			return nil
		}
	}
	ctx = actx.WithTenantId(ctx, event.GetTenantID())

	// 获取事件的参数
	params, err := e.subsampling.GetParameters(ctx, event.GetType(), channel)
	if err != nil {
		hlog.CtxErrorf(ctx, "failed to get subscription parameters: %v", err)
		return fmt.Errorf("failed to get subscription parameters: %w", err)
	}
	eventCtx := mqevent.NewEventContext(ctx, event, channel, params)

	// 从管理器中获取对应的处理器
	handler, ok := e.handlers.Load(key)
	if !ok {
		hlog.CtxErrorf(ctx, "handler not found for event: %s", event.GetType())
		return fmt.Errorf("handler not found for event: %s", event.GetType())
	}

	// 检查是否需要幂等性处理
	if needIdempotence, ok := e.handlerIdempotence.Load(key); ok && needIdempotence.(bool) {
		check, err := e.idempotenceTool.Check(ctx, event.GetType(), channel, event.GetID())
		if err != nil {
			hlog.CtxErrorf(ctx, "failed to check idempotence: %v", err)
			return fmt.Errorf("failed to check idempotence: %w", err)
		}
		if check {
			// 如果已经处理过，直接返回
			hlog.Infof("Event already processed, skipping: %s", event.GetID())
			return nil
		}

		// 处理
		err = handler.(EventHandler).Handle(eventCtx)
		if err != nil {
			hlog.CtxErrorf(ctx, "failed to process event: %v", err)
			return fmt.Errorf("failed to process event: %w", err)
		}
		// 使用幂等性工具处理
		return e.idempotenceTool.MarkProcessed(ctx, event.GetType(), channel, event.GetID())
	}

	return handler.(EventHandler).Handle(eventCtx)
}

// RetryDeadLetter 重试死信队列
func (e *EventBusManager) RetryDeadLetter(ctx context.Context, dead *mqevent.DeadLetterEvent) error {
	return e.handel(ctx, dead.OriginalEvent, dead.Channel)
}

// Start 启动事件总线管理器
func (e *EventBusManager) Start() error {
	// 获取所有激活状态的订阅
	subscribes, err := e.subsampling.GetByStatus(context.Background(), 1)
	if err != nil {
		return fmt.Errorf("failed to get active subscriptions: %w", err)
	}

	// 创建受数据库控制的订阅映射
	dbControlledSubscriptions := make(map[string]bool)
	for _, subscribe := range subscribes {
		key := getHandlerKey(subscribe.Topic, subscribe.Group)
		dbControlledSubscriptions[key] = subscribe.Status != 2
	}

	// 遍历所有已注册的处理器
	e.handlers.Range(func(key, value interface{}) bool {
		topicChannel := key.(string)
		topic, channel, ok := parseHandlerKey(topicChannel)
		if !ok {
			panic(fmt.Sprintf("invalid handler key: %s", topicChannel))
		}
		// 创建事件处理函数
		eventHandler := mqevent.EventHandlerFunc(func(ctx context.Context, event mqevent.Event) error {
			// 使用 defer 和 recover 捕获 panic
			defer func() {
				if r := recover(); r != nil {
					// 记录 panic 信息
					hlog.CtxErrorf(ctx, "panic recovered in event handler: %v, event: %+v", r, event)
					// 可以在这里添加告警通知等逻辑
				}
			}()
			ctx = actx.WithTenantId(ctx, event.GetTenantID())
			// 处理事件
			return e.handel(ctx, event, channel)
		})
		// 订阅事件
		subscriptionID, err := e.ebs.Subscribe(topic, channel, eventHandler)
		if err != nil {
			hlog.Errorf("failed to subscribe event: %v", err)
			return true // 继续处理其他订阅
		}

		// 存储订阅ID
		e.subscriptions.Store(topicChannel, subscriptionID)

		// 设置订阅状态
		if v, ok := dbControlledSubscriptions[topicChannel]; ok {
			// 如果是受数据库控制的订阅，根据数据库状态设置
			e.subscriptionStatus.Store(topicChannel, v)
		} else {
			// 如果是普通订阅，默认激活
			e.subscriptionStatus.Store(topicChannel, true)
		}

		return true
	})

	// 订阅死信队列
	if err := e.subscribeDeadLetter(); err != nil {
		return fmt.Errorf("failed to subscribe dead letter queue: %w", err)
	}

	return nil
}

// subscribeDeadLetter 订阅死信队列
func (e *EventBusManager) subscribeDeadLetter() error {
	// 创建死信处理函数
	deadLetterHandler := mqevent.DeadLetterEventHandlerFunc(func(ctx context.Context, event *mqevent.DeadLetterEvent) error {
		// 使用 defer 和 recover 捕获 panic
		defer func() {
			if r := recover(); r != nil {
				// 记录 panic 信息
				hlog.CtxErrorf(ctx, "panic recovered in dead letter handler: %v, event: %+v", r, event)
			}
		}()

		ctx = actx.WithTenantId(ctx, event.OriginalEvent.GetTenantID())

		// 保存死信事件
		if err := e.subsampling.DadEventSave(ctx, event); err != nil {
			hlog.CtxErrorf(ctx, "failed to save dead letter event: %v", err)
			return fmt.Errorf("failed to save dead letter event: %w", err)
		}

		return nil
	})

	// 订阅死信队列
	subscriptionID, err := e.ebs.SubscribeDeadLetter(deadLetterHandler)
	if err != nil {
		return fmt.Errorf("failed to subscribe dead letter queue: %w", err)
	}

	// 存储死信队列订阅ID
	e.subscriptions.Store(mqevent.EventDeadLetterTopic, subscriptionID)

	return nil
}

// Stop 停止事件总线管理器
func (e *EventBusManager) Stop() error {
	// 遍历所有订阅并取消
	e.subscriptions.Range(func(key, value interface{}) bool {
		if err := e.ebs.Unsubscribe(value.(string)); err != nil {
			// 记录错误但继续处理其他订阅
			hlog.Errorf("failed to unsubscribe: %v", err)
		}
		return true
	})

	// 清空所有映射
	e.subscriptions = sync.Map{}
	e.handlers = sync.Map{}
	e.handlerIdempotence = sync.Map{}
	e.subscriptionStatus = sync.Map{}

	return nil
}

// Close 关闭事件总线管理器
func (e *EventBusManager) Close() error {
	// 先停止所有订阅
	if err := e.Stop(); err != nil {
		return fmt.Errorf("failed to stop event bus manager: %w", err)
	}

	// 关闭事件总线
	if err := e.ebs.Close(); err != nil {
		return fmt.Errorf("failed to close event bus: %w", err)
	}

	return nil
}

// getHandlerKey 获取处理器键值
func getHandlerKey(topic, channel string) string {
	return fmt.Sprintf("%s:%s", topic, channel)
}

// parseHandlerKey 拆解处理器键值为 topic 和 channel
func parseHandlerKey(key string) (topic, channel string, ok bool) {
	parts := strings.SplitN(key, ":", 2)
	if len(parts) != 2 {
		return "", "", false
	}
	return parts[0], parts[1], true
}
