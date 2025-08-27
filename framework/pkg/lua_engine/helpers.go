package lua_engine

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/flare-admin/flare-server-go/framework/pkg/utils"

	lua "github.com/yuin/gopher-lua"
)

// getOrCreateContextTable 获取或创建context全局表
// 这是一个通用的方法，用于确保context表存在
//
// 使用示例：
//
//	// 在Lua辅助函数中使用
//	contextTable := e.getOrCreateContextTable(L)
//	contextTable.RawSetString("key", lua.LString("value"))
//
// 功能：
//  1. 检查全局context表是否存在
//  2. 如果不存在，创建新的表并设置为全局变量
//  3. 返回可用的context表
//
// 优势：
//   - 避免重复代码
//   - 统一错误处理
//   - 提高代码可维护性
func (e *RuleExecutor) getOrCreateContextTable(L *lua.LState) *lua.LTable {
	// 获取context全局表
	contextTable := L.GetGlobal("context")
	if contextTable == lua.LNil {
		// 如果contextTable为nil，创建一个新的表
		contextTable = L.NewTable()
		L.SetGlobal("context", contextTable)
	}
	return contextTable.(*lua.LTable)
}

// registerHelperFunctions 注册辅助函数到Lua状态机
func (e *RuleExecutor) registerHelperFunctions(L *lua.LState) {
	// 成功函数
	// success(action, variables)
	L.SetGlobal("success", L.NewFunction(func(L *lua.LState) int {
		action := L.CheckString(1)               // 第一个参数：action 字符串
		variables := L.OptTable(2, L.NewTable()) // 第二个参数：变量 table

		// 使用通用方法获取或创建context表
		contextTable := e.getOrCreateContextTable(L)

		// 将variables设置到context中
		variables.ForEach(func(k, v lua.LValue) {
			contextTable.RawSetString(k.String(), v)
		})

		L.SetGlobal("valid", lua.LBool(true))
		L.SetGlobal("action", lua.LString(action))
		L.SetGlobal("error", lua.LNil)
		L.SetGlobal("error_reason", lua.LNil)
		return 0
	}))
	// 失败函数
	L.SetGlobal("error", L.NewFunction(func(L *lua.LState) int {
		message := L.CheckString(1)
		reason := L.OptString(2, "lua_rule_execution_failed")
		L.SetGlobal("valid", lua.LBool(false))
		L.SetGlobal("error", lua.LString(message))
		L.SetGlobal("error_reason", lua.LString(reason))
		L.SetGlobal("action", lua.LNil)
		return 0
	}))

	// 添加修改上下文参数值的函数
	L.SetGlobal("set_context_value", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(1)
		value := L.Get(2)

		// 使用通用方法获取或创建context表
		contextTable := e.getOrCreateContextTable(L)

		// 设置新值
		L.SetTable(contextTable, lua.LString(key), value)
		L.Push(lua.LBool(true))
		return 1
	}))

	// 添加修改嵌套对象属性的函数
	L.SetGlobal("set_object_property", L.NewFunction(func(L *lua.LState) int {
		objectKey := L.CheckString(1)
		propertyPath := L.CheckString(2)
		value := L.Get(3)

		// 使用通用方法获取或创建context表
		contextTable := e.getOrCreateContextTable(L)

		// 获取对象
		objectValue := L.GetTable(contextTable, lua.LString(objectKey))
		if objectValue == lua.LNil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("object not found: " + objectKey))
			return 2
		}

		// 检查对象是否为表
		objectTable, ok := objectValue.(*lua.LTable)
		if !ok {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("object is not a table: " + objectKey))
			return 2
		}

		// 设置属性值
		L.SetTable(objectTable, lua.LString(propertyPath), value)
		L.Push(lua.LBool(true))
		return 1
	}))

	// 添加获取嵌套对象属性的函数
	L.SetGlobal("get_object_property", L.NewFunction(func(L *lua.LState) int {
		objectKey := L.CheckString(1)
		propertyPath := L.CheckString(2)

		// 使用通用方法获取或创建context表
		contextTable := e.getOrCreateContextTable(L)

		// 获取对象
		objectValue := L.GetTable(contextTable, lua.LString(objectKey))
		if objectValue == lua.LNil {
			L.Push(lua.LNil)
			L.Push(lua.LString("object not found: " + objectKey))
			return 2
		}

		// 检查对象是否为表
		objectTable, ok := objectValue.(*lua.LTable)
		if !ok {
			L.Push(lua.LNil)
			L.Push(lua.LString("object is not a table: " + objectKey))
			return 2
		}

		// 获取属性值
		propertyValue := L.GetTable(objectTable, lua.LString(propertyPath))
		L.Push(propertyValue)
		return 1
	}))

	// 添加深度修改对象属性的函数
	L.SetGlobal("set_nested_property", L.NewFunction(func(L *lua.LState) int {
		objectKey := L.CheckString(1)
		propertyPath := L.CheckString(2)
		value := L.Get(3)

		// 使用通用方法获取或创建context表
		contextTable := e.getOrCreateContextTable(L)

		// 获取对象
		objectValue := L.GetTable(contextTable, lua.LString(objectKey))
		if objectValue == lua.LNil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("object not found: " + objectKey))
			return 2
		}

		// 检查对象是否为表
		objectTable, ok := objectValue.(*lua.LTable)
		if !ok {
			L.Push(lua.LBool(false))
			L.Push(lua.LString("object is not a table: " + objectKey))
			return 2
		}

		// 解析属性路径（支持点分隔符）
		parts := strings.Split(propertyPath, ".")
		currentTable := objectTable

		// 遍历路径，创建或获取嵌套表
		for i, part := range parts {
			if i == len(parts)-1 {
				// 最后一个部分，设置值
				L.SetTable(currentTable, lua.LString(part), value)
			} else {
				// 中间部分，获取或创建嵌套表
				nextValue := L.GetTable(currentTable, lua.LString(part))
				if nextValue == lua.LNil {
					// 创建新的嵌套表
					nextTable := L.NewTable()
					L.SetTable(currentTable, lua.LString(part), nextTable)
					currentTable = nextTable
				} else if nextTable, ok := nextValue.(*lua.LTable); ok {
					currentTable = nextTable
				} else {
					L.Push(lua.LBool(false))
					L.Push(lua.LString("property path is not a table: " + part))
					return 2
				}
			}
		}

		L.Push(lua.LBool(true))
		return 1
	}))

	// 添加调试和验证context变化的辅助函数
	L.SetGlobal("get_context_size", L.NewFunction(func(L *lua.LState) int {
		// 使用通用方法获取或创建context表
		contextTable := e.getOrCreateContextTable(L)

		// 计算context表的大小
		size := 0
		contextTable.ForEach(func(k, v lua.LValue) {
			size++
		})

		L.Push(lua.LNumber(size))
		return 1
	}))

	L.SetGlobal("print_context_keys", L.NewFunction(func(L *lua.LState) int {
		// 使用通用方法获取或创建context表
		contextTable := e.getOrCreateContextTable(L)

		// 打印context中的所有键
		fmt.Println("Context keys:")
		contextTable.ForEach(func(k, v lua.LValue) {
			fmt.Printf("  - %s: %v\n", k.String(), v)
		})

		return 0
	}))

	L.SetGlobal("verify_context_value", L.NewFunction(func(L *lua.LState) int {
		key := L.CheckString(1)

		// 使用通用方法获取或创建context表
		contextTable := e.getOrCreateContextTable(L)

		// 获取指定键的值
		value := L.GetTable(contextTable, lua.LString(key))

		// 返回值和是否存在标志
		L.Push(value)
		L.Push(lua.LBool(value != lua.LNil))
		return 2
	}))

	// JSON 处理函数
	L.SetGlobal("json_encode", L.NewFunction(func(L *lua.LState) int {
		value := luaValueToGo(L.Get(1))
		if jsonStr, err := json.Marshal(value); err == nil {
			L.Push(lua.LString(string(jsonStr)))
		} else {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
		}
		return 2
	}))

	L.SetGlobal("json_decode", L.NewFunction(func(L *lua.LState) int {
		jsonStr := L.CheckString(1)
		var value interface{}
		if err := json.Unmarshal([]byte(jsonStr), &value); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}
		L.Push(goValueToLua(L, value))
		return 1
	}))

	// 时间处理函数
	L.SetGlobal("now", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LNumber(utils.GetDateUnix()))
		return 1
	}))

	L.SetGlobal("format_time", L.NewFunction(func(L *lua.LState) int {
		timestamp := L.CheckNumber(1)
		layout := L.CheckString(2)
		t := time.Unix(int64(timestamp), 0)
		L.Push(lua.LString(t.Format(layout)))
		return 1
	}))

	// 字符串处理函数
	L.SetGlobal("contains", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		substr := L.CheckString(2)
		L.Push(lua.LBool(contains(str, substr)))
		return 1
	}))

	L.SetGlobal("starts_with", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		prefix := L.CheckString(2)
		L.Push(lua.LBool(startsWith(str, prefix)))
		return 1
	}))

	L.SetGlobal("ends_with", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		suffix := L.CheckString(2)
		L.Push(lua.LBool(endsWith(str, suffix)))
		return 1
	}))

	// 数字和字符串转换函数
	L.SetGlobal("to_number", L.NewFunction(func(L *lua.LState) int {
		value := L.Get(1)
		switch v := value.(type) {
		case lua.LNumber:
			L.Push(v)
		case lua.LString:
			if num, err := strconv.ParseFloat(v.String(), 64); err == nil {
				L.Push(lua.LNumber(num))
			} else {
				L.Push(lua.LNil)
				L.Push(lua.LString(fmt.Sprintf("无法将字符串转换为数字: %v", err)))
				return 2
			}
		default:
			L.Push(lua.LNil)
			L.Push(lua.LString("不支持的类型转换"))
			return 2
		}
		return 1
	}))

	L.SetGlobal("to_string", L.NewFunction(func(L *lua.LState) int {
		value := L.Get(1)
		switch v := value.(type) {
		case lua.LNumber:
			L.Push(lua.LString(fmt.Sprintf("%v", v)))
		case lua.LString:
			L.Push(v)
		case lua.LBool:
			L.Push(lua.LString(fmt.Sprintf("%v", v)))
		case *lua.LTable:
			if jsonStr, err := json.Marshal(luaValueToGo(v)); err == nil {
				L.Push(lua.LString(string(jsonStr)))
			} else {
				L.Push(lua.LNil)
				L.Push(lua.LString(fmt.Sprintf("无法将表转换为字符串: %v", err)))
				return 2
			}
		default:
			L.Push(lua.LString(fmt.Sprintf("%v", v)))
		}
		return 1
	}))

	L.SetGlobal("to_int", L.NewFunction(func(L *lua.LState) int {
		value := L.Get(1)
		switch v := value.(type) {
		case lua.LNumber:
			L.Push(lua.LNumber(float64(int64(v))))
		case lua.LString:
			if num, err := strconv.ParseInt(v.String(), 10, 64); err == nil {
				L.Push(lua.LNumber(float64(num)))
			} else {
				L.Push(lua.LNil)
				L.Push(lua.LString(fmt.Sprintf("无法将字符串转换为整数: %v", err)))
				return 2
			}
		default:
			L.Push(lua.LNil)
			L.Push(lua.LString("不支持的类型转换"))
			return 2
		}
		return 1
	}))

	L.SetGlobal("to_float", L.NewFunction(func(L *lua.LState) int {
		value := L.Get(1)
		switch v := value.(type) {
		case lua.LNumber:
			L.Push(v)
		case lua.LString:
			if num, err := strconv.ParseFloat(v.String(), 64); err == nil {
				L.Push(lua.LNumber(num))
			} else {
				L.Push(lua.LNil)
				L.Push(lua.LString(fmt.Sprintf("无法将字符串转换为浮点数: %v", err)))
				return 2
			}
		default:
			L.Push(lua.LNil)
			L.Push(lua.LString("不支持的类型转换"))
			return 2
		}
		return 1
	}))

	L.SetGlobal("to_bool", L.NewFunction(func(L *lua.LState) int {
		value := L.Get(1)
		switch v := value.(type) {
		case lua.LBool:
			L.Push(v)
		case lua.LString:
			switch strings.ToLower(v.String()) {
			case "true", "1", "yes", "on":
				L.Push(lua.LBool(true))
			case "false", "0", "no", "off":
				L.Push(lua.LBool(false))
			default:
				L.Push(lua.LNil)
				L.Push(lua.LString("无法将字符串转换为布尔值"))
				return 2
			}
		case lua.LNumber:
			if v == 0 {
				L.Push(lua.LBool(false))
			} else {
				L.Push(lua.LBool(true))
			}
		default:
			L.Push(lua.LNil)
			L.Push(lua.LString("不支持的类型转换"))
			return 2
		}
		return 1
	}))

	// 新增：解析JSON字符串为Lua表
	L.SetGlobal("parse_json", L.NewFunction(func(L *lua.LState) int {
		jsonStr := L.CheckString(1)
		var value interface{}
		if err := json.Unmarshal([]byte(jsonStr), &value); err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(fmt.Sprintf("JSON解析错误: %v", err)))
			return 2
		}
		L.Push(goValueToLua(L, value))
		return 1
	}))

	// 新增：检查值是否为表
	L.SetGlobal("is_table", L.NewFunction(func(L *lua.LState) int {
		value := L.Get(1)
		L.Push(lua.LBool(value.Type() == lua.LTTable))
		return 1
	}))

	// 新增：获取表的长度
	L.SetGlobal("table_length", L.NewFunction(func(L *lua.LState) int {
		value := L.Get(1)
		if table, ok := value.(*lua.LTable); ok {
			L.Push(lua.LNumber(float64(table.Len())))
		} else {
			L.Push(lua.LNumber(0))
		}
		return 1
	}))

	// 新增：深度复制表
	L.SetGlobal("deep_copy", L.NewFunction(func(L *lua.LState) int {
		value := L.Get(1)
		if table, ok := value.(*lua.LTable); ok {
			newTable := L.NewTable()
			table.ForEach(func(k, v lua.LValue) {
				if subTable, ok := v.(*lua.LTable); ok {
					// 递归复制子表
					L.Push(subTable)
					L.Call(1, 1)
					L.SetTable(newTable, k, L.Get(-1))
					L.Pop(1)
				} else {
					L.SetTable(newTable, k, v)
				}
			})
			L.Push(newTable)
		} else {
			L.Push(value)
		}
		return 1
	}))

	// 新增：高性能字符串操作
	L.SetGlobal("fast_contains", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		substr := L.CheckString(2)
		L.Push(lua.LBool(strings.Contains(str, substr)))
		return 1
	}))

	L.SetGlobal("fast_starts_with", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		prefix := L.CheckString(2)
		L.Push(lua.LBool(strings.HasPrefix(str, prefix)))
		return 1
	}))

	L.SetGlobal("fast_ends_with", L.NewFunction(func(L *lua.LState) int {
		str := L.CheckString(1)
		suffix := L.CheckString(2)
		L.Push(lua.LBool(strings.HasSuffix(str, suffix)))
		return 1
	}))

	// 新增：高性能数字操作
	L.SetGlobal("fast_math", L.NewFunction(func(L *lua.LState) int {
		op := L.CheckString(1)
		a := L.CheckNumber(2)
		b := L.CheckNumber(3)

		var result float64
		switch op {
		case "add":
			result = float64(a) + float64(b)
		case "sub":
			result = float64(a) - float64(b)
		case "mul":
			result = float64(a) * float64(b)
		case "div":
			if b == 0 {
				L.Push(lua.LNil)
				L.Push(lua.LString("除零错误"))
				return 2
			}
			result = float64(a) / float64(b)
		case "mod":
			result = float64(int64(a) % int64(b))
		case "pow":
			result = math.Pow(float64(a), float64(b))
		default:
			L.Push(lua.LNil)
			L.Push(lua.LString("不支持的操作"))
			return 2
		}

		L.Push(lua.LNumber(result))
		return 1
	}))
	// 添加日志函数
	L.SetGlobal("log", L.NewFunction(func(L *lua.LState) int {
		level := L.CheckString(1)
		message := L.CheckString(2)
		switch level {
		case "debug":
			hlog.Debug(message)
		case "warn":
			hlog.Warn(message)
		case "error":
			hlog.Error(message)
		default:
			hlog.Info(message)
		}
		return 0
	}))
}

