package handler

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/template/application/command"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/constants"
	template_err "github.com/flare-admin/flare-server-go/framework/support/template/domain/err"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/service"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/valueobject"
)

// TemplateCommandHandler 模板命令处理器
type TemplateCommandHandler struct {
	templateService *service.TemplateService
}

// NewTemplateCommandHandler 创建模板命令处理器
func NewTemplateCommandHandler(templateService *service.TemplateService) *TemplateCommandHandler {
	return &TemplateCommandHandler{
		templateService: templateService,
	}
}

// HandleCreateTemplate 处理创建模板命令
func (h *TemplateCommandHandler) HandleCreateTemplate(ctx context.Context, cmd *command.CreateTemplateCommand) *herrors.HError {
	// 1. 验证命令
	if err := h.validateCreateTemplateCommand(cmd); err != nil {
		return err
	}

	// 2. 转换命令对象
	valueCmd := &valueobject.CreateTemplateCommand{
		Code:        cmd.Code,
		Name:        cmd.Name,
		Description: cmd.Description,
		CategoryID:  cmd.CategoryID,
		Attributes:  convertToValueObjectAttributes(cmd.Attributes),
	}

	// 3. 调用领域服务
	if err := h.templateService.CreateTemplate(ctx, valueCmd); herrors.HaveError(err) {
		hlog.CtxErrorf(ctx, "创建模板失败:%v", err)
		return err
	}

	hlog.CtxInfof(ctx, "创建模板成功,编码:%s", cmd.Code)
	return nil
}

// HandleUpdateTemplate 处理更新模板命令
func (h *TemplateCommandHandler) HandleUpdateTemplate(ctx context.Context, cmd *command.UpdateTemplateCommand) *herrors.HError {
	// 1. 验证命令
	if err := h.validateUpdateTemplateCommand(cmd); err != nil {
		return err
	}

	// 2. 转换命令对象
	valueCmd := &valueobject.UpdateTemplateCommand{
		ID:          cmd.ID,
		Code:        cmd.Code,
		Name:        cmd.Name,
		Description: cmd.Description,
		CategoryID:  cmd.CategoryID,
		Attributes:  convertToValueObjectAttributes(cmd.Attributes),
	}

	// 3. 调用领域服务
	if err := h.templateService.UpdateTemplate(ctx, valueCmd); herrors.HaveError(err) {
		hlog.CtxErrorf(ctx, "更新模板失败:%v", err)
		return err
	}

	hlog.CtxInfof(ctx, "更新模板成功,编码:%s", cmd.Code)
	return nil
}

// HandleDeleteTemplate 处理删除模板命令
func (h *TemplateCommandHandler) HandleDeleteTemplate(ctx context.Context, cmd *command.DeleteTemplateCommand) *herrors.HError {
	// 1. 验证命令
	if err := h.validateDeleteTemplateCommand(cmd); err != nil {
		return err
	}

	// 2. 转换命令对象
	valueCmd := &valueobject.DeleteTemplateCommand{
		ID: cmd.ID,
	}

	// 3. 调用领域服务
	if err := h.templateService.DeleteTemplate(ctx, valueCmd); herrors.HaveError(err) {
		hlog.CtxErrorf(ctx, "删除模板失败:%v", err)
		return err
	}

	hlog.CtxInfof(ctx, "删除模板成功,ID:%s", cmd.ID)
	return nil
}

// HandleUpdateTemplateStatus 处理更新模板状态命令
func (h *TemplateCommandHandler) HandleUpdateTemplateStatus(ctx context.Context, cmd *command.UpdateTemplateStatusCommand) *herrors.HError {
	// 1. 验证命令
	if err := h.validateUpdateTemplateStatusCommand(cmd); err != nil {
		return err
	}

	// 2. 转换命令对象
	valueCmd := &valueobject.UpdateTemplateStatusCommand{
		ID:     cmd.ID,
		Status: cmd.Status,
	}

	// 3. 调用领域服务
	if cmd.Status == 1 {
		if err := h.templateService.EnableTemplate(ctx, valueCmd); herrors.HaveError(err) {
			hlog.CtxErrorf(ctx, "启用模板失败:%v", err)
			return err
		}
		hlog.CtxInfof(ctx, "启用模板成功,ID:%s", cmd.ID)
	} else {
		if err := h.templateService.DisableTemplate(ctx, valueCmd); herrors.HaveError(err) {
			hlog.CtxErrorf(ctx, "禁用模板失败:%v", err)
			return err
		}
		hlog.CtxInfof(ctx, "禁用模板成功,ID:%s", cmd.ID)
	}

	return nil
}

