package lua_engine

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"testing"
	"time"

	lua "github.com/yuin/gopher-lua"
)

// BenchmarkExecute 基准测试：规则执行性能
func BenchmarkExecute(b *testing.B) {
	executor := NewRuleExecutor()

	// 准备测试数据
	context := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   12345,
			"name": "张三",
			"age":  25,
			"tags": []string{"vip", "active", "premium"},
			"profile": map[string]interface{}{
				"email": "zhangsan@example.com",
				"phone": "13800138000",
				"address": map[string]interface{}{
					"city":     "北京",
					"district": "朝阳区",
					"street":   "建国路",
				},
			},
		},
		"order": map[string]interface{}{
			"id":      "ORD-2024-001",
			"amount":  999.99,
			"status":  "pending",
			"items":   []string{"商品A", "商品B", "商品C"},
			"created": utils.GetDateUnix(),
		},
		"settings": map[string]interface{}{
			"max_amount": 10000,
			"min_amount": 100,
			"enabled":    true,
			"features":   []string{"feature1", "feature2", "feature3"},
		},
	}

	script := `
		-- 复杂规则逻辑
		local user = context.user
		local order = context.order
		local settings = context.settings
		
		-- 用户验证
		local userValid = user.age >= 18 and user.age <= 65
		local userActive = false
		for i, tag in ipairs(user.tags) do
			if tag == "active" then
				userActive = true
				break
			end
		end
		
		-- 订单验证
		local orderValid = order.amount >= settings.min_amount and order.amount <= settings.max_amount
		local orderStatusValid = order.status == "pending" or order.status == "confirmed"
		
		-- 地址验证
		local addressValid = user.profile.address.city == "北京"
		
		-- 综合判断
		valid = userValid and userActive and orderValid and orderStatusValid and addressValid
		
		-- 设置动作
		if valid then
			action = "approve"
			action_params = {
				reason = "所有条件满足",
				timestamp = now(),
				user_id = user.id,
				order_id = order.id
			}
		else
			action = "reject"
			action_params = {
				reason = "条件不满足",
				timestamp = now(),
				user_id = user.id,
				order_id = order.id
			}
		end
		
		-- 输出变量
		variables = {
			user_info = {
				id = user.id,
				name = user.name,
				age = user.age,
				tags = user.tags
			},
			order_info = {
				id = order.id,
				amount = order.amount,
				status = order.status,
				items = order.items
			},
			validation_result = {
				user_valid = userValid,
				user_active = userActive,
				order_valid = orderValid,
				order_status_valid = orderStatusValid,
				address_valid = addressValid
			}
		}
	`

	opts := &ExecuteOptions{
		Timeout: 5 * time.Second,
		Context: context,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := executor.Execute(script, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkConvertToLua 基准测试：类型转换性能
func BenchmarkConvertToLua(b *testing.B) {
	executor := NewRuleExecutor()
	L := executor.pool.Get().(*lua.LState)
	defer executor.pool.Put(L)

	// 测试数据
	testData := map[string]interface{}{
		"string":      "test string",
		"int":         12345,
		"float":       123.45,
		"bool":        true,
		"slice":       []int{1, 2, 3, 4, 5},
		"map":         map[string]int{"a": 1, "b": 2, "c": 3},
		"mixed_slice": []interface{}{"a", 1, true, 2.5},
		"mixed_map":   map[string]interface{}{"a": "string", "b": 123, "c": true},
		"nested": map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"level3": []interface{}{1, 2, 3, "string", true},
				},
			},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, value := range testData {
			executor.convertToLuaValue(L, value)
		}
	}
}

// BenchmarkContextInjection 基准测试：上下文注入性能
func BenchmarkContextInjection(b *testing.B) {
	executor := NewRuleExecutor()
	L := executor.pool.Get().(*lua.LState)
	defer executor.pool.Put(L)

	// 大型上下文数据
	context := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		context[fmt.Sprintf("key_%d", i)] = map[string]interface{}{
			"id":    i,
			"name":  fmt.Sprintf("item_%d", i),
			"value": i * 15,
			"tags":  []string{fmt.Sprintf("tag_%d", i), fmt.Sprintf("tag_%d_2", i)},
		}
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		contextTable := L.NewTable()
		executor.injectContextData(L, contextTable, context)
	}
}

