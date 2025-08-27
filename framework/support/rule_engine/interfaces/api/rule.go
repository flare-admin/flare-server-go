package ruleapi

import (
	"context"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/service"
)

type IRuleEngineService interface {
	// Validate 验证规则
	Validate(ctx context.Context, params *model.RuleContext) *herrors.HError
	// Execute 执行规则，返回计算结果
	Execute(ctx context.Context, params *model.RuleContext) (*model.RuleResult, *herrors.HError)
	// Exists 检查规则是否存在
	Exists(ctx context.Context, ruleID string) (bool, *herrors.HError)
	// GetById 获取规则详情
	GetById(ctx context.Context, ruleID string) (*model.Rule, *herrors.HError)
}

type RuleEngineService struct {
	ruleService          *service.RuleService
	ruleExecutionService *service.RuleExecutionService
}

func NewRuleEngineService(
	ruleService *service.RuleService,
	ruleExecutionService *service.RuleExecutionService,
) IRuleEngineService {
	return &RuleEngineService{
		ruleService:          ruleService,
		ruleExecutionService: ruleExecutionService,
	}
}

// Validate 验证规则上下文
func (r RuleEngineService) Validate(ctx context.Context, params *model.RuleContext) *herrors.HError {
	// 验证上下文参数
	if err := params.Validate(); err != nil {
		return herrors.NewServerError("RuleContextInvalid")(err)
	}

	// 尝试执行规则来验证上下文是否有效
	// 这里我们使用一个轻量级的验证方式，只检查是否存在匹配的规则
	var hasMatchingRules bool

	// 根据业务类型查找规则
	if params.Scope != "" {
		rules, err := r.ruleService.GetRulesByScope(ctx, params.Scope, params.ScopeID)
		if err == nil && len(rules) > 0 {
			hasMatchingRules = true
		}
	}

	// 根据触发动作查找规则
	if !hasMatchingRules && params.Trigger != "" {
		rules, err := r.ruleService.GetRulesByTrigger(ctx, params.Trigger)
		if err == nil && len(rules) > 0 {
			hasMatchingRules = true
		}
	}

	// 如果没有找到匹配的规则，返回验证失败
	if !hasMatchingRules {
		return herrors.NewBusinessServerError("RuleEngineNoRulesFound")
	}

	return nil
}

// Execute 执行规则，返回计算结果
func (r RuleEngineService) Execute(ctx context.Context, params *model.RuleContext) (*model.RuleResult, *herrors.HError) {
	// 使用规则执行服务执行规则
	result, err := r.ruleExecutionService.ExecuteRules(ctx, params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Exists 检查规则是否存在
func (r RuleEngineService) Exists(ctx context.Context, ruleID string) (bool, *herrors.HError) {
	// 使用规则服务检查规则是否存在
	rule, err := r.ruleService.GetRule(ctx, ruleID)
	if err != nil {
		return false, err
	}

	return rule != nil, nil
}

// GetById 获取规则详情
func (r RuleEngineService) GetById(ctx context.Context, ruleID string) (*model.Rule, *herrors.HError) {
	// 使用规则服务获取规则详情
	rule, err := r.ruleService.GetRule(ctx, ruleID)
	if err != nil {
		return nil, err
	}

	return rule, nil
}