// validateCreateTemplateCommand 验证创建模板命令
func (h *TemplateCommandHandler) validateCreateTemplateCommand(cmd *command.CreateTemplateCommand) *herrors.HError {
	if cmd.Code == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("模板编码不能为空"))
	}
	if cmd.Name == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("模板名称不能为空"))
	}
	if cmd.CategoryID == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("分类ID不能为空"))
	}
	if len(cmd.Attributes) == 0 {
		return template_err.TemplateCreateFailed(fmt.Errorf("模板属性不能为空"))
	}
	for _, attr := range cmd.Attributes {
		if err := h.validateTemplateAttribute(attr); err != nil {
			return err
		}
	}
	return nil
}

// validateUpdateTemplateCommand 验证更新模板命令
func (h *TemplateCommandHandler) validateUpdateTemplateCommand(cmd *command.UpdateTemplateCommand) *herrors.HError {
	if cmd.ID == "" {
		return template_err.TemplateUpdateFailed(fmt.Errorf("模板ID不能为空"))
	}
	if cmd.Code == "" {
		return template_err.TemplateUpdateFailed(fmt.Errorf("模板编码不能为空"))
	}
	if cmd.Name == "" {
		return template_err.TemplateUpdateFailed(fmt.Errorf("模板名称不能为空"))
	}
	if cmd.CategoryID == "" {
		return template_err.TemplateUpdateFailed(fmt.Errorf("分类ID不能为空"))
	}
	if len(cmd.Attributes) == 0 {
		return template_err.TemplateUpdateFailed(fmt.Errorf("模板属性不能为空"))
	}
	for _, attr := range cmd.Attributes {
		if err := h.validateTemplateAttribute(attr); err != nil {
			return err
		}
	}
	return nil
}

// validateDeleteTemplateCommand 验证删除模板命令
func (h *TemplateCommandHandler) validateDeleteTemplateCommand(cmd *command.DeleteTemplateCommand) *herrors.HError {
	if cmd.ID == "" {
		return template_err.TemplateDeleteFailed(fmt.Errorf("模板ID不能为空"))
	}
	return nil
}

// validateUpdateTemplateStatusCommand 验证更新模板状态命令
func (h *TemplateCommandHandler) validateUpdateTemplateStatusCommand(cmd *command.UpdateTemplateStatusCommand) *herrors.HError {
	if cmd.ID == "" {
		return template_err.TemplateUpdateFailed(fmt.Errorf("模板ID不能为空"))
	}
	if cmd.Status != 1 && cmd.Status != 2 {
		return template_err.TemplateUpdateFailed(fmt.Errorf("无效的模板状态"))
	}
	return nil
}

// validateTemplateAttribute 验证模板属性
func (h *TemplateCommandHandler) validateTemplateAttribute(attr command.TemplateAttribute) *herrors.HError {
	if attr.Key == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("属性键不能为空"))
	}
	if attr.Name == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("属性名称不能为空"))
	}
	if attr.Type == "" {
		return template_err.TemplateCreateFailed(fmt.Errorf("属性类型不能为空"))
	}
	if !constants.IsValidAttributeType(attr.Type) {
		return template_err.TemplateCreateFailed(fmt.Errorf("无效的属性类型: %s", attr.Type))
	}
	if attr.Type == constants.AttributeTypeSelect || attr.Type == constants.AttributeTypeSwitch {
		if len(attr.Options) == 0 {
			return template_err.TemplateCreateFailed(fmt.Errorf("选项类型必须提供选项列表"))
		}
		for _, opt := range attr.Options {
			if opt.Label == "" {
				return template_err.TemplateCreateFailed(fmt.Errorf("选项标签不能为空"))
			}
			if opt.Value == nil {
				return template_err.TemplateCreateFailed(fmt.Errorf("选项值不能为空"))
			}
		}
	}
	return nil
}

// convertToValueObjectAttributes 将应用层属性转换为值对象属性
func convertToValueObjectAttributes(attrs []command.TemplateAttribute) []valueobject.TemplateAttribute {
	valueAttrs := make([]valueobject.TemplateAttribute, len(attrs))
	for i, attr := range attrs {
		options := make([]valueobject.Option, len(attr.Options))
		for j, opt := range attr.Options {
			options[j] = valueobject.Option{
				Label: opt.Label,
				Value: opt.Value,
				Sort:  opt.Sort,
			}
		}

		valueAttrs[i] = valueobject.TemplateAttribute{
			Key:      attr.Key,
			Name:     attr.Name,
			Type:     attr.Type,
			Required: attr.Required,
			I18nKey:  attr.I18nKey,
			Options:  options,
			IsQuery:  attr.IsQuery,
			Default:  attr.Default,
			Validation: valueobject.Validation{
				Min:     attr.Validation.Min,
				Max:     attr.Validation.Max,
				Pattern: attr.Validation.Pattern,
				Length:  attr.Validation.Length,
			},
			Description: attr.Description,
		}
	}
	return valueAttrs
}
