package err

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
)

var (
	// TemplateCreateFailed 创建模板失败
	TemplateCreateFailed = herrors.NewServerError("TemplateCreateFailed")
	// TemplateUpdateFailed 更新模板失败
	TemplateUpdateFailed = herrors.NewServerError("TemplateUpdateFailed")
	// TemplateDeleteFailed 删除模板失败
	TemplateDeleteFailed = herrors.NewServerError("TemplateDeleteFailed")
	// TemplateGetFailed 获取模板失败
	TemplateGetFailed = herrors.NewServerError("TemplateGetFailed")
	// TemplateDisabled 模板已禁用
	TemplateDisabled = herrors.NewBusinessServerError("TemplateDisabled")
	// TemplateNotExist 模板不存在
	TemplateNotExist = herrors.NewBusinessServerError("TemplateNotExist")
	// CategoryNotExist 分类不存在
	CategoryNotExist = herrors.NewBusinessServerError("CategoryNotExist")

	// CategoryCreateFailed 创建分类失败
	CategoryCreateFailed = herrors.NewServerError("CategoryCreateFailed")
	// CategoryUpdateFailed 更新分类失败
	CategoryUpdateFailed = herrors.NewServerError("CategoryUpdateFailed")
	// CategoryDeleteFailed 删除分类失败
	CategoryDeleteFailed = herrors.NewServerError("CategoryDeleteFailed")
	// CategoryGetFailed 获取分类失败
	CategoryGetFailed = herrors.NewServerError("CategoryGetFailed")
	// CategoryDisabled 分类已禁用
	CategoryDisabled = herrors.NewBusinessServerError("CategoryDisabled")
	// CategoryNotExist 分类不存在
	CategoryCodeExist = herrors.NewBusinessServerError("CategoryCodeExist")

	// TemplateValidationFailed 模板数据验证失败
	TemplateValidationFailed = herrors.NewServerError("TemplateValidationFailed")

	// TemplateCodeExists 模板编码已存在
	TemplateCodeExists = herrors.NewBusinessServerError("TemplateCodeExists")
)
