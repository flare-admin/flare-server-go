package handlers

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	chacmodel "github.com/flare-admin/flare-server-go/framework/support/cache/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/cache/interfaces/api"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/commands"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/domain/errors"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/infrastructure/entity"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/infrastructure/repository"
)

// ConfigCommandHandler 配置命令处理器
type ConfigCommandHandler struct {
	configRepo repository.IConfigRepository
	cacheSvc   cache_api.InternalCacheService
}

// NewConfigCommandHandler 创建配置命令处理器
func NewConfigCommandHandler(
	configRepo repository.IConfigRepository,
	cacheSvc cache_api.InternalCacheService,
) *ConfigCommandHandler {
	return &ConfigCommandHandler{
		configRepo: configRepo,
		cacheSvc:   cacheSvc,
	}
}

// HandleCreate 处理创建配置命令
func (h *ConfigCommandHandler) HandleCreate(ctx context.Context, cmd commands.CreateConfigCommand) herrors.Herr {
	// 检查配置键是否已存在
	exist, err := h.configRepo.FindByKey(ctx, cmd.Key)
	if err != nil {
		hlog.CtxErrorf(ctx, "Check config key exist error: %v", err)
		return errors.GetConfigFail(err)
	}
	if exist != nil {
		return errors.ConfigKeyExistFail
	}

	// 创建配置实体
	config := &entity.Config{
		Name:        cmd.Name,
		Key:         cmd.Key,
		Value:       cmd.Value,
		Type:        string(cmd.Type),
		Group:       cmd.Group,
		Description: cmd.Description,
		I18nKey:     cmd.I18nKey,
		IsSystem:    cmd.IsSystem,
		IsEnabled:   cmd.IsEnabled,
		Sort:        cmd.Sort,
	}

	// 保存到数据库
	if _, err = h.configRepo.Add(ctx, config); err != nil {
		hlog.CtxErrorf(ctx, "Save config error: %v", err)
		return errors.AddConfigFail(err)
	}

	// 清理缓存
	h.clearConfigCache(ctx, config)

	return nil
}

// HandleUpdate 处理更新配置命令
func (h *ConfigCommandHandler) HandleUpdate(ctx context.Context, cmd commands.UpdateConfigCommand) herrors.Herr {
	// 查找配置
	config, err := h.configRepo.FindById(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config error: %v", err)
		return errors.GetConfigFail(err)
	}
	if config == nil {
		return errors.ConfigNotExistFail
	}

	// 检查配置键是否已存在
	if config.Key != cmd.Key {
		exist, err := h.configRepo.FindByKey(ctx, cmd.Key)
		if err != nil {
			hlog.CtxErrorf(ctx, "Check config key exist error: %v", err)
			return errors.GetConfigFail(err)
		}
		if exist != nil {
			return errors.ConfigKeyExistFail
		}
	}

	// 更新字段
	config.Name = cmd.Name
	config.Key = cmd.Key
	config.Value = cmd.Value
	config.Type = string(cmd.Type)
	config.Group = cmd.Group
	config.Description = cmd.Description
	config.I18nKey = cmd.I18nKey
	config.IsSystem = cmd.IsSystem
	config.IsEnabled = cmd.IsEnabled
	config.Sort = cmd.Sort
	config.UpdatedAt = utils.GetDateUnix()

	// 保存到数据库
	if err = h.configRepo.EditById(ctx, cmd.ID, config); err != nil {
		hlog.CtxErrorf(ctx, "Save config error: %v", err)
		return errors.EditConfigFail(err)
	}

	// 清理缓存
	h.clearConfigCache(ctx, config)

	return nil
}

// HandleDelete 处理删除配置命令
func (h *ConfigCommandHandler) HandleDelete(ctx context.Context, cmd commands.DeleteConfigCommand) herrors.Herr {
	// 查找配置
	config, err := h.configRepo.FindById(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config error: %v", err)
		return errors.GetConfigFail(err)
	}
	if config == nil {
		return errors.ConfigNotExistFail
	}

	// 删除配置
	if err = h.configRepo.DelByIdUnScoped(ctx, cmd.ID); err != nil {
		hlog.CtxErrorf(ctx, "Delete config error: %v", err)
		return errors.DeleteConfigFail(err)
	}

	// 清理缓存
	h.clearConfigCache(ctx, config)

	return nil
}

// HandleUpdateStatus 处理更新配置状态命令
func (h *ConfigCommandHandler) HandleUpdateStatus(ctx context.Context, cmd commands.UpdateConfigStatusCommand) herrors.Herr {
	// 查找配置
	config, err := h.configRepo.FindById(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config error: %v", err)
		return errors.GetConfigFail(err)
	}
	if config == nil {
		return errors.ConfigNotExistFail
	}

	// 更新状态
	config.IsEnabled = cmd.IsEnabled
	config.UpdatedAt = utils.GetDateUnix()

	// 保存到数据库
	if err := h.configRepo.EditById(ctx, config.ID, config); err != nil {
		hlog.CtxErrorf(ctx, "Save config error: %v", err)
		return errors.EditConfigFail(err)
	}

	// 清理缓存
	h.clearConfigCache(ctx, config)

	return nil
}

// HandleBatchUpdate 处理批量更新配置命令
func (h *ConfigCommandHandler) HandleBatchUpdate(ctx context.Context, cmd commands.BatchUpdateConfigCommand) herrors.Herr {
	// 准备批量更新的配置列表
	configs := make([]*entity.Config, 0, len(cmd.Configs))
	for _, config := range cmd.Configs {
		// 查找配置
		exist, err := h.configRepo.FindById(ctx, config.ID)
		if err != nil {
			hlog.CtxErrorf(ctx, "Find config error: %v", err)
			return errors.GetConfigFail(err)
		}
		if exist == nil {
			return errors.ConfigNotExistFail
		}

		// 更新字段
		exist.Value = config.Value
		exist.UpdatedAt = utils.GetDateUnix()
		configs = append(configs, exist)
	}

	// 批量更新配置
	if err := h.configRepo.BatchUpdate(ctx, configs); err != nil {
		hlog.CtxErrorf(ctx, "Batch update config error: %v", err)
		return errors.EditConfigFail(err)
	}

	// 清理缓存
	for _, config := range configs {
		h.clearConfigCache(ctx, config)
	}

	return nil
}

// setConfigCache 设置配置缓存
func (h *ConfigCommandHandler) setConfigCache(ctx context.Context, config *entity.Config) error {
	key := model.CacheKey(model.CacheTypeConfig, config.Key)
	return h.cacheSvc.SetWithGroup(ctx, actx.GetTenantId(ctx), chacmodel.CacheGroupConfig, key, config, 0)
}

// clearConfigCache 清理配置缓存
func (h *ConfigCommandHandler) clearConfigCache(ctx context.Context, config *entity.Config) {
	// 清理配置值缓存
	_ = h.deleteConfigCache(ctx, model.CacheKey(model.CacheTypeConfig, config.Key))
	// 清理配置详情缓存
	_ = h.deleteConfigCache(ctx, model.CacheKey(model.CacheTypeConfig, config.ID))
	// 清理分组配置缓存
	_ = h.deleteConfigCache(ctx, model.CacheKey(model.CacheTypeConfig, config.Group))
}

// deleteConfigCache 删除配置缓存
func (h *ConfigCommandHandler) deleteConfigCache(ctx context.Context, key string) herrors.Herr {
	if err := h.cacheSvc.DeleteWithGroup(ctx, actx.GetTenantId(ctx), chacmodel.CacheGroupConfig, key); err != nil {
		return errors.DeleteConfigCacheFail(err)
	}
	return nil
}
