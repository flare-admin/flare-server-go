package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"golang.org/x/exp/slices"
	"strconv"
	"strings"
	"time"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/lua_engine"
	ruleengineerr "github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/err"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/repository"
)

// sqlExecutor SQL执行器实现

// RuleExecutionService 规则执行服务
type RuleExecutionService struct {
	ruleRepo     repository.IRuleRepository
	ruleExecutor *lua_engine.RuleExecutor
}

// NewRuleExecutionService 创建规则执行服务
func NewRuleExecutionService(
	ruleRepo repository.IRuleRepository,
	ruleExecutor *lua_engine.RuleExecutor,
) *RuleExecutionService {
	return &RuleExecutionService{
		ruleRepo:     ruleRepo,
		ruleExecutor: ruleExecutor,
	}
}

// ExecuteRules 执行多个规则
// 根据上下文中的业务类型和触发动作自动匹配并执行所有相关规则
// 按照优先级排序执行，执行失败时中断，上一个规则的执行结果是下一个规则的输入
func (s *RuleExecutionService) ExecuteRules(ctx context.Context, context *model.RuleContext) (*model.RuleResult, *herrors.HError) {
	// 验证上下文
	if err := context.Validate(); err != nil {
		return nil, ruleengineerr.RuleContextInvalid(err)
	}

	// 查找匹配的规则
	rules, err := s.findMatchingRules(ctx, context)
	if err != nil {
		return nil, err
	}
	// 创建最终结果对象
	finalResult := model.NewRuleResult()
	finalResult.SetSuccess(true, "allow")
	// 如果没有找到匹配的规则，返回空结果
	if len(rules) == 0 {
		return finalResult, nil
	}

	// 按照优先级排序规则（优先级数字越大优先级越高）
	s.sortRulesByPriority(rules)

	// 执行规则链，上一个规则的执行结果是下一个规则的输入
	var currentContext = context

	for i, rule := range rules {
		// 创建执行步骤
		step := model.NewRuleExecutionStep(rule.ID, rule.Code, rule.Name, rule.Priority)

		// 记录输入数据
		step.SetInput(currentContext.Data)

		// 执行单个规则
		result, err := s.executeSingleRule(ctx, rule, currentContext)
		if err != nil {
			// 执行失败，记录失败步骤并中断执行链
			step.SetFailure("deny", err.Error())
			finalResult.AddExecutionStep(step)
			finalResult.SetFailure("deny", err.Reason, err.Error())
			return finalResult, nil
		}

		// 记录执行结果到步骤中
		step.SetExecuteTime(result.ExecuteTime)
		if result.IsSuccess() {
			step.SetSuccess(result.Valid, result.Action)
			// 记录输出数据
			if result.Context != nil {
				step.SetOutput(result.Context)
			}
		} else {
			step.SetFailure(result.Action, result.Error)
		}

		// 添加执行步骤到结果中
		finalResult.AddExecutionStep(step)

		// 如果规则执行失败，中断执行链
		if !result.IsSuccess() {
			finalResult.SetFailure(result.Action, result.ErrorReason, result.Error)
			break
		}

		// 将当前规则的执行结果作为下一个规则的输入
		if i < len(rules)-1 {
			// 更新上下文数据，将当前规则的输出变量合并到上下文中
			currentContext = s.updateContextWithRuleResult(currentContext, result)
		}
	}

	// 设置最终结果的总执行时间
	finalResult.SetExecuteTime(finalResult.GetTotalExecuteTime())

	return finalResult, nil
}

// ExecuteRuleByCode 根据编码执行规则
func (s *RuleExecutionService) ExecuteRuleByCode(ctx context.Context, code string, context *model.RuleContext) (*model.RuleResult, *herrors.HError) {
	rule, err := s.ruleRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, ruleengineerr.RuleGetFailed(err)
	}

	return s.executeSingleRule(ctx, rule, context)
}

