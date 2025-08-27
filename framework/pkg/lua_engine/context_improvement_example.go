package lua_engine

import (
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

// ContextImprovementExample 展示context通用方法的改进效果
func ContextImprovementExample() {
	// 创建规则执行器
	executor := NewRuleExecutor()

	// 创建上下文数据
	context := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   "user_001",
			"name": "张三",
			"age":  25,
		},
		"order": map[string]interface{}{
			"id":     "order_123",
			"amount": 100.50,
			"status": "pending",
		},
	}

	// 创建执行选项
	opts := NewExecuteOptions().WithContext(context)

	// Lua脚本：测试各种context操作
	script := `
		-- 测试1：直接设置context值（使用通用方法）
		set_context_value("test_key", "test_value")
		
		-- 测试2：修改对象属性（使用通用方法）
		set_object_property("user", "age", 30)
		set_object_property("user", "level", 2)
		
		-- 测试3：修改嵌套属性（使用通用方法）
		set_nested_property("order", "status", "completed")
		set_nested_property("order", "processed_at", now())
		
		-- 测试4：获取对象属性（使用通用方法）
		user_age = get_object_property("user", "age")
		order_status = get_object_property("order", "status")
		
		-- 测试5：使用success函数（内部使用通用方法）
		success("update_user", {
			old_age = 25,
			new_age = 30,
			reason = "age_update"
		})
		
		-- 设置执行结果
		valid = true
		action = "context_improvement_test"
	`

	// 执行规则
	result, err := executor.Execute(script, opts)
	if err != nil {
		log.Fatalf("规则执行失败: %v", err)
	}

	// 验证执行结果
	fmt.Println("=== Context通用方法改进效果 ===")
	fmt.Printf("执行成功: %t\n", result.Valid)
	fmt.Printf("执行动作: %s\n", result.Action)
	fmt.Printf("执行时间: %dms\n", result.ExecuteTime)

	// 验证context修改
	fmt.Println("\n=== Context修改验证 ===")
	fmt.Printf("测试键值: %s\n", result.GetContextString("test_key"))
	fmt.Printf("用户年龄: %d\n", result.GetContextInt("user_age"))
	fmt.Printf("订单状态: %s\n", result.GetContextString("order_status"))

	// 验证嵌套对象修改
	user := result.GetContextMap("user")
	if user != nil {
		fmt.Printf("用户年龄(修改后): %v\n", user["age"])
		fmt.Printf("用户等级: %v\n", user["level"])
	}

	order := result.GetContextMap("order")
	if order != nil {
		fmt.Printf("订单状态(修改后): %v\n", order["status"])
		fmt.Printf("处理时间: %v\n", order["processed_at"])
	}

	// 验证success函数设置的变量
	fmt.Println("\n=== Success函数变量验证 ===")
	fmt.Printf("旧年龄: %v\n", result.GetContextValue("old_age"))
	fmt.Printf("新年龄: %v\n", result.GetContextValue("new_age"))
	fmt.Printf("更新原因: %v\n", result.GetContextValue("reason"))
}

// TestContextTableCreation 测试context表创建功能
func TestContextTableCreation() {
	executor := NewRuleExecutor()
	L := lua.NewState()
	defer L.Close()

	// 测试1：当context不存在时，应该创建新表
	contextTable1 := executor.getOrCreateContextTable(L)
	if contextTable1 == nil {
		log.Fatal("context表创建失败")
	}

	// 测试2：再次调用应该返回同一个表
	contextTable2 := executor.getOrCreateContextTable(L)
	if contextTable1 != contextTable2 {
		log.Fatal("应该返回同一个context表")
	}

	// 测试3：设置值后应该能正确获取
	contextTable1.RawSetString("test", lua.LString("value"))
	value := L.GetGlobal("context")
	if value == lua.LNil {
		log.Fatal("全局context表应该存在")
	}

	fmt.Println("✓ Context表创建测试通过")
}

// DemonstrateContextImprovement 演示context改进的效果
func DemonstrateContextImprovement() {
	fmt.Println("=== Context通用方法改进演示 ===")

	// 运行测试
	TestContextTableCreation()

	// 运行示例
	ContextImprovementExample()

	fmt.Println("\n=== 改进总结 ===")
	fmt.Println("1. 统一了context表的获取和创建逻辑")
	fmt.Println("2. 减少了重复代码")
	fmt.Println("3. 提高了代码可维护性")
	fmt.Println("4. 增强了错误处理的健壮性")
	fmt.Println("5. 所有Lua辅助函数现在都使用统一的context处理方法")
}
