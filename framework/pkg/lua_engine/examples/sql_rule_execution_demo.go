package examples

import (
	"fmt"
	"log"

	"github.com/flare-admin/flare-server-go/framework/pkg/lua_engine"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/service"
)

func main() {
	// 创建规则执行服务
	ruleService := createRuleExecutionService()

	// 创建规则上下文
	context := &model.RuleContext{
		RequestID:    "req_123",
		BusinessType: "order_processing",
		Trigger:      "order_created",
		Scope:        "user_balance",
		Data: map[string]interface{}{
			"user_id":  "user_001",
			"amount":   100.50,
			"order_id": "order_123",
		},
	}

	// 创建Lua规则
	rule := &model.Rule{
		ID:     "rule_001",
		Code:   "sql_balance_check",
		Name:   "SQL余额检查规则",
		Type:   "lua",
		Action: "allow",
		LuaScript: `
-- SQL操作示例
local user_id = context.user_id
local amount = context.amount

-- 查询用户余额
local user_sql = "SELECT balance FROM users WHERE id = ?"
local user_result, error_msg = sql_query_one(user_sql, user_id)

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
            local log_affected, log_error = sql_insert("user_balance_logs", log_data)
            
            if log_error ~= "" then
                print("记录日志失败: " .. log_error)
            else
                print("记录日志成功，影响行数: " .. log_affected)
            end
            
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
`,
	}

	// 执行规则
	result, err := ruleService.ExecuteRule(context.Background(), context)
	if err != nil {
		log.Fatalf("执行规则失败: %v", err)
	}

	// 输出结果
	fmt.Printf("规则执行结果:\n")
	fmt.Printf("  请求ID: %s\n", context.RequestID)
	fmt.Printf("  执行结果: %t\n", result.Valid)
	fmt.Printf("  动作: %s\n", result.Action)
	fmt.Printf("  执行时间: %dms\n", result.ExecuteTime)

	// 输出执行链路
	if len(result.ExecutionChain) > 0 {
		fmt.Printf("  执行链路:\n")
		for i, step := range result.ExecutionChain {
			fmt.Printf("    步骤 %d:\n", i+1)
			fmt.Printf("      规则ID: %s\n", step.RuleID)
			fmt.Printf("      规则编码: %s\n", step.RuleCode)
			fmt.Printf("      规则名称: %s\n", step.RuleName)
			fmt.Printf("      优先级: %d\n", step.Priority)
			fmt.Printf("      执行结果: %t\n", step.Valid)
			fmt.Printf("      动作: %s\n", step.Action)
			fmt.Printf("      执行时间: %dms\n", step.ExecuteTime)

			if step.Error != "" {
				fmt.Printf("      错误信息: %s\n", step.Error)
			}

			if step.Output != nil {
				fmt.Printf("      输出变量:\n")
				for key, value := range step.Output {
					fmt.Printf("        %s: %v\n", key, value)
				}
			}
		}
	}

	if result.Variables != nil {
		fmt.Printf("  最终输出变量:\n")
		for key, value := range result.Variables {
			fmt.Printf("    %s: %v\n", key, value)
		}
	}
}

// createRuleExecutionService 创建规则执行服务
func createRuleExecutionService() *service.RuleExecutionService {
	// 这里需要根据实际项目结构创建依赖
	// 为了演示，我们创建一个模拟的服务

	// 创建规则执行器
	ruleExecutor := lua_engine.NewRuleExecutor()

	// 创建数据库连接（这里需要根据实际项目配置）
	// db := database.NewDatabase(config)

	// 创建规则仓储（这里需要根据实际项目配置）
	// ruleRepo := repository.NewRuleRepository(db)

	// 为了演示，我们返回一个模拟的服务
	// 实际使用时需要传入真实的依赖
	return &service.RuleExecutionService{
		RuleExecutor: ruleExecutor,
		// SQLExecutor: service.NewSQLExecutor(db),
		// RuleRepo: ruleRepo,
	}
}
