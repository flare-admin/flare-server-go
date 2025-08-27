# Lua 规则引擎

这是一个基于 Lua 的规则引擎，支持自定义辅助函数，用于执行复杂的业务规则。

## 功能特性

- 支持自定义辅助函数
- 线程安全的函数管理
- 完整的错误处理
- 高性能的对象池复用
- 支持超时控制
- 支持内存限制
- 支持数据库操作
- 支持复杂SQL和占位符
- 支持QueryBuilder查询构建

## 基本使用

### 创建规则执行器

```go
executor := lua_engine.NewRuleExecutor()
```

### 添加自定义辅助函数

```go
// 添加数学计算函数
err := executor.AddCustomHelper("add", func(L *lua.LState) int {
    a := L.CheckNumber(1)
    b := L.CheckNumber(2)
    L.Push(a + b)
    return 1
})
```

**注意**: 自定义辅助函数必须符合 `lua.LGFunction` 类型，即 `func(L *lua.LState) int` 签名。

// 添加字符串处理函数
err = executor.AddCustomHelper("concat", func(L *lua.LState) int {
    str1 := L.CheckString(1)
    str2 := L.CheckString(2)
    L.Push(lua.LString(str1 + str2))
    return 1
})

// 添加业务逻辑函数
err = executor.AddCustomHelper("check_vip", func(L *lua.LState) int {
    userLevel := L.CheckString(1)
    if userLevel == "vip" {
        L.Push(lua.LBool(true))
    } else {
        L.Push(lua.LBool(false))
    }
    return 1
})
```

### 执行规则脚本

```go
script := `
    -- 使用自定义辅助函数
    local sum = add(10, 20)
    local message = concat("Hello, ", "World!")
    local isVip = check_vip("vip")
    
    -- 设置结果到上下文
    set_context_value("sum", sum)
    set_context_value("message", message)
    set_context_value("isVip", isVip)
    
    -- 返回验证结果
    return true, "success", {sum = sum, message = message, isVip = isVip}
`

opts := lua_engine.NewExecuteOptions().
    WithTimeout(5 * time.Second).
    WithContext(map[string]interface{}{
        "user": map[string]interface{}{
            "id":   123,
            "name": "张三",
            "level": "vip",
        },
    })

result, err := executor.Execute(script, opts)
if err != nil {
    log.Printf("执行失败: %v", err)
    return
}

if result.Valid {
    log.Printf("执行成功: %s", result.Action)
    log.Printf("执行时间: %dms", result.ExecuteTime)
}
```

## 自定义辅助函数管理

### 添加函数

```go
// 添加单个函数
err := executor.AddCustomHelper("my_func", func(L *lua.LState) int {
    // 函数实现
    return 0
})

// 批量添加函数
helpers := map[string]lua_engine.HelperFunction{
    "func1": func(L *lua.LState) int { return 0 },
    "func2": func(L *lua.LState) int { return 0 },
}

for name, fn := range helpers {
    err := executor.AddCustomHelper(name, fn)
    if err != nil {
        log.Printf("添加函数 %s 失败: %v", name, err)
    }
}
```

### 移除函数

```go
// 移除单个函数
err := executor.RemoveCustomHelper("my_func")
if err != nil {
    log.Printf("移除函数失败: %v", err)
}

// 清空所有函数
executor.ClearCustomHelpers()
```

### 查询函数

```go
// 获取函数
fn, exists := executor.GetCustomHelper("my_func")
if exists {
    log.Printf("函数存在")
}

// 列出所有函数
names := executor.ListCustomHelpers()
for _, name := range names {
    log.Printf("函数: %s", name)
}
```

## 高级用法

### 使用 ExecuteOptions 传入临时函数

```go
// 创建执行选项
opts := lua_engine.NewExecuteOptions().
    WithTimeout(5 * time.Second).
    WithCustomHelper("temp_func", func(L *lua.LState) int {
        L.Push(lua.LString("临时函数"))
        return 1
    })