// findMatchingRules 查找匹配的规则
func (s *RuleExecutionService) findMatchingRules(ctx context.Context, context *model.RuleContext) ([]*model.Rule, *herrors.HError) {
	var rules []*model.Rule
	// 获取全局规则
	globalRules, err := s.ruleRepo.FindByScope(ctx, "global")
	if err != nil {
		hlog.CtxErrorf(ctx, "get global rules failed: %v", err)
		return nil, ruleengineerr.RuleGetFailed(err)
	}
	if len(globalRules) > 0 {
		rules = append(rules, globalRules...)
	}
	// 根据业务类型查找规则
	if context.Scope != "" {
		businessRules, err := s.ruleRepo.FindByScope(ctx, context.Scope)
		if err != nil {
			return nil, ruleengineerr.RuleGetFailed(err)
		}
		rules = append(rules, businessRules...)
	}

	// 去重并过滤启用的规则
	ruleMap := make(map[string]*model.Rule)
	for _, rule := range rules {
		// 不是指定范围的跳过
		if rule.ScopeID != "" && !slices.Contains(strings.Split(rule.ScopeID, ","), context.ScopeID) {
			continue
		}

		if rule.IsEnabled() &&
			(rule.ExecutionTiming == context.ExecutionTiming || rule.ExecutionTiming == "both") &&
			slices.Contains(rule.Triggers, context.Trigger) {
			ruleMap[rule.ID] = rule
		}
	}

	// 转换为切片
	var uniqueRules []*model.Rule
	for _, rule := range ruleMap {
		uniqueRules = append(uniqueRules, rule)
	}

	return uniqueRules, nil
}

// executeSingleRule 执行单个规则
func (s *RuleExecutionService) executeSingleRule(_ context.Context, rule *model.Rule, context *model.RuleContext) (*model.RuleResult, *herrors.HError) {
	// 创建结果对象
	result := model.NewRuleResult()
	context.AddData("scopeId", context.ScopeID)
	// 执行规则
	execResult, err := s.executeRuleWithExecutor(rule, context)
	if err != nil {
		result.SetFailure("deny", "rule_execute_err", err.Error())
		return result, nil
	}

	// 设置执行结果
	if execResult.Valid {
		result.SetSuccess(execResult.Valid, execResult.Action)
	} else {
		result.SetFailure(execResult.Action, execResult.ErrorReason, execResult.Error)
	}

	result.SetExecuteTime(execResult.ExecuteTime)

	return result, nil
}

// executeRuleWithExecutor 使用规则执行器执行规则
func (s *RuleExecutionService) executeRuleWithExecutor(rule *model.Rule, context *model.RuleContext) (*lua_engine.ExecuteResult, error) {
	switch rule.Type {
	case "condition":
		return s.executeConditionRule(rule, context)
	case "lua":
		return s.executeLuaRule(rule, context)
	case "formula":
		return s.executeFormulaRule(rule, context)
	default:
		return nil, fmt.Errorf("unsupported rule type: %s", rule.Type)
	}
}

// executeLuaRule 执行Lua脚本规则
func (s *RuleExecutionService) executeLuaRule(rule *model.Rule, context *model.RuleContext) (*lua_engine.ExecuteResult, error) {
	context.AddData("tenantId", context.TenantID)

	// 使用自定义执行器执行Lua脚本
	execResult, err := s.ruleExecutor.Execute(rule.LuaScript, &lua_engine.ExecuteOptions{
		Timeout:   5 * time.Second,
		Context:   context.Data,
		MaxMemory: 10 * 1024 * 1024, // 10MB
	})
	if err != nil {
		return nil, fmt.Errorf("lua script execution failed: %w", err)
	}

	return execResult, nil
}

// executeFormulaRule 执行公式规则
func (s *RuleExecutionService) executeFormulaRule(rule *model.Rule, context *model.RuleContext) (*lua_engine.ExecuteResult, error) {
	// 替换公式中的变量
	formula := s.replaceFormulaVariables(rule.Formula, context.Data)

	// 计算公式结果
	formulaResult, err := s.evaluateFormula(formula)
	if err != nil {
		return nil, fmt.Errorf("formula evaluation failed: %w", err)
	}

	return &lua_engine.ExecuteResult{
		Valid:       formulaResult.Valid,
		Action:      rule.Action,
		Context:     map[string]interface{}{"result": formulaResult.Value},
		ExecuteTime: 0,
	}, nil
}

