package service

import (
	"context"
	"fmt"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/snowflake_id"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	ruleengineerr "github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/err"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/repository"
)

// RuleCategoryService 规则分类领域服务
type RuleCategoryService struct {
	categoryRepo repository.ICategoryRepository
	templateRepo repository.ITemplateRepository
	ruleRepo     repository.IRuleRepository
	ig           snowflake_id.IIdGenerate
}

// NewRuleCategoryService 创建规则分类服务
func NewRuleCategoryService(
	categoryRepo repository.ICategoryRepository,
	templateRepo repository.ITemplateRepository,
	ruleRepo repository.IRuleRepository,
	ig snowflake_id.IIdGenerate,
) *RuleCategoryService {
	return &RuleCategoryService{
		categoryRepo: categoryRepo,
		templateRepo: templateRepo,
		ruleRepo:     ruleRepo,
		ig:           ig,
	}
}

// CreateCategory 创建分类
func (s *RuleCategoryService) CreateCategory(ctx context.Context, category *model.RuleCategory) *herrors.HError {
	// 验证分类数据
	if err := category.Validate(); err != nil {
		return ruleengineerr.RuleCategoryValidationFailed(err)
	}

	// 检查编码是否已存在
	exists, err := s.categoryRepo.ExistsByCode(ctx, category.Code)
	if err != nil {
		return ruleengineerr.RuleCategoryGetFailed(err)
	}
	if exists {
		return ruleengineerr.RuleCategoryCodeExists
	}

	// 如果有父分类，检查父分类是否存在
	if category.ParentID != "" {
		parentCategory, err := s.categoryRepo.FindByID(ctx, category.ParentID)
		if err != nil {
			return ruleengineerr.RuleCategoryGetFailed(err)
		}
		if !parentCategory.IsEnabled() {
			return ruleengineerr.RuleCategoryDisabled
		}
		// 设置路径和层级
		category.SetParent(parentCategory.ID, parentCategory.Path)
	} else {
		// 根分类
		category.SetParent("", "")
	}

	// 生成ID
	category.ID = s.ig.GenStringId()

	// 创建分类
	if err := s.categoryRepo.Create(ctx, category); err != nil {
		return ruleengineerr.RuleCategoryCreateFailed(err)
	}

	return nil
}

// UpdateCategory 更新分类
func (s *RuleCategoryService) UpdateCategory(ctx context.Context, category *model.RuleCategory) *herrors.HError {
	// 验证分类数据
	if err := category.Validate(); err != nil {
		return ruleengineerr.RuleCategoryValidationFailed(err)
	}

	// 检查分类是否存在
	existingCategory, err := s.categoryRepo.FindByID(ctx, category.ID)
	if err != nil {
		return ruleengineerr.RuleCategoryGetFailed(err)
	}
	if !existingCategory.IsEnabled() {
		return ruleengineerr.RuleCategoryDisabled
	}

	// 检查编码是否重复（排除自己）
	if category.Code != existingCategory.Code {
		exists, err := s.categoryRepo.ExistsByCode(ctx, category.Code)
		if err != nil {
			return ruleengineerr.RuleCategoryGetFailed(err)
		}
		if exists {
			return ruleengineerr.RuleCategoryCodeExists
		}
	}

	// 如果有父分类，检查父分类是否存在且不能是自己或自己的子分类
	if category.ParentID != "" {
		if category.ParentID == category.ID {
			return ruleengineerr.RuleCategoryValidationFailed(fmt.Errorf("parent category cannot be itself"))
		}

		parentCategory, err := s.categoryRepo.FindByID(ctx, category.ParentID)
		if err != nil {
			return ruleengineerr.RuleCategoryGetFailed(err)
		}
		if !parentCategory.IsEnabled() {
			return ruleengineerr.RuleCategoryDisabled
		}

		// 检查是否将分类设置为自己的子分类
		if parentCategory.IsDescendantOf(existingCategory.Path) {
			return ruleengineerr.RuleCategoryValidationFailed(fmt.Errorf("cannot set category as child of its descendant"))
		}

		// 更新路径和层级
		category.SetParent(parentCategory.ID, parentCategory.Path)
	} else {
		// 根分类
		category.SetParent("", "")
	}

	// 更新分类
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return ruleengineerr.RuleCategoryUpdateFailed(err)
	}

	return nil
}

