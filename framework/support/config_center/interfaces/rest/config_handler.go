package rest

import (
	"context"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/dto"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	_ "github.com/flare-admin/flare-server-go/framework/pkg/hserver/base_info"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/casbin"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/jwt"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/oplog"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/commands"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/handlers"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/queries"
)

// ConfigHandler 配置处理器
type ConfigHandler struct {
	cmdHandler        *handlers.ConfigCommandHandler
	queryHandler      *handlers.ConfigQueryHandler
	groupCmdHandler   *handlers.ConfigGroupCommandHandler
	groupQueryHandler *handlers.ConfigGroupQueryHandler
	ef                *casbin.Enforcer
	moduleName        string
}

// NewConfigHandler 创建配置处理器
func NewConfigHandler(
	cmdHandler *handlers.ConfigCommandHandler,
	queryHandler *handlers.ConfigQueryHandler,
	groupCmdHandler *handlers.ConfigGroupCommandHandler,
	groupQueryHandler *handlers.ConfigGroupQueryHandler,
	ef *casbin.Enforcer,
) *ConfigHandler {
	return &ConfigHandler{
		cmdHandler:        cmdHandler,
		queryHandler:      queryHandler,
		groupCmdHandler:   groupCmdHandler,
		groupQueryHandler: groupQueryHandler,
		ef:                ef,
		moduleName:        "配置管理",
	}
}

// RegisterRouter 注册路由
func (h *ConfigHandler) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	v1 := g.Group("/v1")
	cr := v1.Group("/config", jwt.Handler(t))
	{
		// 配置管理
		cr.POST("", casbin.Handler(h.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      h.moduleName,
			Action:      "新增",
		}), hserver.NewHandlerFu[commands.CreateConfigCommand](h.CreateConfig))
		cr.GET("", casbin.Handler(h.ef), hserver.NewHandlerFu[queries.ListConfigsQuery](h.ListConfigs))
		cr.GET("/:id", casbin.Handler(h.ef), hserver.NewHandlerFu[models.StringIdReq](h.GetConfig))
		cr.PUT("", casbin.Handler(h.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      h.moduleName,
			Action:      "修改",
		}), hserver.NewHandlerFu[commands.UpdateConfigCommand](h.UpdateConfig))
		cr.DELETE("/:id", casbin.Handler(h.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      h.moduleName,
			Action:      "删除",
		}), hserver.NewHandlerFu[models.StringIdReq](h.DeleteConfig))
		cr.PUT("/status", casbin.Handler(h.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      h.moduleName,
			Action:      "更新状态",
		}), hserver.NewHandlerFu[commands.UpdateConfigStatusCommand](h.UpdateConfigStatus))

		// 配置分组管理
		cr.POST("/group", casbin.Handler(h.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      h.moduleName,
			Action:      "新增分组",
		}), hserver.NewHandlerFu[commands.CreateConfigGroupCommand](h.CreateConfigGroup))
		cr.GET("/group", casbin.Handler(h.ef), hserver.NewHandlerFu[queries.ListConfigGroupsQuery](h.ListConfigGroups))
		cr.GET("/group/:id", casbin.Handler(h.ef), hserver.NewHandlerFu[models.StringIdReq](h.GetConfigGroup))
		cr.PUT("/group", casbin.Handler(h.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      h.moduleName,
			Action:      "修改分组",
		}), hserver.NewHandlerFu[commands.UpdateConfigGroupCommand](h.UpdateConfigGroup))
		cr.DELETE("/group/:id", casbin.Handler(h.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      h.moduleName,
			Action:      "删除分组",
		}), hserver.NewHandlerFu[models.StringIdReq](h.DeleteConfigGroup))
		cr.PUT("/group/status", casbin.Handler(h.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      h.moduleName,
			Action:      "更新分组状态",
		}), hserver.NewHandlerFu[commands.UpdateConfigGroupStatusCommand](h.UpdateConfigGroupStatus))

		// 根据分组ID查询配置
		cr.GET("/group/:id/configs", casbin.Handler(h.ef), hserver.NewHandlerFu[models.StringIdReq](h.GetConfigsByGroupId))
		// 批量更新配置
		cr.PUT("/batch", casbin.Handler(h.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      h.moduleName,
			Action:      "批量更新配置",
		}), hserver.NewHandlerFu[commands.BatchUpdateConfigCommand](h.BatchUpdateConfig))
	}
}