// replaceFormulaVariables 替换公式中的变量
func (s *RuleExecutionService) replaceFormulaVariables(formula string, data map[string]interface{}) string {
	// 简单的变量替换，实际项目中可以使用更复杂的表达式解析器
	for key, value := range data {
		placeholder := "${" + key + "}"
		formula = strings.ReplaceAll(formula, placeholder, fmt.Sprintf("%v", value))
	}
	return formula
}

// evaluateFormula 评估公式
func (s *RuleExecutionService) evaluateFormula(formula string) (*FormulaResult, error) {
	// 基本语法检查
	if formula == "" {
		return &FormulaResult{Valid: false, Value: nil}, fmt.Errorf("公式不能为空")
	}

	// 检查括号匹配
	if !s.checkBracketMatch(formula) {
		return &FormulaResult{Valid: false, Value: nil}, fmt.Errorf("公式括号不匹配")
	}

	// 尝试解析为数学表达式
	result, err := s.evaluateMathExpression(formula)
	if err != nil {
		// 如果不是数学表达式，尝试解析为逻辑表达式
		result, err = s.evaluateLogicalExpression(formula)
		if err != nil {
			return &FormulaResult{Valid: false, Value: nil}, fmt.Errorf("公式解析失败: %w", err)
		}
	}

	return &FormulaResult{Valid: true, Value: result}, nil
}

// evaluateMathExpression 评估数学表达式
func (s *RuleExecutionService) evaluateMathExpression(formula string) (interface{}, error) {
	// 简单的数学表达式计算
	// 支持基本的四则运算和比较运算

	// 移除所有空格
	formula = strings.ReplaceAll(formula, " ", "")

	// 检查是否包含数学运算符
	mathOperators := []string{"+", "-", "*", "/", "(", ")", ">", "<", ">=", "<=", "==", "!="}
	hasMathOperator := false
	for _, op := range mathOperators {
		if strings.Contains(formula, op) {
			hasMathOperator = true
			break
		}
	}

	if !hasMathOperator {
		return nil, fmt.Errorf("不是数学表达式")
	}

	// 尝试解析为数字
	if num, err := strconv.ParseFloat(formula, 64); err == nil {
		return num, nil
	}

	// 简单的表达式计算
	return s.calculateSimpleExpression(formula)
}

// executeConditionRule 执行条件规则
func (s *RuleExecutionService) executeConditionRule(rule *model.Rule, context *model.RuleContext) (*lua_engine.ExecuteResult, error) {
	// 获取条件配置
	conditions, err := rule.GetConditions()
	if err != nil {
		return nil, fmt.Errorf("解析条件配置失败: %w", err)
	}

	// 检查所有条件是否匹配
	matched := true
	var failedCondition string

	for _, condition := range conditions {
		if conditionMap, ok := condition.(map[string]interface{}); ok {
			if !s.evaluateCondition(conditionMap, context.Data) {
				matched = false
				if field, ok := conditionMap["field"].(string); ok {
					failedCondition = field
				}
				break
			}
		}
	}

	// 构建输出变量
	variables := make(map[string]interface{})
	if !matched {
		variables["failed_field"] = failedCondition
		variables["condition_result"] = false
	} else {
		variables["condition_result"] = true
	}

	return &lua_engine.ExecuteResult{
		Valid:       matched,
		Action:      rule.Action,
		Context:     variables,
		ExecuteTime: 0,
	}, nil
}

// evaluateCondition 评估单个条件
func (s *RuleExecutionService) evaluateCondition(condition map[string]interface{}, data map[string]interface{}) bool {
	// 获取条件字段
	field, ok := condition["field"].(string)
	if !ok || field == "" {
		return false
	}

	// 获取操作符
	operator, ok := condition["operator"].(string)
	if !ok || operator == "" {
		return false
	}

	// 获取目标值
	targetValue, ok := condition["value"]
	if !ok {
		return false
	}

	// 从数据中获取实际值
	actualValue := s.getFieldValue(data, field)
	if actualValue == nil {
		return false
	}

	// 使用条件评估器评估条件
	return s.evaluateConditionWithOperator(operator, targetValue, actualValue)
}