// 执行脚本
result, err := executor.Execute(script, opts)
```

### 并发安全

```go
// 多个 goroutine 可以安全地添加和移除函数
for i := 0; i < 10; i++ {
    go func(index int) {
        name := fmt.Sprintf("func_%d", index)
        err := executor.AddCustomHelper(name, func(L *lua.LState) int {
            L.Push(lua.LNumber(index))
            return 1
        })
        if err != nil {
            log.Printf("添加函数失败: %v", err)
        }
    }(i)
}
```

### 错误处理

```go
// 添加空名称的函数
err := executor.AddCustomHelper("", func(L *lua.LState) int { return 0 })
if err != nil {
    log.Printf("错误: %v", err) // 输出: 辅助函数名称不能为空
}

// 添加 nil 函数
err = executor.AddCustomHelper("nil_func", nil)
if err != nil {
    log.Printf("错误: %v", err) // 输出: 辅助函数不能为空
}

// 移除不存在的函数
err = executor.RemoveCustomHelper("non_existent")
if err != nil {
    log.Printf("错误: %v", err) // 输出: 辅助函数 'non_existent' 不存在
}
```

## 内置辅助函数

除了自定义函数，规则引擎还提供了许多内置的辅助函数：

### 上下文操作
- `set_context_value(key, value)` - 设置上下文值

### JSON 处理
- `json_encode(value)` - 编码为 JSON 字符串
- `json_decode(json_str)` - 解码 JSON 字符串

### 时间处理
- `now()` - 获取当前时间戳
- `format_time(timestamp, layout)` - 格式化时间

### 字符串处理
- `contains(str, substr)` - 检查字符串包含
- `starts_with(str, prefix)` - 检查字符串开头
- `ends_with(str, suffix)` - 检查字符串结尾

### 类型转换
- `to_number(value)` - 转换为数字
- `to_string(value)` - 转换为字符串
- `to_int(value)` - 转换为整数
- `to_float(value)` - 转换为浮点数
- `to_bool(value)` - 转换为布尔值

## 性能优化

### 对象池复用

规则执行器使用对象池来复用 Lua 状态机，避免频繁创建和销毁：

```go
// 自动复用，无需手动管理
executor := lua_engine.NewRuleExecutor()
```

### 内存限制

可以设置最大内存使用量来防止内存泄漏：

```go
opts := lua_engine.NewExecuteOptions().
    WithMaxMemory(10 * 1024 * 1024) // 10MB
```

### 超时控制

设置执行超时来防止脚本无限执行：

```go
opts := lua_engine.NewExecuteOptions().
    WithTimeout(5 * time.Second)
```

## 数据库操作

### 创建带数据库支持的规则执行器

```go
// 创建数据库服务
dbService := lua_engine.NewDBOperationService(data)

// 创建带数据库支持的规则执行器
executor := lua_engine.NewRuleExecutorWithDB(dbService)
```

### 基本数据库操作

```lua
-- 插入数据
local data = {
    name = "张三",
    age = 25,
    email = "zhangsan@example.com"
}
local affected = db_insert("users", data)

-- 更新数据（使用SQL条件）
local updateData = {age = 26, status = "active"}
local affected = db_update("users", updateData, "id = ?", 1)

-- 更新数据（使用复杂条件）
local updateData = {status = "active", updated_at = now()}
local affected = db_update("users", updateData, "age > ? AND status = ?", 18, "pending")

-- 删除数据（使用SQL条件）
local affected = db_delete("users", "id = ?", 1)

-- 删除数据（使用复杂条件）
local affected = db_delete("users", "age < ? AND status = ?", 18, "inactive")

-- 更新数据（使用map条件，向后兼容）
local where = {id = 1}
local updateData = {age = 26}
local affected = db_update_map("users", where, updateData)

