package lua_engine

import (
	"context"
	"fmt"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	lua "github.com/yuin/gopher-lua"
)

// registerDBHelperFunctions 注册数据库操作辅助函数
func (e *RuleExecutor) registerDBHelperFunctions(L *lua.LState, dbService *DBOperationService) {
	if dbService == nil {
		return
	}

	// 数据库插入操作
	L.SetGlobal("db_insert", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckString(1)
		dataTable := L.CheckTable(2)

		// 将Lua表转换为Go map
		data := make(map[string]interface{})
		dataTable.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				data[key.String()] = luaValueToGo(value)
			}
		})

		// 执行插入操作
		ctx := context.Background()
		affected, err := dbService.Insert(ctx, table, data)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LNumber(affected))
		return 1
	}))

	// 数据库更新操作
	L.SetGlobal("db_update", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckString(1)
		dataTable := L.CheckTable(2)
		whereSQL := L.CheckString(3)
		whereArgs := make([]interface{}, 0)

		// 获取WHERE条件参数
		top := L.GetTop()
		for i := 4; i <= top; i++ {
			whereArgs = append(whereArgs, luaValueToGo(L.Get(i)))
		}

		// 将Lua表转换为Go map
		data := make(map[string]interface{})
		dataTable.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				data[key.String()] = luaValueToGo(value)
			}
		})

		// 执行更新操作
		ctx := context.Background()
		affected, err := dbService.Update(ctx, table, data, whereSQL, whereArgs...)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LNumber(affected))
		return 1
	}))

	// 数据库更新操作（使用map条件，向后兼容）
	L.SetGlobal("db_update_map", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckString(1)
		whereTable := L.CheckTable(2)
		dataTable := L.CheckTable(3)

		// 将Lua表转换为Go map
		where := make(map[string]interface{})
		whereTable.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				where[key.String()] = luaValueToGo(value)
			}
		})

		data := make(map[string]interface{})
		dataTable.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				data[key.String()] = luaValueToGo(value)
			}
		})

		// 执行更新操作
		ctx := context.Background()
		affected, err := dbService.UpdateWithMap(ctx, table, where, data)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LNumber(affected))
		return 1
	}))

	// 数据库删除操作
	L.SetGlobal("db_delete", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckString(1)
		whereSQL := L.CheckString(2)
		whereArgs := make([]interface{}, 0)

		// 获取WHERE条件参数
		top := L.GetTop()
		for i := 3; i <= top; i++ {
			whereArgs = append(whereArgs, luaValueToGo(L.Get(i)))
		}

		// 执行删除操作
		ctx := context.Background()
		affected, err := dbService.Delete(ctx, table, whereSQL, whereArgs...)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LNumber(affected))
		return 1
	}))

	// 数据库删除操作（使用map条件，向后兼容）
	L.SetGlobal("db_delete_map", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckString(1)
		whereTable := L.CheckTable(2)

		// 将Lua表转换为Go map
		where := make(map[string]interface{})
		whereTable.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				where[key.String()] = luaValueToGo(value)
			}
		})

		// 执行删除操作
		ctx := context.Background()
		affected, err := dbService.DeleteWithMap(ctx, table, where)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LNumber(affected))
		return 1
	}))

	// 数据库查询操作
	L.SetGlobal("db_query", L.NewFunction(func(L *lua.LState) int {
		sql := L.CheckString(1)
		args := make([]interface{}, 0)

		// 获取参数
		top := L.GetTop()
		for i := 2; i <= top; i++ {
			args = append(args, luaValueToGo(L.Get(i)))
		}

		// 执行查询操作
		ctx := context.Background()
		results, err := dbService.Query(ctx, sql, args...)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		// 将结果转换为Lua表
		resultTable := L.NewTable()
		for i, row := range results {
			rowTable := L.NewTable()
			for key, value := range row {
				rowTable.RawSetString(key, e.convertToLuaValue(L, value))
			}
			resultTable.RawSetInt(i+1, rowTable)
		}

		L.Push(resultTable)
		return 1
	}))

	// 数据库查询单条数据
	L.SetGlobal("db_query_one", L.NewFunction(func(L *lua.LState) int {
		sql := L.CheckString(1)
		args := make([]interface{}, 0)

		// 获取参数
		top := L.GetTop()
		for i := 2; i <= top; i++ {
			args = append(args, luaValueToGo(L.Get(i)))
		}

		// 执行查询操作
		ctx := context.Background()
		result, err := dbService.QueryOne(ctx, sql, args...)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		// 将结果转换为Lua表
		resultTable := L.NewTable()
		for key, value := range result {
			resultTable.RawSetString(key, e.convertToLuaValue(L, value))
		}

		L.Push(resultTable)
		return 1
	}))

	// 数据库执行操作
	L.SetGlobal("db_execute", L.NewFunction(func(L *lua.LState) int {
		sql := L.CheckString(1)
		args := make([]interface{}, 0)

		// 获取参数
		top := L.GetTop()
		for i := 2; i <= top; i++ {
			args = append(args, luaValueToGo(L.Get(i)))
		}

		// 执行SQL操作
		ctx := context.Background()
		affected, err := dbService.Execute(ctx, sql, args...)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LNumber(affected))
		return 1
	}))

	// 使用QueryBuilder查询
	L.SetGlobal("db_query_builder", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckString(1)
		builderTable := L.CheckTable(2)

		// 构建QueryBuilder
		builder := db_query.NewQueryBuilder()

		// 解析条件
		builderTable.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				keyStr := key.String()
				switch keyStr {
				case "where":
					if value.Type() == lua.LTTable {
						whereTable := value.(*lua.LTable)
						whereTable.ForEach(func(field, condition lua.LValue) {
							if field.Type() == lua.LTString && condition.Type() == lua.LTTable {
								conditionTable := condition.(*lua.LTable)
								operator := conditionTable.RawGetString("operator")
								val := conditionTable.RawGetString("value")

								if operator.Type() == lua.LTString {
									op := db_query.Operator(operator.String())
									builder.Where(field.String(), op, luaValueToGo(val))
								}
							}
						})
					}
				case "order_by":
					if value.Type() == lua.LTTable {
						orderTable := value.(*lua.LTable)
						orderTable.ForEach(func(field, direction lua.LValue) {
							if field.Type() == lua.LTString && direction.Type() == lua.LTString {
								asc := direction.String() == "ASC"
								builder.OrderBy(field.String(), asc)
							}
						})
					}
				case "page":
					if value.Type() == lua.LTTable {
						pageTable := value.(*lua.LTable)
						pageNum := pageTable.RawGetString("pageNum")
						pageSize := pageTable.RawGetString("pageSize")

						if pageNum.Type() == lua.LTNumber && pageSize.Type() == lua.LTNumber {
							page := &db_query.Page{
								Current: int(pageNum.(lua.LNumber)),
								Size:    int(pageSize.(lua.LNumber)),
							}
							builder.WithPage(page)
						}
					}
				}
			}
		})

		// 执行查询
		ctx := context.Background()
		results, err := dbService.QueryWithBuilder(ctx, table, builder)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		// 将结果转换为Lua表
		resultTable := L.NewTable()
		for i, row := range results {
			rowTable := L.NewTable()
			for key, value := range row {
				rowTable.RawSetString(key, e.convertToLuaValue(L, value))
			}
			resultTable.RawSetInt(i+1, rowTable)
		}

		L.Push(resultTable)
		return 1
	}))

	// 使用QueryBuilder统计数量
	L.SetGlobal("db_count_builder", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckString(1)
		builderTable := L.CheckTable(2)

		// 构建QueryBuilder
		builder := db_query.NewQueryBuilder()

		// 解析WHERE条件
		whereTable := builderTable.RawGetString("where")
		if whereTable.Type() == lua.LTTable {
			whereTable.(*lua.LTable).ForEach(func(field, condition lua.LValue) {
				if field.Type() == lua.LTString && condition.Type() == lua.LTTable {
					conditionTable := condition.(*lua.LTable)
					operator := conditionTable.RawGetString("operator")
					val := conditionTable.RawGetString("value")

					if operator.Type() == lua.LTString {
						op := db_query.Operator(operator.String())
						builder.Where(field.String(), op, luaValueToGo(val))
					}
				}
			})
		}

		// 执行统计
		ctx := context.Background()
		count, err := dbService.CountWithBuilder(ctx, table, builder)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LNumber(count))
		return 1
	}))

	// 构建SQL语句
	L.SetGlobal("db_build_sql", L.NewFunction(func(L *lua.LState) int {
		table := L.CheckString(1)
		builderTable := L.CheckTable(2)

		// 构建QueryBuilder
		builder := db_query.NewQueryBuilder()

		// 解析条件
		builderTable.ForEach(func(key, value lua.LValue) {
			if key.Type() == lua.LTString {
				keyStr := key.String()
				switch keyStr {
				case "where":
					if value.Type() == lua.LTTable {
						whereTable := value.(*lua.LTable)
						whereTable.ForEach(func(field, condition lua.LValue) {
							if field.Type() == lua.LTString && condition.Type() == lua.LTTable {
								conditionTable := condition.(*lua.LTable)
								operator := conditionTable.RawGetString("operator")
								val := conditionTable.RawGetString("value")

								if operator.Type() == lua.LTString {
									op := db_query.Operator(operator.String())
									builder.Where(field.String(), op, luaValueToGo(val))
								}
							}
						})
					}
				case "order_by":
					if value.Type() == lua.LTTable {
						orderTable := value.(*lua.LTable)
						orderTable.ForEach(func(field, direction lua.LValue) {
							if field.Type() == lua.LTString && direction.Type() == lua.LTString {
								asc := direction.String() == "ASC"
								builder.OrderBy(field.String(), asc)
							}
						})
					}
				case "page":
					if value.Type() == lua.LTTable {
						pageTable := value.(*lua.LTable)
						pageNum := pageTable.RawGetString("pageNum")
						pageSize := pageTable.RawGetString("pageSize")

						if pageNum.Type() == lua.LTNumber && pageSize.Type() == lua.LTNumber {
							page := &db_query.Page{
								Current: int(pageNum.(lua.LNumber)),
								Size:    int(pageSize.(lua.LNumber)),
							}
							builder.WithPage(page)
						}
					}
				}
			}
		})

		// 构建SQL
		sql, args, err := dbService.BuildSQL(table, builder)
		if err != nil {
			L.Push(lua.LNil)
			L.Push(lua.LString(err.Error()))
			return 2
		}

		// 返回SQL和参数
		L.Push(lua.LString(sql))

		argsTable := L.NewTable()
		for i, arg := range args {
			argsTable.RawSetInt(i+1, e.convertToLuaValue(L, arg))
		}
		L.Push(argsTable)

		return 2
	}))

	// 事务操作
	L.SetGlobal("db_transaction", L.NewFunction(func(L *lua.LState) int {
		fn := L.CheckFunction(1)

		// 执行事务
		ctx := context.Background()
		err := dbService.Transaction(ctx, func(ctx context.Context) error {
			// 调用Lua函数
			L.Push(fn)
			if err := L.PCall(0, 0, nil); err != nil {
				return fmt.Errorf("事务中的Lua函数执行失败: %v", err)
			}
			return nil
		})

		if err != nil {
			L.Push(lua.LBool(false))
			L.Push(lua.LString(err.Error()))
			return 2
		}

		L.Push(lua.LBool(true))
		return 1
	}))
}
