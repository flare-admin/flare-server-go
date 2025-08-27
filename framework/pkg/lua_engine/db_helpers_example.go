package lua_engine

import (
	"fmt"
	"log"

	lua "github.com/yuin/gopher-lua"
)

// ExampleDBQueryWithConvertToLuaValue 展示使用convertToLuaValue方法的数据库查询示例
func ExampleDBQueryWithConvertToLuaValue() {
	// 创建Lua状态机
	L := lua.NewState()
	defer L.Close()

	// 创建数据库服务（这里只是示例，实际使用时需要真实的数据库连接）
	dbService := &DBOperationService{}

	// 创建规则执行器
	executor := &RuleExecutor{}
	executor.registerDBHelperFunctions(L, dbService)

	// 示例1: 使用db_query执行COUNT查询
	// SQL: SELECT count(*) FROM users WHERE from_uid = '711083979805036544' and is_real = 1
	script := `
		-- 执行COUNT查询
		local result = db_query("SELECT count(*) as total FROM users WHERE from_uid = ? and is_real = ?", "711083979805036544", 1)
		
		if result then
			-- 获取第一行结果
			local row = result[1]
			if row then
				-- 使用convertToLuaValue转换后的结果
				local count = row.total
				print("用户数量: " .. tostring(count))
				
				-- 检查类型
				if type(count) == "number" then
					print("类型: number")
				else
					print("类型: " .. type(count))
				end
				
				-- 返回结果
				success("查询成功", {
					count = count,
					user_id = "711083979805036544",
					is_real = 1
				})
			else
				error("未找到结果")
			end
		else
			error("查询失败")
		end
	`

	// 执行脚本
	if err := L.DoString(script); err != nil {
		log.Printf("执行脚本失败: %v", err)
		return
	}

	// 获取执行结果
	valid := L.GetGlobal("valid")
	action := L.GetGlobal("action")
	errorMsg := L.GetGlobal("error")

	fmt.Printf("执行结果:\n")
	fmt.Printf("  valid: %v\n", valid)
	fmt.Printf("  action: %v\n", action)
	if errorMsg != lua.LNil {
		fmt.Printf("  error: %v\n", errorMsg)
	}

	// 示例2: 使用db_query_one查询单条记录
	script2 := `
		-- 查询单个用户信息
		local user = db_query_one("SELECT id, username, email, created_at FROM users WHERE id = ?", "711083979805036544")
		
		if user then
			print("用户信息:")
			print("  ID: " .. tostring(user.id))
			print("  用户名: " .. tostring(user.username))
			print("  邮箱: " .. tostring(user.email))
			print("  创建时间: " .. tostring(user.created_at))
			
			-- 检查时间戳类型转换
			if type(user.created_at) == "number" then
				print("  创建时间类型: number (时间戳)")
			end
			
			success("查询成功", {
				user = user
			})
		else
			error("用户不存在")
		end
	`

	fmt.Printf("\n=== 示例2: 查询单条记录 ===\n")
	if err := L.DoString(script2); err != nil {
		log.Printf("执行脚本失败: %v", err)
		return
	}

	// 示例3: 使用db_query_builder进行复杂查询
	script3 := `
		-- 使用QueryBuilder进行分页查询
		local builder = {
			where = {
				from_uid = {
					operator = "eq",
					value = "711083979805036544"
				},
				is_real = {
					operator = "eq", 
					value = 1
				}
			},
			order_by = {
				created_at = "DESC"
			},
			page = {
				pageNum = 1,
				pageSize = 10
			}
		}
		
		local results = db_query_builder("users", builder)
		
		if results then
			print("查询结果数量: " .. tostring(#results))
			
			for i, user in ipairs(results) do
				print("用户 " .. i .. ":")
				print("  ID: " .. tostring(user.id))
				print("  用户名: " .. tostring(user.username))
				print("  创建时间: " .. tostring(user.created_at))
				print("  是否实名: " .. tostring(user.is_real))
			end
			
			success("查询成功", {
				users = results,
				total = #results
			})
		else
			error("查询失败")
		end
	`

	fmt.Printf("\n=== 示例3: 使用QueryBuilder查询 ===\n")
	if err := L.DoString(script3); err != nil {
		log.Printf("执行脚本失败: %v", err)
		return
	}
}

// ExampleCountQueryResult 展示COUNT查询的返回结果示例
func ExampleCountQueryResult() {
	fmt.Printf("\n=== COUNT查询返回结果示例 ===\n")
	fmt.Printf("SQL: SELECT count(*) FROM users WHERE from_uid = '711083979805036544' and is_real = 1\n")
	fmt.Printf("\n使用convertToLuaValue转换后的Lua表结构:\n")
	fmt.Printf("{\n")
	fmt.Printf("  [1] = {\n")
	fmt.Printf("    total = 42  -- 数字类型，由convertToLuaValue自动转换\n")
	fmt.Printf("  }\n")
	fmt.Printf("}\n")
	fmt.Printf("\n在Lua中的使用方式:\n")
	fmt.Printf("local result = db_query(\"SELECT count(*) as total FROM users WHERE from_uid = ? and is_real = ?\", \"711083979805036544\", 1)\n")
	fmt.Printf("local count = result[1].total  -- 获取数量\n")
	fmt.Printf("print(\"用户数量: \" .. tostring(count))  -- 输出: 用户数量: 42\n")
	fmt.Printf("\n类型转换说明:\n")
	fmt.Printf("- 数据库返回的count(*)是int64类型\n")
	fmt.Printf("- convertToLuaValue将其转换为lua.LNumber\n")
	fmt.Printf("- 在Lua中表现为number类型\n")
	fmt.Printf("- 可以直接进行数学运算\n")
}

// ExampleTypeConversion 展示各种数据类型的转换示例
func ExampleTypeConversion() {
	fmt.Printf("\n=== 数据类型转换示例 ===\n")
	fmt.Printf("convertToLuaValue支持的Go类型转换:\n")
	fmt.Printf("\n1. 基本类型:\n")
	fmt.Printf("   - bool -> lua.LBool\n")
	fmt.Printf("   - int/int8/int16/int32/int64 -> lua.LNumber\n")
	fmt.Printf("   - uint/uint8/uint16/uint32/uint64 -> lua.LNumber\n")
	fmt.Printf("   - float32/float64 -> lua.LNumber\n")
	fmt.Printf("   - string -> lua.LString\n")
	fmt.Printf("\n2. 切片类型:\n")
	fmt.Printf("   - []int -> lua.LTable (数字索引)\n")
	fmt.Printf("   - []string -> lua.LTable (数字索引)\n")
	fmt.Printf("   - []bool -> lua.LTable (数字索引)\n")
	fmt.Printf("   - []interface{} -> lua.LTable (数字索引)\n")
	fmt.Printf("\n3. 映射类型:\n")
	fmt.Printf("   - map[string]int -> lua.LTable (字符串键)\n")
	fmt.Printf("   - map[string]string -> lua.LTable (字符串键)\n")
	fmt.Printf("   - map[string]bool -> lua.LTable (字符串键)\n")
	fmt.Printf("   - map[string]interface{} -> lua.LTable (字符串键)\n")
	fmt.Printf("\n4. 指针类型:\n")
	fmt.Printf("   - 自动解引用指针\n")
	fmt.Printf("   - nil指针转换为lua.LNil\n")
	fmt.Printf("\n5. 结构体类型:\n")
	fmt.Printf("   - 使用反射转换为lua.LTable\n")
	fmt.Printf("   - 支持json标签\n")
	fmt.Printf("   - 忽略私有字段\n")
}