// luaValueToGo 将Lua值转换为Go值
func luaValueToGo(lv lua.LValue) interface{} {
	if lv == lua.LNil {
		return nil
	}
	switch v := lv.(type) {
	case lua.LBool:
		return bool(v)
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return string(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn > 0 { // table is an array
			arr := make([]interface{}, 0, maxn)
			for i := 1; i <= maxn; i++ {
				arr = append(arr, luaValueToGo(v.RawGetInt(i)))
			}
			return arr
		}
		// table is a map
		m := make(map[string]interface{})
		v.ForEach(func(k, v lua.LValue) {
			m[k.String()] = luaValueToGo(v)
		})
		return m
	default:
		return nil
	}
}

// goValueToLua 将Go值转换为Lua值
func goValueToLua(L *lua.LState, value interface{}) lua.LValue {
	if value == nil {
		return lua.LNil
	}
	switch v := value.(type) {
	case bool:
		return lua.LBool(v)
	case int:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []interface{}:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), goValueToLua(L, item))
		}
		return table
	case map[string]interface{}:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), goValueToLua(L, item))
		}
		return table
	default:
		return lua.LNil
	}
}

// convertToLuaValue 智能转换Go值到Lua值
// 支持所有可转换的类型，不支持的类型返回nil
func (e *RuleExecutor) convertToLuaValue(L *lua.LState, value interface{}) lua.LValue {
	if value == nil {
		return lua.LNil
	}

	// 快速类型检查，避免反射开销
	switch v := value.(type) {
	case bool:
		return lua.LBool(v)
	case int:
		return lua.LNumber(v)
	case int8:
		return lua.LNumber(v)
	case int16:
		return lua.LNumber(v)
	case int32:
		return lua.LNumber(v)
	case int64:
		return lua.LNumber(v)
	case uint:
		return lua.LNumber(v)
	case uint8:
		return lua.LNumber(v)
	case uint16:
		return lua.LNumber(v)
	case uint32:
		return lua.LNumber(v)
	case uint64:
		return lua.LNumber(v)
	case float32:
		return lua.LNumber(v)
	case float64:
		return lua.LNumber(v)
	case string:
		return lua.LString(v)
	case []interface{}:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), e.convertToLuaValue(L, item))
		}
		return table
	case []string:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LString(item))
		}
		return table
	case []int:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []int8:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []int16:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []int32:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []int64:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []uint:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []uint8:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []uint16:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []uint32:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []uint64:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []float32:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []float64:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LNumber(item))
		}
		return table
	case []bool:
		table := L.NewTable()
		for i, item := range v {
			L.RawSet(table, lua.LNumber(i+1), lua.LBool(item))
		}
		return table
	case map[string]interface{}:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), e.convertToLuaValue(L, item))
		}
		return table
	case map[string]string:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LString(item))
		}
		return table
	case map[string]int:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]int8:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]int16:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]int32:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]int64:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]uint:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]uint8:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]uint16:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]uint32:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]uint64:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]float32:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]float64:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LNumber(item))
		}
		return table
	case map[string]bool:
		table := L.NewTable()
		for k, item := range v {
			L.RawSet(table, lua.LString(k), lua.LBool(item))
		}
		return table
	default:
		// 处理指针类型
		if reflect.TypeOf(value).Kind() == reflect.Ptr {
			if reflect.ValueOf(value).IsNil() {
				return lua.LNil
			}
			// 解引用指针
			value = reflect.ValueOf(value).Elem().Interface()
			return e.convertToLuaValue(L, value)
		}

		// 使用反射处理结构体
		return e.convertStructToLua(L, value)
	}
}

