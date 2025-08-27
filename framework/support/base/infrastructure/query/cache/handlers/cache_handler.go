package handlers

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/casbin"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query/cache"
)

// CacheHandler 缓存处理器
type CacheHandler struct {
	userCache     *cache.UserQueryCache
	roleCache     *cache.RoleQueryCache
	deptCache     *cache.DepartmentQueryCache
	permCache     *cache.PermissionsQueryCache
	dataPermCache *cache.DataPermissionQueryCache
	tenantCache   *cache.TenantQueryCache
	en            *casbin.Enforcer
}

// NewCacheHandler 创建缓存处理器
func NewCacheHandler(
	userCache *cache.UserQueryCache,
	roleCache *cache.RoleQueryCache,
	deptCache *cache.DepartmentQueryCache,
	permCache *cache.PermissionsQueryCache,
	dataPermCache *cache.DataPermissionQueryCache,
	tenantCache *cache.TenantQueryCache,
	en *casbin.Enforcer,
) *CacheHandler {
	return &CacheHandler{
		userCache:     userCache,
		roleCache:     roleCache,
		deptCache:     deptCache,
		permCache:     permCache,
		dataPermCache: dataPermCache,
		tenantCache:   tenantCache,
		en:            en,
	}
}

// 用户缓存相关
func (h *CacheHandler) InvalidateUserAllCache(ctx context.Context, userID string) error {
	hlog.CtxDebugf(ctx, "清除用户所有相关缓存: userID=%s", userID)

	// 1. 清除用户基本信息
	if err := h.userCache.InvalidateUserCache(ctx, userID); err != nil {
		return fmt.Errorf("清除用户基本信息缓存失败: %w", err)
	}

	// 2. 清除用户权限相关
	if err := h.userCache.InvalidateUserPermissionCache(ctx, userID); err != nil {
		return fmt.Errorf("清除用户权限缓存失败: %w", err)
	}
	if err := h.userCache.InvalidateUserMenuCache(ctx, userID); err != nil {
		return fmt.Errorf("清除用户菜单缓存失败: %w", err)
	}

	// 3. 清除用户部门关系
	if err := h.userCache.InvalidateUserDepartmentCache(ctx, userID); err != nil {
		return fmt.Errorf("清除用户部门缓存失败: %w", err)
	}

	// 4. 清除用户角色关系
	roles, err := h.userCache.GetUserRoles(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户角色列表失败: %w", err)
	}
	for _, role := range roles {
		if err := h.roleCache.InvalidateRoleUserCache(ctx, role.ID); err != nil {
			return fmt.Errorf("清除角色用户列表缓存失败: %w", err)
		}
	}

	// 5. 清除用户所在部门的用户列表
	depts, err := h.userCache.GetUserDepartments(ctx, userID)
	if err != nil {
		return fmt.Errorf("获取用户部门列表失败: %w", err)
	}
	for _, dept := range depts {
		if err := h.deptCache.InvalidateDepartmentUserCache(ctx, dept.ID); err != nil {
			return fmt.Errorf("清除部门用户列表缓存失败: %w", err)
		}
	}

	return nil
}

// 角色缓存相关
func (h *CacheHandler) InvalidateRoleAllCache(ctx context.Context, roleID int64) error {
	hlog.CtxDebugf(ctx, "清除角色所有相关缓存: roleID=%d", roleID)

	// 1. 清除角色基本信息
	if err := h.roleCache.InvalidateRoleCache(ctx, roleID); err != nil {
		return fmt.Errorf("清除角色缓存失败: %w", err)
	}

	// 2. 清除角色权限关系
	if err := h.roleCache.InvalidateRolePermissionCache(ctx, roleID); err != nil {
		return fmt.Errorf("清除角色权限缓存失败: %w", err)
	}

	// 3. 清除角色下所有用户的权限缓存
	users, err := h.roleCache.GetRoleUsers(ctx, roleID)
	if err != nil {
		return fmt.Errorf("获取角色用户列表失败: %w", err)
	}
	for _, user := range users {
		if err := h.InvalidateUserAllCache(ctx, user.ID); err != nil {
			return fmt.Errorf("清除用户缓存失败: %w", err)
		}
	}

	// 4. 重新加载权限
	if err := h.en.PublishUpdate(ctx); err != nil {
		return fmt.Errorf("重新加载权限失败: %w", err)
	}

	return nil
}

