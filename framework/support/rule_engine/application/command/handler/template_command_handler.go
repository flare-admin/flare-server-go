package handler

import (
	"context"
	"fmt"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/command"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/err"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/service"
)

// TemplateCommandHandler 模板命令处理器
type TemplateCommandHandler struct {
	templateService *service.RuleTemplateService
}

// NewTemplateCommandHandler 创建模板命令处理器
func NewTemplateCommandHandler(templateService *service.RuleTemplateService) *TemplateCommandHandler {
	return &TemplateCommandHandler{
		templateService: templateService,
	}
}

// HandleCreateTemplate 处理创建模板命令
func (h *TemplateCommandHandler) HandleCreateTemplate(ctx context.Context, cmd *command.CreateTemplateCommand) *herrors.HError {
	// 验证命令参数
	if err := h.validateCreateTemplateCommand(cmd); err != nil {
		return err
	}

	// 创建模板领域模型
	template := model.NewRuleTemplate(cmd.Code, cmd.Name, cmd.Description, cmd.CategoryID, cmd.Type)
	template.Conditions = cmd.Conditions
	template.LuaScript = cmd.LuaScript
	template.Formula = cmd.Formula
	template.FormulaVars = cmd.FormulaVars
	template.Parameters = cmd.Parameters
	template.Priority = cmd.Priority
	template.Sorting = cmd.Sorting

	// 调用领域服务创建模板
	return h.templateService.CreateTemplate(ctx, template)
}

// HandleUpdateTemplate 处理更新模板命令
func (h *TemplateCommandHandler) HandleUpdateTemplate(ctx context.Context, cmd *command.UpdateTemplateCommand) *herrors.HError {
	// 验证命令参数
	if err := h.validateUpdateTemplateCommand(cmd); err != nil {
		return err
	}

	// 创建模板领域模型
	template := &model.RuleTemplate{
		ID:          cmd.ID,
		Code:        cmd.Code,
		Name:        cmd.Name,
		Description: cmd.Description,
		CategoryID:  cmd.CategoryID,
		Type:        cmd.Type,
		Conditions:  cmd.Conditions,
		LuaScript:   cmd.LuaScript,
		Formula:     cmd.Formula,
		FormulaVars: cmd.FormulaVars,
		Parameters:  cmd.Parameters,
		Priority:    cmd.Priority,
		Sorting:     cmd.Sorting,
		UpdatedAt:   utils.GetDateUnix(),
	}

	// 调用领域服务更新模板
	return h.templateService.UpdateTemplate(ctx, template)
}

// HandleUpdateTemplateStatus 处理更新模板状态命令
func (h *TemplateCommandHandler) HandleUpdateTemplateStatus(ctx context.Context, cmd *command.UpdateTemplateStatusCommand) *herrors.HError {
	// 验证命令参数
	if cmd.ID == "" {
		return err.RuleTemplateValidationFailed(fmt.Errorf("模板ID不能为空"))
	}
	if cmd.Status != 1 && cmd.Status != 2 {
		return err.RuleTemplateValidationFailed(fmt.Errorf("无效的状态值"))
	}

	// 根据状态调用相应的服务方法
	if cmd.Status == 1 {
		return h.templateService.EnableTemplate(ctx, cmd.ID)
	} else {
		return h.templateService.DisableTemplate(ctx, cmd.ID)
	}
}

// HandleDeleteTemplate 处理删除模板命令
func (h *TemplateCommandHandler) HandleDeleteTemplate(ctx context.Context, cmd *command.DeleteTemplateCommand) *herrors.HError {
	// 验证命令参数
	if cmd.ID == "" {
		return err.RuleTemplateValidationFailed(fmt.Errorf("模板ID不能为空"))
	}

	// 调用领域服务删除模板
	return h.templateService.DeleteTemplate(ctx, cmd.ID)
}

// validateCreateTemplateCommand 验证创建模板命令
func (h *TemplateCommandHandler) validateCreateTemplateCommand(cmd *command.CreateTemplateCommand) *herrors.HError {
	if cmd.Code == "" {
		return err.RuleTemplateValidationFailed(fmt.Errorf("模板编码不能为空"))
	}
	if cmd.Name == "" {
		return err.RuleTemplateValidationFailed(fmt.Errorf("模板名称不能为空"))
	}
	//if cmd.CategoryID == "" {
	//	return err.RuleTemplateValidationFailed(fmt.Errorf("分类ID不能为空"))
	//}
	if cmd.Type == "" {
		return err.RuleTemplateValidationFailed(fmt.Errorf("模板类型不能为空"))
	}
	if !h.isValidTemplateType(cmd.Type) {
		return err.RuleTemplateValidationFailed(fmt.Errorf("无效的模板类型: %s", cmd.Type))
	}

	// 根据模板类型验证必填字段
	switch cmd.Type {
	case "condition":
		if cmd.Conditions == "" {
			return err.RuleTemplateValidationFailed(fmt.Errorf("条件模板必须提供条件表达式"))
		}
	case "lua":
		if cmd.LuaScript == "" {
			return err.RuleTemplateValidationFailed(fmt.Errorf("Lua模板必须提供Lua脚本"))
		}
	case "formula":
		if cmd.Formula == "" {
			return err.RuleTemplateValidationFailed(fmt.Errorf("公式模板必须提供计算公式"))
		}
	}

	return nil
}

// validateUpdateTemplateCommand 验证更新模板命令
func (h *TemplateCommandHandler) validateUpdateTemplateCommand(cmd *command.UpdateTemplateCommand) *herrors.HError {
	if cmd.ID == "" {
		return err.RuleTemplateValidationFailed(fmt.Errorf("模板ID不能为空"))
	}
	if cmd.Code == "" {
		return err.RuleTemplateValidationFailed(fmt.Errorf("模板编码不能为空"))
	}
	if cmd.Name == "" {
		return err.RuleTemplateValidationFailed(fmt.Errorf("模板名称不能为空"))
	}
	//if cmd.CategoryID == "" {
	//	return err.RuleTemplateValidationFailed(fmt.Errorf("分类ID不能为空"))
	//}
	if cmd.Type == "" {
		return err.RuleTemplateValidationFailed(fmt.Errorf("模板类型不能为空"))
	}
	if !h.isValidTemplateType(cmd.Type) {
		return err.RuleTemplateValidationFailed(fmt.Errorf("无效的模板类型: %s", cmd.Type))
	}

	// 根据模板类型验证必填字段
	switch cmd.Type {
	case "condition":
		if cmd.Conditions == "" {
			return err.RuleTemplateValidationFailed(fmt.Errorf("条件模板必须提供条件表达式"))
		}
	case "lua":
		if cmd.LuaScript == "" {
			return err.RuleTemplateValidationFailed(fmt.Errorf("Lua模板必须提供Lua脚本"))
		}
	case "formula":
		if cmd.Formula == "" {
			return err.RuleTemplateValidationFailed(fmt.Errorf("公式模板必须提供计算公式"))
		}
	}

	return nil
}

// isValidTemplateType 验证模板类型是否有效
func (h *TemplateCommandHandler) isValidTemplateType(templateType string) bool {
	validTypes := []string{"condition", "lua", "formula"}
	for _, validType := range validTypes {
		if templateType == validType {
			return true
		}
	}
	return false
}
