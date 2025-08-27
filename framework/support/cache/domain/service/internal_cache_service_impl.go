package service

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/errors"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/repository"
	"time"
)

// InternalCacheServiceImpl 内部缓存服务实现
type InternalCacheServiceImpl struct {
	repo repository.CacheRepository
}

// NewInternalCacheService 创建内部缓存服务
func NewInternalCacheService(repo repository.CacheRepository) *InternalCacheServiceImpl {
	return &InternalCacheServiceImpl{
		repo: repo,
	}
}

// SetCache 设置缓存
// cache: 缓存对象
// 返回值：错误
func (s *InternalCacheServiceImpl) SetCache(ctx context.Context, cache *model.Cache) error {
	if cache == nil {
		return cache_errors.ErrInvalidKey
	}
	if cache.Key == "" {
		return cache_errors.ErrInvalidKey
	}
	if cache.GroupID == "" {
		cache.GroupID = model.CacheGroupDefault
	}
	if !model.IsValidCacheGroup(cache.GroupID) {
		return cache_errors.ErrInvalidKey
	}
	return s.repo.Set(ctx, cache)
}

// Set 设置缓存
func (s *InternalCacheServiceImpl) Set(ctx context.Context, tenantID string, key string, value interface{}, expireSeconds int64) error {
	return s.SetWithGroup(ctx, tenantID, model.CacheGroupDefault, key, value, expireSeconds)
}

// Get 获取缓存
func (s *InternalCacheServiceImpl) Get(ctx context.Context, tenantID string, key string) (*model.Cache, error) {
	return s.GetWithGroup(ctx, tenantID, model.CacheGroupDefault, key)
}

// Delete 删除缓存
func (s *InternalCacheServiceImpl) Delete(ctx context.Context, tenantID string, key string) error {
	return s.DeleteWithGroup(ctx, tenantID, model.CacheGroupDefault, key)
}

// SetWithGroup 设置带分组的缓存
func (s *InternalCacheServiceImpl) SetWithGroup(ctx context.Context, tenantID string, groupID, key string, value interface{}, expireSeconds int64) error {
	if key == "" {
		return cache_errors.ErrInvalidKey
	}

	if !model.IsValidCacheGroup(groupID) {
		return cache_errors.ErrInvalidKey
	}

	cache := &model.Cache{
		Key:       key,
		TenantID:  tenantID,
		GroupID:   groupID,
		ExpireAt:  utils.GetTimeNow().Add(time.Duration(expireSeconds) * time.Second).Unix(),
		CreatedAt: utils.GetDateUnix(),
		UpdatedAt: utils.GetDateUnix(),
	}
	if err := cache.SetValue(value); err != nil {
		return cache_errors.ErrInvalidValue
	}
	return s.repo.Set(ctx, cache)
}

// GetWithGroup 获取带分组的缓存
func (s *InternalCacheServiceImpl) GetWithGroup(ctx context.Context, tenantID string, groupID, key string) (*model.Cache, error) {
	if key == "" {
		return nil, cache_errors.ErrInvalidKey
	}

	if !model.IsValidCacheGroup(groupID) {
		return nil, cache_errors.ErrInvalidKey
	}

	// 构建带租户ID的缓存键
	cache := &model.Cache{
		Key:      key,
		GroupID:  groupID,
		TenantID: tenantID,
	}

	err := s.repo.Get(ctx, cache)
	if err != nil {
		return nil, err
	}
	if cache.IsExpired() {
		if err := s.repo.Delete(ctx, cache); err != nil {
			return nil, err
		}
		return nil, cache_errors.ErrKeyExpired
	}
	return cache, nil
}

// DeleteWithGroup 删除带分组的缓存
func (s *InternalCacheServiceImpl) DeleteWithGroup(ctx context.Context, tenantID string, groupID, key string) error {
	if key == "" {
		return cache_errors.ErrInvalidKey
	}

	if !model.IsValidCacheGroup(groupID) {
		return cache_errors.ErrInvalidKey
	}

	// 构建带租户ID的缓存键
	cache := &model.Cache{
		Key:      key,
		GroupID:  groupID,
		TenantID: tenantID,
	}
	return s.repo.Delete(ctx, cache)
}

// DeleteGroup 删除分组下的所有缓存
func (s *InternalCacheServiceImpl) DeleteGroup(ctx context.Context, tenantID string, groupID string) error {
	if !model.IsValidCacheGroup(groupID) {
		return cache_errors.ErrInvalidKey
	}
	return s.repo.DeleteByGroup(ctx, tenantID, groupID)
}
