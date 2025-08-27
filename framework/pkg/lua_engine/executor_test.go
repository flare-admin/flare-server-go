package lua_engine

import (
	"fmt"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func TestRuleExecutor_CustomHelpers(t *testing.T) {
	executor := NewRuleExecutor()

	// 测试添加自定义辅助函数
	err := executor.AddCustomHelper("test_func", func(L *lua.LState) int {
		L.Push(lua.LString("test_result"))
		return 1
	})
	if err != nil {
		t.Fatalf("添加自定义辅助函数失败: %v", err)
	}

	// 测试添加空名称的函数
	err = executor.AddCustomHelper("", func(L *lua.LState) int {
		return 0
	})
	if err == nil {
		t.Error("应该返回错误，名称为空")
	}

	// 测试添加nil函数
	err = executor.AddCustomHelper("nil_func", nil)
	if err == nil {
		t.Error("应该返回错误，函数为nil")
	}

	// 测试获取自定义辅助函数
	fn, exists := executor.GetCustomHelper("test_func")
	if !exists {
		t.Error("应该能找到自定义辅助函数")
	}
	if fn == nil {
		t.Error("返回的函数不应该为nil")
	}

	// 测试获取不存在的函数
	fn, exists = executor.GetCustomHelper("non_existent")
	if exists {
		t.Error("不应该找到不存在的函数")
	}

	// 测试列出自定义辅助函数
	names := executor.ListCustomHelpers()
	if len(names) != 1 {
		t.Errorf("期望1个函数，实际有%d个", len(names))
	}
	if names[0] != "test_func" {
		t.Errorf("期望函数名为'test_func'，实际为'%s'", names[0])
	}

	// 测试移除自定义辅助函数
	err = executor.RemoveCustomHelper("test_func")
	if err != nil {
		t.Fatalf("移除自定义辅助函数失败: %v", err)
	}

	// 测试移除不存在的函数
	err = executor.RemoveCustomHelper("non_existent")
	if err == nil {
		t.Error("移除不存在的函数应该返回错误")
	}

	// 验证函数已被移除
	names = executor.ListCustomHelpers()
	if len(names) != 0 {
		t.Errorf("期望0个函数，实际有%d个", len(names))
	}
}

func TestRuleExecutor_ExecuteWithCustomHelpers(t *testing.T) {
	executor := NewRuleExecutor()

	// 添加自定义辅助函数
	err := executor.AddCustomHelper("add_numbers", func(L *lua.LState) int {
		a := L.CheckNumber(1)
		b := L.CheckNumber(2)
		L.Push(a + b)
		return 1
	})
	if err != nil {
		t.Fatalf("添加自定义辅助函数失败: %v", err)
	}

	// 添加字符串处理函数
	err = executor.AddCustomHelper("concat_strings", func(L *lua.LState) int {
		str1 := L.CheckString(1)
		str2 := L.CheckString(2)
		L.Push(lua.LString(str1 + str2))
		return 1
	})
	if err != nil {
		t.Fatalf("添加自定义辅助函数失败: %v", err)
	}

	// 测试脚本
	script := `
		local result1 = add_numbers(10, 20)
		local result2 = concat_strings("Hello, ", "World!")
		
		-- 设置结果到上下文
		set_context_value("sum", result1)
		set_context_value("message", result2)
		
		-- 返回验证结果
		return true, "success", {sum = result1, message = result2}
	`

	opts := NewExecuteOptions().WithTimeout(5 * time.Second)
	result, err := executor.Execute(script, opts)
	if err != nil {
		t.Fatalf("执行脚本失败: %v", err)
	}

	if !result.Valid {
		t.Errorf("脚本执行应该成功，但返回了错误: %s", result.Error)
	}

	if result.Action != "success" {
		t.Errorf("期望action为'success'，实际为'%s'", result.Action)
	}
}

func TestRuleExecutor_ValidateScriptWithCustomHelpers(t *testing.T) {
	executor := NewRuleExecutor()

	// 添加自定义辅助函数
	err := executor.AddCustomHelper("test_helper", func(L *lua.LState) int {
		return 0
	})
	if err != nil {
		t.Fatalf("添加自定义辅助函数失败: %v", err)
	}

	// 测试有效的脚本
	validScript := `
		local result = test_helper()
		return true
	`

	err = executor.ValidateScript(validScript)
	if err != nil {
		t.Errorf("脚本应该有效，但验证失败: %v", err)
	}

	// 测试无效的脚本
	invalidScript := `
		local result = test_helper(
		return true
	`

	err = executor.ValidateScript(invalidScript)
	if err == nil {
		t.Error("脚本应该无效，但验证通过了")
	}
}

func TestRuleExecutor_ClearCustomHelpers(t *testing.T) {
	executor := NewRuleExecutor()

	// 添加多个自定义辅助函数
	helpers := []string{"func1", "func2", "func3"}
	for _, name := range helpers {
		err := executor.AddCustomHelper(name, func(L *lua.LState) int {
			return 0
		})
		if err != nil {
			t.Fatalf("添加自定义辅助函数失败: %v", err)
		}
	}

	// 验证函数已添加
	names := executor.ListCustomHelpers()
	if len(names) != len(helpers) {
		t.Errorf("期望%d个函数，实际有%d个", len(helpers), len(names))
	}

	// 清空所有自定义辅助函数
	executor.ClearCustomHelpers()

	// 验证函数已清空
	names = executor.ListCustomHelpers()
	if len(names) != 0 {
		t.Errorf("期望0个函数，实际有%d个", len(names))
	}
}

func TestRuleExecutor_ConcurrentCustomHelpers(t *testing.T) {
	executor := NewRuleExecutor()

	// 并发添加自定义辅助函数
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			name := fmt.Sprintf("func_%d", index)
			err := executor.AddCustomHelper(name, func(L *lua.LState) int {
				L.Push(lua.LNumber(index))
				return 1
			})
			if err != nil {
				t.Errorf("并发添加函数失败: %v", err)
			}
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	// 验证所有函数都已添加
	names := executor.ListCustomHelpers()
	if len(names) != 10 {
		t.Errorf("期望10个函数，实际有%d个", len(names))
	}
}