// getFieldValue 从数据中获取字段值，支持嵌套字段
func (s *RuleExecutionService) getFieldValue(data map[string]interface{}, field string) interface{} {
	// 处理嵌套字段，如 "user.age"
	parts := strings.Split(field, ".")

	current := data
	for i, part := range parts {
		if i == len(parts)-1 {
			// 最后一个部分，返回值
			if value, exists := current[part]; exists {
				return value
			}
			return nil
		}

		// 中间部分，继续遍历
		if next, ok := current[part].(map[string]interface{}); ok {
			current = next
		} else {
			return nil
		}
	}

	return nil
}

// evaluateConditionWithOperator 使用操作符评估条件
func (s *RuleExecutionService) evaluateConditionWithOperator(operator string, targetValue, actualValue interface{}) bool {
	switch operator {
	case "eq":
		return s.evaluateEqual(targetValue, actualValue)
	case "neq":
		return !s.evaluateEqual(targetValue, actualValue)
	case "gt":
		return s.compareValues(actualValue, targetValue) > 0
	case "gte":
		return s.compareValues(actualValue, targetValue) >= 0
	case "lt":
		return s.compareValues(actualValue, targetValue) < 0
	case "lte":
		return s.compareValues(actualValue, targetValue) <= 0
	case "in":
		return s.evaluateIn(targetValue, actualValue)
	case "between":
		return s.evaluateBetween(targetValue, actualValue)
	case "contains":
		return s.evaluateContains(targetValue, actualValue)
	case "not_contains":
		return !s.evaluateContains(targetValue, actualValue)
	case "regex":
		return s.evaluateRegex(targetValue, actualValue)
	default:
		return false
	}
}

// evaluateEqual 评估等于
func (s *RuleExecutionService) evaluateEqual(targetValue, actualValue interface{}) bool {
	return fmt.Sprintf("%v", actualValue) == fmt.Sprintf("%v", targetValue)
}

// evaluateIn 评估包含
func (s *RuleExecutionService) evaluateIn(targetValue, actualValue interface{}) bool {
	// 将目标值转换为字符串列表
	targetStr := fmt.Sprintf("%v", targetValue)
	values := strings.Split(targetStr, ",")
	actualStr := fmt.Sprintf("%v", actualValue)

	for _, v := range values {
		if strings.TrimSpace(v) == actualStr {
			return true
		}
	}
	return false
}

// evaluateBetween 评估区间
func (s *RuleExecutionService) evaluateBetween(targetValue, actualValue interface{}) bool {
	// 将目标值解析为区间
	targetStr := fmt.Sprintf("%v", targetValue)
	values := strings.Split(targetStr, ",")
	if len(values) != 2 {
		return false
	}

	min, err1 := strconv.ParseFloat(strings.TrimSpace(values[0]), 64)
	max, err2 := strconv.ParseFloat(strings.TrimSpace(values[1]), 64)
	if err1 != nil || err2 != nil {
		return false
	}

	actualFloat, err := s.toFloat64(actualValue)
	if err != nil {
		return false
	}

	return actualFloat >= min && actualFloat <= max
}

// evaluateContains 评估包含关系
func (s *RuleExecutionService) evaluateContains(targetValue, actualValue interface{}) bool {
	actualStr := fmt.Sprintf("%v", actualValue)
	targetStr := fmt.Sprintf("%v", targetValue)

	return strings.Contains(actualStr, targetStr)
}

// evaluateRegex 评估正则表达式
func (s *RuleExecutionService) evaluateRegex(targetValue, actualValue interface{}) bool {
	// 这里可以添加正则表达式评估逻辑
	// 为了简化，暂时返回false
	return false
}

// compareValues 比较值
func (s *RuleExecutionService) compareValues(actual, target interface{}) int {
	actualFloat, err1 := s.toFloat64(actual)
	if err1 != nil {
		return 0
	}

	targetFloat, err2 := s.toFloat64(target)
	if err2 != nil {
		return 0
	}

	if actualFloat < targetFloat {
		return -1
	} else if actualFloat > targetFloat {
		return 1
	}
	return 0
}

