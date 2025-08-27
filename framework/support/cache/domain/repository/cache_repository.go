package repository

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/model"
)

// CacheRepository 缓存仓储接口
type CacheRepository interface {
	// Set 设置缓存
	Set(ctx context.Context, cache *model.Cache) error

	// Get 获取缓存
	Get(ctx context.Context, cache *model.Cache) error

	// Delete 删除缓存
	Delete(ctx context.Context, cache *model.Cache) error

	// DeleteByGroup 删除分组下的所有缓存
	DeleteByGroup(ctx context.Context, tenantID string, groupID string) error

	// GetGroupStats 获取分组统计信息
	GetGroupStats(ctx context.Context, tenantID string, groupID string) (*model.CacheGroupInfo, error)

	// ListGroupsWithStats 获取所有分组及其统计信息
	ListGroupsWithStats(ctx context.Context, tenantID string) ([]*model.CacheGroupInfo, error)

	// ListGroupKeys 获取分组下的所有键
	ListGroupKeys(ctx context.Context, tenantID string, groupID string) ([]string, error)

	// ListByGroup 获取分组下的所有缓存
	ListByGroup(ctx context.Context, tenantID string, groupID string) ([]*model.Cache, error)

	// Clear 清空所有缓存
	Clear(ctx context.Context, tenantID string) error
}
