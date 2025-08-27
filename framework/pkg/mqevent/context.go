package mqevent

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"reflect"
	"time"
)

// EventContext 事件上下文
type EventContext struct {
	// 基础信息
	ctx        context.Context
	event      Event
	channel    string
	parameters map[string]interface{}
	// 时间信息
	receivedTime  time.Time
	processedTime time.Time
}

// NewEventContext 创建事件上下文
func NewEventContext(ctx context.Context, event Event, channel string, parameters map[string]interface{}) *EventContext {
	return &EventContext{
		ctx:          ctx,
		event:        event,
		channel:      channel,
		parameters:   parameters,
		receivedTime: utils.GetTimeNow(),
	}
}

// Context 获取基础上下文
func (c *EventContext) Context() context.Context {
	return c.ctx
}

// Event 获取事件
func (c *EventContext) Event() Event {
	return c.event
}

// Topic 获取主题
func (c *EventContext) Topic() string {
	return c.event.GetType()
}

// Channel 获取消费通道
func (c *EventContext) Channel() string {
	return c.channel
}

// GetData 获取事件数据
func (c *EventContext) GetData() interface{} {
	return c.event.GetData()
}

// GetMetadata 获取元数据
func (c *EventContext) GetMetadata() map[string]string {
	return c.event.GetMetadata()
}

// GetMetadataValue 获取元数据值
func (c *EventContext) GetMetadataValue(key string) string {
	return c.event.GetMetadata()[key]
}

// GetParameter 获取参数
func (c *EventContext) GetParameter(key string) interface{} {
	return c.parameters[key]
}

// SetParameter 设置参数
func (c *EventContext) SetParameter(key string, value interface{}) {
	c.parameters[key] = value
}

// ReceivedTime 获取接收时间
func (c *EventContext) ReceivedTime() time.Time {
	return c.receivedTime
}

// ProcessedTime 获取处理时间
func (c *EventContext) ProcessedTime() time.Time {
	return c.processedTime
}

// SetProcessedTime 设置处理时间
func (c *EventContext) SetProcessedTime(t time.Time) {
	c.processedTime = t
}

// GetParameterAs 泛型函数：根据类型获取参数
func GetParameterAs[T any](ctx *EventContext, key string) (T, error) {
	var result T
	param := ctx.parameters[key]
	if param == nil {
		return result, fmt.Errorf("parameter not found: %s", key)
	}
	// 处理基础类型
	switch any(result).(type) {
	case string:
		if str, ok := param.(string); ok {
			return any(str).(T), nil
		}
	case int:
		if num, ok := param.(int); ok {
			return any(num).(T), nil
		}
	case float64:
		if num, ok := param.(float64); ok {
			return any(num).(T), nil
		}
	case bool:
		if b, ok := param.(bool); ok {
			return any(b).(T), nil
		}
	}
	// 处理结构体和map
	switch v := param.(type) {
	case map[string]interface{}:
		bytes, err := json.Marshal(v)
		if err != nil {
			return result, fmt.Errorf("failed to marshal map: %w", err)
		}
		if err := json.Unmarshal(bytes, &result); err != nil {
			return result, fmt.Errorf("failed to unmarshal to target type: %w", err)
		}
		return result, nil
	default:
		if reflect.TypeOf(v) == reflect.TypeOf(result) {
			return v.(T), nil
		}
		bytes, err := json.Marshal(v)
		if err != nil {
			return result, fmt.Errorf("failed to marshal parameter: %w", err)
		}
		if err := json.Unmarshal(bytes, &result); err != nil {
			return result, fmt.Errorf("failed to unmarshal to target type: %w", err)
		}
		return result, nil
	}
}

// GetDataAs 泛型函数：根据类型获取事件数据
func GetDataAs[T any](ctx *EventContext) (T, error) {
	var result T
	data := ctx.event.GetData()
	// 处理基础类型
	switch any(result).(type) {
	case string:
		if str, ok := data.(string); ok {
			return any(str).(T), nil
		}
	case int:
		if num, ok := data.(int); ok {
			return any(num).(T), nil
		}
	case float64:
		if num, ok := data.(float64); ok {
			return any(num).(T), nil
		}
	case bool:
		if b, ok := data.(bool); ok {
			return any(b).(T), nil
		}
	}
	// 处理结构体和map
	switch v := data.(type) {
	case map[string]interface{}:
		bytes, err := json.Marshal(v)
		if err != nil {
			return result, fmt.Errorf("failed to marshal map: %w", err)
		}
		if err := json.Unmarshal(bytes, &result); err != nil {
			return result, fmt.Errorf("failed to unmarshal to target type: %w", err)
		}
		return result, nil
	default:
		if reflect.TypeOf(v) == reflect.TypeOf(result) {
			return v.(T), nil
		}
		bytes, err := json.Marshal(v)
		if err != nil {
			return result, fmt.Errorf("failed to marshal data: %w", err)
		}
		if err := json.Unmarshal(bytes, &result); err != nil {
			return result, fmt.Errorf("failed to unmarshal to target type: %w", err)
		}
		return result, nil
	}
}

// GetAllParametersAs 泛型函数：将所有参数转换为指定结构体或 map 类型
func GetAllParametersAs[T any](ctx *EventContext) (T, error) {
	var result T
	if ctx == nil || ctx.parameters == nil {
		return result, fmt.Errorf("parameters is nil")
	}
	bytes, err := json.Marshal(ctx.parameters)
	if err != nil {
		return result, fmt.Errorf("failed to marshal parameters: %w", err)
	}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return result, fmt.Errorf("failed to unmarshal parameters to target type: %w", err)
	}
	return result, nil
}
