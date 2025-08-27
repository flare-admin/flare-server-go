package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/lua_engine"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
)

// RuleContext 规则执行上下文
type RuleContext struct {
	Scope           string `json:"scope"`           // 作用域：global(全局) product(商品) user(用户) order(订单) withdraw(提现) declare(申报) payment(支付)
	Trigger         string `json:"trigger"`         // 触发动作
	ScopeID         string `json:"scopeId"`         // 作用域ID（商品ID、用户ID、订单ID等）
	ExecutionTiming string `json:"executionTiming"` // 执行时机：before(前置) after(后置) both(前后都执行)
	// 执行数据
	Data map[string]interface{} `json:"data"` // 执行数据
	// 租户信息
	TenantID string `json:"tenantId"` // 租户ID
}

// NewRuleContext 创建规则执行上下文
func NewRuleContext(scope, trigger, executionTiming, scopeId string) *RuleContext {
	return &RuleContext{
		Trigger:         trigger,
		Scope:           scope,
		ScopeID:         scopeId,
		ExecutionTiming: executionTiming,
		Data:            make(map[string]interface{}),
		TenantID:        "",
	}
}

// SetData 设置数据
func (rc *RuleContext) SetData(data map[string]interface{}) {
	rc.Data = data
}

// AddData 添加数据
func (rc *RuleContext) AddData(key string, value interface{}) {
	rc.Data[key] = value
}

// GetData 获取数据
func (rc *RuleContext) GetData(key string) interface{} {
	return rc.Data[key]
}

// WithTenantID 设置租户ID
func (rc *RuleContext) WithTenantID(tenantID string) {
	rc.TenantID = tenantID
}
func (rc *RuleContext) GetCtx(ctx context.Context) context.Context {
	return actx.WithTenantId(ctx, rc.TenantID)
}

// Validate 验证上下文
func (rc *RuleContext) Validate() error {
	if rc.Trigger == "" {
		return fmt.Errorf("trigger cannot be empty")
	}

	if rc.Scope == "" {
		return fmt.Errorf("business type cannot be empty")
	}

	if rc.ExecutionTiming == "" {
		return fmt.Errorf("scope id cannot be empty")
	}
	return nil
}

// RuleResult 规则执行结果
type RuleResult struct {
	// 执行结果
	lua_engine.ExecuteResult

	// 执行统计
	ExecuteTime int64 `json:"executeTime"` // 执行时间(毫秒)
	ExecuteAt   int64 `json:"executeAt"`   // 执行时间戳

	// 规则执行链路
	ExecutionChain []*RuleExecutionStep `json:"executionChain"` // 规则执行链路

	// 租户信息
	TenantID string `json:"tenantId"` // 租户ID
}

// RuleExecutionStep 规则执行步骤
type RuleExecutionStep struct {
	// 规则信息
	RuleID   string `json:"ruleId"`   // 规则ID
	RuleCode string `json:"ruleCode"` // 规则编码
	RuleName string `json:"ruleName"` // 规则名称
	Priority int32  `json:"priority"` // 优先级

	// 执行输入
	Input map[string]interface{} `json:"input"` // 执行输入数据

	// 执行输出
	Output map[string]interface{} `json:"output"` // 执行输出数据

	// 执行结果
	Valid       bool   `json:"valid"`       // 是否验证通过
	Action      string `json:"action"`      // 执行动作
	Error       string `json:"error"`       // 错误信息
	ExecuteTime int64  `json:"executeTime"` // 执行时间(毫秒)

	// 执行时间
	ExecuteAt int64 `json:"executeAt"` // 执行时间戳
}

// NewRuleResult 创建规则执行结果
func NewRuleResult() *RuleResult {
	return &RuleResult{
		ExecuteTime: 0,
		ExecuteAt:   utils.GetDateUnix(),
		TenantID:    "",
	}
}

// SetSuccess 设置成功结果
func (rr *RuleResult) SetSuccess(valid bool, action string) {
	rr.Valid = valid
	rr.Action = action
	rr.Error = ""
	rr.ExecuteAt = utils.GetDateUnix()
}

// SetFailure 设置失败结果
func (rr *RuleResult) SetFailure(action, reason, error string) {
	rr.Valid = false
	rr.Action = action
	rr.Error = error
	rr.ErrorReason = reason
	rr.ExecuteAt = utils.GetDateUnix()
}