// 部门缓存相关
func (h *CacheHandler) InvalidateDepartmentAllCache(ctx context.Context, deptID string) error {
	hlog.CtxDebugf(ctx, "清除部门所有相关缓存: deptID=%s", deptID)

	// 1. 清除部门基本信息
	if err := h.deptCache.InvalidateCache(ctx, deptID); err != nil {
		return fmt.Errorf("清除部门缓存失败: %w", err)
	}

	// 2. 清除部门树结构
	if err := h.deptCache.InvalidateDepartmentTreeCache(ctx); err != nil {
		return fmt.Errorf("清除部门树缓存失败: %w", err)
	}

	// 3. 清除子部门缓存
	if err := h.deptCache.InvalidateChildrenCache(ctx, deptID); err != nil {
		return fmt.Errorf("清除子部门缓存失败: %w", err)
	}

	// 4. 清除部门用户关系
	if err := h.deptCache.InvalidateDepartmentUserCache(ctx, deptID); err != nil {
		return fmt.Errorf("清除部门用户缓存失败: %w", err)
	}

	return nil
}

// 权限缓存相关
func (h *CacheHandler) InvalidatePermissionAllCache(ctx context.Context, permID int64) error {
	hlog.CtxDebugf(ctx, "清除权限所有相关缓存: permID=%d", permID)

	// 1. 清除权限基本信息
	if err := h.permCache.InvalidatePermissionCache(ctx, permID); err != nil {
		return fmt.Errorf("清除权限缓存失败: %w", err)
	}

	// 2. 清除权限树结构
	if err := h.permCache.InvalidatePermissionTreeCache(ctx); err != nil {
		return fmt.Errorf("清除权限树缓存失败: %w", err)
	}

	// 3. 清除拥有该权限的角色缓存
	roles, err := h.permCache.GetPermissionRoles(ctx, permID)
	if err != nil {
		return fmt.Errorf("获取权限角色列表失败: %w", err)
	}
	for _, role := range roles {
		if err := h.InvalidateRoleAllCache(ctx, role.ID); err != nil {
			return fmt.Errorf("清除角色缓存失败: %w", err)
		}
	}

	// 4. 重新加载权限
	if err := h.en.PublishUpdate(ctx); err != nil {
		return fmt.Errorf("重新加载权限失败: %w", err)
	}

	return nil
}

// 租户缓存相关
func (h *CacheHandler) InvalidateTenantAllCache(ctx context.Context, tenantID string) error {
	hlog.CtxDebugf(ctx, "清除租户所有相关缓存: tenantID=%s", tenantID)

	// 1. 清除租户基本信息
	if err := h.tenantCache.InvalidateTenantCache(ctx, tenantID); err != nil {
		return fmt.Errorf("清除租户缓存失败: %w", err)
	}

	// 2. 清除租户状态
	if err := h.tenantCache.InvalidateTenantStatusCache(ctx, tenantID); err != nil {
		return fmt.Errorf("清除租户状态缓存失败: %w", err)
	}

	// 3. 清除租户权限
	if err := h.tenantCache.InvalidateTenantPermissionCache(ctx, tenantID); err != nil {
		return fmt.Errorf("清除租户权限缓存失败: %w", err)
	}

	// 4. 清除租户下所有用户缓存
	if err := h.userCache.InvalidateTenantUserCache(ctx, tenantID); err != nil {
		return fmt.Errorf("清除租户用户缓存失败: %w", err)
	}

	// 5. 清除租户下所有角色缓存
	if err := h.roleCache.InvalidateTenantRoleCache(ctx, tenantID); err != nil {
		return fmt.Errorf("清除租户角色缓存失败: %w", err)
	}

	// 6. 清除租户下所有部门缓存
	if err := h.deptCache.InvalidateTenantDepartmentCache(ctx, tenantID); err != nil {
		return fmt.Errorf("清除租户部门缓存失败: %w", err)
	}

	// 7. 重新加载权限
	if err := h.en.PublishUpdate(ctx); err != nil {
		return fmt.Errorf("重新加载权限失败: %w", err)
	}

	return nil
}
