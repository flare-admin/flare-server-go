package templateapi

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	template_err "github.com/flare-admin/flare-server-go/framework/support/template/domain/err"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/service"
)

type ITemplateService interface {
	// ValidateContent 验证模板内容
	ValidateContent(ctx context.Context, templateID string, content map[string]interface{}) *herrors.HError
	GetById(ctx context.Context, templateID string) (*model.Template, *herrors.HError)
	// Exists 检查模板是否存在
	Exists(ctx context.Context, templateID string) (bool, *herrors.HError)
}

type TemplateService struct {
	templateService *service.TemplateService
}

func NewTemplateService(templateService *service.TemplateService) ITemplateService {
	return &TemplateService{
		templateService: templateService,
	}
}

func (t TemplateService) ValidateContent(ctx context.Context, templateID string, content map[string]interface{}) *herrors.HError {
	return t.templateService.ValidateTemplateContent(ctx, templateID, content)
}

func (t TemplateService) GetById(ctx context.Context, templateID string) (*model.Template, *herrors.HError) {
	return t.templateService.FindByID(ctx, templateID)
}

// Exists 检查模板是否存在
func (t TemplateService) Exists(ctx context.Context, templateID string) (bool, *herrors.HError) {
	template, err := t.templateService.FindByID(ctx, templateID)
	if herrors.HaveError(err) {
		if herrors.Is(err, template_err.TemplateCodeExists) {
			return false, nil
		}
		return false, err
	}
	return template != nil, nil
}
