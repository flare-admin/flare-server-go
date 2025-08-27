package handlers

import (
	"context"
	"encoding/json"
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
	"time"
)

// ConfigQueryHandler 配置查询处理器
type ConfigQueryHandler struct {
	groupRepo  repository.IConfigGroupRepository
	configRepo repository.IConfigRepository
	cacheSvc   cache_api.InternalCacheService
}

// NewConfigQueryHandler 创建配置查询处理器
func NewConfigQueryHandler(
	groupRepo repository.IConfigGroupRepository,
	configRepo repository.IConfigRepository,
	cacheSvc cache_api.InternalCacheService,
) *ConfigQueryHandler {
	return &ConfigQueryHandler{
		configRepo: configRepo,
		groupRepo:  groupRepo,
		cacheSvc:   cacheSvc,
	}
}

// HandleGet 处理获取配置查询
func (h *ConfigQueryHandler) HandleGet(ctx context.Context, query queries.GetConfigQuery) (*dto.ConfigDTO, herrors.Herr) {
	// 先从缓存获取
	cache, err := h.cacheSvc.Get(ctx, actx.GetTenantId(ctx), model.CacheKey(model.CacheTypeConfig, query.ID))
	if err == nil && cache != nil {
		var dto *dto.ConfigDTO
		if err := cache.GetValue(&dto); err != nil {
			hlog.CtxErrorf(ctx, "Get cache value error: %v", err)
		} else if dto != nil {
			return dto, nil
		}
	}

	// 缓存未命中，从数据库获取
	config, err := h.configRepo.FindById(ctx, query.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config error: %v", err)
		return nil, errors.GetConfigFail(err)
	}
	if config == nil {
		return nil, errors.ConfigNotExistFail
	}

	// 转换为DTO
	dto := converter.ToConfigDTO(config)

	// 设置缓存
	cache = &chacmodel.Cache{
		Key:      model.CacheKey(model.CacheTypeConfig, config.ID),
		TenantID: actx.GetTenantId(ctx),
		GroupID:  chacmodel.CacheGroupConfig,
		ExpireAt: 0,
	}
	if err := cache.SetValue(dto); err != nil {
		hlog.CtxErrorf(ctx, "Set cache value error: %v", err)
		return nil, errors.SetConfigCacheFail(err)
	}
	if err := h.cacheSvc.SetCache(ctx, cache); err != nil {
		hlog.CtxErrorf(ctx, "Set config cache error: %v", err)
		return nil, errors.SetConfigCacheFail(err)
	}

	return dto, nil
}

// HandleList 处理获取配置列表查询
func (h *ConfigQueryHandler) HandleList(ctx context.Context, query queries.ListConfigsQuery) ([]*dto.ConfigDTO, int64, herrors.Herr) {
	quer := db_query.NewQueryBuilder()
	if query.Key != "" {
		quer.Where("key", db_query.Like, "%"+query.Key+"%")
	}
	if query.Type != "" {
		quer.Where("type", db_query.Eq, query.Type)
	}
	if query.Group != "" {
		quer.Where("group_id", db_query.Eq, query.Group)
	}
	if query.IsSystem != nil {
		quer.Where("is_system", db_query.Eq, *query.IsSystem)
	}
	if query.IsEnabled != nil {
		quer.Where("is_enabled", db_query.Eq, *query.IsEnabled)
	}
	if query.Group != "" {
		quer.Where("group_id", db_query.Eq, query.Group)
	}

	// 查询总数
	total, err := h.configRepo.Count(ctx, quer)
	if err != nil {
		hlog.CtxErrorf(ctx, "Count configs error: %v", err)
		return nil, 0, errors.GetConfigFail(err)
	}

	// 分页查询
	quer.WithPage(&query.Page)
	quer.OrderBy("sort", true)
	configs, err := h.configRepo.Find(ctx, quer)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find configs error: %v", err)
		return nil, 0, errors.GetConfigFail(err)
	}

	// 转换为DTO列表
	dtos := make([]*dto.ConfigDTO, 0, len(configs))
	for _, config := range configs {
		dtos = append(dtos, converter.ToConfigDTO(config))
	}

	return dtos, total, nil
}

