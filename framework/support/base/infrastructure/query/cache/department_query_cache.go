package cache

import (
	"context"
	"fmt"

	dCache "github.com/flare-admin/flare-server-go/framework/infrastructure/database/cache"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/dto"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query/cache/keys"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query/impl"
)

type DepartmentQueryCache struct {
	next      *impl.DepartmentQueryService
	decorator *dCache.CacheDecorator
}

func NewDepartmentQueryCache(
	next *impl.DepartmentQueryService,
	decorator *dCache.CacheDecorator,
) *DepartmentQueryCache {
	return &DepartmentQueryCache{
		next:      next,
		decorator: decorator,
	}
}

// GetDepartment 获取部门信息(带缓存)
func (c *DepartmentQueryCache) GetDepartment(ctx context.Context, id string) (*dto.DepartmentDto, error) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.DepartmentKey(tenantID, id)
	var dept *dto.DepartmentDto
	err := c.decorator.Cached(ctx, key, &dept, func() error {
		var err error
		dept, err = c.next.GetDepartment(ctx, id)
		return err
	})
	return dept, err
}

// GetDepartmentTree 获取部门树(带缓存)
func (c *DepartmentQueryCache) GetDepartmentTree(ctx context.Context, parentID string) ([]*dto.DepartmentTreeDto, error) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.DepartmentTreeKey(tenantID, parentID)
	var tree []*dto.DepartmentTreeDto
	err := c.decorator.Cached(ctx, key, &tree, func() error {
		var err error
		tree, err = c.next.GetDepartmentTree(ctx, parentID)
		return err
	})
	return tree, err
}

// GetUserDepartments 获取用户部门列表(带缓存)
func (c *DepartmentQueryCache) GetUserDepartments(ctx context.Context, userID string) ([]*dto.DepartmentDto, error) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.UserDepartmentsKey(tenantID, userID)
	var depts []*dto.DepartmentDto
	err := c.decorator.Cached(ctx, key, &depts, func() error {
		var err error
		depts, err = c.next.GetUserDepartments(ctx, userID)
		return err
	})
	return depts, err
}

// 列表查询不缓存,直接透传
func (c *DepartmentQueryCache) FindDepartments(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.DepartmentDto, error) {
	return c.next.FindDepartments(ctx, qb)
}

func (c *DepartmentQueryCache) CountDepartments(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return c.next.CountDepartments(ctx, qb)
}

func (c *DepartmentQueryCache) GetDepartmentUsers(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) ([]*dto.UserDto, error) {
	return c.next.GetDepartmentUsers(ctx, deptID, excludeAdminID, qb)
}

func (c *DepartmentQueryCache) CountDepartmentUsers(ctx context.Context, deptID string, excludeAdminID string, qb *db_query.QueryBuilder) (int64, error) {
	return c.next.CountDepartmentUsers(ctx, deptID, excludeAdminID, qb)
}

func (c *DepartmentQueryCache) GetUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.UserDto, error) {
	return c.next.GetUnassignedUsers(ctx, qb)
}

func (c *DepartmentQueryCache) CountUnassignedUsers(ctx context.Context, qb *db_query.QueryBuilder) (int64, error) {
	return c.next.CountUnassignedUsers(ctx, qb)
}

// WarmupCache 预热缓存
func (c *DepartmentQueryCache) WarmupCache(ctx context.Context, deptID string) error {

	// 1. 预热部门详情
	if _, err := c.GetDepartment(ctx, deptID); err != nil {
		return fmt.Errorf("预热部门详情缓存失败: %w", err)
	}

	// 2. 预热部门树
	if _, err := c.GetDepartmentTree(ctx, ""); err != nil {
		return fmt.Errorf("预热部门树缓存失败: %w", err)
	}

	return nil
}

// InvalidateDepartmentUserCache 清除部门用户列表缓存
func (c *DepartmentQueryCache) InvalidateDepartmentUserCache(ctx context.Context, deptID string) error {
	tenantID := actx.GetTenantId(ctx)
	return c.decorator.InvalidateCache(ctx, keys.DepartmentUsersKey(tenantID, deptID))
}

// InvalidateChildrenCache 清除部门子节点列表缓存
func (c *DepartmentQueryCache) InvalidateChildrenCache(ctx context.Context, parentID string) error {
	tenantID := actx.GetTenantId(ctx)
	return c.decorator.InvalidateCache(ctx, keys.DepartmentChildrenKey(tenantID, parentID))
}

// InvalidateDepartmentTreeCache 清除部门树缓存
func (c *DepartmentQueryCache) InvalidateDepartmentTreeCache(ctx context.Context) error {
	tenantID := actx.GetTenantId(ctx)
	return c.decorator.InvalidateCache(ctx, keys.DepartmentTreeKey(tenantID, ""))
}

// InvalidateCache 清除部门缓存
func (c *DepartmentQueryCache) InvalidateCache(ctx context.Context, deptID string) error {
	tenantID := actx.GetTenantId(ctx)
	return c.decorator.InvalidateCache(ctx, keys.DepartmentKeys(tenantID, deptID)...)
}

// InvalidateTenantDepartmentCache 清除租户下所有部门缓存
func (c *DepartmentQueryCache) InvalidateTenantDepartmentCache(ctx context.Context, tenantID string) error {
	// 使用租户前缀清除所有相关缓存
	return c.decorator.InvalidateTenantTypeCache(ctx, tenantID, keys.DeptPrefix)
}
