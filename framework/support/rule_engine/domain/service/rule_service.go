package service

import (
	"context"
	"fmt"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/snowflake_id"
	"github.com/flare-admin/flare-server-go/framework/pkg/lua_engine"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	ruleengineerr "github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/err"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/repository"
)

// RuleService 规则领域服务
type RuleService struct {
	ruleRepo     repository.IRuleRepository
	templateRepo repository.ITemplateRepository
	categoryRepo repository.ICategoryRepository
	ruleExecutor *lua_engine.RuleExecutor
	ig           snowflake_id.IIdGenerate
}

// NewRuleService 创建规则服务
func NewRuleService(
	ruleRepo repository.IRuleRepository,
	templateRepo repository.ITemplateRepository,
	categoryRepo repository.ICategoryRepository,
	ruleExecutor *lua_engine.RuleExecutor,
	ig snowflake_id.IIdGenerate,
) *RuleService {
	return &RuleService{
		ruleRepo:     ruleRepo,
		templateRepo: templateRepo,
		categoryRepo: categoryRepo,
		ruleExecutor: ruleExecutor,
		ig:           ig,
	}
}

// CreateRule 创建规则
func (s *RuleService) CreateRule(ctx context.Context, rule *model.Rule) *herrors.HError {
	// 验证规则数据
	if err := rule.Validate(); err != nil {
		return ruleengineerr.RuleValidationFailed(err)
	}

	// 检查编码是否已存在
	exists, err := s.ruleRepo.ExistsByCode(ctx, rule.Code)
	if err != nil {
		return ruleengineerr.RuleGetFailed(err)
	}
	if exists {
		return ruleengineerr.RuleCodeExists
	}

	// 检查分类是否存在
	if rule.CategoryID != "" {
		category, err := s.categoryRepo.FindByID(ctx, rule.CategoryID)
		if err != nil {
			return ruleengineerr.RuleCategoryGetFailed(err)
		}
		if !category.IsEnabled() {
			return ruleengineerr.RuleCategoryDisabled
		}
	}

	// 检查模板是否存在
	if rule.TemplateID != "" {
		template, err := s.templateRepo.FindByID(ctx, rule.TemplateID)
		if err != nil {
			return ruleengineerr.RuleTemplateGetFailed(err)
		}
		if !template.IsEnabled() {
			return ruleengineerr.RuleTemplateDisabled
		}
	}

	// 验证规则内容
	if err := s.validateRuleContent(ctx, rule); err != nil {
		return ruleengineerr.RuleContentInvalid
	}

	// 生成ID
	rule.ID = s.ig.GenStringId()
	rule.Completion()
	// 创建规则
	if err := s.ruleRepo.Create(ctx, rule); err != nil {
		return ruleengineerr.RuleCreateFailed(err)
	}

	return nil
}

// UpdateRule 更新规则
func (s *RuleService) UpdateRule(ctx context.Context, rule *model.Rule) *herrors.HError {
	// 验证规则数据
	if err := rule.Validate(); err != nil {
		return ruleengineerr.RuleValidationFailed(err)
	}

	// 检查规则是否存在
	existingRule, err := s.ruleRepo.FindByID(ctx, rule.ID)
	if err != nil {
		return ruleengineerr.RuleGetFailed(err)
	}
	if !existingRule.IsEnabled() {
		return ruleengineerr.RuleDisabled
	}

	// 检查编码是否重复（排除自己）
	if rule.Code != existingRule.Code {
		exists, err := s.ruleRepo.ExistsByCode(ctx, rule.Code)
		if err != nil {
			return ruleengineerr.RuleGetFailed(err)
		}
		if exists {
			return ruleengineerr.RuleCodeExists
		}
	}

	// 检查分类是否存在
	if rule.CategoryID != "" {
		category, err := s.categoryRepo.FindByID(ctx, rule.CategoryID)
		if err != nil {
			return ruleengineerr.RuleCategoryGetFailed(err)
		}
		if !category.IsEnabled() {
			return ruleengineerr.RuleCategoryDisabled
		}
	}

	// 检查模板是否存在
	if rule.TemplateID != "" {
		template, err := s.templateRepo.FindByID(ctx, rule.TemplateID)
		if err != nil {
			return ruleengineerr.RuleTemplateGetFailed(err)
		}
		if !template.IsEnabled() {
			return ruleengineerr.RuleTemplateDisabled
		}
	}

	// 验证规则内容
	if err := s.validateRuleContent(ctx, rule); err != nil {
		return ruleengineerr.RuleContentInvalid
	}
	rule.Completion()
	// 更新规则
	if err := s.ruleRepo.Update(ctx, rule); err != nil {
		return ruleengineerr.RuleUpdateFailed(err)
	}

	return nil
}

