package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/hredis"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"time"

	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/errors"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/repository"
	"github.com/redis/go-redis/v9"
)

// RedisCacheService Redis缓存服务实现
type RedisCacheService struct {
	redisClient *redis.Client
}

// NewRedisCacheService 创建Redis缓存服务实例
func NewRedisCacheService(redisClient *hredis.RedisClient) repository.CacheRepository {
	return &RedisCacheService{
		redisClient: redisClient.GetClient(),
	}
}

// Set 设置缓存
func (s *RedisCacheService) Set(ctx context.Context, cache *model.Cache) error {
	// 序列化值
	bytes, err := json.Marshal(cache.Value)
	if err != nil {
		return cache_errors.ErrInvalidValue
	}

	// 生成缓存键
	key := cache.BuildKey()
	if key == "" {
		return cache_errors.ErrInvalidKey
	}

	// 设置缓存
	if cache.ExpireAt > 0 {
		expireDuration := time.Duration(cache.ExpireAt-utils.GetDateUnix()) * time.Second
		if err := s.redisClient.Set(ctx, key, bytes, expireDuration).Err(); err != nil {
			return cache_errors.ErrConnectionFailed
		}
		return nil
	}

	if err := s.redisClient.Set(ctx, key, bytes, 0).Err(); err != nil {
		return cache_errors.ErrConnectionFailed
	}
	return nil
}

// Get 获取缓存
func (s *RedisCacheService) Get(ctx context.Context, cache *model.Cache) error {
	// 生成缓存键
	key := cache.BuildKey()
	if key == "" {
		return cache_errors.ErrInvalidKey
	}

	// 获取缓存
	bytes, err := s.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return cache_errors.ErrNil
		}
		return cache_errors.ErrConnectionFailed
	}

	// 反序列化值
	if err := json.Unmarshal(bytes, &cache.Value); err != nil {
		return cache_errors.ErrInvalidValue
	}

	// 获取过期时间
	ttl, err := s.redisClient.TTL(ctx, key).Result()
	if err != nil {
		return cache_errors.ErrConnectionFailed
	}

	if ttl > 0 {
		cache.ExpireAt = utils.GetTimeNow().Add(ttl).Unix()
	}

	return nil
}

// Delete 删除缓存
func (s *RedisCacheService) Delete(ctx context.Context, cache *model.Cache) error {
	// 生成缓存键
	key := cache.BuildKey()
	if key == "" {
		return cache_errors.ErrInvalidKey
	}

	if err := s.redisClient.Del(ctx, key).Err(); err != nil {
		return cache_errors.ErrConnectionFailed
	}
	return nil
}

// DeleteByGroup 删除分组下的所有缓存
func (s *RedisCacheService) DeleteByGroup(ctx context.Context, tenantID string, groupID string) error {
	// 验证分组ID
	if !model.IsValidCacheGroup(groupID) {
		return cache_errors.ErrInvalidKey
	}

	// 获取分组下的所有缓存
	caches, err := s.ListByGroup(ctx, tenantID, groupID)
	if err != nil {
		return err
	}

	// 删除所有缓存
	for _, cache := range caches {
		if err := s.Delete(ctx, cache); err != nil {
			return err
		}
	}

	return nil
}

// GetGroupStats 获取分组统计信息
func (s *RedisCacheService) GetGroupStats(ctx context.Context, tenantID string, groupID string) (*model.CacheGroupInfo, error) {
	// 验证分组ID
	if !model.IsValidCacheGroup(groupID) {
		return nil, fmt.Errorf("无效的缓存分组ID: %s", groupID)
	}

	// 获取分组下的所有缓存
	caches, err := s.ListByGroup(ctx, tenantID, groupID)
	if err != nil {
		return nil, err
	}

	// 获取默认分组信息
	groups := model.GetDefaultCacheGroups()
	var groupInfo *model.CacheGroupInfo
	for _, g := range groups {
		if g.GroupID == groupID {
			groupInfo = g
			break
		}
	}
	if groupInfo == nil {
		return nil, fmt.Errorf("未找到分组信息: %s", groupID)
	}

	// 更新统计信息
	groupInfo.KeyCount = int64(len(caches))

	// 计算总内存使用量
	var totalMemory int64
	for _, cache := range caches {
		key := cache.BuildKey()
		memory, err := s.redisClient.MemoryUsage(ctx, key).Result()
		if err != nil {
			continue
		}
		totalMemory += memory
	}
	groupInfo.MemoryUsage = totalMemory

	return groupInfo, nil
}

// ListGroupsWithStats 获取所有分组及其统计信息
func (s *RedisCacheService) ListGroupsWithStats(ctx context.Context, tenantID string) ([]*model.CacheGroupInfo, error) {
	// 获取默认分组列表
	groups := model.GetDefaultCacheGroups()
	result := make([]*model.CacheGroupInfo, len(groups))

	// 获取每个分组的统计信息
	for i, group := range groups {
		stats, err := s.GetGroupStats(ctx, tenantID, group.GroupID)
		if err != nil {
			// 如果获取统计信息失败，使用默认值
			result[i] = group
			continue
		}
		result[i] = stats
	}

	return result, nil
}

// ListGroupKeys 获取分组下的所有键
func (s *RedisCacheService) ListGroupKeys(ctx context.Context, tenantID string, groupID string) ([]string, error) {
	// 获取分组下的所有缓存
	caches, err := s.ListByGroup(ctx, tenantID, groupID)
	if err != nil {
		return nil, err
	}

	// 提取键
	keys := make([]string, len(caches))
	for i, cache := range caches {
		keys[i] = cache.BuildKey()
	}

	return keys, nil
}

// ListByGroup 获取分组下的所有缓存
func (s *RedisCacheService) ListByGroup(ctx context.Context, tenantID string, groupID string) ([]*model.Cache, error) {
	// 验证分组ID
	if !model.IsValidCacheGroup(groupID) {
		return nil, cache_errors.ErrInvalidKey
	}

	// 构建缓存对象用于查询
	cache := &model.Cache{
		TenantID: tenantID,
		GroupID:  groupID,
	}

	// 获取所有键
	pattern := cache.BuildGroupPattern()
	keys, err := s.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, cache_errors.ErrConnectionFailed
	}

	// 获取所有缓存值
	caches := make([]*model.Cache, 0, len(keys))
	for _, key := range keys {
		// 创建缓存对象
		cache := &model.Cache{
			Key:      key,
			TenantID: tenantID,
			GroupID:  groupID,
		}

		// 获取缓存值
		if err := s.Get(ctx, cache); err != nil {
			if errors.Is(err, cache_errors.ErrNil) {
				continue
			}
			return nil, err
		}

		caches = append(caches, cache)
	}

	return caches, nil
}

// Clear 清空所有缓存
func (s *RedisCacheService) Clear(ctx context.Context, tenantID string) error {
	// 获取所有分组
	groups := model.GetDefaultCacheGroups()

	// 清空每个分组
	for _, group := range groups {
		if err := s.DeleteByGroup(ctx, tenantID, group.GroupID); err != nil {
			return err
		}
	}

	return nil
}
