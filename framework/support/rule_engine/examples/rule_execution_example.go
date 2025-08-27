package examples

import (
	"fmt"

	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
)

// RuleExecutionExample 规则执行示例
func RuleExecutionExample() {
	// 创建规则执行上下文
	context := model.NewRuleContext(
		"req_123456",        // 请求ID
		"user_register",     // 触发动作
		"user_verification", // 业务类型
	)

	// 设置作用域
	context.SetScope("user", "user_123")

	// 添加执行数据
	context.SetData(map[string]interface{}{
		"user": map[string]interface{}{
			"id":         "user_123",
			"name":       "张三",
			"age":        16,
			"isRealName": false,
			"level":      1,
		},
		"product": map[string]interface{}{
			"id":    "prod_001",
			"name":  "测试商品",
			"price": 100.0,
		},
	})

	// 模拟规则执行服务
	// 在实际使用中，需要通过依赖注入获取服务实例
	// ruleExecutionService := wire.GetRuleExecutionService()

	fmt.Println("=== 规则执行示例 ===")
	fmt.Printf("请求ID: %s\n", context.RequestID)
	fmt.Printf("触发动作: %s\n", context.Trigger)
	fmt.Printf("业务类型: %s\n", context.BusinessType)
	fmt.Printf("作用域: %s\n", context.Scope)

	// 执行单个规则
	// result, err := ruleExecutionService.ExecuteRule(context.Background(), context)
	// if err != nil {
	//     fmt.Printf("执行规则失败: %v\n", err)
	//     return
	// }

	// 处理执行结果
	// fmt.Printf("规则执行结果:\n")
	// fmt.Printf("  规则ID: %s\n", result.RuleID)
	// fmt.Printf("  规则编码: %s\n", result.RuleCode)
	// fmt.Printf("  规则名称: %s\n", result.RuleName)
	// fmt.Printf("  是否成功: %t\n", result.Success)
	// fmt.Printf("  是否有效: %t\n", result.Valid)
	// fmt.Printf("  执行动作: %s\n", result.Action)
	// fmt.Printf("  执行时间: %dms\n", result.ExecuteTime)
	// fmt.Printf("  错误信息: %s\n", result.Error)

	// 检查结果类型
	// if result.IsAllowed() {
	//     fmt.Println("✅ 规则验证通过，允许执行")
	// } else if result.IsDenied() {
	//     fmt.Println("❌ 规则验证失败，拒绝执行")
	// } else if result.IsModified() {
	//     fmt.Println("🔄 规则验证通过，需要修改")
	// } else if result.IsNotified() {
	//     fmt.Println("📢 规则验证通过，需要通知")
	// } else if result.IsRedirected() {
	//     fmt.Println("🔄 规则验证通过，需要重定向")
	// }

	// 获取输出变量
	// if result.Variables != nil {
	//     fmt.Println("输出变量:")
	//     for key, value := range result.Variables {
	//         fmt.Printf("  %s: %v\n", key, value)
	//     }
	// }

	fmt.Println("\n=== 执行多个规则示例 ===")

	// 执行多个规则
	// results, err := ruleExecutionService.ExecuteRules(context.Background(), context)
	// if err != nil {
	//     fmt.Printf("执行规则失败: %v\n", err)
	//     return
	// }

	// 处理多个规则结果
	// fmt.Printf("共执行 %d 个规则:\n", len(results))
	// for i, result := range results {
	//     fmt.Printf("规则 %d:\n", i+1)
	//     fmt.Printf("  规则ID: %s\n", result.RuleID)
	//     fmt.Printf("  规则名称: %s\n", result.RuleName)
	//     fmt.Printf("  是否成功: %t\n", result.Success)
	//     fmt.Printf("  是否有效: %t\n", result.Valid)
	//     fmt.Printf("  执行动作: %s\n", result.Action)
	//     fmt.Printf("  执行时间: %dms\n", result.ExecuteTime)
	//     if result.Error != "" {
	//         fmt.Printf("  错误信息: %s\n", result.Error)
	//     }
	//     fmt.Println()
	// }

	fmt.Println("=== 根据编码执行规则示例 ===")

	// 根据编码执行特定规则
	// result, err := ruleExecutionService.ExecuteRuleByCode(context.Background(), "USER_VERIFICATION", context)
	// if err != nil {
	//     fmt.Printf("执行规则失败: %v\n", err)
	//     return
	// }

	// fmt.Printf("规则执行结果:\n")
	// fmt.Printf("  规则编码: %s\n", result.RuleCode)
	// fmt.Printf("  是否成功: %t\n", result.Success)
	// fmt.Printf("  是否有效: %t\n", result.Valid)
	// fmt.Printf("  执行动作: %s\n", result.Action)

	fmt.Println("示例完成")
}

