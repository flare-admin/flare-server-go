package lua_engine

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// TestHelperFunctionType 测试 HelperFunction 类型是否正确
func TestHelperFunctionType(t *testing.T) {
	// 创建一个符合 HelperFunction 类型的函数
	var fn HelperFunction = func(L *lua.LState) int {
		L.Push(lua.LString("test"))
		return 1
	}

	// 验证函数不为空
	if fn == nil {
		t.Error("HelperFunction 不应该为 nil")
	}

	// 验证类型兼容性
	var _ lua.LGFunction = fn
}

// TestExecuteOptionsWithCustomHelper 测试 ExecuteOptions 的 WithCustomHelper 方法
func TestExecuteOptionsWithCustomHelper(t *testing.T) {
	opts := NewExecuteOptions()

	// 添加自定义辅助函数
	fn := func(L *lua.LState) int {
		L.Push(lua.LString("test"))
		return 1
	}

	opts = opts.WithCustomHelper("test_func", fn)

	// 验证函数已添加
	if opts.CustomHelpers == nil {
		t.Error("CustomHelpers 不应该为 nil")
	}

	if _, exists := opts.CustomHelpers["test_func"]; !exists {
		t.Error("应该能找到添加的自定义函数")
	}
}

// TestRuleExecutorAddCustomHelper 测试 RuleExecutor 的 AddCustomHelper 方法
func TestRuleExecutorAddCustomHelper(t *testing.T) {
	executor := NewRuleExecutor()

	// 添加自定义辅助函数
	fn := func(L *lua.LState) int {
		L.Push(lua.LString("test"))
		return 1
	}

	err := executor.AddCustomHelper("test_func", fn)
	if err != nil {
		t.Fatalf("添加自定义辅助函数失败: %v", err)
	}

	// 验证函数已添加
	fn2, exists := executor.GetCustomHelper("test_func")
	if !exists {
		t.Error("应该能找到添加的自定义函数")
	}

	if fn2 == nil {
		t.Error("返回的函数不应该为 nil")
	}
}
