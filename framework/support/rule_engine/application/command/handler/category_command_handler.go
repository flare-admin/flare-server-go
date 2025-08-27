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

// CategoryCommandHandler 分类命令处理器
type CategoryCommandHandler struct {
	categoryService *service.RuleCategoryService
}

// NewCategoryCommandHandler 创建分类命令处理器
func NewCategoryCommandHandler(categoryService *service.RuleCategoryService) *CategoryCommandHandler {
	return &CategoryCommandHandler{
		categoryService: categoryService,
	}
}

// HandleCreateCategory 处理创建分类命令
func (h *CategoryCommandHandler) HandleCreateCategory(ctx context.Context, cmd *command.CreateCategoryCommand) *herrors.HError {
	// 验证命令参数
	if err := h.validateCreateCategoryCommand(cmd); err != nil {
		return err
	}

	// 创建分类领域模型
	category := model.NewRuleCategory(cmd.Code, cmd.Name, cmd.Description, cmd.Type, cmd.BusinessType)

	// 设置父分类ID
	if cmd.ParentID != "" {
		category.ParentID = cmd.ParentID
	}

	// 设置排序
	category.SetSorting(cmd.Sorting)

	// 调用领域服务创建分类
	return h.categoryService.CreateCategory(ctx, category)
}

// HandleUpdateCategory 处理更新分类命令
func (h *CategoryCommandHandler) HandleUpdateCategory(ctx context.Context, cmd *command.UpdateCategoryCommand) *herrors.HError {
	// 验证命令参数
	if err := h.validateUpdateCategoryCommand(cmd); err != nil {
		return err
	}

	// 创建分类领域模型
	category := &model.RuleCategory{
		ID:           cmd.ID,
		Code:         cmd.Code,
		Name:         cmd.Name,
		Description:  cmd.Description,
		ParentID:     cmd.ParentID,
		Type:         cmd.Type,
		BusinessType: cmd.BusinessType,
		Sorting:      cmd.Sorting,
		UpdatedAt:    utils.GetDateUnix(),
	}

	// 调用领域服务更新分类
	return h.categoryService.UpdateCategory(ctx, category)
}

// HandleUpdateCategoryStatus 处理更新分类状态命令
func (h *CategoryCommandHandler) HandleUpdateCategoryStatus(ctx context.Context, cmd *command.UpdateCategoryStatusCommand) *herrors.HError {
	// 验证命令参数
	if cmd.ID == "" {
		return err.RuleCategoryValidationFailed(fmt.Errorf("分类ID不能为空"))
	}
	if cmd.Status != 1 && cmd.Status != 2 {
		return err.RuleCategoryValidationFailed(fmt.Errorf("无效的状态值"))
	}

	// 根据状态调用相应的服务方法
	if cmd.Status == 1 {
		return h.categoryService.EnableCategory(ctx, cmd.ID)
	} else {
		return h.categoryService.DisableCategory(ctx, cmd.ID)
	}
}

// HandleDeleteCategory 处理删除分类命令
func (h *CategoryCommandHandler) HandleDeleteCategory(ctx context.Context, cmd *command.DeleteCategoryCommand) *herrors.HError {
	// 验证命令参数
	if cmd.ID == "" {
		return err.RuleCategoryValidationFailed(fmt.Errorf("分类ID不能为空"))
	}

	// 调用领域服务删除分类
	return h.categoryService.DeleteCategory(ctx, cmd.ID)
}

// validateCreateCategoryCommand 验证创建分类命令
func (h *CategoryCommandHandler) validateCreateCategoryCommand(cmd *command.CreateCategoryCommand) *herrors.HError {
	if cmd.Code == "" {
		return err.RuleCategoryValidationFailed(fmt.Errorf("分类编码不能为空"))
	}
	if cmd.Name == "" {
		return err.RuleCategoryValidationFailed(fmt.Errorf("分类名称不能为空"))
	}
	if cmd.Type == "" {
		return err.RuleCategoryValidationFailed(fmt.Errorf("分类类型不能为空"))
	}
	if cmd.BusinessType == "" {
		return err.RuleCategoryValidationFailed(fmt.Errorf("业务类型不能为空"))
	}

	// 验证分类类型
	if !h.isValidCategoryType(cmd.Type) {
		return err.RuleCategoryValidationFailed(fmt.Errorf("无效的分类类型: %s", cmd.Type))
	}

	return nil
}

// validateUpdateCategoryCommand 验证更新分类命令
func (h *CategoryCommandHandler) validateUpdateCategoryCommand(cmd *command.UpdateCategoryCommand) *herrors.HError {
	if cmd.ID == "" {
		return err.RuleCategoryValidationFailed(fmt.Errorf("分类ID不能为空"))
	}
	if cmd.Code == "" {
		return err.RuleCategoryValidationFailed(fmt.Errorf("分类编码不能为空"))
	}
	if cmd.Name == "" {
		return err.RuleCategoryValidationFailed(fmt.Errorf("分类名称不能为空"))
	}
	if cmd.Type == "" {
		return err.RuleCategoryValidationFailed(fmt.Errorf("分类类型不能为空"))
	}
	if cmd.BusinessType == "" {
		return err.RuleCategoryValidationFailed(fmt.Errorf("业务类型不能为空"))
	}

	// 验证分类类型
	if !h.isValidCategoryType(cmd.Type) {
		return err.RuleCategoryValidationFailed(fmt.Errorf("无效的分类类型: %s", cmd.Type))
	}

	return nil
}

// isValidCategoryType 验证分类类型是否有效
func (h *CategoryCommandHandler) isValidCategoryType(categoryType string) bool {
	validTypes := []string{"system", "business", "custom"}
	for _, validType := range validTypes {
		if categoryType == validType {
			return true
		}
	}
	return false
}
