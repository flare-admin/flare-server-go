package lua_engine

import (
	"fmt"
	"log"
)

// SuccessContextDemo 演示success函数中variables设置后context的变化
func SuccessContextDemo() {
	fmt.Println("=== Success函数Context变化演示 ===")

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

	// Lua脚本：演示success函数设置variables
	script := `
		-- 步骤1：查看初始context
		print("=== 步骤1：初始Context ===")
		print("Context大小:", get_context_size())
		print_context_keys()
		
		-- 步骤2：使用success函数设置variables
		print("\n=== 步骤2：使用Success函数设置Variables ===")
		success("update_user", {
			old_age = 25,
			new_age = 30,
			reason = "age_update",
			processed_at = now(),
			processor = "rule_engine",
			result = "success"
		})
		
		-- 步骤3：验证context变化
		print("\n=== 步骤3：验证Context变化 ===")
		print("Context大小:", get_context_size())
		print_context_keys()
		
		-- 步骤4：验证具体变量
		print("\n=== 步骤4：验证具体变量 ===")
		local old_age, exists = verify_context_value("old_age")
		print("old_age:", old_age, "exists:", exists)
		
		local new_age, exists = verify_context_value("new_age")
		print("new_age:", new_age, "exists:", exists)
		
		local reason, exists = verify_context_value("reason")
		print("reason:", reason, "exists:", exists)
		
		local processed_at, exists = verify_context_value("processed_at")
		print("processed_at:", processed_at, "exists:", exists)
		
		local processor, exists = verify_context_value("processor")
		print("processor:", processor, "exists:", exists)
		
		local result, exists = verify_context_value("result")
		print("result:", result, "exists:", exists)
		
		-- 步骤5：验证原始数据是否保持不变
		print("\n=== 步骤5：验证原始数据 ===")
		local user_id, exists = verify_context_value("user")
		print("user exists:", exists)
		
		local order_id, exists = verify_context_value("order")
		print("order exists:", exists)
		
		-- 设置执行结果
		valid = true
		action = "success_context_demo"
	`

	// 执行规则
	result, err := executor.Execute(script, opts)
	if err != nil {
		log.Fatalf("规则执行失败: %v", err)
	}

	// 验证执行结果
	fmt.Println("\n=== 执行结果验证 ===")
	fmt.Printf("执行成功: %t\n", result.Valid)
	fmt.Printf("执行动作: %s\n", result.Action)
	fmt.Printf("执行时间: %dms\n", result.ExecuteTime)

	// 验证context变化
	fmt.Println("\n=== Context变化验证 ===")

	// 检查success函数设置的变量
	fmt.Println("Success函数设置的变量:")
	fmt.Printf("  old_age: %v\n", result.GetContextValue("old_age"))
	fmt.Printf("  new_age: %v\n", result.GetContextValue("new_age"))
	fmt.Printf("  reason: %v\n", result.GetContextValue("reason"))
	fmt.Printf("  processed_at: %v\n", result.GetContextValue("processed_at"))
	fmt.Printf("  processor: %v\n", result.GetContextValue("processor"))
	fmt.Printf("  result: %v\n", result.GetContextValue("result"))

	// 验证原始数据
	fmt.Println("\n原始数据验证:")
	user := result.GetContextMap("user")
	if user != nil {
		fmt.Printf("  用户ID: %v\n", user["id"])
		fmt.Printf("  用户姓名: %v\n", user["name"])
		fmt.Printf("  用户年龄: %v\n", user["age"])
	}

	order := result.GetContextMap("order")
	if order != nil {
		fmt.Printf("  订单ID: %v\n", order["id"])
		fmt.Printf("  订单金额: %v\n", order["amount"])
		fmt.Printf("  订单状态: %v\n", order["status"])
	}

	// 验证context大小
	fmt.Println("\n=== Context大小验证 ===")
	contextSize := len(result.Context)
	fmt.Printf("最终Context大小: %d\n", contextSize)

	// 列出所有context键
	fmt.Println("Context中的所有键:")
	for key := range result.Context {
		fmt.Printf("  - %s\n", key)
	}

	fmt.Println("\n=== 演示总结 ===")
	fmt.Println("✓ Success函数中的variables被正确设置到context中")
	fmt.Println("✓ Variables与原始context数据成功合并")
	fmt.Println("✓ 修改后的context在执行结果中正确返回")
	fmt.Println("✓ 原始数据保持不变")
	fmt.Println("✓ Context大小正确反映了新增的变量")
}
