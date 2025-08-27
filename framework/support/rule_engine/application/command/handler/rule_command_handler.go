package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/command"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/err"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/service"
)

// RuleCommandHandler 规则命令处理器
type RuleCommandHandler struct {
	ruleService *service.RuleService
}

// NewRuleCommandHandler 创建规则命令处理器
func NewRuleCommandHandler(ruleService *service.RuleService) *RuleCommandHandler {
	return &RuleCommandHandler{
		ruleService: ruleService,
	}
}

// HandleCreateRule 处理创建规则命令
func (h *RuleCommandHandler) HandleCreateRule(ctx context.Context, cmd *command.CreateRuleCommand) *herrors.HError {
	// 验证命令参数
	if err := h.validateCreateRuleCommand(cmd); err != nil {
		return err
	}

	// 创建规则领域模型
	rule := model.NewRule(cmd.Code, cmd.Name, cmd.Description, cmd.CategoryID, cmd.Type)

	// 设置模板ID
	if cmd.TemplateID != "" {
		rule.SetTemplate(cmd.TemplateID)
	}

	err2 := rule.SetTriggers(cmd.Triggers)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "Failed to set trigger: %v", err2)
		return herrors.NewBadReqHError(err2)
	}

	// 设置作用域
	err2 = rule.SetScope(cmd.Scope, cmd.ScopeID)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "Failed to set scope: %v", err2)
		return herrors.NewBadReqHError(err2)
	}

	// 设置条件表达式
	if cmd.Condition != nil {
		conditions := map[string]interface{}{
			"type":       cmd.Condition.Type,
			"expression": cmd.Condition.Expression,
			"parameters": cmd.Condition.Parameters,
		}
		rule.SetConditions(conditions)
	}

	// 设置Lua脚本
	if cmd.LuaScript != "" {
		rule.SetLuaScript(cmd.LuaScript)
	}
	// 设置执行时机
	err2 = rule.SetExecutionTiming(cmd.ExecutionTiming)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "Failed to set execution timing: %v", err2)
		return herrors.NewBadReqHError(err2)
	}

	// 设置计算公式
	if cmd.Formula != "" {
		rule.SetFormula(cmd.Formula, make(map[string]interface{}))
	}

	// 设置动作
	if cmd.Action != "" {
		rule.SetAction(cmd.Action)
	}

	// 设置优先级和排序
	rule.SetPriority(cmd.Priority)
	rule.SetSorting(cmd.Sorting)

	// 调用领域服务创建规则
	return h.ruleService.CreateRule(ctx, rule)
}

// HandleUpdateRule 处理更新规则命令
func (h *RuleCommandHandler) HandleUpdateRule(ctx context.Context, cmd *command.UpdateRuleCommand) *herrors.HError {
	// 验证命令参数
	if err := h.validateUpdateRuleCommand(cmd); err != nil {
		return err
	}

	// 先获取现有规则
	existingRule, err := h.ruleService.GetRule(ctx, cmd.ID)
	if err != nil {
		return err
	}

	// 更新规则字段
	existingRule.Update(cmd.Name, cmd.Description)

	// 设置模板ID
	if cmd.TemplateID != "" {
		existingRule.SetTemplate(cmd.TemplateID)
	}

	// 设置触发动作
	err2 := existingRule.SetTriggers(cmd.Triggers)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "Failed to set trigger: %v", err2)
		return herrors.NewBadReqHError(err2)
	}
	// 设置作用域
	err2 = existingRule.SetScope(cmd.Scope, cmd.ScopeID)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "Failed to set scope: %v", err2)
		return herrors.NewBadReqHError(err2)
	}

	// 设置执行时机
	err2 = existingRule.SetExecutionTiming(cmd.ExecutionTiming)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "Failed to set execution timing: %v", err2)
		return herrors.NewBadReqHError(err2)
	}

	// 设置条件表达式
	if cmd.Condition != nil {
		conditions := map[string]interface{}{
			"type":       cmd.Condition.Type,
			"expression": cmd.Condition.Expression,
			"parameters": cmd.Condition.Parameters,
		}
		existingRule.SetConditions(conditions)
	}

	// 设置Lua脚本
	if cmd.LuaScript != "" {
		existingRule.SetLuaScript(cmd.LuaScript)
	}

	// 设置计算公式
	if cmd.Formula != "" {
	}

	// 设置动作
	if cmd.Action != "" {
		existingRule.SetAction(cmd.Action)
	}

	// 设置优先级和排序
	existingRule.SetPriority(cmd.Priority)
	existingRule.SetSorting(cmd.Sorting)

	// 调用领域服务更新规则
	return h.ruleService.UpdateRule(ctx, existingRule)
}

