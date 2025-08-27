package lua_engine

import (
	"encoding/json"
	"strings"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// HelperFunction Lua辅助函数类型
type HelperFunction = lua.LGFunction

// ContextModificationCallback 上下文修改回调函数
type ContextModificationCallback func(modifiedContext map[string]interface{})

// ExecuteOptions 规则执行选项
type ExecuteOptions struct {
	Timeout                     time.Duration               // 执行超时时间
	MaxMemory                   uint64                      // 最大内存使用(bytes)
	Context                     map[string]interface{}      // 上下文数据
	RequireFields               []string                    // 必需的返回字段
	CustomHelpers               map[string]HelperFunction   // 自定义辅助函数
	DBService                   *DBOperationService         // 数据库操作服务
	ContextModificationCallback ContextModificationCallback // 上下文修改回调函数
}

// DefaultOptions 默认选项
var DefaultOptions = &ExecuteOptions{
	Timeout:   5 * time.Second,
	MaxMemory: 10 * 1024 * 1024, // 10MB
}

// NewExecuteOptions 创建执行选项
func NewExecuteOptions() *ExecuteOptions {
	return &ExecuteOptions{
		Timeout:       5 * time.Second,
		MaxMemory:     10 * 1024 * 1024, // 10MB
		CustomHelpers: make(map[string]HelperFunction),
	}
}

// WithCustomHelper 添加自定义辅助函数到执行选项
func (opts *ExecuteOptions) WithCustomHelper(name string, fn HelperFunction) *ExecuteOptions {
	if opts.CustomHelpers == nil {
		opts.CustomHelpers = make(map[string]HelperFunction)
	}
	opts.CustomHelpers[name] = fn
	return opts
}

// WithTimeout 设置超时时间
func (opts *ExecuteOptions) WithTimeout(timeout time.Duration) *ExecuteOptions {
	opts.Timeout = timeout
	return opts
}

// WithMaxMemory 设置最大内存使用
func (opts *ExecuteOptions) WithMaxMemory(maxMemory uint64) *ExecuteOptions {
	opts.MaxMemory = maxMemory
	return opts
}

// WithContext 设置上下文数据
func (opts *ExecuteOptions) WithContext(context map[string]interface{}) *ExecuteOptions {
	opts.Context = context
	return opts
}

// WithDBService 设置数据库操作服务
func (opts *ExecuteOptions) WithDBService(dbService *DBOperationService) *ExecuteOptions {
	opts.DBService = dbService
	return opts
}

// WithContextModificationCallback 设置上下文修改回调函数
func (opts *ExecuteOptions) WithContextModificationCallback(callback ContextModificationCallback) *ExecuteOptions {
	opts.ContextModificationCallback = callback
	return opts
}

// ExecuteResult 规则执行结果
type ExecuteResult struct {
	Valid       bool                   `json:"valid"`        // 验证结果
	Action      string                 `json:"action"`       // 动作
	Error       string                 `json:"error"`        // 错误信息
	ErrorReason string                 `json:"error_reason"` // 错误原因（国际化键）
	ExecuteTime int64                  `json:"execute_time"` // 执行时间(ms)
	Context     map[string]interface{} `json:"context"`      // 上下文数据（包含修改后的值）
}

// GetContextValue 获取上下文中的值（泛型方法）
func (r *ExecuteResult) GetContextValue(key string) interface{} {
	if r.Context == nil {
		return nil
	}
	return r.Context[key]
}

// GetContextStruct 从上下文中获取结构体（通用方法）
func (r *ExecuteResult) GetContextStruct(key string, target interface{}) bool {
	if r.Context == nil {
		return false
	}

	value, exists := r.Context[key]
	if !exists {
		return false
	}

	// 尝试直接类型转换
	if converted, ok := value.(map[string]interface{}); ok {
		return r.mapToStruct(converted, target)
	}

	return false
}

// GetNestedContextStruct 从嵌套上下文中获取结构体（通用方法）
func (r *ExecuteResult) GetNestedContextStruct(path string, target interface{}) bool {
	if r.Context == nil {
		return false
	}

	keys := strings.Split(path, ".")
	current := r.Context

	for i, key := range keys {
		if current == nil {
			return false
		}

		if i == len(keys)-1 {
			// 最后一个键，返回值
			value, exists := current[key]
			if !exists {
				return false
			}

			// 尝试直接类型转换
			if converted, ok := value.(map[string]interface{}); ok {
				return r.mapToStruct(converted, target)
			}

			return false
		}

		// 中间键，继续遍历
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return false
		}
	}

	return false
}

// mapToStruct 将map转换为结构体
func (r *ExecuteResult) mapToStruct(data map[string]interface{}, target interface{}) bool {
	// 使用JSON序列化和反序列化来实现map到结构体的转换
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false
	}

	err = json.Unmarshal(jsonData, target)
	return err == nil
}

