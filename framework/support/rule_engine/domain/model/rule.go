package model

import (
	"encoding/json"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
)

// Rule 规则领域模型
type Rule struct {
	// 基础信息
	ID          string `json:"id"`          // 规则ID
	Code        string `json:"code"`        // 规则编码
	Name        string `json:"name"`        // 规则名称
	Description string `json:"description"` // 规则描述
	CategoryID  string `json:"categoryId"`  // 分类ID
	TemplateID  string `json:"templateId"`  // 模板ID（可选）

	// 规则配置
	Type    string `json:"type"`    // 规则类型：condition(条件规则) lua(lua脚本规则) formula(公式规则)
	Version string `json:"version"` // 规则版本
	Status  int32  `json:"status"`  // 状态：1-启用 2-禁用

	// 触发配置
	Triggers        []string `json:"triggers"`        // 触发动作列表
	Scope           string   `json:"scope"`           // 作用域：global(全局) product(商品) user(用户) order(订单) withdraw(提现) declare(申报) payment(支付)
	ScopeID         string   `json:"scopeId"`         // 作用域ID（商品ID、用户ID、订单ID等）
	ExecutionTiming string   `json:"executionTiming"` // 执行时机：before(前置) after(后置) both(前后都执行)

	// 规则内容（从模板继承或自定义）
	Conditions  string `json:"conditions"`  // 条件表达式(JSON格式)
	LuaScript   string `json:"luaScript"`   // Lua脚本代码
	Formula     string `json:"formula"`     // 计算公式
	FormulaVars string `json:"formulaVars"` // 公式变量映射(JSON格式)

	// 动作配置
	Action string `json:"action"` // 规则动作：allow(允许) deny(拒绝) modify(修改) notify(通知) redirect(重定向)

	// 优先级和排序
	Priority int32 `json:"priority"` // 优先级，数字越大优先级越高
	Sorting  int32 `json:"sorting"`  // 排序权重

	// 统计信息
	ExecuteCount  int64 `json:"executeCount"`  // 执行次数
	SuccessCount  int64 `json:"successCount"`  // 成功次数
	LastExecuteAt int64 `json:"lastExecuteAt"` // 最后执行时间

	// 时间信息
	CreatedAt int64 `json:"createdAt"` // 创建时间
	UpdatedAt int64 `json:"updatedAt"` // 更新时间

	// 租户信息
	TenantID string `json:"tenantId"` // 租户ID
}

// NewRule 创建规则
func NewRule(code, name, description, categoryID, ruleType string) *Rule {
	now := utils.GetDateUnix()
	return &Rule{
		ID:              "",
		Code:            code,
		Name:            name,
		Description:     description,
		CategoryID:      categoryID,
		TemplateID:      "",
		Type:            ruleType,
		Version:         "1.0.0",
		Status:          1,
		Triggers:        []string{},
		Scope:           "global",
		ScopeID:         "",
		ExecutionTiming: "before",
		Conditions:      "{}",
		LuaScript:       "",
		Formula:         "",
		FormulaVars:     "{}",
		Action:          "allow",
		Priority:        0,
		Sorting:         0,
		ExecuteCount:    0,
		SuccessCount:    0,
		LastExecuteAt:   0,
		CreatedAt:       now,
		UpdatedAt:       now,
		TenantID:        "",
	}
}

func (r *Rule) Completion() {
	if len(r.Triggers) == 0 {
		r.Triggers = []string{}
	}
	if r.Conditions == "" {
		r.Conditions = "{}"
	}
	if r.Formula == "" {
		r.Formula = ""
	}
	if r.FormulaVars == "" {
		r.FormulaVars = "{}"
	}
	if r.ExecutionTiming == "" {
		r.ExecutionTiming = "before"
	}
}

// SetTemplate 设置模板
func (r *Rule) SetTemplate(templateID string) {
	r.TemplateID = templateID
	r.UpdatedAt = utils.GetDateUnix()
}

// SetTriggers 设置触发动作
func (r *Rule) SetTriggers(triggers []string) error {
	if len(triggers) == 0 {
		r.Triggers = []string{}
		return nil
	}

	// 验证触发动作
	validTriggers := []string{"create", "update", "delete", "approve", "reject", "placeOrder", "pay", "withdraw", "declare"}
	validTriggerMap := make(map[string]bool)
	for _, vt := range validTriggers {
		validTriggerMap[vt] = true
	}

	for _, trigger := range triggers {
		if !validTriggerMap[trigger] {
			return fmt.Errorf("invalid trigger: %s", trigger)
		}
	}

	r.Triggers = triggers
	r.UpdatedAt = utils.GetDateUnix()
	return nil
}

