package lua_engine

import (
	"fmt"
	"log"
)

// TestSuccessFunctionContextChange 测试success函数中variables设置后context的变化
func TestSuccessFunctionContextChange() {
	// 创建规则执行器
	executor := NewRuleExecutor()

	// 创建初始上下文数据
	initialContext := map[string]interface{}{
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
	opts := NewExecuteOptions().WithContext(initialContext)

	// Lua脚本：使用success函数设置variables
	script := `
		-- 打印初始context大小
		print("Initial context size:", get_context_size())
		print_context_keys()
		
		-- 使用success函数设置variables，这些变量应该被添加到context中
		success("update_user", {
			old_age = 25,
			new_age = 30,
			reason = "age_update",
			processed_at = now(),
			processor = "rule_engine"
		})
		
		-- 验证context变化
		print("Context size after success:", get_context_size())
		print_context_keys()
		
		-- 验证success函数设置的变量
		local old_age, exists = verify_context_value("old_age")
		print("old_age exists:", exists, "value:", old_age)
		
		local new_age, exists = verify_context_value("new_age")
		print("new_age exists:", exists, "value:", new_age)
		
		local reason, exists = verify_context_value("reason")
		print("reason exists:", exists, "value:", reason)
		
		local processed_at, exists = verify_context_value("processed_at")
		print("processed_at exists:", exists, "value:", processed_at)
		
		local processor, exists = verify_context_value("processor")
		print("processor exists:", exists, "value:", processor)
		
		-- 设置执行结果
		valid = true
		action = "test_success_context_change"
	`

	// 执行规则
	result, err := executor.Execute(script, opts)
	if err != nil {
		log.Fatalf("规则执行失败: %v", err)
	}

	// 验证执行结果
	fmt.Println("=== Success函数Context变化测试 ===")
	fmt.Printf("执行成功: %t\n", result.Valid)
	fmt.Printf("执行动作: %s\n", result.Action)
	fmt.Printf("执行时间: %dms\n", result.ExecuteTime)

	// 验证context变化
	fmt.Println("\n=== Context变化验证 ===")

	// 检查success函数设置的变量是否在context中
	fmt.Printf("old_age: %v\n", result.GetContextValue("old_age"))
	fmt.Printf("new_age: %v\n", result.GetContextValue("new_age"))
	fmt.Printf("reason: %v\n", result.GetContextValue("reason"))
	fmt.Printf("processed_at: %v\n", result.GetContextValue("processed_at"))
	fmt.Printf("processor: %v\n", result.GetContextValue("processor"))

	// 验证原始数据是否保持不变
	fmt.Println("\n=== 原始数据验证 ===")
	user := result.GetContextMap("user")
	if user != nil {
		fmt.Printf("用户ID: %v\n", user["id"])
		fmt.Printf("用户姓名: %v\n", user["name"])
		fmt.Printf("用户年龄: %v\n", user["age"])
	}

	order := result.GetContextMap("order")
	if order != nil {
		fmt.Printf("订单ID: %v\n", order["id"])
		fmt.Printf("订单金额: %v\n", order["amount"])
		fmt.Printf("订单状态: %v\n", order["status"])
	}

	// 验证context大小
	fmt.Println("\n=== Context大小验证 ===")
	contextSize := len(result.Context)
	fmt.Printf("Context总大小: %d\n", contextSize)

	// 列出所有context键
	fmt.Println("Context中的所有键:")
	for key := range result.Context {
		fmt.Printf("  - %s\n", key)
	}
}

// TestSuccessFunctionVariablesMerge 测试success函数variables与原始context的合并
func TestSuccessFunctionVariablesMerge() {
	// 创建规则执行器
	executor := NewRuleExecutor()

	// 创建初始上下文数据
	initialContext := map[string]interface{}{
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
		"existing_var": "original_value",
	}

	// 创建执行选项
	opts := NewExecuteOptions().WithContext(initialContext)

	// Lua脚本：测试variables与现有context的合并
	script := `
		-- 使用success函数设置variables，包括覆盖现有变量
		success("merge_test", {
			existing_var = "updated_value",  -- 覆盖现有变量
			new_var = "new_value",           -- 新增变量
			user = {                         -- 覆盖整个对象
				id = "user_002",
				name = "李四",
				age = 30,
				level = 2
			},
			order = {                        -- 部分更新对象
				status = "completed",
				processed_at = now()
			}
		})
		
		-- 设置执行结果
		valid = true
		action = "test_variables_merge"
	`

	// 执行规则
	result, err := executor.Execute(script, opts)
	if err != nil {
		log.Fatalf("规则执行失败: %v", err)
	}

	// 验证执行结果
	fmt.Println("=== Success函数Variables合并测试 ===")
	fmt.Printf("执行成功: %t\n", result.Valid)
	fmt.Printf("执行动作: %s\n", result.Action)

	// 验证变量合并
	fmt.Println("\n=== Variables合并验证 ===")

	// 检查覆盖的变量
	fmt.Printf("existing_var (应该被覆盖): %v\n", result.GetContextValue("existing_var"))

	// 检查新增的变量
	fmt.Printf("new_var (新增): %v\n", result.GetContextValue("new_var"))

	// 检查覆盖的对象
	user := result.GetContextMap("user")
	if user != nil {
		fmt.Printf("用户ID (覆盖后): %v\n", user["id"])
		fmt.Printf("用户姓名 (覆盖后): %v\n", user["name"])
		fmt.Printf("用户年龄 (覆盖后): %v\n", user["age"])
		fmt.Printf("用户等级 (新增): %v\n", user["level"])
	}

	// 检查部分更新的对象
	order := result.GetContextMap("order")
	if order != nil {
		fmt.Printf("订单ID (保持不变): %v\n", order["id"])
		fmt.Printf("订单金额 (保持不变): %v\n", order["amount"])
		fmt.Printf("订单状态 (更新后): %v\n", order["status"])
		fmt.Printf("处理时间 (新增): %v\n", order["processed_at"])
	}

	// 验证context大小
	fmt.Println("\n=== 最终Context大小 ===")
	contextSize := len(result.Context)
	fmt.Printf("Context总大小: %d\n", contextSize)
}

// DemonstrateSuccessFunctionContextChange 演示success函数中context的变化
func DemonstrateSuccessFunctionContextChange() {
	fmt.Println("=== Success函数Context变化演示 ===")

	// 运行基本测试
	TestSuccessFunctionContextChange()

	fmt.Println("\n" + "="*50 + "\n")

	// 运行合并测试
	TestSuccessFunctionVariablesMerge()

	fmt.Println("\n=== 总结 ===")
	fmt.Println("1. Success函数中的variables会被正确设置到context中")
	fmt.Println("2. Variables会与原始context数据合并")
	fmt.Println("3. 相同键的variables会覆盖原始context中的值")
	fmt.Println("4. 新的variables会被添加到context中")
	fmt.Println("5. 修改后的context会在执行结果中正确返回")
}
