package service

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/model"
)

// CacheService 缓存服务接口
type CacheService interface {
	// GetWithGroup 获取带分组的缓存
	// tenantID: 租户ID
	// groupID: 分组ID（必须是预定义的分组之一）
	// key: 缓存键
	// target: 目标对象指针，用于存储获取到的值
	GetWithGroup(ctx context.Context, tenantID string, groupID, key string, target interface{}) error

	// DeleteWithGroup 删除带分组的缓存
	// tenantID: 租户ID
	// groupID: 分组ID（必须是预定义的分组之一）
	// key: 缓存键
	DeleteWithGroup(ctx context.Context, tenantID string, groupID, key string) error

	// DeleteGroup 删除分组下的所有缓存
	// tenantID: 租户ID
	// groupID: 分组ID（必须是预定义的分组之一）
	DeleteGroup(ctx context.Context, tenantID string, groupID string) error

	// GetGroupStats 获取分组统计信息
	// tenantID: 租户ID
	// groupID: 分组ID（必须是预定义的分组之一）
	GetGroupStats(ctx context.Context, tenantID string, groupID string) (*model.CacheGroupInfo, error)

	// ListGroupsWithStats 获取所有分组及其统计信息
	// tenantID: 租户ID
	ListGroupsWithStats(ctx context.Context, tenantID string) ([]*model.CacheGroupInfo, error)

	// ListGroupKeys 获取分组下的所有键
	// tenantID: 租户ID
	// groupID: 分组ID（必须是预定义的分组之一）
	ListGroupKeys(ctx context.Context, tenantID string, groupID string) ([]string, error)
}