// GetTriggers 获取触发动作
func (r *Rule) GetTriggers() []string {
	if len(r.Triggers) == 0 {
		return make([]string, 0)
	}
	return r.Triggers
}

// AddTrigger 添加触发动作
func (r *Rule) AddTrigger(trigger string) error {
	if trigger == "" {
		return fmt.Errorf("trigger cannot be empty")
	}

	triggers := r.GetTriggers()

	// 检查是否已存在
	for _, t := range triggers {
		if t == trigger {
			return nil // 已存在，不重复添加
		}
	}

	triggers = append(triggers, trigger)
	return r.SetTriggers(triggers)
}

// RemoveTrigger 移除触发动作
func (r *Rule) RemoveTrigger(trigger string) error {
	triggers := r.GetTriggers()

	for i, t := range triggers {
		if t == trigger {
			triggers = append(triggers[:i], triggers[i+1:]...)
			return r.SetTriggers(triggers)
		}
	}

	return nil
}

// SetScope 设置作用域
func (r *Rule) SetScope(scope, scopeID string) error {
	r.Scope = scope
	r.ScopeID = scopeID
	r.UpdatedAt = utils.GetDateUnix()
	return nil
}

// SetExecutionTiming 设置执行时机
func (r *Rule) SetExecutionTiming(timing string) error {
	// 验证执行时机
	validTimings := []string{"before", "after", "both"}
	isValidTiming := false
	for _, validTiming := range validTimings {
		if timing == validTiming {
			isValidTiming = true
			break
		}
	}

	if !isValidTiming {
		return fmt.Errorf("invalid execution timing: %s", timing)
	}

	r.ExecutionTiming = timing
	r.UpdatedAt = utils.GetDateUnix()
	return nil
}

// GetExecutionTiming 获取执行时机
func (r *Rule) GetExecutionTiming() string {
	if r.ExecutionTiming == "" {
		return "before"
	}
	return r.ExecutionTiming
}

// SetConditions 设置条件表达式
func (r *Rule) SetConditions(conditions map[string]interface{}) error {
	// 验证JSON格式
	conditionsJSON, err := json.Marshal(conditions)
	if err != nil {
		return fmt.Errorf("invalid conditions format: %v", err)
	}

	r.Conditions = string(conditionsJSON)
	r.UpdatedAt = utils.GetDateUnix()
	return nil
}

// GetConditions 获取条件表达式
func (r *Rule) GetConditions() (map[string]interface{}, error) {
	var conditions map[string]interface{}
	if r.Conditions == "" {
		return make(map[string]interface{}), nil
	}

	err := json.Unmarshal([]byte(r.Conditions), &conditions)
	if err != nil {
		return nil, fmt.Errorf("invalid conditions format: %v", err)
	}

	return conditions, nil
}

// SetLuaScript 设置Lua脚本
func (r *Rule) SetLuaScript(script string) error {
	if r.Type != "lua" {
		return fmt.Errorf("rule type is not lua")
	}

	r.LuaScript = script
	r.UpdatedAt = utils.GetDateUnix()
	return nil
}

// SetFormula 设置计算公式
func (r *Rule) SetFormula(formula string, vars map[string]interface{}) error {
	if r.Type != "formula" {
		return fmt.Errorf("rule type is not formula")
	}

	r.Formula = formula

	// 序列化变量映射
	varsJSON, err := json.Marshal(vars)
	if err != nil {
		return fmt.Errorf("invalid formula vars format: %v", err)
	}

	r.FormulaVars = string(varsJSON)
	r.UpdatedAt = utils.GetDateUnix()
	return nil
}

// GetFormulaVars 获取公式变量
func (r *Rule) GetFormulaVars() (map[string]interface{}, error) {
	var vars map[string]interface{}
	if r.FormulaVars == "" {
		return make(map[string]interface{}), nil
	}

	err := json.Unmarshal([]byte(r.FormulaVars), &vars)
	if err != nil {
		return nil, fmt.Errorf("invalid formula vars format: %v", err)
	}

	return vars, nil
}

// SetAction 设置规则动作
func (r *Rule) SetAction(action string) error {
	r.Action = action
	r.UpdatedAt = utils.GetDateUnix()
	return nil
}

// GetAction 获取规则动作
func (r *Rule) GetAction() string {
	return r.Action
}

// Enable 启用规则
func (r *Rule) Enable() {
	r.Status = 1
	r.UpdatedAt = utils.GetDateUnix()
}

// Disable 禁用规则
func (r *Rule) Disable() {
	r.Status = 2
	r.UpdatedAt = utils.GetDateUnix()
}

// IsEnabled 是否启用
func (r *Rule) IsEnabled() bool {
	return r.Status == 1
}

