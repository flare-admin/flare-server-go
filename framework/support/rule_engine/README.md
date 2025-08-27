# 规则引擎 (Rule Engine)

## 概述

规则引擎是一个基于DDD和CQRS架构的灵活规则执行系统，支持多种规则类型和SQL操作。

## 功能特性

- **多种规则类型**: 支持条件规则、Lua脚本规则、公式规则
- **SQL执行**: 在Lua脚本中直接执行SQL操作
- **CQRS模式**: 命令和查询职责分离
- **DDD架构**: 领域驱动设计
- **灵活配置**: 支持动态规则配置

## 规则类型

### 1. 条件规则 (Condition)

基于条件判断的规则，支持多种操作符：

```json
{
  "type": "condition",
  "conditions": [
    {
      "field": "user.age",
      "operator": "gte",
      "value": 18
    },
    {
      "field": "order.amount",
      "operator": "between",
      "value": "100,1000"
    }
  ],
  "action": "allow"
}
```

### 2. Lua脚本规则 (Lua)

支持复杂业务逻辑的Lua脚本规则，包含SQL执行功能：

```lua
-- 获取上下文数据
local user_id = context.user_id
local amount = context.amount

-- 查询用户余额
local user_sql = "SELECT balance FROM users WHERE id = ?"
local user_result, error_msg = sql_query_one(user_sql, user_id)

if user_result and user_result.balance >= amount then
    -- 更新余额
    local update_where = {id = user_id}
    local update_data = {balance = user_result.balance - amount}
    local affected_rows, update_error = sql_update("users", update_where, update_data)
    
    if update_error == "" then
        set_context_value("balance_after", user_result.balance - amount)
        return true
    end
end

return false
```

### 3. 公式规则 (Formula)

基于数学和逻辑表达式的规则：

```
${user.age} >= 18 and ${order.amount} <= 1000
```

## SQL执行功能

在Lua脚本中，可以使用以下SQL函数：

### sql_insert(table, data)
插入数据到指定表

```lua
local data = {
    user_id = "user_001",
    action = "login",
    created_at = now()
}
local affected_rows, error_msg = sql_insert("user_logs", data)
```

### sql_update(table, where, data)
更新指定表的数据

```lua
local where = {id = "user_001"}
local data = {status = "active"}
local affected_rows, error_msg = sql_update("users", where, data)
```

### sql_delete(table, where)
删除指定表的数据

```lua
local where = {created_at = {["<"] = now() - 86400}}
local affected_rows, error_msg = sql_delete("temp_data", where)
```

### sql_query(sql, ...args)
执行查询SQL，返回结果数组

```lua
local sql = "SELECT * FROM orders WHERE user_id = ? AND status = ?"
local results, error_msg = sql_query(sql, user_id, "pending")
```

### sql_query_one(sql, ...args)
执行查询SQL，返回单条结果

```lua
local sql = "SELECT balance FROM users WHERE id = ?"
local result, error_msg = sql_query_one(sql, user_id)
```

### sql_execute(sql, ...args)
执行任意SQL语句

```lua
local sql = "UPDATE user_stats SET total_orders = total_orders + 1 WHERE user_id = ?"
local affected_rows, error_msg = sql_execute(sql, user_id)
```

## 使用示例

### 基本使用

```go
// 创建规则执行服务
ruleService := service.NewRuleExecutionService(ruleRepo, ruleExecutor, db)

// 创建规则上下文
context := &model.RuleContext{
    RequestID:    "req_123",
    BusinessType: "order_processing",
    Trigger:      "order_created",
    Data: map[string]interface{}{
        "user_id":  "user_001",
        "amount":   100.50,
        "order_id": "order_123",
    },
}

// 执行规则
result, err := ruleService.ExecuteRule(ctx, context)
if err != nil {
    log.Fatalf("执行规则失败: %v", err)
}

fmt.Printf("执行结果: %t\n", result.Success)
fmt.Printf("动作: %s\n", result.Action)
```

### 余额检查规则示例

```lua
-- 用户余额检查规则
local user_id = context.user_id
local amount = context.amount

-- 查询用户余额
local user_result, error_msg = sql_query_one("SELECT balance FROM users WHERE id = ?", user_id)

if error_msg ~= "" then
    print("查询失败: " .. error_msg)
    return false
elseif user_result == nil then
    print("用户不存在")
    return false
else
    local balance = user_result.balance
    print("用户余额: " .. balance)
    
    -- 检查余额是否足够
    if balance >= amount then
        -- 更新余额
        local update_where = {id = user_id}
        local update_data = {balance = balance - amount}
        local affected_rows, update_error = sql_update("users", update_where, update_data)
        
        if update_error ~= "" then
            print("更新失败: " .. update_error)
            return false
        else
            print("更新成功，影响行数: " .. affected_rows)
            
            -- 记录操作日志
            local log_data = {
                user_id = user_id,
                action = "balance_deducted",
                amount = amount,
                balance_before = balance,
                balance_after = balance - amount,
                created_at = now()
            }
            sql_insert("user_balance_logs", log_data)
            
            -- 设置上下文变量
            set_context_value("balance_before", balance)
            set_context_value("balance_after", balance - amount)
            set_context_value("deducted_amount", amount)
            
            return true
        end
    else
        print("余额不足")
        set_context_value("insufficient_balance", true)
        return false
    end
end
```

## 安全注意事项

1. **SQL注入防护**: 使用参数化查询，避免直接拼接SQL字符串
2. **权限控制**: 确保Lua脚本只能访问授权的表和字段
3. **事务管理**: 在需要原子性的操作中使用数据库事务
4. **资源限制**: 设置合理的超时时间和内存限制
5. **日志记录**: 记录所有SQL操作的日志，便于审计

## 性能优化

1. **连接池**: 使用数据库连接池提高性能
2. **缓存**: 对频繁查询的数据进行缓存
3. **索引**: 为查询字段添加适当的数据库索引
4. **批量操作**: 对于大量数据操作，使用批量SQL语句
5. **异步处理**: 对于非关键操作，考虑异步处理

## 错误处理

所有SQL函数都返回两个值：结果和错误信息

```lua
local result, error_msg = sql_query("SELECT * FROM users WHERE id = ?", user_id)
if error_msg ~= "" then
    print("查询失败: " .. error_msg)
    return false
end
```

## 配置说明

### 数据库配置

确保数据库连接配置正确，包括：
- 数据库类型 (PostgreSQL)
- 连接字符串
- 用户名和密码
- 连接池配置

### 规则配置

规则支持以下配置项：
- `id`: 规则唯一标识
- `code`: 规则编码
- `name`: 规则名称
- `type`: 规则类型 (condition/lua/formula)
- `action`: 规则动作
- `business_type`: 业务类型
- `trigger`: 触发条件
- `scope`: 作用域
- `status`: 规则状态

## 扩展功能

### 自定义函数

可以在Lua脚本中注册自定义函数：

```go
// 注册自定义函数
L.SetGlobal("custom_function", L.NewFunction(func(L *lua.LState) int {
    // 自定义逻辑
    return 1
}))
```

### 事件处理

支持规则执行后的事件处理：

```go
// 注册事件处理器
ruleService.OnRuleExecuted(func(result *model.RuleResult) {
    // 处理规则执行结果
})
```

## 最佳实践

1. **规则设计**: 保持规则简单明确，避免过于复杂的逻辑
2. **测试覆盖**: 为每个规则编写充分的测试用例
3. **版本管理**: 对规则进行版本控制，支持回滚
4. **监控告警**: 监控规则执行性能和错误率
5. **文档维护**: 及时更新规则文档和使用说明 