// HandleGetValue 处理获取配置值查询
func (h *ConfigQueryHandler) HandleGetValue(ctx context.Context, query queries.GetConfigValueQuery) (interface{}, herrors.Herr) {
	// 先从缓存获取
	cache, err := h.cacheSvc.Get(ctx, actx.GetTenantId(ctx), model.CacheKey(model.CacheTypeConfig, query.Key))
	if err == nil && cache != nil {
		var value interface{}
		if err := cache.GetValue(&value); err != nil {
			hlog.CtxErrorf(ctx, "Get cache value error: %v", err)
		} else if value != nil {
			return value, nil
		}
	}

	// 缓存未命中，从数据库获取
	config, err := h.configRepo.FindByKey(ctx, query.Key)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config error: %v", err)
		return query.DefaultValue, errors.GetConfigFail(err)
	}
	if config == nil {
		return query.DefaultValue, nil
	}
	if !config.IsEnabled {
		return query.DefaultValue, errors.ConfigNotEnableFail
	}

	// 解析配置值
	value, err := h.parseConfigValue(config, query.DefaultValue)
	if err != nil {
		return query.DefaultValue, errors.GetConfigFail(err)
	}

	// 设置缓存
	cache = &chacmodel.Cache{
		Key:      model.CacheKey(model.CacheTypeConfig, config.Key),
		TenantID: actx.GetTenantId(ctx),
		GroupID:  chacmodel.CacheGroupConfig,
		ExpireAt: 0,
	}
	if err := cache.SetValue(value); err != nil {
		hlog.CtxErrorf(ctx, "Set cache value error: %v", err)
		return nil, errors.SetConfigCacheFail(err)
	}
	if err := h.cacheSvc.SetCache(ctx, cache); err != nil {
		hlog.CtxErrorf(ctx, "Set config cache error: %v", err)
		return nil, errors.SetConfigCacheFail(err)
	}

	return value, nil
}

// HandleGetValueMap 处理获取配置值映射查询
func (h *ConfigQueryHandler) HandleGetValueMap(ctx context.Context, query queries.GetConfigValueMapQuery) (map[string]interface{}, herrors.Herr) {
	result := make(map[string]interface{})
	for _, key := range query.Keys {
		value, err := h.HandleGetValue(ctx, queries.GetConfigValueQuery{
			Key:          key,
			DefaultValue: query.DefaultValue,
		})
		if err != nil {
			continue
		}
		result[key] = value
	}
	return result, nil
}

// HandleGetValueByGroupCode 处理获取配置值查询
func (h *ConfigQueryHandler) HandleGetValueByGroupCode(ctx context.Context, code string) ([]*dto.ConfigDTO, herrors.Herr) {
	// 先从缓存获取
	cache, err := h.cacheSvc.Get(ctx, actx.GetTenantId(ctx), model.CacheKey(model.CacheTypeConfig, code))
	if err == nil && cache != nil {
		var dtos []*dto.ConfigDTO
		if err := cache.GetValue(&dtos); err != nil {
			hlog.CtxErrorf(ctx, "Get cache value error: %v", err)
		} else if dtos != nil {
			return dtos, nil
		}
	}

	// 缓存未命中，从数据库获取
	// 1. 先获取分组信息
	group, err := h.groupRepo.FindByCode(ctx, code)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config group error: %v", err)
		return nil, errors.GetConfigGroupFail(err)
	}
	if group == nil {
		return nil, errors.ConfigGroupNotExistFail
	}
	if !group.IsEnabled {
		return nil, errors.ConfigGroupNotEnableFail
	}

	// 2. 获取分组下的所有配置
	query := db_query.NewQueryBuilder()
	query = query.Where("group_id", db_query.Eq, group.ID)
	query = query.Where("is_enabled", db_query.Eq, true)
	configs, err := h.configRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find configs by group error: %v", err)
		return nil, errors.GetConfigFail(err)
	}

	// 3. 转换为DTO列表
	dtos := make([]*dto.ConfigDTO, 0, len(configs))
	for _, config := range configs {
		dtos = append(dtos, converter.ToConfigDTO(config))
	}

	// 4. 设置缓存
	cache = &chacmodel.Cache{
		Key:      model.CacheKey(model.CacheTypeConfig, code),
		TenantID: actx.GetTenantId(ctx),
		GroupID:  chacmodel.CacheGroupConfig,
		ExpireAt: 0,
	}
	if err := cache.SetValue(dtos); err != nil {
		hlog.CtxErrorf(ctx, "Set cache value error: %v", err)
		return nil, errors.SetConfigCacheFail(err)
	}
	if err := h.cacheSvc.SetCache(ctx, cache); err != nil {
		hlog.CtxErrorf(ctx, "Set config cache error: %v", err)
		return nil, errors.SetConfigCacheFail(err)
	}

	return dtos, nil
}

