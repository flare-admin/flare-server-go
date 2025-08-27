-- SQL操作示例脚本
-- 这个脚本展示了如何在Lua规则中使用SQL执行函数

-- 获取上下文数据
local user_id = context.user_id
local amount = context.amount
local order_id = context.order_id

-- 1. 插入数据示例
-- 插入用户操作日志
local insert_data = {
    user_id = user_id,
    action = "order_created",
    amount = amount,
    order_id = order_id,
    created_at = now()
}

local affected_rows, error_msg = sql_insert("user_operation_logs", insert_data)
if error_msg ~= "" then
    print("插入失败: " .. error_msg)
    return false
else
    print("插入成功，影响行数: " .. affected_rows)
end

-- 2. 查询数据示例
-- 查询用户信息
local user_sql = "SELECT * FROM users WHERE id = ?"
local user_result, error_msg = sql_query_one(user_sql, user_id)
if error_msg ~= "" then
    print("查询失败: " .. error_msg)
    return false
elseif user_result == nil then
    print("用户不存在")
    return false
else
    print("用户余额: " .. user_result.balance)
end

-- 3. 更新数据示例
-- 更新用户余额
local update_where = {id = user_id}
local update_data = {balance = user_result.balance - amount}
local affected_rows, error_msg = sql_update("users", update_where, update_data)
if error_msg ~= "" then
    print("更新失败: " .. error_msg)
    return false
else
    print("更新成功，影响行数: " .. affected_rows)
end

-- 4. 复杂查询示例
-- 查询用户最近的订单
local recent_orders_sql = [[
    SELECT * FROM orders 
    WHERE user_id = ? 
    AND created_at > ? 
    ORDER BY created_at DESC 
    LIMIT 5
]]
local recent_time = now() - 86400 -- 24小时前
local orders_result, error_msg = sql_query(recent_orders_sql, user_id, recent_time)
if error_msg ~= "" then
    print("查询订单失败: " .. error_msg)
    return false
else
    print("最近订单数量: " .. #orders_result)
    for i, order in ipairs(orders_result) do
        print("订单ID: " .. order.id .. ", 金额: " .. order.amount)
    end
end

-- 5. 删除数据示例（谨慎使用）
-- 删除过期的临时数据
local delete_where = {created_at = {["<"] = now() - 604800}} -- 7天前
local affected_rows, error_msg = sql_delete("temp_data", delete_where)
if error_msg ~= "" then
    print("删除失败: " .. error_msg)
else
    print("删除过期数据，影响行数: " .. affected_rows)
end

-- 6. 执行自定义SQL示例
-- 执行复杂的业务逻辑SQL
local custom_sql = [[
    UPDATE user_statistics 
    SET total_orders = total_orders + 1,
        total_amount = total_amount + ?,
        last_order_time = ?
    WHERE user_id = ?
]]
local affected_rows, error_msg = sql_execute(custom_sql, amount, now(), user_id)
if error_msg ~= "" then
    print("执行自定义SQL失败: " .. error_msg)
    return false
else
    print("更新用户统计成功，影响行数: " .. affected_rows)
end

-- 设置上下文变量
set_context_value("sql_affected_rows", affected_rows)
set_context_value("user_balance", user_result.balance - amount)

-- 返回执行结果
return true 