// DeleteCategory 删除分类
func (s *RuleCategoryService) DeleteCategory(ctx context.Context, categoryID string) *herrors.HError {
	// 检查分类是否存在
	_, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return ruleengineerr.RuleCategoryGetFailed(err)
	}

	// 检查是否有子分类
	children, err := s.categoryRepo.FindByParentID(ctx, categoryID)
	if err != nil {
		return ruleengineerr.RuleCategoryGetFailed(err)
	}
	if len(children) > 0 {
		return ruleengineerr.RuleCategoryHasChildren
	}

	// 检查是否有模板
	templates, err := s.templateRepo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return ruleengineerr.RuleTemplateGetFailed(err)
	}
	if len(templates) > 0 {
		return ruleengineerr.RuleCategoryHasTemplates
	}

	// 检查是否有规则
	rules, err := s.ruleRepo.FindByCategoryID(ctx, categoryID)
	if err != nil {
		return ruleengineerr.RuleGetFailed(err)
	}
	if len(rules) > 0 {
		return ruleengineerr.RuleCategoryHasRules
	}

	// 删除分类
	if err := s.categoryRepo.Delete(ctx, categoryID); err != nil {
		return ruleengineerr.RuleCategoryDeleteFailed(err)
	}

	return nil
}

// GetCategory 获取分类
func (s *RuleCategoryService) GetCategory(ctx context.Context, categoryID string) (*model.RuleCategory, *herrors.HError) {
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, ruleengineerr.RuleCategoryGetFailed(err)
	}

	return category, nil
}

// GetCategoryByCode 根据编码获取分类
func (s *RuleCategoryService) GetCategoryByCode(ctx context.Context, code string) (*model.RuleCategory, *herrors.HError) {
	category, err := s.categoryRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, ruleengineerr.RuleCategoryGetFailed(err)
	}

	return category, nil
}

// GetCategoriesByParent 根据父分类获取子分类列表
func (s *RuleCategoryService) GetCategoriesByParent(ctx context.Context, parentID string) ([]*model.RuleCategory, *herrors.HError) {
	if parentID != "" {
		// 检查父分类是否存在
		_, err := s.categoryRepo.FindByID(ctx, parentID)
		if err != nil {
			return nil, ruleengineerr.RuleCategoryGetFailed(err)
		}
	}

	categories, err := s.categoryRepo.FindByParentID(ctx, parentID)
	if err != nil {
		return nil, ruleengineerr.RuleCategoryGetFailed(err)
	}

	return categories, nil
}

// GetCategoriesByBusinessType 根据业务类型获取分类列表
func (s *RuleCategoryService) GetCategoriesByBusinessType(ctx context.Context, businessType string) ([]*model.RuleCategory, *herrors.HError) {
	categories, err := s.categoryRepo.FindByBusinessType(ctx, businessType)
	if err != nil {
		return nil, ruleengineerr.RuleCategoryGetFailed(err)
	}

	return categories, nil
}

// GetCategoriesByType 根据分类类型获取分类列表
func (s *RuleCategoryService) GetCategoriesByType(ctx context.Context, categoryType string) ([]*model.RuleCategory, *herrors.HError) {
	categories, err := s.categoryRepo.FindByType(ctx, categoryType)
	if err != nil {
		return nil, ruleengineerr.RuleCategoryGetFailed(err)
	}

	return categories, nil
}

// GetRootCategories 获取根分类列表
func (s *RuleCategoryService) GetRootCategories(ctx context.Context) ([]*model.RuleCategory, *herrors.HError) {
	categories, err := s.categoryRepo.FindRootCategories(ctx)
	if err != nil {
		return nil, ruleengineerr.RuleCategoryGetFailed(err)
	}

	return categories, nil
}

// GetCategoryTree 获取分类树
func (s *RuleCategoryService) GetCategoryTree(ctx context.Context, categoryID string) ([]*model.RuleCategory, *herrors.HError) {
	if categoryID == "" {
		// 获取所有根分类
		return s.GetRootCategories(ctx)
	}

	// 获取指定分类及其所有后代
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return nil, ruleengineerr.RuleCategoryGetFailed(err)
	}

	descendants, err := s.categoryRepo.FindDescendants(ctx, category.Path)
	if err != nil {
		return nil, ruleengineerr.RuleCategoryGetFailed(err)
	}

	// 构建树结构
	result := []*model.RuleCategory{category}
	result = append(result, descendants...)

	return result, nil
}

// EnableCategory 启用分类
func (s *RuleCategoryService) EnableCategory(ctx context.Context, categoryID string) *herrors.HError {
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return ruleengineerr.RuleCategoryGetFailed(err)
	}

	category.Enable()
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return ruleengineerr.RuleCategoryUpdateFailed(err)
	}

	return nil
}

// DisableCategory 禁用分类
func (s *RuleCategoryService) DisableCategory(ctx context.Context, categoryID string) *herrors.HError {
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		return ruleengineerr.RuleCategoryGetFailed(err)
	}

	category.Disable()
	if err := s.categoryRepo.Update(ctx, category); err != nil {
		return ruleengineerr.RuleCategoryUpdateFailed(err)
	}

	return nil
}

// ValidateCategory 验证分类
func (s *RuleCategoryService) ValidateCategory(ctx context.Context, category *model.RuleCategory) *herrors.HError {
	if err := category.Validate(); err != nil {
		return ruleengineerr.RuleCategoryValidationFailed(err)
	}

	return nil
}
