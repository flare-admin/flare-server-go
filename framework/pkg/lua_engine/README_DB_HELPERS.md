# 数据库辅助函数类型转换优化

## 修改概述

在 `registerDBHelperFunctions` 函数中，将所有的 `goValueToLua` 方法调用替换为 `convertToLuaValue` 方法，以提供更好的类型转换支持。

## 修改内容

### 1. 修改的函数

- `db_query` - 数据库查询操作
- `db_query_one` - 数据库查询单条数据
- `db_query_builder` - 使用QueryBuilder查询
- `db_build_sql` - 构建SQL语句

### 2. 修改位置

```go
// 修改前
rowTable.RawSetString(key, goValueToLua(L, value))

// 修改后
rowTable.RawSetString(key, e.convertToLuaValue(L, value))
```

## convertToLuaValue 方法优势

### 1. 更全面的类型支持

`convertToLuaValue` 方法支持更多的Go类型转换：

- **基本类型**: `bool`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `string`
- **切片类型**: `[]int`, `[]string`, `[]bool`, `[]interface{}` 等
- **映射类型**: `map[string]int`, `map[string]string`, `map[string]bool`, `map[string]interface{}` 等
- **指针类型**: 自动解引用
- **结构体类型**: 使用反射转换，支持json标签

### 2. 更好的性能

- 使用快速类型检查，避免反射开销
- 预分配容量以提高性能
- 支持批量转换

### 3. 更安全的转换

- 处理nil指针
- 忽略不支持的类型
- 提供错误处理

## 使用示例

### 1. COUNT查询示例

```lua
-- SQL: SELECT count(*) FROM users WHERE from_uid = '711083979805036544' and is_real = 1
local result = db_query("SELECT count(*) as total FROM users WHERE from_uid = ? and is_real = ?", "711083979805036544", 1)

if result then
    local row = result[1]
    if row then
        local count = row.total  -- 自动转换为number类型
        print("用户数量: " .. tostring(count))
        
        -- 可以直接进行数学运算
        local doubleCount = count * 2
        print("双倍数量: " .. tostring(doubleCount))
        
        success("查询成功", {
            count = count,
            user_id = "711083979805036544",
            is_real = 1
        })
    end
end
```

### 2. 用户信息查询示例

```lua
-- 查询单个用户信息
local user = db_query_one("SELECT id, username, email, created_at, is_real, balance, active FROM users WHERE id = ?", "711083979805036544")

if user then
    print("用户信息:")
    print("  ID: " .. tostring(user.id))           -- string类型
    print("  用户名: " .. tostring(user.username)) -- string类型
    print("  邮箱: " .. tostring(user.email))      -- string类型
    print("  创建时间: " .. tostring(user.created_at)) -- number类型(时间戳)
    print("  是否实名: " .. tostring(user.is_real))    -- number类型
    print("  余额: " .. tostring(user.balance))        -- number类型
    print("  是否激活: " .. tostring(user.active))     -- boolean类型
    
    -- 时间戳转换示例
    if type(user.created_at) == "number" then
        local timestamp = user.created_at
        local date = os.date("%Y-%m-%d %H:%M:%S", timestamp)
        print("  格式化时间: " .. date)
    end
    
    success("查询成功", {
        user = user
    })
end
```

### 3. 复杂查询示例

```lua
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
        
        -- 类型检查
        if type(user.is_real) == "number" then
            if user.is_real == 1 then
                print("  状态: 已实名")
            else
                print("  状态: 未实名")
            end
        end
    end
    
    success("查询成功", {
        users = results,
        total = #results
    })
end
```

## 返回结果示例

### COUNT查询返回结果

```lua
-- 数据库返回: {"total": 42}
-- convertToLuaValue转换后的Lua表结构:
{
  [1] = {
    total = 42  -- number类型，由convertToLuaValue自动转换
  }
}

-- 在Lua中的使用:
local count = result[1].total
print("用户数量: " .. tostring(count))  -- 输出: 用户数量: 42
```

### 用户查询返回结果

```lua
-- 数据库返回: {
--   "id": "711083979805036544",
--   "username": "testuser", 
--   "email": "test@example.com",
--   "created_at": 1640995200,
--   "is_real": 1,
--   "balance": 100.50,
--   "active": true
-- }

-- convertToLuaValue转换后的Lua表结构:
{
  id = "711083979805036544",      -- string类型
  username = "testuser",          -- string类型
  email = "test@example.com",     -- string类型
  created_at = 1640995200,        -- number类型(时间戳)
  is_real = 1,                    -- number类型
  balance = 100.5,                -- number类型
  active = true                   -- boolean类型
}
```

## 类型转换说明

### 1. 数字类型转换

- `int`, `int8`, `int16`, `int32`, `int64` → `lua.LNumber`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64` → `lua.LNumber`
- `float32`, `float64` → `lua.LNumber`

### 2. 字符串类型转换

- `string` → `lua.LString`

### 3. 布尔类型转换

- `bool` → `lua.LBool`

### 4. 时间戳处理

- 数据库中的时间戳字段（如 `created_at`）会被转换为 `lua.LNumber`
- 在Lua中可以直接使用 `os.date()` 函数格式化

### 5. 空值处理

- `nil` → `lua.LNil`
- 空指针 → `lua.LNil`

## 性能优化

### 1. 快速类型检查

`convertToLuaValue` 使用类型断言进行快速检查，避免反射开销：

```go
switch v := value.(type) {
case bool:
    return lua.LBool(v)
case int:
    return lua.LNumber(v)
case string:
    return lua.LString(v)
// ... 更多类型
}
```

### 2. 预分配容量

对于切片和映射类型，预分配容量以提高性能：

```go
case []string:
    table := L.NewTable()
    for i, item := range v {
        L.RawSet(table, lua.LNumber(i+1), lua.LString(item))
    }
    return table
```

### 3. 批量转换

支持批量转换多个值，减少函数调用开销。

## 测试验证

运行测试以验证修改的正确性：

```bash
go test ./framework/pkg/lua_engine -v -run TestConvertToLuaValueInDBQuery
```

测试包括：
- 基本类型转换测试
- COUNT查询结果测试
- 用户查询结果测试
- 查询结果数组测试
- 性能基准测试

## 注意事项

1. **向后兼容性**: 修改保持了API的向后兼容性
2. **性能提升**: 新的转换方法提供了更好的性能
3. **类型安全**: 更全面的类型支持确保了数据转换的安全性
4. **错误处理**: 不支持的类型会被忽略，不会导致程序崩溃 