-- 删除数据（使用map条件，向后兼容）
local where = {id = 1}
local affected = db_delete_map("users", where)

-- 查询数据
local sql = "SELECT * FROM users WHERE age > ?"
local results = db_query(sql, 18)

-- 查询单条数据
local sql = "SELECT * FROM users WHERE id = ?"
local result = db_query_one(sql, 1)

-- 执行SQL
local sql = "UPDATE users SET status = ? WHERE age > ?"
local affected = db_execute(sql, "active", 18)
```

### 使用QueryBuilder查询

```lua
-- 构建复杂查询
local builder = {
    where = {
        age = {
            operator = ">",
            value = 18
        },
        status = {
            operator = "=",
            value = "active"
        }
    },
    order_by = {
        age = "DESC",
        name = "ASC"
    },
    page = {
        pageNum = 1,
        pageSize = 10
    }
}

-- 执行查询
local results = db_query_builder("users", builder)

-- 统计数量
local count = db_count_builder("users", builder)

-- 构建SQL语句
local sql, args = db_build_sql("users", builder)
```

### 事务操作

```lua
-- 事务操作
local success = db_transaction(function()
    -- 在事务中执行操作
    local data1 = {name = "张三", age = 25}
    local affected1 = db_insert("users", data1)
    
    local data2 = {name = "李四", age = 30}
    local affected2 = db_insert("users", data2)
    
    -- 如果任何操作失败，事务会自动回滚
    if affected1 > 0 and affected2 > 0 then
        return true
    else
        return false
    end
end)
```

### 支持的查询操作符

- `=` - 等于
- `!=` - 不等于
- `>` - 大于
- `>=` - 大于等于
- `<` - 小于
- `<=` - 小于等于
- `LIKE` - 模糊匹配
- `IN` - 包含
- `NOT IN` - 不包含
- `IS NULL` - 为空
- `IS NOT NULL` - 不为空

### 动态参数支持

更新和删除操作支持动态SQL条件和参数：

```lua
-- 复杂条件更新
local updateData = {status = "active"}
local affected = db_update("users", updateData, "age > ? AND status = ? AND created_at > ?", 18, "pending", "2024-01-01")

-- 复杂条件删除
local affected = db_delete("users", "age < ? OR (status = ? AND updated_at < ?)", 18, "inactive", "2024-01-01")

-- 使用LIKE条件
local affected = db_update("users", updateData, "name LIKE ? AND email LIKE ?", "%张%", "%@example.com%")

-- 使用IN条件
local affected = db_delete("users", "id IN (?, ?, ?)", 1, 2, 3)
```

### 复杂SQL示例

```lua
-- 复杂查询
local sql = [[
    SELECT u.id, u.name, u.age, 
           COUNT(o.id) as order_count,
           SUM(o.amount) as total_amount
    FROM users u
    LEFT JOIN orders o ON u.id = o.user_id
    WHERE u.age > ? AND u.status = ?
    GROUP BY u.id, u.name, u.age
    HAVING total_amount > ?
    ORDER BY total_amount DESC
    LIMIT ?, ?
]]

local results = db_query(sql, 18, "active", 1000, 0, 10)
```

## 最佳实践

1. **函数命名**: 使用有意义的函数名，避免与内置函数冲突
2. **错误处理**: 在自定义函数中正确处理错误情况
3. **性能考虑**: 避免在函数中执行耗时操作
4. **线程安全**: 自定义函数应该是线程安全的
5. **内存管理**: 及时清理不再使用的函数
6. **SQL安全**: 使用参数化查询避免SQL注入
7. **事务管理**: 合理使用事务确保数据一致性
8. **查询优化**: 使用QueryBuilder构建高效查询

## 注意事项

1. 自定义函数名称不能为空
2. 自定义函数不能为 nil
3. 函数应该返回正确的参数数量
4. 避免在函数中修改全局状态
5. 注意内存使用，避免内存泄漏 