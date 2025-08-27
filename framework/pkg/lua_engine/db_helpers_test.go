package lua_engine

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// TestConvertToLuaValueInDBQuery 测试在数据库查询中使用convertToLuaValue方法
func TestConvertToLuaValueInDBQuery(t *testing.T) {
	// 创建Lua状态机
	L := lua.NewState()
	defer L.Close()

	// 创建规则执行器
	executor := &RuleExecutor{}

	// 测试各种数据类型的转换
	testCases := []struct {
		name     string
		input    interface{}
		expected lua.LValueType
	}{
		{
			name:     "int64 to number",
			input:    int64(42),
			expected: lua.LTNumber,
		},
		{
			name:     "string to string",
			input:    "test string",
			expected: lua.LTString,
		},
		{
			name:     "bool to bool",
			input:    true,
			expected: lua.LTBool,
		},
		{
			name:     "float64 to number",
			input:    float64(3.14),
			expected: lua.LTNumber,
		},
		{
			name:     "nil to nil",
			input:    nil,
			expected: lua.LTNil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := executor.convertToLuaValue(L, tc.input)
			if result.Type() != tc.expected {
				t.Errorf("期望类型 %v, 实际类型 %v", tc.expected, result.Type())
			}
		})
	}
}

// TestCountQueryResult 测试COUNT查询结果的转换
func TestCountQueryResult(t *testing.T) {
	// 创建Lua状态机
	L := lua.NewState()
	defer L.Close()

	// 创建规则执行器
	executor := &RuleExecutor{}

	// 模拟COUNT查询结果
	countResult := map[string]interface{}{
		"total": int64(42),
	}

	// 转换为Lua表
	resultTable := L.NewTable()
	for key, value := range countResult {
		resultTable.RawSetString(key, executor.convertToLuaValue(L, value))
	}

	// 验证转换结果
	totalValue := resultTable.RawGetString("total")
	if totalValue.Type() != lua.LTNumber {
		t.Errorf("期望total字段为number类型，实际为 %v", totalValue.Type())
	}

	if totalValue.(lua.LNumber) != 42 {
		t.Errorf("期望total值为42，实际为 %v", totalValue)
	}
}

// TestUserQueryResult 测试用户查询结果的转换
func TestUserQueryResult(t *testing.T) {
	// 创建Lua状态机
	L := lua.NewState()
	defer L.Close()

	// 创建规则执行器
	executor := &RuleExecutor{}

	// 模拟用户查询结果
	userResult := map[string]interface{}{
		"id":         "711083979805036544",
		"username":   "testuser",
		"email":      "test@example.com",
		"created_at": int64(1640995200), // 2022-01-01 00:00:00
		"is_real":    int8(1),
		"balance":    float64(100.50),
		"active":     true,
	}

	// 转换为Lua表
	resultTable := L.NewTable()
	for key, value := range userResult {
		resultTable.RawSetString(key, executor.convertToLuaValue(L, value))
	}

	// 验证各个字段的类型转换
	testCases := []struct {
		field    string
		expected lua.LValueType
	}{
		{"id", lua.LTString},
		{"username", lua.LTString},
		{"email", lua.LTString},
		{"created_at", lua.LTNumber},
		{"is_real", lua.LTNumber},
		{"balance", lua.LTNumber},
		{"active", lua.LTBool},
	}

	for _, tc := range testCases {
		t.Run(tc.field, func(t *testing.T) {
			value := resultTable.RawGetString(tc.field)
			if value.Type() != tc.expected {
				t.Errorf("字段 %s 期望类型 %v, 实际类型 %v", tc.field, tc.expected, value.Type())
			}
		})
	}
}

// TestQueryResultsArray 测试查询结果数组的转换
func TestQueryResultsArray(t *testing.T) {
	// 创建Lua状态机
	L := lua.NewState()
	defer L.Close()

	// 创建规则执行器
	executor := &RuleExecutor{}

	// 模拟查询结果数组
	results := []map[string]interface{}{
		{
			"id":       "1",
			"username": "user1",
			"count":    int64(10),
		},
		{
			"id":       "2",
			"username": "user2",
			"count":    int64(20),
		},
	}

	// 转换为Lua表
	resultTable := L.NewTable()
	for i, row := range results {
		rowTable := L.NewTable()
		for key, value := range row {
			rowTable.RawSetString(key, executor.convertToLuaValue(L, value))
		}
		resultTable.RawSetInt(i+1, rowTable)
	}

	// 验证数组长度
	if resultTable.Len() != 2 {
		t.Errorf("期望数组长度为2，实际为 %d", resultTable.Len())
	}

	// 验证第一行数据
	firstRow := resultTable.RawGetInt(1).(*lua.LTable)
	if firstRow.RawGetString("id").String() != "1" {
		t.Errorf("期望第一行id为'1'，实际为 %s", firstRow.RawGetString("id").String())
	}

	// 验证第二行数据
	secondRow := resultTable.RawGetInt(2).(*lua.LTable)
	if secondRow.RawGetString("count").(lua.LNumber) != 20 {
		t.Errorf("期望第二行count为20，实际为 %v", secondRow.RawGetString("count"))
	}
}

// BenchmarkConvertToLuaValue 基准测试convertToLuaValue方法
func BenchmarkConvertToLuaValue(b *testing.B) {
	L := lua.NewState()
	defer L.Close()

	executor := &RuleExecutor{}

	// 测试数据
	testData := map[string]interface{}{
		"id":         "711083979805036544",
		"username":   "testuser",
		"email":      "test@example.com",
		"created_at": int64(1640995200),
		"is_real":    int8(1),
		"balance":    float64(100.50),
		"active":     true,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, value := range testData {
			executor.convertToLuaValue(L, value)
		}
	}
}