// UserVerificationExample 用户实名验证示例
func UserVerificationExample() {
	fmt.Println("\n=== 用户实名验证示例 ===")

	// 创建用户验证上下文
	context := model.NewRuleContext(
		"user_verify_001",
		"user_verification",
		"user_management",
	)

	// 设置用户数据
	context.SetData(map[string]interface{}{
		"user": map[string]interface{}{
			"id":         "user_123",
			"name":       "张三",
			"age":        16,
			"isRealName": false,
			"idCard":     "",
			"phone":      "13800138000",
			"level":      1,
			"isTest":     false,
		},
	})

	fmt.Printf("用户信息:\n")
	fmt.Printf("  用户ID: %s\n", context.GetData("user").(map[string]interface{})["id"])
	fmt.Printf("  用户姓名: %s\n", context.GetData("user").(map[string]interface{})["name"])
	fmt.Printf("  用户年龄: %d\n", context.GetData("user").(map[string]interface{})["age"])
	fmt.Printf("  是否实名: %t\n", context.GetData("user").(map[string]interface{})["isRealName"])

	// 模拟执行实名验证规则
	// result, err := ruleExecutionService.ExecuteRule(context.Background(), context)
	// if err != nil {
	//     fmt.Printf("实名验证失败: %v\n", err)
	//     return
	// }

	// 处理验证结果
	// if result.IsAllowed() {
	//     fmt.Println("✅ 实名验证通过")
	// } else {
	//     fmt.Println("❌ 实名验证失败")
	//     fmt.Printf("失败原因: %s\n", result.Error)
	// }
}

// BusinessRuleExample 业务规则示例
func BusinessRuleExample() {
	fmt.Println("\n=== 业务规则示例 ===")

	// 创建订单验证上下文
	context := model.NewRuleContext(
		"order_check_001",
		"order_create",
		"order_management",
	)

	// 设置订单数据
	context.SetData(map[string]interface{}{
		"order": map[string]interface{}{
			"id":     "order_123",
			"amount": 1000.0,
			"items":  []string{"item1", "item2"},
		},
		"user": map[string]interface{}{
			"id":         "user_123",
			"level":      2,
			"isRealName": true,
		},
		"product": map[string]interface{}{
			"id":    "prod_001",
			"price": 100.0,
			"stock": 50,
		},
	})

	fmt.Printf("订单信息:\n")
	fmt.Printf("  订单ID: %s\n", context.GetData("order").(map[string]interface{})["id"])
	fmt.Printf("  订单金额: %.2f\n", context.GetData("order").(map[string]interface{})["amount"])
	fmt.Printf("  用户等级: %d\n", context.GetData("user").(map[string]interface{})["level"])

	// 模拟执行订单验证规则
	// result, err := ruleExecutionService.ExecuteRule(context.Background(), context)
	// if err != nil {
	//     fmt.Printf("订单验证失败: %v\n", err)
	//     return
	// }

	// 处理验证结果
	// if result.IsAllowed() {
	//     fmt.Println("✅ 订单验证通过，允许创建")
	// } else if result.IsDenied() {
	//     fmt.Println("❌ 订单验证失败，拒绝创建")
	//     fmt.Printf("失败原因: %s\n", result.Error)
	// } else if result.IsModified() {
	//     fmt.Println("🔄 订单验证通过，需要修改")
	//     // 处理修改后的数据
	//     if modifiedData, ok := result.GetVariable("modified_data").(map[string]interface{}); ok {
	//         fmt.Println("修改后的数据:")
	//         for key, value := range modifiedData {
	//             fmt.Printf("  %s: %v\n", key, value)
	//         }
	//     }
	// }
}