// DeleteRule 删除规则
func (s *RuleService) DeleteRule(ctx context.Context, ruleID string) *herrors.HError {
	// 检查规则是否存在
	_, err := s.ruleRepo.FindByID(ctx, ruleID)
	if err != nil {
		return ruleengineerr.RuleGetFailed(err)
	}

	// 删除规则
	if err := s.ruleRepo.Delete(ctx, ruleID); err != nil {
		return ruleengineerr.RuleDeleteFailed(err)
	}

	return nil
}

// GetRule 获取规则
func (s *RuleService) GetRule(ctx context.Context, ruleID string) (*model.Rule, *herrors.HError) {
	rule, err := s.ruleRepo.FindByID(ctx, ruleID)
	if err != nil {
		return nil, ruleengineerr.RuleGetFailed(err)
	}

	return rule, nil
}

// GetRuleByCode 根据编码获取规则
func (s *RuleService) GetRuleByCode(ctx context.Context, code string) (*model.Rule, *herrors.HError) {
	rule, err := s.ruleRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, ruleengineerr.RuleGetFailed(err)
	}

	return rule, nil
}

// GetRulesByCategory 根据分类获取规则列表
func (s *RuleService) GetRulesByCategory(ctx context.Context, categoryID string) ([]*model.Rule, *herrors.HError) {
	// 检查分类是否存在
	if categoryID != "" {
		_, err := s.categoryRepo.FindByID(ctx, categoryID)
		if err != nil {
			return nil, ruleengineerr.RuleCategoryGetFailed(err)
		}
	}

	rules, err := s.ruleRepo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, ruleengineerr.RuleGetFailed(err)
	}

	return rules, nil
}

// GetRulesByTemplate 根据模板获取规则列表
func (s *RuleService) GetRulesByTemplate(ctx context.Context, templateID string) ([]*model.Rule, *herrors.HError) {
	// 检查模板是否存在
	if templateID != "" {
		_, err := s.templateRepo.FindByID(ctx, templateID)
		if err != nil {
			return nil, ruleengineerr.RuleTemplateGetFailed(err)
		}
	}

	rules, err := s.ruleRepo.FindByTemplateID(ctx, templateID)
	if err != nil {
		return nil, ruleengineerr.RuleGetFailed(err)
	}

	return rules, nil
}

// GetRulesByType 根据类型获取规则列表
func (s *RuleService) GetRulesByType(ctx context.Context, ruleType string) ([]*model.Rule, *herrors.HError) {
	rules, err := s.ruleRepo.FindByType(ctx, ruleType)
	if err != nil {
		return nil, ruleengineerr.RuleGetFailed(err)
	}

	return rules, nil
}

// GetRulesByTrigger 根据触发条件获取规则列表
func (s *RuleService) GetRulesByTrigger(ctx context.Context, trigger string) ([]*model.Rule, *herrors.HError) {
	rules, err := s.ruleRepo.FindByTrigger(ctx, trigger)
	if err != nil {
		return nil, ruleengineerr.RuleGetFailed(err)
	}

	return rules, nil
}

// GetRulesByScope 根据作用域获取规则列表
func (s *RuleService) GetRulesByScope(ctx context.Context, scope string, scopeID string) ([]*model.Rule, *herrors.HError) {
	rules, err := s.ruleRepo.FindByScope(ctx, scope)
	if err != nil {
		return nil, ruleengineerr.RuleGetFailed(err)
	}

	return rules, nil
}

// GetRulesByBusinessType 根据业务类型获取规则列表
func (s *RuleService) GetRulesByBusinessType(ctx context.Context, businessType string) ([]*model.Rule, *herrors.HError) {
	rules, err := s.ruleRepo.FindByBusinessType(ctx, businessType)
	if err != nil {
		return nil, ruleengineerr.RuleGetFailed(err)
	}

	return rules, nil
}

// EnableRule 启用规则
func (s *RuleService) EnableRule(ctx context.Context, ruleID string) *herrors.HError {
	rule, err := s.ruleRepo.FindByID(ctx, ruleID)
	if err != nil {
		return ruleengineerr.RuleGetFailed(err)
	}

	rule.Enable()
	if err := s.ruleRepo.Update(ctx, rule); err != nil {
		return ruleengineerr.RuleUpdateFailed(err)
	}

	return nil
}

// DisableRule 禁用规则
func (s *RuleService) DisableRule(ctx context.Context, ruleID string) *herrors.HError {
	rule, err := s.ruleRepo.FindByID(ctx, ruleID)
	if err != nil {
		return ruleengineerr.RuleGetFailed(err)
	}

	rule.Disable()
	if err := s.ruleRepo.Update(ctx, rule); err != nil {
		return ruleengineerr.RuleUpdateFailed(err)
	}

	return nil
}

// ValidateRule 验证规则
func (s *RuleService) ValidateRule(ctx context.Context, rule *model.Rule) *herrors.HError {
	if err := rule.Validate(); err != nil {
		return ruleengineerr.RuleValidationFailed(err)
	}

	// 验证规则内容
	if err := s.validateRuleContent(ctx, rule); err != nil {
		return ruleengineerr.RuleContentInvalid
	}

	return nil
}