// convertStructToLua 将结构体转换为Lua表
func (e *RuleExecutor) convertStructToLua(L *lua.LState, value interface{}) lua.LValue {
	val := reflect.ValueOf(value)
	typ := val.Type()

	// 如果不是结构体，返回nil
	if val.Kind() != reflect.Struct {
		return lua.LNil
	}

	table := L.NewTable()

	// 遍历结构体字段
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 获取字段名（优先使用json标签）
		fieldName := fieldType.Name
		if jsonTag := fieldType.Tag.Get("json"); jsonTag != "" {
			// 处理json标签，取第一个逗号前的部分
			if commaIndex := strings.Index(jsonTag, ","); commaIndex != -1 {
				fieldName = jsonTag[:commaIndex]
			} else {
				fieldName = jsonTag
			}
			// 如果标签是"-"，则跳过该字段
			if fieldName == "-" {
				continue
			}
		}

		// 如果字段名为空，跳过
		if fieldName == "" {
			continue
		}

		// 转换字段值
		fieldValue := e.convertToLuaValue(L, field.Interface())
		if fieldValue != lua.LNil {
			L.RawSet(table, lua.LString(fieldName), fieldValue)
		}
	}

	return table
}

// resetLuaState 高效重置Lua状态机
func (e *RuleExecutor) resetLuaState(L *lua.LState, maxMemory uint64) {
	// 清空全局表而不是重新创建状态机
	L.SetGlobal("context", lua.LNil)
	L.SetGlobal("valid", lua.LNil)
	L.SetGlobal("action", lua.LNil)
	L.SetGlobal("action_params", lua.LNil)
	L.SetGlobal("variables", lua.LNil)

	// 设置内存限制
	if maxMemory > 0 {
		L.SetMx(int(maxMemory))
	}
}

