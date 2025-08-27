package service

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/cache/application/commands"
	"github.com/flare-admin/flare-server-go/framework/support/cache/application/queries"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/service"
)

// CacheService 缓存应用服务
type CacheService struct {
	cacheService service.CacheService
}

// NewCacheService 创建缓存应用服务
func NewCacheService(cacheService service.CacheService) *CacheService {
	return &CacheService{
		cacheService: cacheService,
	}
}

// GetCache 获取缓存
func (s *CacheService) GetCache(ctx context.Context, query *queries.GetCacheQuery) (interface{}, herrors.Herr) {
	if err := query.Validate(); err != nil {
		return nil, herrors.QueryFail(err)
	}

	var value interface{}
	err := s.cacheService.GetWithGroup(ctx, actx.GetTenantId(ctx), query.GroupID, query.Key, &value)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	return value, nil
}

// DeleteCache 删除缓存
func (s *CacheService) DeleteCache(ctx context.Context, cmd *commands.DeleteCacheCommand) herrors.Herr {
	if err := cmd.Validate(); err != nil {
		return herrors.DeleteFail(err)
	}

	err := s.cacheService.DeleteWithGroup(ctx, actx.GetTenantId(ctx), cmd.GroupID, cmd.Key)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// GetGroup 获取分组信息
func (s *CacheService) GetGroup(ctx context.Context, query *queries.GetGroupQuery) (*model.CacheGroupInfo, herrors.Herr) {
	if err := query.Validate(); err != nil {
		return nil, herrors.QueryFail(err)
	}

	group, err := s.cacheService.GetGroupStats(ctx, actx.GetTenantId(ctx), query.GroupID)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	return group, nil
}

// ListGroups 获取所有分组
func (s *CacheService) ListGroups(ctx context.Context) ([]*model.CacheGroupInfo, herrors.Herr) {
	groups, err := s.cacheService.ListGroupsWithStats(ctx, actx.GetTenantId(ctx))
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	return groups, nil
}

// DeleteGroup 删除分组
func (s *CacheService) DeleteGroup(ctx context.Context, cmd *commands.DeleteGroupCommand) herrors.Herr {
	if err := cmd.Validate(); err != nil {
		return herrors.DeleteFail(err)
	}

	err := s.cacheService.DeleteGroup(ctx, actx.GetTenantId(ctx), cmd.GroupID)
	if err != nil {
		return herrors.NewServerHError(err)
	}
	return nil
}

// ListGroupKeys 获取分组下的所有键
func (s *CacheService) ListGroupKeys(ctx context.Context, query *queries.ListGroupKeysQuery) ([]string, herrors.Herr) {
	if err := query.Validate(); err != nil {
		return nil, herrors.QueryFail(err)
	}

	keys, err := s.cacheService.ListGroupKeys(ctx, actx.GetTenantId(ctx), query.GroupID)
	if err != nil {
		return nil, herrors.NewServerHError(err)
	}
	return keys, nil
}