// SetExecuteTime 设置执行时间
func (rr *RuleResult) SetExecuteTime(executeTime int64) {
	rr.ExecuteTime = executeTime
}

// IsSuccess 判断规则执行是否成功
func (rr *RuleResult) IsSuccess() bool {
	return rr.Valid && rr.Error == ""
}

// AddExecutionStep 添加执行步骤
func (rr *RuleResult) AddExecutionStep(step *RuleExecutionStep) {
	if rr.ExecutionChain == nil {
		rr.ExecutionChain = make([]*RuleExecutionStep, 0)
	}
	rr.ExecutionChain = append(rr.ExecutionChain, step)
}

// GetExecutionChain 获取执行链路
func (rr *RuleResult) GetExecutionChain() []*RuleExecutionStep {
	return rr.ExecutionChain
}

// GetLastExecutionStep 获取最后一个执行步骤
func (rr *RuleResult) GetLastExecutionStep() *RuleExecutionStep {
	if len(rr.ExecutionChain) == 0 {
		return nil
	}
	return rr.ExecutionChain[len(rr.ExecutionChain)-1]
}

// GetExecutionStepCount 获取执行步骤数量
func (rr *RuleResult) GetExecutionStepCount() int {
	return len(rr.ExecutionChain)
}

// GetSuccessfulSteps 获取成功执行的步骤
func (rr *RuleResult) GetSuccessfulSteps() []*RuleExecutionStep {
	var successfulSteps []*RuleExecutionStep
	for _, step := range rr.ExecutionChain {
		if step.IsSuccess() {
			successfulSteps = append(successfulSteps, step)
		}
	}
	return successfulSteps
}

// GetFailedSteps 获取失败执行的步骤
func (rr *RuleResult) GetFailedSteps() []*RuleExecutionStep {
	var failedSteps []*RuleExecutionStep
	for _, step := range rr.ExecutionChain {
		if !step.IsSuccess() {
			failedSteps = append(failedSteps, step)
		}
	}
	return failedSteps
}

// GetTotalExecuteTime 获取总执行时间
func (rr *RuleResult) GetTotalExecuteTime() int64 {
	var totalTime int64
	for _, step := range rr.ExecutionChain {
		totalTime += step.ExecuteTime
	}
	return totalTime
}

// ToJSON 转换为JSON字符串
func (rr *RuleResult) ToJSON() (string, error) {
	data, err := json.Marshal(rr)
	if err != nil {
		return "", fmt.Errorf("marshal rule result error: %v", err)
	}
	return string(data), nil
}

// NewRuleExecutionStep 创建规则执行步骤
func NewRuleExecutionStep(ruleID, ruleCode, ruleName string, priority int32) *RuleExecutionStep {
	return &RuleExecutionStep{
		RuleID:      ruleID,
		RuleCode:    ruleCode,
		RuleName:    ruleName,
		Priority:    priority,
		Input:       make(map[string]interface{}),
		Output:      make(map[string]interface{}),
		Valid:       false,
		Action:      "",
		Error:       "",
		ExecuteTime: 0,
		ExecuteAt:   utils.GetDateUnix(),
	}
}

// SetInput 设置输入数据
func (step *RuleExecutionStep) SetInput(input map[string]interface{}) {
	step.Input = input
}

// SetOutput 设置输出数据
func (step *RuleExecutionStep) SetOutput(output map[string]interface{}) {
	step.Output = output
}

// SetSuccess 设置成功结果
func (step *RuleExecutionStep) SetSuccess(valid bool, action string) {
	step.Valid = valid
	step.Action = action
	step.Error = ""
	step.ExecuteAt = utils.GetDateUnix()
}

// SetFailure 设置失败结果
func (step *RuleExecutionStep) SetFailure(action, error string) {
	step.Valid = false
	step.Action = action
	step.Error = error
	step.ExecuteAt = utils.GetDateUnix()
}

// SetExecuteTime 设置执行时间
func (step *RuleExecutionStep) SetExecuteTime(executeTime int64) {
	step.ExecuteTime = executeTime
}

// IsSuccess 判断执行是否成功
func (step *RuleExecutionStep) IsSuccess() bool {
	return step.Valid && step.Error == ""
}
