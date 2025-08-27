package lua_engine

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"sync"
	"time"

	"github.com/flare-admin/flare-server-go/framework/pkg/database"

	lua "github.com/yuin/gopher-lua"
)

// RuleExecutor Lua规则执行器
type RuleExecutor struct {
	pool               sync.Pool
	customHelpers      map[string]HelperFunction // 自定义辅助函数
	customHelpersMutex sync.RWMutex              // 自定义辅助函数读写锁
	dbService          *DBOperationService       // 数据库操作服务
}

// NewRuleExecutor 创建规则执行器
func NewRuleExecutor() *RuleExecutor {
	return &RuleExecutor{
		pool: sync.Pool{
			New: func() interface{} {
				return lua.NewState()
			},
		},
		customHelpers: make(map[string]HelperFunction),
	}
}

// NewRuleExecutorWithDB 创建带数据库支持的规则执行器
func NewRuleExecutorWithDB(db database.IDataBase) *RuleExecutor {
	return &RuleExecutor{
		pool: sync.Pool{
			New: func() interface{} {
				return lua.NewState()
			},
		},
		customHelpers: make(map[string]HelperFunction),
		dbService:     NewDBOperationService(db),
	}
}

// Execute 执行规则
func (e *RuleExecutor) Execute(script string, opts *ExecuteOptions) (*ExecuteResult, error) {
	if opts == nil {
		opts = DefaultOptions
	}

	// 从池中获取Lua状态机
	L := e.pool.Get().(*lua.LState)
	defer e.pool.Put(L)

	// 高效重置状态机（避免重新创建）
	e.resetLuaState(L, opts.MaxMemory)

	// 注入上下文数据
	contextTable := L.NewTable()
	e.injectContextData(L, contextTable, opts.Context)
	L.SetGlobal("context", contextTable)

	// 注入内置辅助函数
	e.registerHelperFunctions(L)

	// 注入自定义辅助函数
	e.registerCustomHelperFunctions(L, opts.CustomHelpers)

	// 注入数据库操作辅助函数
	dbService := e.dbService
	if opts.DBService != nil {
		dbService = opts.DBService
	}
	e.registerDBHelperFunctions(L, dbService)

	// 执行脚本
	start := utils.GetTimeNow()
	done := make(chan error, 1)
	go func() {
		done <- L.DoString(script)
	}()

	// 等待执行完成或超时
	var err error
	select {
	case err = <-done:
	case <-time.After(opts.Timeout):
		err = fmt.Errorf("执行超时")
	}

	if err != nil {
		return nil, err
	}

	// 获取执行结果
	result := &ExecuteResult{
		ExecuteTime: time.Since(start).Milliseconds(),
	}

	// 高效获取结果
	e.extractResults(L, result)

	// 提取修改后的上下文
	modifiedContext := e.extractModifiedContext(L, opts.Context)
	result.Context = modifiedContext

	// 调用上下文修改回调函数
	if opts.ContextModificationCallback != nil {
		opts.ContextModificationCallback(modifiedContext)
	}

	return result, nil
}

// ValidateScript 验证脚本语法
func (e *RuleExecutor) ValidateScript(script string) error {
	L := e.pool.Get().(*lua.LState)
	defer e.pool.Put(L)

	// 重置状态机
	L.Close()
	L = lua.NewState()

	// 注入空的上下文和辅助函数
	contextTable := L.NewTable()
	L.SetGlobal("context", contextTable)
	e.registerHelperFunctions(L)

	// 注册自定义辅助函数
	e.registerCustomHelperFunctions(L, nil)

	// 注册数据库操作辅助函数
	e.registerDBHelperFunctions(L, e.dbService)

	// 解析脚本
	_, err := L.LoadString(script)
	if err != nil {
		return fmt.Errorf("脚本语法错误: %v", err)
	}

	return nil
}