// toFloat64 转换为float64
func (s *RuleExecutionService) toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}

// evaluateLogicalExpression 评估逻辑表达式
func (s *RuleExecutionService) evaluateLogicalExpression(formula string) (interface{}, error) {
	// 支持基本的逻辑表达式
	// 如: true, false, and, or, not

	formula = strings.ToLower(strings.TrimSpace(formula))

	switch formula {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "1":
		return true, nil
	case "0":
		return false, nil
	default:
		// 尝试解析逻辑表达式
		return s.calculateLogicalExpression(formula)
	}
}

// calculateSimpleExpression 计算简单表达式
func (s *RuleExecutionService) calculateSimpleExpression(formula string) (interface{}, error) {
	// 这里实现一个简单的表达式计算器
	// 实际项目中建议使用成熟的表达式解析库

	// 处理比较运算
	if strings.Contains(formula, ">=") {
		return s.evaluateComparison(formula, ">=")
	} else if strings.Contains(formula, "<=") {
		return s.evaluateComparison(formula, "<=")
	} else if strings.Contains(formula, "==") {
		return s.evaluateComparison(formula, "==")
	} else if strings.Contains(formula, "!=") {
		return s.evaluateComparison(formula, "!=")
	} else if strings.Contains(formula, ">") {
		return s.evaluateComparison(formula, ">")
	} else if strings.Contains(formula, "<") {
		return s.evaluateComparison(formula, "<")
	}

	// 处理数学运算
	if strings.Contains(formula, "+") {
		return s.evaluateMathOperation(formula, "+")
	} else if strings.Contains(formula, "-") {
		return s.evaluateMathOperation(formula, "-")
	} else if strings.Contains(formula, "*") {
		return s.evaluateMathOperation(formula, "*")
	} else if strings.Contains(formula, "/") {
		return s.evaluateMathOperation(formula, "/")
	}

	return nil, fmt.Errorf("不支持的表达式: %s", formula)
}

// evaluateComparison 评估比较运算
func (s *RuleExecutionService) evaluateComparison(formula, operator string) (interface{}, error) {
	parts := strings.Split(formula, operator)
	if len(parts) != 2 {
		return nil, fmt.Errorf("无效的比较表达式: %s", formula)
	}

	left, err1 := s.parseValue(strings.TrimSpace(parts[0]))
	right, err2 := s.parseValue(strings.TrimSpace(parts[1]))

	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("解析比较值失败")
	}

	switch operator {
	case ">=":
		return s.compareValues(left, right) >= 0, nil
	case "<=":
		return s.compareValues(left, right) <= 0, nil
	case "==":
		return s.evaluateEqual(left, right), nil
	case "!=":
		return !s.evaluateEqual(left, right), nil
	case ">":
		return s.compareValues(left, right) > 0, nil
	case "<":
		return s.compareValues(left, right) < 0, nil
	default:
		return nil, fmt.Errorf("不支持的比较运算符: %s", operator)
	}
}

// evaluateMathOperation 评估数学运算
func (s *RuleExecutionService) evaluateMathOperation(formula, operator string) (interface{}, error) {
	parts := strings.Split(formula, operator)
	if len(parts) != 2 {
		return nil, fmt.Errorf("无效的数学表达式: %s", formula)
	}

	left, err1 := s.parseValue(strings.TrimSpace(parts[0]))
	right, err2 := s.parseValue(strings.TrimSpace(parts[1]))

	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("解析数学值失败")
	}

	leftFloat, err1 := s.toFloat64(left)
	rightFloat, err2 := s.toFloat64(right)

	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("转换为数字失败")
	}

	switch operator {
	case "+":
		return leftFloat + rightFloat, nil
	case "-":
		return leftFloat - rightFloat, nil
	case "*":
		return leftFloat * rightFloat, nil
	case "/":
		if rightFloat == 0 {
			return nil, fmt.Errorf("除零错误")
		}
		return leftFloat / rightFloat, nil
	default:
		return nil, fmt.Errorf("不支持的数学运算符: %s", operator)
	}
}

