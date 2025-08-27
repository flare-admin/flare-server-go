package service

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	template_err "github.com/flare-admin/flare-server-go/framework/support/template/domain/err"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/repository"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/valueobject"
)

// CategoryService 分类领域服务
type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryService 创建分类服务实例
func NewCategoryService(
	categoryRepo repository.CategoryRepository,
) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

// CreateCategory 创建分类
func (s *CategoryService) CreateCategory(ctx context.Context, cmd *valueobject.CreateCategoryCommand) *herrors.HError {
	// 1. 验证命令
	if err := cmd.Validate(); err != nil {
		return err
	}

	// 2. 验证编码是否已存在
	existCategory, err := s.categoryRepo.FindByCode(ctx, cmd.Code)
	if !database.IfErrorNotFound(err) {
		hlog.CtxErrorf(ctx, "查询分类失败:%v", err)
		return template_err.CategoryGetFailed(err)
	}
	if existCategory != nil && existCategory.ID != "" {
		return template_err.CategoryCodeExist
	}

	// 3. 创建分类
	category := model.NewCategory(cmd.Name, cmd.Code, cmd.Description)
	category.Sort = cmd.Sort

	if err := s.categoryRepo.Create(ctx, category); err != nil {
		hlog.CtxErrorf(ctx, "创建分类失败:%v", err)
		return template_err.CategoryCreateFailed(err)
	}

	hlog.CtxInfof(ctx, "创建分类成功,分类ID:%s,名称:%s", category.ID, category.Name)
	return nil
}

// UpdateCategory 更新分类
func (s *CategoryService) UpdateCategory(ctx context.Context, cmd *valueobject.UpdateCategoryCommand) *herrors.HError {
	// 1. 验证命令
	if err := cmd.Validate(); err != nil {
		return err
	}

	// 2. 获取原分类信息
	oldCategory, err := s.categoryRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取分类失败:%v", err)
		return template_err.CategoryGetFailed(err)
	}
	if oldCategory == nil {
		return template_err.CategoryNotExist
	}

	// 3. 验证分类状态
	if oldCategory.Status != 1 {
		hlog.CtxErrorf(ctx, "分类已禁用,状态:%d", oldCategory.Status)
		return template_err.CategoryDisabled
	}

	// 4. 验证编码是否已存在
	if cmd.Code != oldCategory.Code {
		existCategory, err := s.categoryRepo.FindByCode(ctx, cmd.Code)
		if err != nil {
			hlog.CtxErrorf(ctx, "查询分类失败:%v", err)
			return template_err.CategoryGetFailed(err)
		}
		if existCategory != nil {
			return template_err.CategoryCodeExist
		}
	}

	// 5. 更新分类
	oldCategory.Name = cmd.Name
	oldCategory.Code = cmd.Code
	oldCategory.Description = cmd.Description
	oldCategory.Sort = cmd.Sort
	oldCategory.UpdatedAt = utils.GetDateUnix()

	if err := s.categoryRepo.Update(ctx, oldCategory); err != nil {
		hlog.CtxErrorf(ctx, "更新分类失败:%v", err)
		return template_err.CategoryUpdateFailed(err)
	}

	hlog.CtxInfof(ctx, "更新分类成功,分类ID:%s,名称:%s", oldCategory.ID, oldCategory.Name)
	return nil
}

// DeleteCategory 删除分类
func (s *CategoryService) DeleteCategory(ctx context.Context, cmd *valueobject.DeleteCategoryCommand) *herrors.HError {
	// 1. 验证命令
	if err := cmd.Validate(); err != nil {
		return err
	}

	// 2. 获取分类信息
	category, err := s.categoryRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取分类失败:%v", err)
		return template_err.CategoryGetFailed(err)
	}
	if category == nil {
		return template_err.CategoryNotExist
	}

	// 3. 验证分类状态
	if category.Status != 1 {
		hlog.CtxErrorf(ctx, "分类已禁用,状态:%d", category.Status)
		return template_err.CategoryDisabled
	}

	// 4. 删除分类
	if err := s.categoryRepo.Delete(ctx, cmd.ID); err != nil {
		hlog.CtxErrorf(ctx, "删除分类失败:%v", err)
		return template_err.CategoryDeleteFailed(err)
	}

	hlog.CtxInfof(ctx, "删除分类成功,分类ID:%s", cmd.ID)
	return nil
}

// EnableCategory 启用分类
func (s *CategoryService) EnableCategory(ctx context.Context, cmd *valueobject.UpdateCategoryStatusCommand) *herrors.HError {
	// 1. 验证命令
	if err := cmd.Validate(); err != nil {
		return err
	}

	// 2. 获取分类信息
	category, err := s.categoryRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取分类失败:%v", err)
		return template_err.CategoryGetFailed(err)
	}
	if category == nil {
		return template_err.CategoryNotExist
	}

	// 3. 更新状态
	category.Status = 1
	category.UpdatedAt = utils.GetDateUnix()

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		hlog.CtxErrorf(ctx, "更新分类状态失败:%v", err)
		return template_err.CategoryUpdateFailed(err)
	}

	hlog.CtxInfof(ctx, "启用分类成功,分类ID:%s", cmd.ID)
	return nil
}

// DisableCategory 禁用分类
func (s *CategoryService) DisableCategory(ctx context.Context, cmd *valueobject.UpdateCategoryStatusCommand) *herrors.HError {
	// 1. 验证命令
	if err := cmd.Validate(); err != nil {
		return err
	}

	// 2. 获取分类信息
	category, err := s.categoryRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取分类失败:%v", err)
		return template_err.CategoryGetFailed(err)
	}
	if category == nil {
		return template_err.CategoryNotExist
	}

	// 3. 更新状态
	category.Status = 2
	category.UpdatedAt = utils.GetDateUnix()

	if err := s.categoryRepo.Update(ctx, category); err != nil {
		hlog.CtxErrorf(ctx, "更新分类状态失败:%v", err)
		return template_err.CategoryUpdateFailed(err)
	}

	hlog.CtxInfof(ctx, "禁用分类成功,分类ID:%s", cmd.ID)
	return nil
}
