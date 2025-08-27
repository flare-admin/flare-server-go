package cache

import (
	"context"

	dCache "github.com/flare-admin/flare-server-go/framework/infrastructure/database/cache"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/dto"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query/cache/keys"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query/impl"
)

type DataPermissionQueryCache struct {
	next      *impl.DataPermissionQueryService
	decorator *dCache.CacheDecorator
}

func NewDataPermissionQueryCache(
	next *impl.DataPermissionQueryService,
	decorator *dCache.CacheDecorator,
) *DataPermissionQueryCache {
	return &DataPermissionQueryCache{
		next:      next,
		decorator: decorator,
	}
}

// GetByRoleID 获取角色的数据权限(带缓存)
func (c *DataPermissionQueryCache) GetByRoleID(ctx context.Context, roleID int64) (*dto.DataPermissionDto, herrors.Herr) {
	tenantID := actx.GetTenantId(ctx)
	key := keys.DataPermissionKey(tenantID, roleID)
	var permission *dto.DataPermissionDto
	err := c.decorator.Cached(ctx, key, &permission, func() error {
		var err error
		permission, err = c.next.GetByRoleID(ctx, roleID)
		return err
	})
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	return permission, nil
}

// InvalidateCache 使缓存失效
func (c *DataPermissionQueryCache) InvalidateCache(ctx context.Context, roleID int64) error {
	tenantID := actx.GetTenantId(ctx)
	return c.decorator.InvalidateCache(ctx, keys.DataPermissionKey(tenantID, roleID))
}