// GetContextString 获取上下文中的字符串值
func (r *ExecuteResult) GetContextString(key string) string {
	if value := r.GetContextValue(key); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetContextInt 获取上下文中的整数值
func (r *ExecuteResult) GetContextInt(key string) int {
	if value := r.GetContextValue(key); value != nil {
		switch v := value.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case int64:
			return int(v)
		case int32:
			return int(v)
		case int16:
			return int(v)
		case int8:
			return int(v)
		case uint:
			return int(v)
		case uint64:
			return int(v)
		case uint32:
			return int(v)
		case uint16:
			return int(v)
		case uint8:
			return int(v)
		}
	}
	return 0
}

// GetContextInt64 获取上下文中的int64值
func (r *ExecuteResult) GetContextInt64(key string) int64 {
	if value := r.GetContextValue(key); value != nil {
		switch v := value.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case float64:
			return int64(v)
		case uint64:
			return int64(v)
		case uint:
			return int64(v)
		}
	}
	return 0
}

// GetContextFloat64 获取上下文中的float64值
func (r *ExecuteResult) GetContextFloat64(key string) float64 {
	if value := r.GetContextValue(key); value != nil {
		switch v := value.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case int64:
			return float64(v)
		case float32:
			return float64(v)
		}
	}
	return 0.0
}

// GetContextBool 获取上下文中的布尔值
func (r *ExecuteResult) GetContextBool(key string) bool {
	if value := r.GetContextValue(key); value != nil {
		switch v := value.(type) {
		case bool:
			return v
		case int:
			return v != 0
		case float64:
			return v != 0
		case string:
			return v == "true" || v == "1" || v == "yes"
		}
	}
	return false
}

// GetContextMap 获取上下文中的map值
func (r *ExecuteResult) GetContextMap(key string) map[string]interface{} {
	if value := r.GetContextValue(key); value != nil {
		if m, ok := value.(map[string]interface{}); ok {
			return m
		}
	}
	return nil
}

// GetContextSlice 获取上下文中的slice值
func (r *ExecuteResult) GetContextSlice(key string) []interface{} {
	if value := r.GetContextValue(key); value != nil {
		if s, ok := value.([]interface{}); ok {
			return s
		}
	}
	return nil
}

// GetNestedContextValue 获取嵌套上下文中的值（支持点分隔符）
func (r *ExecuteResult) GetNestedContextValue(path string) interface{} {
	if r.Context == nil {
		return nil
	}

	keys := strings.Split(path, ".")
	current := r.Context

	for i, key := range keys {
		if current == nil {
			return nil
		}

		if i == len(keys)-1 {
			// 最后一个键，返回值
			return current[key]
		}

		// 中间键，继续遍历
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			return nil
		}
	}

	return nil
}

// GetNestedContextString 获取嵌套上下文中的字符串值
func (r *ExecuteResult) GetNestedContextString(path string) string {
	if value := r.GetNestedContextValue(path); value != nil {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetNestedContextInt 获取嵌套上下文中的整数值
func (r *ExecuteResult) GetNestedContextInt(path string) int {
	if value := r.GetNestedContextValue(path); value != nil {
		switch v := value.(type) {
		case int:
			return v
		case float64:
			return int(v)
		case int64:
			return int(v)
		}
	}
	return 0
}

// GetNestedContextFloat64 获取嵌套上下文中的float64值
func (r *ExecuteResult) GetNestedContextFloat64(path string) float64 {
	if value := r.GetNestedContextValue(path); value != nil {
		switch v := value.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case int64:
			return float64(v)
		}
	}
	return 0.0
}

// GetNestedContextBool 获取嵌套上下文中的布尔值
func (r *ExecuteResult) GetNestedContextBool(path string) bool {
	if value := r.GetNestedContextValue(path); value != nil {
		switch v := value.(type) {
		case bool:
			return v
		case int:
			return v != 0
		case float64:
			return v != 0
		case string:
			return v == "true" || v == "1" || v == "yes"
		}
	}
	return false
}

// HasContextKey 检查上下文中是否存在指定的键
func (r *ExecuteResult) HasContextKey(key string) bool {
	if r.Context == nil {
		return false
	}
	_, exists := r.Context[key]
	return exists
}

// HasNestedContextKey 检查嵌套上下文中是否存在指定的路径
func (r *ExecuteResult) HasNestedContextKey(path string) bool {
	return r.GetNestedContextValue(path) != nil
}

// GetContextKeys 获取上下文中所有的键
func (r *ExecuteResult) GetContextKeys() []string {
	if r.Context == nil {
		return nil
	}
	keys := make([]string, 0, len(r.Context))
	for key := range r.Context {
		keys = append(keys, key)
	}
	return keys
}