// Validate 验证规则
func (r *Rule) Validate() error {
	if r.Code == "" {
		return fmt.Errorf("rule code cannot be empty")
	}

	if r.Name == "" {
		return fmt.Errorf("rule name cannot be empty")
	}

	if r.CategoryID == "" {
		return fmt.Errorf("category ID cannot be empty")
	}

	// 验证规则类型
	validTypes := []string{"condition", "lua", "formula"}
	isValidType := false
	for _, validType := range validTypes {
		if r.Type == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("invalid rule type: %s", r.Type)
	}

	// 验证执行时机
	validTimings := []string{"before", "after", "both"}
	isValidTiming := false
	for _, validTiming := range validTimings {
		if r.ExecutionTiming == validTiming {
			isValidTiming = true
			break
		}
	}

	if !isValidTiming {
		return fmt.Errorf("invalid execution timing: %s", r.ExecutionTiming)
	}

	// 根据类型验证内容
	switch r.Type {
	case "condition":
		if _, err := r.GetConditions(); err != nil {
			return fmt.Errorf("invalid conditions: %v", err)
		}
	case "lua":
		if r.LuaScript == "" {
			return fmt.Errorf("lua script cannot be empty for lua rule")
		}
	case "formula":
		if r.Formula == "" {
			return fmt.Errorf("formula cannot be empty for formula rule")
		}
		if _, err := r.GetFormulaVars(); err != nil {
			return fmt.Errorf("invalid formula vars: %v", err)
		}
	}

	// 验证触发动作
	triggers := r.GetTriggers()
	validTriggers := []string{"create", "update", "delete", "approve", "reject", "placeOrder", "pay", "withdraw", "declare"}
	validTriggerMap := make(map[string]bool)
	for _, vt := range validTriggers {
		validTriggerMap[vt] = true
	}

	for _, trigger := range triggers {
		if !validTriggerMap[trigger] {
			return fmt.Errorf("invalid trigger: %s", trigger)
		}
	}

	return nil
}

// Update 更新规则
func (r *Rule) Update(name, description string) {
	r.Name = name
	r.Description = description
	r.UpdatedAt = utils.GetDateUnix()
}

// SetPriority 设置优先级
func (r *Rule) SetPriority(priority int32) {
	r.Priority = priority
	r.UpdatedAt = utils.GetDateUnix()
}

// SetSorting 设置排序
func (r *Rule) SetSorting(sorting int32) {
	r.Sorting = sorting
	r.UpdatedAt = utils.GetDateUnix()
}

// RecordExecution 记录执行
func (r *Rule) RecordExecution(success bool) {
	r.ExecuteCount++
	if success {
		r.SuccessCount++
	}
	r.LastExecuteAt = utils.GetDateUnix()
	r.UpdatedAt = utils.GetDateUnix()
}

// GetSuccessRate 获取成功率
func (r *Rule) GetSuccessRate() float64 {
	if r.ExecuteCount == 0 {
		return 0.0
	}
	return float64(r.SuccessCount) / float64(r.ExecuteCount) * 100
}

// ApplyTemplate 应用模板内容
func (r *Rule) ApplyTemplate(template *RuleTemplate) error {
	if template == nil {
		return fmt.Errorf("template cannot be nil")
	}

	// 设置模板ID
	r.TemplateID = template.ID
	r.Type = template.Type

	// 应用模板内容
	switch template.Type {
	case "condition":
		if conditions, err := template.GetConditions(); err == nil {
			r.SetConditions(conditions)
		}
	case "lua":
		r.LuaScript = template.LuaScript
	case "formula":
		r.Formula = template.Formula
		if vars, err := template.GetFormulaVars(); err == nil {
			r.SetFormula(template.Formula, vars)
		}
	}

	r.UpdatedAt = utils.GetDateUnix()
	return nil
}

// IsTriggered 检查是否触发
func (r *Rule) IsTriggered(trigger string) bool {
	triggers := r.GetTriggers()

	for _, t := range triggers {
		if t == trigger {
			return true
		}
	}

	return false
}

// IsInScope 检查是否在作用域内
func (r *Rule) IsInScope(scope, scopeID string) bool {
	if r.Scope == "global" {
		return true
	}

	if r.Scope == scope && r.ScopeID == scopeID {
		return true
	}

	return false
}

// IsExecutionTiming 检查是否匹配执行时机
func (r *Rule) IsExecutionTiming(timing string) bool {
	if r.ExecutionTiming == "both" {
		return true
	}
	return r.ExecutionTiming == timing
}

// SetTenantID 设置租户ID
func (r *Rule) SetTenantID(tenantID string) {
	r.TenantID = tenantID
	r.UpdatedAt = utils.GetDateUnix()
}

// GetTenantID 获取租户ID
func (r *Rule) GetTenantID() string {
	return r.TenantID
}