// calculateLogicalExpression 计算逻辑表达式
func (s *RuleExecutionService) calculateLogicalExpression(formula string) (interface{}, error) {
	// 简单的逻辑表达式计算
	// 支持: and, or, not

	if strings.Contains(formula, " and ") {
		return s.evaluateLogicalOperation(formula, " and ")
	} else if strings.Contains(formula, " or ") {
		return s.evaluateLogicalOperation(formula, " or ")
	} else if strings.HasPrefix(formula, "not ") {
		return s.evaluateNotOperation(formula)
	}

	return nil, fmt.Errorf("不支持的逻辑表达式: %s", formula)
}

// evaluateLogicalOperation 评估逻辑运算
func (s *RuleExecutionService) evaluateLogicalOperation(formula, operator string) (interface{}, error) {
	parts := strings.Split(formula, operator)
	if len(parts) != 2 {
		return nil, fmt.Errorf("无效的逻辑表达式: %s", formula)
	}

	left, err1 := s.parseValue(strings.TrimSpace(parts[0]))
	right, err2 := s.parseValue(strings.TrimSpace(parts[1]))

	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("解析逻辑值失败")
	}

	leftBool := s.toBool(left)
	rightBool := s.toBool(right)

	switch operator {
	case " and ":
		return leftBool && rightBool, nil
	case " or ":
		return leftBool || rightBool, nil
	default:
		return nil, fmt.Errorf("不支持的逻辑运算符: %s", operator)
	}
}

// evaluateNotOperation 评估非运算
func (s *RuleExecutionService) evaluateNotOperation(formula string) (interface{}, error) {
	if !strings.HasPrefix(formula, "not ") {
		return nil, fmt.Errorf("无效的非运算表达式: %s", formula)
	}

	operand := strings.TrimSpace(strings.TrimPrefix(formula, "not "))
	value, err := s.parseValue(operand)
	if err != nil {
		return nil, err
	}

	return !s.toBool(value), nil
}

// parseValue 解析值
func (s *RuleExecutionService) parseValue(value string) (interface{}, error) {
	// 尝试解析为数字
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num, nil
	}

	// 尝试解析为布尔值
	switch strings.ToLower(value) {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	}

	// 返回字符串
	return value, nil
}

// toBool 转换为布尔值
func (s *RuleExecutionService) toBool(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0
	case string:
		switch strings.ToLower(v) {
		case "true", "1", "yes", "on":
			return true
		case "false", "0", "no", "off":
			return false
		default:
			return len(v) > 0
		}
	default:
		return value != nil
	}
}

// checkBracketMatch 检查括号匹配
func (s *RuleExecutionService) checkBracketMatch(text string) bool {
	stack := make([]rune, 0)
	bracketPairs := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, char := range text {
		switch char {
		case '(', '{', '[':
			stack = append(stack, char)
		case ')', '}', ']':
			if len(stack) == 0 {
				return false
			}
			if stack[len(stack)-1] != bracketPairs[char] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

// sortRulesByPriority 按照优先级排序规则（优先级数字越大优先级越高）
func (s *RuleExecutionService) sortRulesByPriority(rules []*model.Rule) {
	// 使用稳定的排序算法，按照优先级降序排列
	for i := 0; i < len(rules)-1; i++ {
		for j := i + 1; j < len(rules); j++ {
			if rules[i].Priority < rules[j].Priority {
				rules[i], rules[j] = rules[j], rules[i]
			}
		}
	}
}

// updateContextWithRuleResult 将规则执行结果更新到上下文中
func (s *RuleExecutionService) updateContextWithRuleResult(context *model.RuleContext, result *model.RuleResult) *model.RuleContext {
	// 创建新的上下文副本
	newContext := &model.RuleContext{
		Scope:           context.Scope,
		ScopeID:         context.ScopeID,
		TenantID:        context.TenantID,
		ExecutionTiming: context.ExecutionTiming,
		Trigger:         context.Trigger,
		Data:            make(map[string]interface{}),
	}

	// 复制原始数据
	for k, v := range context.Data {
		newContext.Data[k] = v
	}
	return newContext
}

// FormulaResult 公式计算结果
type FormulaResult struct {
	Valid bool        `json:"valid"`
	Value interface{} `json:"value"`
}
