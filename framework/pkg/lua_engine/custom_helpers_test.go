package lua_engine

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestCustomHelpers(t *testing.T) {
	executor := NewRuleExecutor()

	// 创建执行选项并添加自定义辅助函数
	opts := NewExecuteOptions().
		WithTimeout(5*1000000000). // 5秒
		WithContext(map[string]interface{}{
			"user_id": 12345,
			"amount":  100.50,
		}).
		WithCustomHelper("test_sum", func(L *lua.LState) int {
			// 获取参数数量
			n := L.GetTop()
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
		WithCustomHelper("test_multiply", func(L *lua.LState) int {
			a := L.CheckNumber(1)
			b := L.CheckNumber(2)
			result := a * b
			L.Push(lua.LNumber(result))
			return 1
		})

	// 定义测试脚本
	script := `
		local user_id = context.user_id
		local amount = context.amount
		
		-- 测试自定义求和函数
		local sum = test_sum(amount, 50, 25.5)
		
		-- 测试自定义乘法函数
		local product = test_multiply(amount, 2)
		
		-- 设置验证结果
		valid = true
		action = "approve"
		variables = {
			user_id = user_id,
			original_amount = amount,
			sum_result = sum,
			product_result = product
		}
	`

	// 执行规则
	result, err := executor.Execute(script, opts)
	if err != nil {
		t.Fatalf("执行规则失败: %v", err)
	}

	// 验证结果
	if !result.Valid {
		t.Errorf("期望验证通过，但实际为 false")
	}

	if result.Action != "approve" {
		t.Errorf("期望动作为 'approve'，但实际为 '%s'", result.Action)
	}

	// 验证输出变量
	variables, ok := result.Variables.(map[string]interface{})
	if !ok {
		t.Fatalf("输出变量类型错误")
	}

	// 验证用户ID
	if userID, ok := variables["user_id"].(float64); !ok || userID != 12345 {
		t.Errorf("用户ID验证失败")
	}

	// 验证原始金额
	if originalAmount, ok := variables["original_amount"].(float64); !ok || originalAmount != 100.50 {
		t.Errorf("原始金额验证失败")
	}

	// 验证求和结果 (100.50 + 50 + 25.5 = 176.0)
	if sumResult, ok := variables["sum_result"].(float64); !ok || sumResult != 176.0 {
		t.Errorf("求和结果验证失败，期望 176.0，实际 %f", sumResult)
	}

	// 验证乘法结果 (100.50 * 2 = 201.0)
	if productResult, ok := variables["product_result"].(float64); !ok || productResult != 201.0 {
		t.Errorf("乘法结果验证失败，期望 201.0，实际 %f", productResult)
	}
}

func TestCustomHelpersWithError(t *testing.T) {
	executor := NewRuleExecutor()

	// 创建执行选项并添加带错误处理的自定义函数
	opts := NewExecuteOptions().
		WithCustomHelper("safe_divide", func(L *lua.LState) int {
			a := L.CheckNumber(1)
			b := L.CheckNumber(2)

			if b == 0 {
				L.Push(lua.LNil)
				L.Push(lua.LString("除零错误"))
				return 2
			}

			result := a / b
			L.Push(lua.LNumber(result))
			return 1
		})

	// 测试正常除法
	script1 := `
		local result = safe_divide(10, 2)
		valid = true
		action = "approve"
		variables = { result = result }
	`

	result1, err := executor.Execute(script1, opts)
	if err != nil {
		t.Fatalf("执行规则失败: %v", err)
	}

	if !result1.Valid {
		t.Errorf("期望验证通过，但实际为 false")
	}

	// 测试除零错误
	script2 := `
		local result, error = safe_divide(10, 0)
		valid = false
		action = "reject"
		error = error
	`

	result2, err := executor.Execute(script2, opts)
	if err != nil {
		t.Fatalf("执行规则失败: %v", err)
	}

	if result2.Valid {
		t.Errorf("期望验证失败，但实际为 true")
	}

	if result2.Action != "reject" {
		t.Errorf("期望动作为 'reject'，但实际为 '%s'", result2.Action)
	}

	if result2.Error != "除零错误" {
		t.Errorf("期望错误为 '除零错误'，但实际为 '%s'", result2.Error)
	}
}

func TestCustomHelpersMultipleFunctions(t *testing.T) {
	executor := NewRuleExecutor()

	// 创建多个自定义辅助函数
	opts := NewExecuteOptions().
		WithContext(map[string]interface{}{
			"name": "张三",
			"age":  25,
		}).
		WithCustomHelper("format_name", func(L *lua.LState) int {
			name := L.CheckString(1)
			formatted := "尊敬的 " + name + " 先生"
			L.Push(lua.LString(formatted))
			return 1
		}).
		WithCustomHelper("check_age", func(L *lua.LState) int {
			age := L.CheckNumber(1)
			isAdult := age >= 18
			L.Push(lua.LBool(isAdult))
			return 1
		}).
		WithCustomHelper("calculate_bmi", func(L *lua.LState) int {
			weight := L.CheckNumber(1)
			height := L.CheckNumber(2)
			bmi := weight / (height * height)
			L.Push(lua.LNumber(bmi))
			return 1
		})

	// 定义测试脚本
	script := `
		local name = context.name
		local age = context.age
		
		-- 使用多个自定义函数
		local formatted_name = format_name(name)
		local is_adult = check_age(age)
		local bmi = calculate_bmi(70, 1.75)
		
		-- 设置验证结果
		valid = is_adult
		action = is_adult and "approve" or "reject"
		variables = {
			formatted_name = formatted_name,
			is_adult = is_adult,
			bmi = bmi
		}
	`

	// 执行规则
	result, err := executor.Execute(script, opts)
	if err != nil {
		t.Fatalf("执行规则失败: %v", err)
	}

	// 验证结果
	if !result.Valid {
		t.Errorf("期望验证通过，但实际为 false")
	}

	if result.Action != "approve" {
		t.Errorf("期望动作为 'approve'，但实际为 '%s'", result.Action)
	}

	// 验证输出变量
	variables, ok := result.Variables.(map[string]interface{})
	if !ok {
		t.Fatalf("输出变量类型错误")
	}

	// 验证格式化姓名
	if formattedName, ok := variables["formatted_name"].(string); !ok || formattedName != "尊敬的 张三 先生" {
		t.Errorf("格式化姓名验证失败")
	}

	// 验证成年状态
	if isAdult, ok := variables["is_adult"].(bool); !ok || !isAdult {
		t.Errorf("成年状态验证失败")
	}

	// 验证BMI (70 / (1.75 * 1.75) ≈ 22.86)
	if bmi, ok := variables["bmi"].(float64); !ok || bmi < 22.8 || bmi > 22.9 {
		t.Errorf("BMI验证失败，期望约22.86，实际 %f", bmi)
	}
}
