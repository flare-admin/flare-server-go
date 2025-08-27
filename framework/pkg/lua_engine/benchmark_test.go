package lua_engine

import (
	"fmt"
	"testing"
	"time"
)

// 定义测试结构体
type BenchmarkUser struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Age      int      `json:"age"`
	IsActive bool     `json:"isActive"`
	Tags     []string `json:"tags"`
	Settings struct {
		Theme    string                 `json:"theme"`
		Language string                 `json:"language"`
		Config   map[string]interface{} `json:"config"`
	} `json:"settings"`
}

type BenchmarkOrder struct {
	ID       string  `json:"id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Items    []struct {
		Name     string   `json:"name"`
		Price    float64  `json:"price"`
		Quantity int      `json:"quantity"`
		Tags     []string `json:"tags"`
	} `json:"items"`
}

func BenchmarkRuleExecutor_OptimizedConversion(b *testing.B) {
	executor := NewRuleExecutor()

	// 创建复杂的测试数据
	user := &BenchmarkUser{
		ID:       123,
		Name:     "张三",
		Age:      25,
		IsActive: true,
		Tags:     []string{"vip", "premium", "verified"},
	}
	user.Settings.Theme = "dark"
	user.Settings.Language = "zh-CN"
	user.Settings.Config = map[string]interface{}{
		"notifications": true,
		"autoSave":      false,
		"theme":         "dark",
	}

	order := &BenchmarkOrder{
		ID:       "ORD-001",
		Amount:   99.99,
		Currency: "CNY",
		Items: []struct {
			Name     string   `json:"name"`
			Price    float64  `json:"price"`
			Quantity int      `json:"quantity"`
			Tags     []string `json:"tags"`
		}{
			{Name: "商品1", Price: 29.99, Quantity: 2, Tags: []string{"hot", "new"}},
			{Name: "商品2", Price: 39.99, Quantity: 1, Tags: []string{"sale"}},
			{Name: "商品3", Price: 49.99, Quantity: 3, Tags: []string{"premium"}},
		},
	}

	context := map[string]interface{}{
		"user":          user,
		"order":         order,
		"simple_string": "hello world",
		"simple_number": 42,
		"simple_bool":   true,
		"numbers":       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		"strings":       []string{"a", "b", "c", "d", "e"},
		"nested": map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]int{
					"level3": 42,
				},
			},
		},
	}

	script := `
		-- 验证基本类型
		if context.simple_string ~= "hello world" then
			valid = false
			return
		end
		
		if context.simple_number ~= 42 then
			valid = false
			return
		end
		
		if context.simple_bool ~= true then
			valid = false
			return
		end
		
		-- 验证结构体字段访问
		if context.user.name ~= "张三" then
			valid = false
			return
		end
		
		if context.user.age ~= 25 then
			valid = false
			return
		end
		
		if context.user.isActive ~= true then
			valid = false
			return
		end
		
		-- 验证结构体中的切片
		if context.user.tags[1] ~= "vip" then
			valid = false
			return
		end
		
		if context.user.tags[2] ~= "premium" then
			valid = false
			return
		end
		
		-- 验证嵌套结构体
		if context.user.settings.theme ~= "dark" then
			valid = false
			return
		end
		
		if context.user.settings.language ~= "zh-CN" then
			valid = false
			return
		end
		
		-- 验证订单结构体
		if context.order.id ~= "ORD-001" then
			valid = false
			return
		end
		
		if context.order.amount ~= 99.99 then
			valid = false
			return
		end
		
		if context.order.currency ~= "CNY" then
			valid = false
			return
		end
		
		-- 验证订单中的切片结构体
		if context.order.items[1].name ~= "商品1" then
			valid = false
			return
		end
		
		if context.order.items[1].price ~= 29.99 then
			valid = false
			return
		end
		
		if context.order.items[1].quantity ~= 2 then
			valid = false
			return
		end
		
		-- 验证嵌套映射
		if context.nested.level1.level2.level3 ~= 42 then
			valid = false
			return
		end
		
		-- 验证切片
		if context.numbers[1] ~= 1 then
			valid = false
			return
		end
		
		if context.numbers[10] ~= 10 then
			valid = false
			return
		end
		
		-- 所有验证通过
		valid = true
		action = "benchmark_test"
		variables = {
			user_name = context.user.name,
			user_age = context.user.age,
			order_amount = context.order.amount,
			tags_count = table_length(context.user.tags),
			items_count = table_length(context.order.items),
			numbers_count = table_length(context.numbers)
		}
	`

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := executor.Execute(script, &ExecuteOptions{
			Context: context,
			Timeout: 5 * time.Second,
		})

		if err != nil {
			b.Fatalf("规则执行失败: %v", err)
		}

		if !result.Valid {
			b.Fatalf("规则验证失败")
		}
	}
}

func BenchmarkRuleExecutor_ComplexDataProcessing(b *testing.B) {
	executor := NewRuleExecutor()

	// 创建大量复杂数据
	largeData := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		user := &BenchmarkUser{
			ID:       i,
			Name:     "用户" + string(rune('A'+i%26)),
			Age:      20 + i%50,
			IsActive: i%2 == 0,
			Tags:     []string{"tag1", "tag2", "tag3"},
		}
		user.Settings.Theme = "theme" + string(rune('A'+i%5))
		user.Settings.Language = "lang" + string(rune('A'+i%3))

		order := &BenchmarkOrder{
			ID:       "ORD-" + string(rune('A'+i)),
			Amount:   float64(i) * 10.5,
			Currency: "CNY",
			Items: []struct {
				Name     string   `json:"name"`
				Price    float64  `json:"price"`
				Quantity int      `json:"quantity"`
				Tags     []string `json:"tags"`
			}{
				{Name: "商品" + string(rune('A'+i)), Price: float64(i) * 5.5, Quantity: i%5 + 1, Tags: []string{"hot"}},
			},
		}

		largeData[fmt.Sprintf("user_%d", i)] = user
		largeData[fmt.Sprintf("order_%d", i)] = order
	}

	script := `
		-- 简单的验证逻辑
		valid = true
		action = "process_large_data"
		action_params = {
			count = table_length(context)
		}
		variables = {
			message = "大数据处理测试完成",
			total_count = table_length(context)
		}
	`

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		result, err := executor.Execute(script, &ExecuteOptions{
			Context: largeData,
			Timeout: 10 * time.Second,
		})

		if err != nil {
			b.Fatalf("规则执行失败: %v", err)
		}

		if !result.Valid {
			b.Fatalf("规则验证失败")
		}
	}
}