// CreateConfig 创建配置
// @Summary 创建配置
// @Description 创建配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param req body commands.CreateConfigCommand true "配置信息"
// @Success 200 {object} base_info.Success
// @Router /v1/config [post]
func (h *ConfigHandler) CreateConfig(ctx context.Context, params *commands.CreateConfigCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := h.cmdHandler.HandleCreate(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// GetConfig 获取配置详情
// @Summary 获取配置详情
// @Description 获取配置详情
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param id path string true "配置ID"
// @Success 200 {object} base_info.Success{data=dto.ConfigDTO}
// @Router /v1/config/{id} [get]
func (h *ConfigHandler) GetConfig(ctx context.Context, params *models.StringIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := h.queryHandler.HandleGet(ctx, queries.GetConfigQuery{ID: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// ListConfigs 获取配置列表
// @Summary 获取配置列表
// @Description 获取配置列表
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param req query queries.ListConfigsQuery true "查询参数"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.ConfigDTO]{list=[]dto.ConfigDTO}}
// @Router /v1/config [get]
func (h *ConfigHandler) ListConfigs(ctx context.Context, params *queries.ListConfigsQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, total, err := h.queryHandler.HandleList(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data).WithData(&models.PageRes[dto.ConfigDTO]{
		Total: total,
		List:  data,
	})
}

// UpdateConfig 更新配置
// @Summary 更新配置
// @Description 更新配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param req body commands.UpdateConfigCommand true "配置信息"
// @Success 200 {object} base_info.Success
// @Router /v1/config [put]
func (h *ConfigHandler) UpdateConfig(ctx context.Context, params *commands.UpdateConfigCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := h.cmdHandler.HandleUpdate(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// DeleteConfig 删除配置
// @Summary 删除配置
// @Description 删除配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param id path string true "配置ID"
// @Success 200 {object} base_info.Success
// @Router /v1/config/{id} [delete]
func (h *ConfigHandler) DeleteConfig(ctx context.Context, params *models.StringIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := h.cmdHandler.HandleDelete(ctx, commands.DeleteConfigCommand{ID: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// UpdateConfigStatus 更新配置状态
// @Summary 更新配置状态
// @Description 更新配置状态
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param req body commands.UpdateConfigStatusCommand true "状态信息"
// @Success 200 {object} base_info.Success
// @Router /v1/config/status [put]
func (h *ConfigHandler) UpdateConfigStatus(ctx context.Context, params *commands.UpdateConfigStatusCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := h.cmdHandler.HandleUpdateStatus(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// CreateConfigGroup 创建配置分组
// @Summary 创建配置分组
// @Description 创建配置分组
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param req body commands.CreateConfigGroupCommand true "分组信息"
// @Success 200 {object} base_info.Success
// @Router /v1/config/group [post]
func (h *ConfigHandler) CreateConfigGroup(ctx context.Context, params *commands.CreateConfigGroupCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := h.groupCmdHandler.HandleCreate(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// GetConfigGroup 获取配置分组详情
// @Summary 获取配置分组详情
// @Description 获取配置分组详情
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param id path string true "分组ID"
// @Success 200 {object} base_info.Success{data=dto.ConfigGroupDTO}
// @Router /v1/config/group/{id} [get]
func (h *ConfigHandler) GetConfigGroup(ctx context.Context, params *models.StringIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := h.groupQueryHandler.HandleGet(ctx, queries.GetConfigGroupQuery{ID: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// ListConfigGroups 获取配置分组列表
// @Summary 获取配置分组列表
// @Description 获取配置分组列表
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param req query queries.ListConfigGroupsQuery true "查询参数"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.ConfigGroupDTO]{list=[]dto.ConfigGroupDTO}}
// @Router /v1/config/group [get]
func (h *ConfigHandler) ListConfigGroups(ctx context.Context, params *queries.ListConfigGroupsQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, total, err := h.groupQueryHandler.HandleList(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data).WithData(&models.PageRes[dto.ConfigGroupDTO]{
		Total: total,
		List:  data,
	})
}

// UpdateConfigGroup 更新配置分组
// @Summary 更新配置分组
// @Description 更新配置分组
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param req body commands.UpdateConfigGroupCommand true "分组信息"
// @Success 200 {object} base_info.Success
// @Router /v1/config/group [put]
func (h *ConfigHandler) UpdateConfigGroup(ctx context.Context, params *commands.UpdateConfigGroupCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := h.groupCmdHandler.HandleUpdate(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// DeleteConfigGroup 删除配置分组
// @Summary 删除配置分组
// @Description 删除配置分组
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param id path string true "分组ID"
// @Success 200 {object} base_info.Success
// @Router /v1/config/group/{id} [delete]
func (h *ConfigHandler) DeleteConfigGroup(ctx context.Context, params *models.StringIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := h.groupCmdHandler.HandleDelete(ctx, commands.DeleteConfigGroupCommand{ID: params.Id})
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// UpdateConfigGroupStatus 更新配置分组状态
// @Summary 更新配置分组状态
// @Description 更新配置分组状态
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param req body commands.UpdateConfigGroupStatusCommand true "状态信息"
// @Success 200 {object} base_info.Success
// @Router /v1/config/group/status [put]
func (h *ConfigHandler) UpdateConfigGroupStatus(ctx context.Context, params *commands.UpdateConfigGroupStatusCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := h.groupCmdHandler.HandleUpdateStatus(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// GetConfigsByGroupId 根据分组ID获取配置列表
// @Summary 根据分组ID获取配置列表
// @Description 根据分组ID获取配置列表
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param id path string true "分组ID"
// @Success 200 {object} base_info.Success{data=[]dto.ConfigDTO}
// @Router /v1/config/group/{id}/configs [get]
func (h *ConfigHandler) GetConfigsByGroupId(ctx context.Context, params *models.StringIdReq) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := h.queryHandler.HandleGetByGroupId(ctx, params.Id)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// BatchUpdateConfig 批量更新配置
// @Summary 批量更新配置
// @Description 批量更新配置
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param req body commands.BatchUpdateConfigCommand true "配置信息"
// @Success 200 {object} base_info.Success
// @Router /v1/config/batch [put]
func (h *ConfigHandler) BatchUpdateConfig(ctx context.Context, params *commands.BatchUpdateConfigCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := h.cmdHandler.HandleBatchUpdate(ctx, *params)
	if err != nil {
		return result.WithError(err)
	}
	return result
}