// BenchmarkLuaStateReset 基准测试：Lua状态机重置性能
func BenchmarkLuaStateReset(b *testing.B) {
	executor := NewRuleExecutor()
	L := executor.pool.Get().(*lua.LState)
	defer executor.pool.Put(L)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		executor.resetLuaState(L, 10*1024*1024) // 10MB
	}
}

// BenchmarkHelperFunctions 基准测试：辅助函数性能
func BenchmarkHelperFunctions(b *testing.B) {
	executor := NewRuleExecutor()
	L := executor.pool.Get().(*lua.LState)
	defer executor.pool.Put(L)

	executor.registerHelperFunctions(L)

	// 测试脚本
	script := `
		local result = {}
		for i = 1, 1000 do
			-- 测试字符串操作
			local str = "test_string_" .. i
			table.insert(result, fast_contains(str, "test"))
			table.insert(result, fast_starts_with(str, "test"))
			table.insert(result, fast_ends_with(str, tostring(i)))
			
			-- 测试数学操作
			local num1 = i * 1.5
			local num2 = i * 2.5
			table.insert(result, fast_math("add", num1, num2))
			table.insert(result, fast_math("mul", num1, num2))
			table.insert(result, fast_math("pow", num1, 2))
		end
		return result
	`

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := L.DoString(script)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkComplexRule 基准测试：复杂规则性能
func BenchmarkComplexRule(b *testing.B) {
	executor := NewRuleExecutor()

	// 复杂上下文
	context := map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": 1, "name": "Alice", "age": 25, "score": 85.5, "vip": true},
			{"id": 2, "name": "Bob", "age": 30, "score": 92.3, "vip": false},
			{"id": 3, "name": "Charlie", "age": 35, "score": 78.9, "vip": true},
		},
		"products": map[string]interface{}{
			"electronics": []map[string]interface{}{
				{"id": "E001", "name": "Laptop", "price": 2999.99, "stock": 50},
				{"id": "E002", "name": "Phone", "price": 1999.99, "stock": 100},
			},
			"books": []map[string]interface{}{
				{"id": "B001", "name": "Go Programming", "price": 59.99, "stock": 200},
				{"id": "B002", "name": "Lua Programming", "price": 49.99, "stock": 150},
			},
		},
		"rules": map[string]interface{}{
			"min_age":         18,
			"max_age":         65,
			"min_score":       80.0,
			"vip_discount":    0.1,
			"stock_threshold": 10,
		},
	}

	script := `
		local users = context.users
		local products = context.products
		local rules = context.rules
		
		local valid_users = {}
		local valid_products = {}
		local total_value = 0
		
		-- 验证用户
		for i, user in ipairs(users) do
			local user_valid = user.age >= rules.min_age and user.age <= rules.max_age
			if user_valid and user.score >= rules.min_score then
				table.insert(valid_users, user)
			end
		end
		
		-- 验证产品
		for category, items in pairs(products) do
			for j, product in ipairs(items) do
				if product.stock >= rules.stock_threshold then
					table.insert(valid_products, product)
					total_value = total_value + product.price * product.stock
				end
			end
		end
		
		-- 计算VIP折扣
		local vip_count = 0
		for i, user in ipairs(valid_users) do
			if user.vip then
				vip_count = vip_count + 1
			end
		end
		
		local discount = vip_count * rules.vip_discount
		local final_value = total_value * (1 - discount)
		
		-- 设置结果
		valid = #valid_users > 0 and #valid_products > 0
		
		if valid then
			action = "process_order"
			action_params = {
				user_count = #valid_users,
				product_count = #valid_products,
				total_value = total_value,
				final_value = final_value,
				discount = discount,
				vip_count = vip_count
			}
		else
			action = "reject_order"
			action_params = {
				reason = "No valid users or products",
				user_count = #valid_users,
				product_count = #valid_products
			}
		end
		
		variables = {
			valid_users = valid_users,
			valid_products = valid_products,
			summary = {
				total_users = #users,
				valid_user_count = #valid_users,
				total_products = #valid_products,
				valid_product_count = #valid_products,
				total_value = total_value,
				final_value = final_value
			}
		}
	`

	opts := &ExecuteOptions{
		Timeout: 5 * time.Second,
		Context: context,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := executor.Execute(script, opts)
		if err != nil {
			b.Fatal(err)
		}
	}
}