// injectContextData 高效注入上下文数据
func (e *RuleExecutor) injectContextData(L *lua.LState, contextTable *lua.LTable, context map[string]interface{}) {
	// 预分配容量以提高性能
	contextTable.RawSetString("__size", lua.LNumber(len(context)))

	for k, v := range context {
		// 直接转换为Lua值，不支持的类型将被忽略
		luaValue := e.convertToLuaValue(L, v)
		if luaValue != lua.LNil {
			L.RawSet(contextTable, lua.LString(k), luaValue)
		}
	}
}

// extractResults 高效提取执行结果
func (e *RuleExecutor) extractResults(L *lua.LState, result *ExecuteResult) {
	// 获取必需字段
	valid := L.GetGlobal("valid")
	if valid == lua.LNil || valid.Type() != lua.LTBool {
		result.Error = "规则必须返回布尔类型的valid变量"
		result.ErrorReason = "rule.missing_valid_field"
		return
	}
	result.Valid = bool(valid.(lua.LBool))

	// 获取动作
	if action := L.GetGlobal("action"); action != lua.LNil {
		result.Action = action.String()
	}
	// 获取错误信息（如果存在）
	if errorMsg := L.GetGlobal("error"); errorMsg != lua.LNil {
		result.Error = errorMsg.String()
	}

	// 获取错误原因（国际化键）
	if errorReason := L.GetGlobal("error_reason"); errorReason != lua.LNil {
		result.ErrorReason = errorReason.String()
	}
}

// extractModifiedContext 提取修改后的上下文
func (e *RuleExecutor) extractModifiedContext(L *lua.LState, originalContext map[string]interface{}) map[string]interface{} {
	// 使用通用方法获取context表
	contextTable := e.getOrCreateContextTable(L)

	// 将Lua表转换为Go map
	modifiedContext := luaValueToGo(contextTable)
	if modifiedMap, ok := modifiedContext.(map[string]interface{}); ok {
		return modifiedMap
	}

	return originalContext
}

// 字符串辅助函数
func contains(s, substr string) bool {
	for i := 0; i < len(s); i++ {
		if startsWith(s[i:], substr) {
			return true
		}
	}
	return false
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[0:len(prefix)] == prefix
}

func endsWith(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}

// testGetOrCreateContextTable 测试getOrCreateContextTable方法
// 这是一个内部测试函数，用于验证通用方法的功能
func (e *RuleExecutor) testGetOrCreateContextTable(L *lua.LState) *lua.LTable {
	return e.getOrCreateContextTable(L)
}
