package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/template/application/command"
	template_err "github.com/flare-admin/flare-server-go/framework/support/template/domain/err"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/service"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/valueobject"
)

// CategoryCommandHandler 分类命令处理器
type CategoryCommandHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryCommandHandler 创建分类命令处理器
func NewCategoryCommandHandler(categoryService *service.CategoryService) *CategoryCommandHandler {
	return &CategoryCommandHandler{
		categoryService: categoryService,
	}
}

// HandleCreateCategory 处理创建分类命令
func (h *CategoryCommandHandler) HandleCreateCategory(ctx context.Context, cmd *command.CreateCategoryCommand) *herrors.HError {
	// 1. 验证命令
	if err := h.validateCreateCategoryCommand(cmd); err != nil {
		return err
	}

	// 2. 转换命令对象
	valueCmd := &valueobject.CreateCategoryCommand{
		Name:        cmd.Name,
		Code:        cmd.Code,
		Description: cmd.Description,
		Sort:        cmd.Sort,
	}

	// 3. 调用领域服务
	if err := h.categoryService.CreateCategory(ctx, valueCmd); herrors.HaveError(err) {
		hlog.CtxErrorf(ctx, "创建分类失败:%v", err)
		return err
	}

	hlog.CtxInfof(ctx, "创建分类成功,名称:%s", cmd.Name)
	return nil
}

// HandleUpdateCategory 处理更新分类命令
func (h *CategoryCommandHandler) HandleUpdateCategory(ctx context.Context, cmd *command.UpdateCategoryCommand) *herrors.HError {
	// 1. 验证命令
	if err := h.validateUpdateCategoryCommand(cmd); err != nil {
		return err
	}

	// 2. 转换命令对象
	valueCmd := &valueobject.UpdateCategoryCommand{
		ID:          cmd.ID,
		Name:        cmd.Name,
		Code:        cmd.Code,
		Description: cmd.Description,
		Sort:        cmd.Sort,
	}

	// 3. 调用领域服务
	if err := h.categoryService.UpdateCategory(ctx, valueCmd); herrors.HaveError(err) {
		hlog.CtxErrorf(ctx, "更新分类失败:%v", err)
		return err
	}

	hlog.CtxInfof(ctx, "更新分类成功,ID:%s", cmd.ID)
	return nil
}

// HandleDeleteCategory 处理删除分类命令
func (h *CategoryCommandHandler) HandleDeleteCategory(ctx context.Context, cmd *command.DeleteCategoryCommand) *herrors.HError {
	// 1. 验证命令
	if err := h.validateDeleteCategoryCommand(cmd); err != nil {
		return err
	}

	// 2. 转换命令对象
	valueCmd := &valueobject.DeleteCategoryCommand{
		ID: cmd.ID,
	}

	// 3. 调用领域服务
	if err := h.categoryService.DeleteCategory(ctx, valueCmd); herrors.HaveError(err) {
		hlog.CtxErrorf(ctx, "删除分类失败:%v", err)
		return err
	}

	hlog.CtxInfof(ctx, "删除分类成功,ID:%s", cmd.ID)
	return nil
}

// HandleUpdateCategoryStatus 处理分类状态变更命令
func (h *CategoryCommandHandler) HandleUpdateCategoryStatus(ctx context.Context, cmd *command.UpdateCategoryStatusCommand) *herrors.HError {
	// 1. 验证命令
	if err := h.validateUpdateCategoryStatusCommand(cmd); err != nil {
		return err
	}

	// 2. 转换命令对象
	valueCmd := &valueobject.UpdateCategoryStatusCommand{
		ID:     cmd.ID,
		Status: cmd.Status,
	}

	// 3. 调用领域服务
	if cmd.Status == 1 {
		if err := h.categoryService.EnableCategory(ctx, valueCmd); herrors.HaveError(err) {
			hlog.CtxErrorf(ctx, "启用分类失败:%v", err)
			return err
		}
	} else {
		if err := h.categoryService.DisableCategory(ctx, valueCmd); herrors.HaveError(err) {
			hlog.CtxErrorf(ctx, "禁用分类失败:%v", err)
			return err
		}
	}

	hlog.CtxInfof(ctx, "启用分类成功,ID:%s", cmd.ID)
	return nil
}

// validateCreateCategoryCommand 验证创建分类命令
func (h *CategoryCommandHandler) validateCreateCategoryCommand(cmd *command.CreateCategoryCommand) *herrors.HError {
	if cmd.Name == "" {
		return template_err.CategoryCreateFailed(fmt.Errorf("分类名称不能为空"))
	}
	if cmd.Code == "" {
		return template_err.CategoryCreateFailed(fmt.Errorf("分类编码不能为空"))
	}
	return nil
}

// validateUpdateCategoryCommand 验证更新分类命令
func (h *CategoryCommandHandler) validateUpdateCategoryCommand(cmd *command.UpdateCategoryCommand) *herrors.HError {
	if cmd.ID == "" {
		return template_err.CategoryUpdateFailed(fmt.Errorf("分类ID不能为空"))
	}
	if cmd.Name == "" {
		return template_err.CategoryUpdateFailed(fmt.Errorf("分类名称不能为空"))
	}
	if cmd.Code == "" {
		return template_err.CategoryUpdateFailed(fmt.Errorf("分类编码不能为空"))
	}
	return nil
}

// validateDeleteCategoryCommand 验证删除分类命令
func (h *CategoryCommandHandler) validateDeleteCategoryCommand(cmd *command.DeleteCategoryCommand) *herrors.HError {
	if cmd.ID == "" {
		return template_err.CategoryDeleteFailed(fmt.Errorf("分类ID不能为空"))
	}
	return nil
}

// validateUpdateCategoryStatusCommand 验证更新分类状态命令
func (h *CategoryCommandHandler) validateUpdateCategoryStatusCommand(cmd *command.UpdateCategoryStatusCommand) *herrors.HError {
	if cmd.ID == "" {
		return template_err.CategoryUpdateFailed(fmt.Errorf("分类ID不能为空"))
	}
	if cmd.Status != 1 && cmd.Status != 2 {
		return template_err.CategoryUpdateFailed(fmt.Errorf("无效的分类状态"))
	}
	return nil
}
