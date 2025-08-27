package handlers

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	chacmodel "github.com/flare-admin/flare-server-go/framework/support/cache/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/cache/interfaces/api"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/converter"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/dto"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/queries"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/domain/errors"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/infrastructure/entity"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/infrastructure/repository"
)

// ConfigGroupQueryHandler 配置分组查询处理器
type ConfigGroupQueryHandler struct {
	configGroupRepo repository.IConfigGroupRepository
	cacheSvc        cache_api.InternalCacheService
}

// NewConfigGroupQueryHandler 创建配置分组查询处理器
func NewConfigGroupQueryHandler(
	configGroupRepo repository.IConfigGroupRepository,
	cacheSvc cache_api.InternalCacheService,
) *ConfigGroupQueryHandler {
	return &ConfigGroupQueryHandler{
		configGroupRepo: configGroupRepo,
		cacheSvc:        cacheSvc,
	}
}

// HandleGet 处理获取配置分组查询
func (h *ConfigGroupQueryHandler) HandleGet(ctx context.Context, query queries.GetConfigGroupQuery) (*dto.ConfigGroupDTO, herrors.Herr) {
	// 先从缓存获取
	cache, err := h.cacheSvc.GetWithGroup(ctx, actx.GetTenantId(ctx), chacmodel.CacheGroupConfig, model.CacheKey(model.CacheTypeConfigGroup, query.ID))
	if err == nil && cache != nil {
		return cache.Value.(*dto.ConfigGroupDTO), nil
	}

	// 缓存未命中，从数据库获取
	group, err := h.configGroupRepo.FindById(ctx, query.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config group error: %v", err)
		return nil, errors.GetConfigGroupFail(err)
	}
	if group == nil {
		return nil, errors.ConfigGroupNotExistFail
	}

	// 转换为DTO
	dto := converter.ToConfigGroupDTO(group)

	// 设置缓存
	cache = &chacmodel.Cache{
		Key:      model.CacheKey(model.CacheTypeConfigGroup, group.ID),
		Value:    dto,
		TenantID: actx.GetTenantId(ctx),
		ExpireAt: 0,
	}
	if err := h.setConfigGroupCache(ctx, group); err != nil {
		hlog.CtxErrorf(ctx, "Set config group cache error: %v", err)
		return nil, errors.SetGroupCacheFail(err)
	}

	return dto, nil
}

// HandleList 处理获取配置分组列表查询
func (h *ConfigGroupQueryHandler) HandleList(ctx context.Context, query queries.ListConfigGroupsQuery) ([]*dto.ConfigGroupDTO, int64, herrors.Herr) {
	quer := db_query.NewQueryBuilder()
	if query.Name != "" {
		quer.Where("name", db_query.Like, "%"+query.Name+"%")
	}
	if query.Code != "" {
		quer.Where("code", db_query.Eq, query.Code)
	}
	if query.IsSystem != nil {
		quer.Where("is_system", db_query.Eq, *query.IsSystem)
	}
	if query.IsEnabled != nil {
		quer.Where("is_enabled", db_query.Eq, *query.IsEnabled)
	}

	// 查询总数
	total, err := h.configGroupRepo.Count(ctx, quer)
	if err != nil {
		hlog.CtxErrorf(ctx, "Count config groups error: %v", err)
		return nil, 0, errors.GetConfigGroupFail(err)
	}

	// 分页查询
	quer.WithPage(&query.Page)
	quer.OrderBy("sort", true)
	groups, err := h.configGroupRepo.Find(ctx, quer)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config groups error: %v", err)
		return nil, 0, errors.GetConfigGroupFail(err)
	}

	// 转换为DTO列表
	dtos := make([]*dto.ConfigGroupDTO, 0, len(groups))
	for _, group := range groups {
		dtos = append(dtos, converter.ToConfigGroupDTO(group))
	}

	return dtos, total, nil
}

// setConfigGroupCache 设置配置分组缓存
func (h *ConfigGroupQueryHandler) setConfigGroupCache(ctx context.Context, group *entity.ConfigGroup) herrors.Herr {
	key := model.CacheKey(model.CacheTypeConfigGroup, group.Code)
	err := h.cacheSvc.SetWithGroup(ctx, actx.GetTenantId(ctx), chacmodel.CacheGroupConfig, key, group, 0)
	if err != nil {
		return errors.SetGroupCacheFail(err)
	}
	return nil
}
