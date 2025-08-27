package examples

import (
	"fmt"
	"log"

	"github.com/flare-admin/flare-server-go/framework/pkg/lua_engine"
	lua "github.com/yuin/gopher-lua"
)

// CustomHelpersExample 自定义辅助函数示例
func CustomHelpersExample() {
	// 创建规则执行器
	executor := lua_engine.NewRuleExecutor()

	// 创建执行选项并添加自定义辅助函数
	opts := lua_engine.NewExecuteOptions().
		WithTimeout(10*1000000000).  // 10秒
		WithMaxMemory(20*1024*1024). // 20MB
		WithContext(map[string]interface{}{
			"user_id": 12345,
			"amount":  100.50,
		}).
		WithCustomHelper("custom_sum", func(L *lua.LState) int {
			// 获取参数数量
			n := L.GetTop()
			if n < 2 {
				L.Push(lua.LNumber(0))
				return 1
			}

			// 计算所有数字参数的和
			sum := 0.0
			for i := 1; i <= n; i++ {
				value := L.Get(i)
				if num, ok := value.(lua.LNumber); ok {
					sum += float64(num)
				}
			}

			L.Push(lua.LNumber(sum))
			return 1
		}).
		WithCustomHelper("custom_format_currency", func(L *lua.LState) int {
			// 获取金额参数
			amount := L.CheckNumber(1)
			currency := L.CheckString(2)

			// 格式化货币
			formatted := fmt.Sprintf("%s %.2f", currency, amount)
			L.Push(lua.LString(formatted))
			return 1
		}).
		WithCustomHelper("custom_check_permission", func(L *lua.LState) int {
			// 获取用户ID和权限
			userID := L.CheckNumber(1)
			permission := L.CheckString(2)

			// 模拟权限检查逻辑
			hasPermission := false
			if userID == 12345 && permission == "admin" {
				hasPermission = true
			}

			L.Push(lua.LBool(hasPermission))
			return 1
		}).
		WithCustomHelper("custom_log", func(L *lua.LState) int {
			// 获取日志消息
			message := L.CheckString(1)
			level := L.OptString(2, "INFO")

			// 输出日志
			log.Printf("[%s] %s", level, message)
			return 0
		})

	// 定义Lua脚本，使用自定义辅助函数
	script := `
		-- 使用自定义辅助函数
		custom_log("开始执行规则验证", "DEBUG")
		
		-- 获取上下文数据
		local user_id = context.user_id
		local amount = context.amount
		
		-- 使用自定义权限检查
		local has_admin = custom_check_permission(user_id, "admin")
		if has_admin then
			custom_log("用户具有管理员权限", "INFO")
		end
		
		-- 使用自定义求和函数
		local total = custom_sum(amount, 50, 25.5)
		custom_log("计算总金额: " .. to_string(total), "INFO")
		
		-- 使用自定义货币格式化
		local formatted = custom_format_currency(total, "CNY")
		custom_log("格式化金额: " .. formatted, "INFO")
		
		-- 设置验证结果
		valid = true
		action = "approve"
		variables = {
			total_amount = total,
			formatted_amount = formatted,
			has_admin = has_admin
		}
		
		custom_log("规则执行完成", "DEBUG")
	`

	// 执行规则
	result, err := executor.Execute(script, opts)
	if err != nil {
		log.Printf("执行规则失败: %v", err)
		return
	}

	// 输出结果
	fmt.Printf("执行结果:\n")
	fmt.Printf("  验证通过: %t\n", result.Valid)
	fmt.Printf("  动作: %s\n", result.Action)
	fmt.Printf("  执行时间: %dms\n", result.ExecuteTime)
	fmt.Printf("  输出变量: %+v\n", result.Variables)
}

// BusinessLogicExample 业务逻辑示例
func BusinessLogicExample() {
	executor := lua_engine.NewRuleExecutor()

	// 添加业务特定的辅助函数
	opts := lua_engine.NewExecuteOptions().
		WithContext(map[string]interface{}{
			"order_amount": 1500.00,
			"user_level":   "vip",
			"region":       "CN",
		}).
		WithCustomHelper("calculate_discount", func(L *lua.LState) int {
			amount := L.CheckNumber(1)
			userLevel := L.CheckString(2)

			discount := 0.0
			switch userLevel {
			case "vip":
				discount = 0.15
			case "premium":
				discount = 0.10
			case "regular":
				discount = 0.05
			}

			finalAmount := amount * (1 - discount)
			L.Push(lua.LNumber(finalAmount))
			return 1
		}).
		WithCustomHelper("check_region_restriction", func(L *lua.LState) int {
			region := L.CheckString(1)
			restrictedRegions := []string{"US", "EU"}

			for _, restricted := range restrictedRegions {
				if region == restricted {
					L.Push(lua.LBool(true))
					return 1
				}
			}

			L.Push(lua.LBool(false))
			return 1
		})

	// 业务规则脚本
	script := `
		local amount = context.order_amount
		local user_level = context.user_level
		local region = context.region
		
		-- 检查地区限制
		local is_restricted = check_region_restriction(region)
		if is_restricted then
			valid = false
			action = "reject"
			error = "该地区暂不支持此服务"
			return
		end
		
		-- 计算折扣后金额
		local final_amount = calculate_discount(amount, user_level)
		
		-- 业务规则验证
		if final_amount > 1000 then
			valid = true
			action = "approve"
			variables = {
				original_amount = amount,
				final_amount = final_amount,
				discount_rate = (amount - final_amount) / amount,
				user_level = user_level
			}
		else
			valid = false
			action = "reject"
			error = "订单金额过低"
		end
	`

	result, err := executor.Execute(script, opts)
	if err != nil {
		log.Printf("执行业务规则失败: %v", err)
		return
	}

	fmt.Printf("业务规则执行结果:\n")
	fmt.Printf("  验证通过: %t\n", result.Valid)
	fmt.Printf("  动作: %s\n", result.Action)
	if result.Error != "" {
		fmt.Printf("  错误: %s\n", result.Error)
	}
	fmt.Printf("  输出变量: %+v\n", result.Variables)
}
