package service

import (
	"context"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/snowflake_id"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	ruleengineerr "github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/err"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/repository"
)

// RuleTemplateService 规则模板领域服务
type RuleTemplateService struct {
	templateRepo repository.ITemplateRepository
	categoryRepo repository.ICategoryRepository
	ig           snowflake_id.IIdGenerate
	ruleRepo     repository.IRuleRepository
}

// NewRuleTemplateService 创建规则模板服务
func NewRuleTemplateService(
	templateRepo repository.ITemplateRepository,
	categoryRepo repository.ICategoryRepository,
	ruleRepo repository.IRuleRepository,
	ig snowflake_id.IIdGenerate,
) *RuleTemplateService {
	return &RuleTemplateService{
		templateRepo: templateRepo,
		categoryRepo: categoryRepo,
		ruleRepo:     ruleRepo,
		ig:           ig,
	}
}

// CreateTemplate 创建模板
func (s *RuleTemplateService) CreateTemplate(ctx context.Context, template *model.RuleTemplate) *herrors.HError {
	// 验证模板数据
	if err := template.Validate(); err != nil {
		return ruleengineerr.RuleTemplateValidationFailed(err)
	}

	//// 检查分类是否存在
	//category, err := s.categoryRepo.FindByID(ctx, template.CategoryID)
	//if err != nil {
	//	return ruleengineerr.RuleCategoryGetFailed(err)
	//}
	//if !category.IsEnabled() {
	//	return ruleengineerr.RuleCategoryDisabled
	//}

	// 检查编码是否已存在
	exists, err := s.templateRepo.ExistsByCode(ctx, template.Code)
	if err != nil {
		return ruleengineerr.RuleTemplateGetFailed(err)
	}
	if exists {
		return ruleengineerr.RuleTemplateCodeExists
	}

	// 生成ID
	template.ID = s.ig.GenStringId()
	template.Completion()
	// 创建模板
	if err := s.templateRepo.Create(ctx, template); err != nil {
		return ruleengineerr.RuleTemplateCreateFailed(err)
	}

	return nil
}

// UpdateTemplate 更新模板
func (s *RuleTemplateService) UpdateTemplate(ctx context.Context, template *model.RuleTemplate) *herrors.HError {
	// 验证模板数据
	if err := template.Validate(); err != nil {
		return ruleengineerr.RuleTemplateValidationFailed(err)
	}

	// 检查模板是否存在
	existingTemplate, err := s.templateRepo.FindByID(ctx, template.ID)
	if err != nil {
		return ruleengineerr.RuleTemplateGetFailed(err)
	}
	if !existingTemplate.IsEnabled() {
		return ruleengineerr.RuleTemplateDisabled
	}

	//// 检查分类是否存在
	//category, err := s.categoryRepo.FindByID(ctx, template.CategoryID)
	//if err != nil {
	//	return ruleengineerr.RuleCategoryGetFailed(err)
	//}
	//if !category.IsEnabled() {
	//	return ruleengineerr.RuleCategoryDisabled
	//}

	// 检查编码是否重复（排除自己）
	if template.Code != existingTemplate.Code {
		exists, err := s.templateRepo.ExistsByCode(ctx, template.Code)
		if err != nil {
			return ruleengineerr.RuleTemplateGetFailed(err)
		}
		if exists {
			return ruleengineerr.RuleTemplateCodeExists
		}
	}
	template.Completion()
	// 更新模板
	if err := s.templateRepo.Update(ctx, template); err != nil {
		return ruleengineerr.RuleTemplateUpdateFailed(err)
	}

	return nil
}

// DeleteTemplate 删除模板
func (s *RuleTemplateService) DeleteTemplate(ctx context.Context, templateID string) *herrors.HError {
	// 检查模板是否存在
	_, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return ruleengineerr.RuleTemplateGetFailed(err)
	}

	// 检查是否有规则在使用此模板
	rules, err := s.ruleRepo.FindByTemplateID(ctx, templateID)
	if err != nil {
		return ruleengineerr.RuleGetFailed(err)
	}
	if len(rules) > 0 {
		return ruleengineerr.RuleCategoryHasTemplates
	}

	// 删除模板
	if err := s.templateRepo.Delete(ctx, templateID); err != nil {
		return ruleengineerr.RuleTemplateDeleteFailed(err)
	}

	return nil
}

// GetTemplate 获取模板
func (s *RuleTemplateService) GetTemplate(ctx context.Context, templateID string) (*model.RuleTemplate, *herrors.HError) {
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return nil, ruleengineerr.RuleTemplateGetFailed(err)
	}

	return template, nil
}

// GetTemplateByCode 根据编码获取模板
func (s *RuleTemplateService) GetTemplateByCode(ctx context.Context, code string) (*model.RuleTemplate, *herrors.HError) {
	template, err := s.templateRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, ruleengineerr.RuleTemplateGetFailed(err)
	}

	return template, nil
}

// GetTemplatesByCategory 根据分类获取模板列表
func (s *RuleTemplateService) GetTemplatesByCategory(ctx context.Context, categoryID string) ([]*model.RuleTemplate, *herrors.HError) {
	// 检查分类是否存在
	_, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, ruleengineerr.RuleCategoryGetFailed(err)
	}

	templates, err := s.templateRepo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, ruleengineerr.RuleTemplateGetFailed(err)
	}

	return templates, nil
}

// GetTemplatesByType 根据类型获取模板列表
func (s *RuleTemplateService) GetTemplatesByType(ctx context.Context, templateType string) ([]*model.RuleTemplate, *herrors.HError) {
	templates, err := s.templateRepo.FindByType(ctx, templateType)
	if err != nil {
		return nil, ruleengineerr.RuleTemplateGetFailed(err)
	}

	return templates, nil
}

// EnableTemplate 启用模板
func (s *RuleTemplateService) EnableTemplate(ctx context.Context, templateID string) *herrors.HError {
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return ruleengineerr.RuleTemplateGetFailed(err)
	}

	template.Enable()
	if err := s.templateRepo.Update(ctx, template); err != nil {
		return ruleengineerr.RuleTemplateUpdateFailed(err)
	}

	return nil
}

// DisableTemplate 禁用模板
func (s *RuleTemplateService) DisableTemplate(ctx context.Context, templateID string) *herrors.HError {
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return ruleengineerr.RuleTemplateGetFailed(err)
	}

	template.Disable()
	if err := s.templateRepo.Update(ctx, template); err != nil {
		return ruleengineerr.RuleTemplateUpdateFailed(err)
	}

	return nil
}

// ValidateTemplate 验证模板
func (s *RuleTemplateService) ValidateTemplate(ctx context.Context, template *model.RuleTemplate) *herrors.HError {
	if err := template.Validate(); err != nil {
		return ruleengineerr.RuleTemplateValidationFailed(err)
	}

	return nil
}

// ApplyTemplateParameters 应用模板参数
func (s *RuleTemplateService) ApplyTemplateParameters(ctx context.Context, templateID string, params map[string]interface{}) (map[string]interface{}, *herrors.HError) {
	template, err := s.templateRepo.FindByID(ctx, templateID)
	if err != nil {
		return nil, ruleengineerr.RuleTemplateGetFailed(err)
	}

	content, err := template.ApplyParameters(params)
	if err != nil {
		return nil, ruleengineerr.RuleTemplateContentInvalid
	}

	return content, nil
}
