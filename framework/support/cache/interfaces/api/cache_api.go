package cache_api

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/model"
)

// InternalCacheService 内部缓存服务接口
type InternalCacheService interface {
	// SetCache 设置缓存
	// cache: 缓存对象
	// 返回值：错误
	SetCache(ctx context.Context, cache *model.Cache) error

	// Set 设置缓存
	// tenantID: 租户ID
	// key: 缓存键
	// value: 缓存值
	// expireSeconds: 过期时间(秒)，0表示永不过期
	Set(ctx context.Context, tenantID string, key string, value interface{}, expireSeconds int64) error

	// Get 获取缓存
	// tenantID: 租户ID
	// key: 缓存键
	Get(ctx context.Context, tenantID string, key string) (*model.Cache, error)

	// Delete 删除缓存
	// tenantID: 租户ID
	// key: 缓存键
	Delete(ctx context.Context, tenantID string, key string) error

	// SetWithGroup 设置带分组的缓存
	// tenantID: 租户ID
	// groupID: 分组ID（必须是预定义的分组之一）
	// key: 缓存键
	// value: 缓存值
	// expireSeconds: 过期时间(秒)，0表示永不过期
	SetWithGroup(ctx context.Context, tenantID string, groupID, key string, value interface{}, expireSeconds int64) error

	// GetWithGroup 获取带分组的缓存
	// tenantID: 租户ID
	// groupID: 分组ID（必须是预定义的分组之一）
	// key: 缓存键
	GetWithGroup(ctx context.Context, tenantID string, groupID, key string) (*model.Cache, error)

	// DeleteWithGroup 删除带分组的缓存
	// tenantID: 租户ID
	// groupID: 分组ID（必须是预定义的分组之一）
	// key: 缓存键
	DeleteWithGroup(ctx context.Context, tenantID string, groupID, key string) error

	// DeleteGroup 删除分组下的所有缓存
	// tenantID: 租户ID
	// groupID: 分组ID（必须是预定义的分组之一）
	DeleteGroup(ctx context.Context, tenantID string, groupID string) error
}