// HandleGetByGroupId 处理根据分组ID获取配置列表查询
func (h *ConfigQueryHandler) HandleGetByGroupId(ctx context.Context, groupId string) ([]*dto.ConfigDTO, herrors.Herr) {
	// 先从缓存获取
	cache, err := h.cacheSvc.Get(ctx, actx.GetTenantId(ctx), model.CacheKey(model.CacheTypeConfig, groupId))
	if err == nil && cache != nil {
		var dtos []*dto.ConfigDTO
		if err := cache.GetValue(&dtos); err != nil {
			hlog.CtxErrorf(ctx, "Get cache value error: %v", err)
		} else if dtos != nil {
			return dtos, nil
		}
	}

	// 缓存未命中，从数据库获取
	// 1. 先获取分组信息
	group, err := h.groupRepo.FindById(ctx, groupId)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config group error: %v", err)
		return nil, errors.GetConfigGroupFail(err)
	}
	if group == nil {
		return nil, errors.ConfigGroupNotExistFail
	}
	if !group.IsEnabled {
		return nil, errors.ConfigGroupNotEnableFail
	}

	// 2. 获取分组下的所有配置
	query := db_query.NewQueryBuilder()
	query = query.Where("group_id", db_query.Eq, groupId)
	configs, err := h.configRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find configs by group error: %v", err)
		return nil, errors.GetConfigFail(err)
	}

	// 3. 转换为DTO列表
	dtos := make([]*dto.ConfigDTO, 0, len(configs))
	for _, config := range configs {
		dtos = append(dtos, converter.ToConfigDTO(config))
	}

	// 4. 设置缓存
	cache = &chacmodel.Cache{
		Key:      model.CacheKey(model.CacheTypeConfig, groupId),
		TenantID: actx.GetTenantId(ctx),
		GroupID:  chacmodel.CacheGroupConfig,
		ExpireAt: 0,
	}
	if err := cache.SetValue(dtos); err != nil {
		hlog.CtxErrorf(ctx, "Set cache value error: %v", err)
		return nil, errors.SetConfigCacheFail(err)
	}
	if err := h.cacheSvc.SetCache(ctx, cache); err != nil {
		hlog.CtxErrorf(ctx, "Set config cache error: %v", err)
		return nil, errors.SetConfigCacheFail(err)
	}

	return dtos, nil
}

// parseConfigValue 解析配置值
func (h *ConfigQueryHandler) parseConfigValue(config *entity.Config, defaultValue interface{}) (interface{}, herrors.Herr) {
	switch model.ConfigType(config.Type) {
	case model.ConfigTypeString, model.ConfigTypeRegex:
		return config.Value, nil
	case model.ConfigTypeInt:
		var value int64
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, errors.ConfigTypeInvalidFail
		}
		return value, nil
	case model.ConfigTypeFloat:
		var value float64
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, errors.ConfigTypeInvalidFail
		}
		return value, nil
	case model.ConfigTypeBool:
		var value bool
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, errors.ConfigTypeInvalidFail
		}
		return value, nil
	case model.ConfigTypeJSON:
		var value interface{}
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, errors.ConfigTypeInvalidFail
		}
		return value, nil
	case model.ConfigTypeArray:
		var value []interface{}
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, errors.ConfigTypeInvalidFail
		}
		return value, nil
	case model.ConfigTypeObject:
		var value map[string]interface{}
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, errors.ConfigTypeInvalidFail
		}
		return value, nil
	case model.ConfigTypeTime:
		var value time.Time
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, errors.ConfigTypeInvalidFail
		}
		return value, nil
	case model.ConfigTypeDate:
		var value time.Time
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, errors.ConfigTypeInvalidFail
		}
		return value.Format("2006-01-02"), nil
	case model.ConfigTypeDateTime:
		var value time.Time
		if err := json.Unmarshal([]byte(config.Value), &value); err != nil {
			return defaultValue, errors.ConfigTypeInvalidFail
		}
		return value.Format("2006-01-02 15:04:05"), nil
	default:
		return defaultValue, errors.ConfigTypeInvalidFail
	}
}