// HandleUpdateRuleStatus 处理更新规则状态命令
func (h *RuleCommandHandler) HandleUpdateRuleStatus(ctx context.Context, cmd *command.UpdateRuleStatusCommand) *herrors.HError {
	// 验证命令参数
	if cmd.ID == "" {
		return err.RuleValidationFailed(fmt.Errorf("规则ID不能为空"))
	}
	if cmd.Status != 1 && cmd.Status != 2 {
		return err.RuleValidationFailed(fmt.Errorf("无效的状态值"))
	}

	// 根据状态调用相应的服务方法
	if cmd.Status == 1 {
		return h.ruleService.EnableRule(ctx, cmd.ID)
	} else {
		return h.ruleService.DisableRule(ctx, cmd.ID)
	}
}

// HandleDeleteRule 处理删除规则命令
func (h *RuleCommandHandler) HandleDeleteRule(ctx context.Context, cmd *command.DeleteRuleCommand) *herrors.HError {
	// 验证命令参数
	if cmd.ID == "" {
		return err.RuleValidationFailed(fmt.Errorf("规则ID不能为空"))
	}

	// 调用领域服务删除规则
	return h.ruleService.DeleteRule(ctx, cmd.ID)
}

// validateCreateRuleCommand 验证创建规则命令
func (h *RuleCommandHandler) validateCreateRuleCommand(cmd *command.CreateRuleCommand) *herrors.HError {
	if cmd.Code == "" {
		return err.RuleValidationFailed(fmt.Errorf("规则编码不能为空"))
	}
	if cmd.Name == "" {
		return err.RuleValidationFailed(fmt.Errorf("规则名称不能为空"))
	}
	if cmd.Type == "" {
		return err.RuleValidationFailed(fmt.Errorf("规则类型不能为空"))
	}

	// 验证规则类型
	if !h.isValidRuleType(cmd.Type) {
		return err.RuleValidationFailed(fmt.Errorf("无效的规则类型: %s", cmd.Type))
	}

	// 根据规则类型验证必填字段
	switch cmd.Type {
	case "condition":
		if cmd.Condition == nil {
			return err.RuleValidationFailed(fmt.Errorf("条件规则必须提供条件表达式"))
		}
	case "lua":
		if cmd.LuaScript == "" {
			return err.RuleValidationFailed(fmt.Errorf("Lua规则必须提供Lua脚本"))
		}
	case "formula":
		if cmd.Formula == "" {
			return err.RuleValidationFailed(fmt.Errorf("公式规则必须提供计算公式"))
		}
	}

	// 验证动作配置
	if len(cmd.Triggers) == 0 {
		return err.RuleValidationFailed(fmt.Errorf("规则必须提供动作配置"))
	}

	return nil
}

// validateUpdateRuleCommand 验证更新规则命令
func (h *RuleCommandHandler) validateUpdateRuleCommand(cmd *command.UpdateRuleCommand) *herrors.HError {
	if cmd.ID == "" {
		return err.RuleValidationFailed(fmt.Errorf("规则ID不能为空"))
	}
	if cmd.Code == "" {
		return err.RuleValidationFailed(fmt.Errorf("规则编码不能为空"))
	}
	if cmd.Name == "" {
		return err.RuleValidationFailed(fmt.Errorf("规则名称不能为空"))
	}
	if cmd.Type == "" {
		return err.RuleValidationFailed(fmt.Errorf("规则类型不能为空"))
	}
	// 验证规则类型
	if !h.isValidRuleType(cmd.Type) {
		return err.RuleValidationFailed(fmt.Errorf("无效的规则类型: %s", cmd.Type))
	}

	// 根据规则类型验证必填字段
	switch cmd.Type {
	case "condition":
		if cmd.Condition == nil {
			return err.RuleValidationFailed(fmt.Errorf("条件规则必须提供条件表达式"))
		}
	case "lua":
		if cmd.LuaScript == "" {
			return err.RuleValidationFailed(fmt.Errorf("Lua规则必须提供Lua脚本"))
		}
	case "formula":
		if cmd.Formula == "" {
			return err.RuleValidationFailed(fmt.Errorf("公式规则必须提供计算公式"))
		}
	}

	// 验证动作配置
	if len(cmd.Triggers) == 0 {
		return err.RuleValidationFailed(fmt.Errorf("规则必须提供动作配置"))
	}

	return nil
}

// isValidRuleType 验证规则类型是否有效
func (h *RuleCommandHandler) isValidRuleType(ruleType string) bool {
	validTypes := []string{"condition", "lua", "formula"}
	for _, validType := range validTypes {
		if ruleType == validType {
			return true
		}
	}
	return false
}
