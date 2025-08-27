package service

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/errors"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/repository"
)

// CacheServiceImpl 缓存服务实现
type CacheServiceImpl struct {
	repo repository.CacheRepository
}

// NewCacheService 创建缓存服务
func NewCacheService(repo repository.CacheRepository) CacheService {
	return &CacheServiceImpl{
		repo: repo,
	}
}

// GetWithGroup 获取带分组的缓存
func (s *CacheServiceImpl) GetWithGroup(ctx context.Context, tenantID string, groupID, key string, target interface{}) error {
	if key == "" {
		return cache_errors.ErrInvalidKey
	}

	if !model.IsValidCacheGroup(groupID) {
		return cache_errors.ErrInvalidKey
	}

	// 构建缓存对象
	cache := &model.Cache{
		Key:      key,
		GroupID:  groupID,
		TenantID: tenantID,
	}

	// 获取缓存
	if err := s.repo.Get(ctx, cache); err != nil {
		return err
	}

	// 检查是否过期
	if cache.IsExpired() {
		if err := s.repo.Delete(ctx, cache); err != nil {
			return err
		}
		return cache_errors.ErrKeyExpired
	}

	// 获取值
	if err := cache.GetValue(target); err != nil {
		return cache_errors.ErrInvalidValue
	}
	return nil
}

// DeleteWithGroup 删除带分组的缓存
func (s *CacheServiceImpl) DeleteWithGroup(ctx context.Context, tenantID string, groupID, key string) error {
	if key == "" {
		return cache_errors.ErrInvalidKey
	}

	if !model.IsValidCacheGroup(groupID) {
		return cache_errors.ErrInvalidKey
	}

	// 构建缓存对象
	cache := &model.Cache{
		Key:      key,
		TenantID: tenantID,
		GroupID:  groupID,
	}
	return s.repo.Delete(ctx, cache)
}

// DeleteGroup 删除分组下的所有缓存
func (s *CacheServiceImpl) DeleteGroup(ctx context.Context, tenantID string, groupID string) error {
	if !model.IsValidCacheGroup(groupID) {
		return cache_errors.ErrInvalidKey
	}
	return s.repo.DeleteByGroup(ctx, tenantID, groupID)
}

// GetGroupStats 获取分组统计信息
func (s *CacheServiceImpl) GetGroupStats(ctx context.Context, tenantID string, groupID string) (*model.CacheGroupInfo, error) {
	if !model.IsValidCacheGroup(groupID) {
		return nil, cache_errors.ErrInvalidKey
	}

	// 获取基础分组信息
	groups := model.GetDefaultCacheGroups()
	var groupInfo *model.CacheGroupInfo
	for _, g := range groups {
		if g.GroupID == groupID {
			groupInfo = g
			break
		}
	}
	if groupInfo == nil {
		return nil, cache_errors.ErrKeyNotFound
	}

	// 获取分组统计信息
	stats, err := s.repo.GetGroupStats(ctx, tenantID, groupID)
	if err != nil {
		return nil, err
	}

	// 更新统计信息
	groupInfo.KeyCount = stats.KeyCount
	groupInfo.MemoryUsage = stats.MemoryUsage

	return groupInfo, nil
}

// ListGroupsWithStats 获取所有分组及其统计信息
func (s *CacheServiceImpl) ListGroupsWithStats(ctx context.Context, tenantID string) ([]*model.CacheGroupInfo, error) {
	// 获取基础分组信息
	groups := model.GetDefaultCacheGroups()
	result := make([]*model.CacheGroupInfo, len(groups))

	// 获取每个分组的统计信息
	for i, group := range groups {
		stats, err := s.repo.GetGroupStats(ctx, tenantID, group.GroupID)
		if err != nil {
			// 如果获取统计信息失败，使用默认值
			result[i] = group
			continue
		}
		// 更新统计信息
		group.KeyCount = stats.KeyCount
		group.MemoryUsage = stats.MemoryUsage
		result[i] = group
	}

	return result, nil
}

// ListGroupKeys 获取分组下的所有键
func (s *CacheServiceImpl) ListGroupKeys(ctx context.Context, tenantID string, groupID string) ([]string, error) {
	// 获取分组下的所有缓存
	caches, err := s.repo.ListByGroup(ctx, tenantID, groupID)
	if err != nil {
		return nil, err
	}
	// 提取键
	var keys []string
	for _, cache := range caches {
		keys = append(keys, cache.Key)
	}
	return keys, nil
}