// CompileTemplate 编译规则模板
func (e *RuleExecutor) CompileTemplate(template string, params map[string]interface{}) (string, error) {
	L := e.pool.Get().(*lua.LState)
	defer e.pool.Put(L)

	// 重置状态机
	L.Close()
	L = lua.NewState()

	// 注入模板参数
	paramsTable := L.NewTable()
	for k, v := range params {
		// 直接转换为Lua值，不支持的类型将被忽略
		luaValue := e.convertToLuaValue(L, v)
		if luaValue != lua.LNil {
			L.SetTable(paramsTable, lua.LString(k), luaValue)
		}
	}
	L.SetGlobal("params", paramsTable)

	// 包装模板
	script := fmt.Sprintf(`
		local template = %q
		local result = template
		for k, v in pairs(params) do
			result = string.gsub(result, "${" .. k .. "}", tostring(v))
		end
		return result
	`, template)

	// 执行模板编译
	if err := L.DoString(script); err != nil {
		return "", fmt.Errorf("模板编译错误: %v", err)
	}

	// 获取编译结果
	result := L.Get(-1)
	if result.Type() != lua.LTString {
		return "", fmt.Errorf("模板必须返回字符串")
	}

	return result.String(), nil
}

// registerCustomHelperFunctions 注册自定义辅助函数到Lua状态机
func (e *RuleExecutor) registerCustomHelperFunctions(L *lua.LState, customHelpers map[string]HelperFunction) {
	// 注册传入的自定义辅助函数
	if customHelpers != nil {
		for name, fn := range customHelpers {
			if name != "" && fn != nil {
				L.SetGlobal(name, L.NewFunction(fn))
			}
		}
	}

	// 注册执行器实例中存储的自定义辅助函数
	e.customHelpersMutex.RLock()
	defer e.customHelpersMutex.RUnlock()

	for name, fn := range e.customHelpers {
		if name != "" && fn != nil {
			L.SetGlobal(name, L.NewFunction(fn))
		}
	}
}

// AddCustomHelper 添加自定义辅助函数
func (e *RuleExecutor) AddCustomHelper(name string, fn HelperFunction) error {
	if name == "" {
		return fmt.Errorf("辅助函数名称不能为空")
	}

	if fn == nil {
		return fmt.Errorf("辅助函数不能为空")
	}

	e.customHelpersMutex.Lock()
	defer e.customHelpersMutex.Unlock()

	e.customHelpers[name] = fn
	return nil
}

// RemoveCustomHelper 移除自定义辅助函数
func (e *RuleExecutor) RemoveCustomHelper(name string) error {
	if name == "" {
		return fmt.Errorf("辅助函数名称不能为空")
	}

	e.customHelpersMutex.Lock()
	defer e.customHelpersMutex.Unlock()

	if _, exists := e.customHelpers[name]; !exists {
		return fmt.Errorf("辅助函数 '%s' 不存在", name)
	}

	delete(e.customHelpers, name)
	return nil
}

// GetCustomHelper 获取自定义辅助函数
func (e *RuleExecutor) GetCustomHelper(name string) (HelperFunction, bool) {
	if name == "" {
		return nil, false
	}

	e.customHelpersMutex.RLock()
	defer e.customHelpersMutex.RUnlock()

	fn, exists := e.customHelpers[name]
	return fn, exists
}

// ListCustomHelpers 列出所有自定义辅助函数名称
func (e *RuleExecutor) ListCustomHelpers() []string {
	e.customHelpersMutex.RLock()
	defer e.customHelpersMutex.RUnlock()

	names := make([]string, 0, len(e.customHelpers))
	for name := range e.customHelpers {
		names = append(names, name)
	}

	return names
}

// ClearCustomHelpers 清空所有自定义辅助函数
func (e *RuleExecutor) ClearCustomHelpers() {
	e.customHelpersMutex.Lock()
	defer e.customHelpersMutex.Unlock()

	e.customHelpers = make(map[string]HelperFunction)
}
