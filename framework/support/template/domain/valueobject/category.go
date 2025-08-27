package valueobject

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	template_err "github.com/flare-admin/flare-server-go/framework/support/template/domain/err"
)

// Category 分类值对象
type Category struct {
	ID          string `json:"id" is_query:"true"`          // 分类ID
	Name        string `json:"name" is_query:"true"`        // 分类名称
	Code        string `json:"code" is_query:"true"`        // 分类编码
	Description string `json:"description" is_query:"true"` // 分类描述
	Sort        int    `json:"sort" is_query:"true"`        // 排序
	Status      int    `json:"status" is_query:"true"`      // 状态
	CreatedAt   int64  `json:"created_at" is_query:"true"`  // 创建时间
	UpdatedAt   int64  `json:"updated_at" is_query:"true"`  // 更新时间
}

// CreateCategoryCommand 创建分类命令
type CreateCategoryCommand struct {
	Name        string `json:"name" binding:"required" comment:"分类名称"`
	Code        string `json:"code" binding:"required" comment:"分类编码"`
	Description string `json:"description" comment:"分类描述"`
	Sort        int    `json:"sort" comment:"排序"`
}

// UpdateCategoryCommand 更新分类命令
type UpdateCategoryCommand struct {
	ID          string `json:"id" binding:"required" comment:"分类ID"`
	Name        string `json:"name" binding:"required" comment:"分类名称"`
	Code        string `json:"code" binding:"required" comment:"分类编码"`
	Description string `json:"description" comment:"分类描述"`
	Sort        int    `json:"sort" comment:"排序"`
}

// UpdateCategoryStatusCommand 更新分类状态命令
type UpdateCategoryStatusCommand struct {
	ID     string `json:"id" binding:"required" comment:"分类ID"`
	Status int    `json:"status" binding:"required" comment:"分类状态"`
}

// DeleteCategoryCommand 删除分类命令
type DeleteCategoryCommand struct {
	ID string `json:"id" binding:"required" comment:"分类ID"`
}

// Validate 验证创建分类命令
func (c *CreateCategoryCommand) Validate() *herrors.HError {
	if c.Name == "" {
		return template_err.CategoryCreateFailed(fmt.Errorf("分类名称不能为空"))
	}
	if c.Code == "" {
		return template_err.CategoryCreateFailed(fmt.Errorf("分类编码不能为空"))
	}
	return nil
}

// Validate 验证更新分类命令
func (c *UpdateCategoryCommand) Validate() *herrors.HError {
	if c.ID == "" {
		return template_err.CategoryUpdateFailed(fmt.Errorf("分类ID不能为空"))
	}
	if c.Name == "" {
		return template_err.CategoryUpdateFailed(fmt.Errorf("分类名称不能为空"))
	}
	if c.Code == "" {
		return template_err.CategoryUpdateFailed(fmt.Errorf("分类编码不能为空"))
	}
	return nil
}

// Validate 验证更新分类状态命令
func (c *UpdateCategoryStatusCommand) Validate() *herrors.HError {
	if c.ID == "" {
		return template_err.CategoryUpdateFailed(fmt.Errorf("分类ID不能为空"))
	}
	if c.Status != 1 && c.Status != 2 {
		return template_err.CategoryUpdateFailed(fmt.Errorf("无效的分类状态"))
	}
	return nil
}

// Validate 验证删除分类命令
func (c *DeleteCategoryCommand) Validate() *herrors.HError {
	if c.ID == "" {
		return template_err.CategoryDeleteFailed(fmt.Errorf("分类ID不能为空"))
	}
	return nil
}
