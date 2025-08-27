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

// ConfigGroupCommandHandler 配置分组命令处理器
type ConfigGroupCommandHandler struct {
	configGroupRepo repository.IConfigGroupRepository
	cacheSvc        cache_api.InternalCacheService
}

// NewConfigGroupCommandHandler 创建配置分组命令处理器
func NewConfigGroupCommandHandler(
	configGroupRepo repository.IConfigGroupRepository,
	cacheSvc cache_api.InternalCacheService,
) *ConfigGroupCommandHandler {
	return &ConfigGroupCommandHandler{
		configGroupRepo: configGroupRepo,
		cacheSvc:        cacheSvc,
	}
}

// HandleCreate 处理创建配置分组命令
func (h *ConfigGroupCommandHandler) HandleCreate(ctx context.Context, cmd commands.CreateConfigGroupCommand) herrors.Herr {
	// 检查配置分组编码是否已存在
	exist, err := h.configGroupRepo.FindByCode(ctx, cmd.Code)
	if err != nil {
		hlog.CtxErrorf(ctx, "Check config group code exist error: %v", err)
		return errors.GetConfigGroupFail(err)
	}
	if exist != nil {
		return errors.ConfigGroupCodeExistFail
	}

	// 创建配置分组实体
	group := &entity.ConfigGroup{
		ID:          cmd.Code,
		Name:        cmd.Name,
		Code:        cmd.Code,
		Description: cmd.Description,
		I18nKey:     cmd.I18nKey,
		IsSystem:    cmd.IsSystem,
		IsEnabled:   cmd.IsEnabled,
		Sort:        cmd.Sort,
	}

	// 保存到数据库
	if _, err = h.configGroupRepo.Add(ctx, group); err != nil {
		hlog.CtxErrorf(ctx, "Save config group error: %v", err)
		return errors.AddConfigGroupFail(err)
	}

	// 设置缓存
	if err := h.setConfigGroupCache(ctx, group); err != nil {
		hlog.CtxErrorf(ctx, "Set config group cache error: %v", err)
	}

	return nil
}

// HandleUpdate 处理更新配置分组命令
func (h *ConfigGroupCommandHandler) HandleUpdate(ctx context.Context, cmd commands.UpdateConfigGroupCommand) herrors.Herr {
	// 查找配置分组
	group, err := h.configGroupRepo.FindById(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config group error: %v", err)
		return errors.GetConfigGroupFail(err)
	}
	if group == nil {
		return errors.ConfigGroupNotExistFail
	}

	// 检查配置分组编码是否已存在
	if group.Code != cmd.Code {
		exist, err := h.configGroupRepo.FindByCode(ctx, cmd.Code)
		if err != nil {
			hlog.CtxErrorf(ctx, "Check config group code exist error: %v", err)
			return errors.GetConfigGroupFail(err)
		}
		if exist != nil {
			return errors.ConfigGroupCodeExistFail
		}
	}

	// 更新字段
	group.Name = cmd.Name
	group.Code = cmd.Code
	group.Description = cmd.Description
	group.I18nKey = cmd.I18nKey
	group.IsSystem = cmd.IsSystem
	group.IsEnabled = cmd.IsEnabled
	group.Sort = cmd.Sort
	group.UpdatedAt = utils.GetDateUnix()

	// 保存到数据库
	if err = h.configGroupRepo.EditById(ctx, cmd.ID, group); err != nil {
		hlog.CtxErrorf(ctx, "Save config group error: %v", err)
		return errors.EditConfigGroupFail(err)
	}

	// 更新缓存
	if err := h.setConfigGroupCache(ctx, group); err != nil {
		hlog.CtxErrorf(ctx, "Set config group cache error: %v", err)
	}

	return nil
}

// HandleDelete 处理删除配置分组命令
func (h *ConfigGroupCommandHandler) HandleDelete(ctx context.Context, cmd commands.DeleteConfigGroupCommand) herrors.Herr {
	// 查找配置分组
	group, err := h.configGroupRepo.FindById(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config group error: %v", err)
		return errors.GetConfigGroupFail(err)
	}
	if group == nil {
		return errors.ConfigGroupNotExistFail
	}

	// 删除配置分组
	if err = h.configGroupRepo.DelByIdUnScoped(ctx, cmd.ID); err != nil {
		hlog.CtxErrorf(ctx, "Delete config group error: %v", err)
		return errors.DeleteConfigGroupFail(err)
	}

	// 删除缓存
	if err = h.deleteConfigGroupCache(ctx, group.Code); err != nil {
		hlog.CtxErrorf(ctx, "Delete config group cache error: %v", err)
		return errors.DeleteGroupCacheFail(err)
	}

	return nil
}

// HandleUpdateStatus 处理更新配置分组状态命令
func (h *ConfigGroupCommandHandler) HandleUpdateStatus(ctx context.Context, cmd commands.UpdateConfigGroupStatusCommand) herrors.Herr {
	// 查找配置分组
	group, err := h.configGroupRepo.FindById(ctx, cmd.ID)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find config group error: %v", err)
		return errors.GetConfigGroupFail(err)
	}
	if group == nil {
		return errors.ConfigGroupNotExistFail
	}

	// 更新状态
	group.IsEnabled = cmd.IsEnabled
	group.UpdatedAt = utils.GetDateUnix()

	// 保存到数据库
	if err := h.configGroupRepo.EditById(ctx, group.ID, group); err != nil {
		hlog.CtxErrorf(ctx, "Save config group error: %v", err)
		return errors.EditConfigGroupFail(err)
	}

	// 更新缓存
	if err := h.setConfigGroupCache(ctx, group); err != nil {
		hlog.CtxErrorf(ctx, "Set config group cache error: %v", err)
		return errors.SetGroupCacheFail(err)
	}

	return nil
}

// setConfigGroupCache 设置配置分组缓存
func (h *ConfigGroupCommandHandler) setConfigGroupCache(ctx context.Context, group *entity.ConfigGroup) herrors.Herr {
	key := model.CacheKey(model.CacheTypeConfigGroup, group.Code)
	err := h.cacheSvc.SetWithGroup(ctx, actx.GetTenantId(ctx), chacmodel.CacheGroupConfig, key, group, 0)
	if err != nil {
		return errors.SetGroupCacheFail(err)
	}
	return nil
}

// deleteConfigGroupCache 删除配置分组缓存
func (h *ConfigGroupCommandHandler) deleteConfigGroupCache(ctx context.Context, code string) herrors.Herr {
	if err := h.cacheSvc.Delete(ctx, actx.GetTenantId(ctx), code); err != nil {
		return errors.DeleteGroupCacheFail(err)
	}
	return nil
}