// validateRuleContent 验证规则内容
func (s *RuleService) validateRuleContent(ctx context.Context, rule *model.Rule) error {
	switch rule.Type {
	case "condition":
		// 验证条件规则
		//if rule.Conditions == nil || len(rule.Conditions) == 0 {
		//	return fmt.Errorf("condition rule must have conditions configuration")
		//}
		//// 验证条件配置
		//for i, condition := range rule.Conditions {
		//	if err := s.validateCondition(condition); err != nil {
		//		return fmt.Errorf("condition %d validation failed: %w", i, err)
		//	}
		//}
		return nil
	case "lua":
		// 验证Lua脚本规则
		if rule.LuaScript == "" {
			return fmt.Errorf("lua rule must have lua script")
		}
		// 简单验证Lua脚本格式，不执行脚本
		if err := s.ruleExecutor.ValidateScript(rule.LuaScript); err != nil {
			return fmt.Errorf("lua script validation failed: %w", err)
		}
		return nil
	case "formula":
		// 验证公式规则
		if rule.Formula == "" {
			return fmt.Errorf("formula rule must have formula")
		}
		// 验证公式语法
		if err := s.validateFormula(rule.Formula); err != nil {
			return fmt.Errorf("formula validation failed: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported rule type: %s", rule.Type)
	}
}

// validateCondition 验证条件配置
func (s *RuleService) validateCondition(condition map[string]interface{}) error {
	// 检查必要字段
	field, ok := condition["field"].(string)
	if !ok || field == "" {
		return fmt.Errorf("condition must have valid field")
	}

	operator, ok := condition["operator"].(string)
	if !ok || operator == "" {
		return fmt.Errorf("condition must have valid operator")
	}

	// 验证操作符
	validOperators := []string{"eq", "ne", "gt", "gte", "lt", "lte", "in", "nin", "contains", "not_contains", "regex"}
	isValidOperator := false
	for _, validOp := range validOperators {
		if operator == validOp {
			isValidOperator = true
			break
		}
	}
	if !isValidOperator {
		return fmt.Errorf("invalid operator: %s", operator)
	}

	// 检查值字段
	if _, ok := condition["value"]; !ok {
		return fmt.Errorf("condition must have value")
	}

	return nil
}

// validateFormula 验证公式语法
func (s *RuleService) validateFormula(formula string) error {
	// 基本语法检查
	if len(formula) == 0 {
		return fmt.Errorf("formula cannot be empty")
	}

	// 检查基本数学运算符
	validOperators := []string{"+", "-", "*", "/", "(", ")", "=", ">", "<", ">=", "<=", "!=", "&&", "||"}
	for _, operator := range validOperators {
		if contains(formula, operator) {
			// 检查运算符前后是否有操作数
			if !s.checkOperatorContext(formula, operator) {
				return fmt.Errorf("invalid operator context for '%s'", operator)
			}
		}
	}

	// 检查括号匹配
	if !s.checkBracketMatch(formula) {
		return fmt.Errorf("formula has unmatched brackets")
	}

	// 检查变量引用格式
	if err := s.checkVariableReferences(formula); err != nil {
		return fmt.Errorf("invalid variable reference: %w", err)
	}

	return nil
}

// checkBracketMatch 检查括号匹配
func (s *RuleService) checkBracketMatch(text string) bool {
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

// checkOperatorContext 检查运算符上下文
func (s *RuleService) checkOperatorContext(formula, operator string) bool {
	// 简单的上下文检查，确保运算符前后有字符
	index := 0
	for {
		pos := indexOf(formula[index:], operator)
		if pos == -1 {
			break
		}
		actualPos := index + pos

		// 检查运算符前是否有字符
		if actualPos > 0 && formula[actualPos-1] == ' ' {
			// 运算符前是空格，需要检查更前面的字符
			prevChar := actualPos - 2
			if prevChar < 0 || formula[prevChar] == ' ' {
				return false
			}
		}

		// 检查运算符后是否有字符
		if actualPos+len(operator) < len(formula) && formula[actualPos+len(operator)] == ' ' {
			// 运算符后是空格，需要检查更后面的字符
			nextChar := actualPos + len(operator) + 1
			if nextChar >= len(formula) || formula[nextChar] == ' ' {
				return false
			}
		}

		index = actualPos + len(operator)
	}

	return true
}

// checkVariableReferences 检查变量引用
func (s *RuleService) checkVariableReferences(formula string) error {
	// 检查变量引用格式，例如 ${varName} 或 $varName
	// 这里实现简单的检查逻辑
	if contains(formula, "${") && !contains(formula, "}") {
		return fmt.Errorf("unclosed variable reference")
	}

	return nil
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return indexOf(s, substr) != -1
}

// indexOf 查找子字符串的位置
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
