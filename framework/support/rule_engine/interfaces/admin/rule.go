package admin

import (
	"context"

	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/casbin"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/jwt"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/oplog"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/command"
	commandhandler "github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/command/handler"
	_ "github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/dto"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/queries"
	queryhandler "github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/queries/handler"
)

// RuleService 规则管理服务
type RuleService struct {
	rh *queryhandler.RuleQueryHandler
	rc *commandhandler.RuleCommandHandler
	ef *casbin.Enforcer
}

// NewRuleService 创建规则管理服务
func NewRuleService(
	rh *queryhandler.RuleQueryHandler,
	rc *commandhandler.RuleCommandHandler,
	ef *casbin.Enforcer,
) *RuleService {
	return &RuleService{
		rh: rh,
		rc: rc,
		ef: ef,
	}
}

// RegisterRouter 注册路由
func (rs *RuleService) RegisterRouter(rg *route.RouterGroup, t token.IToken) {
	g := rg.Group("/v1/rule-engine/rule", jwt.Handler(t))
	{
		g.POST("", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则管理",
			Action:      "新建",
		}), hserver.NewHandlerFu[command.CreateRuleCommand](rs.CreateRule)) // 新建规则

		g.GET("", hserver.NewHandlerFu[queries.GetRuleListReq](rs.GetRuleList)) // 获取规则列表

		g.GET("/:id", hserver.NewHandlerFu[models.StringIdReq](rs.GetRuleInfo)) // 获取规则信息

		g.DELETE("/:id", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则管理",
			Action:      "删除",
		}), hserver.NewHandlerFu[models.StringIdReq](rs.DeleteRule)) // 删除规则

		g.PUT("", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则管理",
			Action:      "修改",
		}), hserver.NewHandlerFu[command.UpdateRuleCommand](rs.UpdateRule)) // 修改规则信息

		g.PUT("/status", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则管理",
			Action:      "更新状态",
		}), hserver.NewHandlerFu[command.UpdateRuleStatusCommand](rs.UpdateRuleStatus)) // 更新规则状态

		g.GET("/enabled", hserver.NewHandlerFu[queries.GetEnabledRulesReq](rs.GetEnabledRules)) // 获取启用的规则列表

		g.GET("/by-category", hserver.NewHandlerFu[queries.GetRulesByCategoryReq](rs.GetRulesByCategory)) // 根据分类获取规则

		g.GET("/by-type", hserver.NewHandlerFu[queries.GetRulesByTypeReq](rs.GetRulesByType)) // 根据类型获取规则
	}
}

// CreateRule 创建规则
// @Summary 创建规则
// @Description 创建新规则
// @Tags 规则引擎
// @ID CreateRule
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.CreateRuleCommand true "规则信息"
// @Success 200 {object} base_info.Success{} "创建成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/rule [post]
func (rs *RuleService) CreateRule(ctx context.Context, req *command.CreateRuleCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := rs.rc.HandleCreateRule(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// GetRuleList 获取规则列表
// @Summary 获取规则列表
// @Description 分页获取规则列表
// @Tags 规则引擎
// @ID GetRuleList
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetRuleListReq true "查询参数"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.RuleDTO]{list=[]dto.RuleDTO}} "规则列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/rule [get]
func (rs *RuleService) GetRuleList(ctx context.Context, req *queries.GetRuleListReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, total, err := rs.rh.HandleGetRuleList(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(models.NewPageRes(total, data))
}

// GetRuleInfo 获取规则信息
// @Summary 获取规则信息
// @Description 获取指定规则的详细信息
// @Tags 规则引擎
// @ID GetRuleInfo
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "规则ID"
// @Success 200 {object} base_info.Success{data=dto.RuleDTO} "规则信息"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/rule/{id} [get]
func (rs *RuleService) GetRuleInfo(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := rs.rh.HandleGetRule(ctx, &queries.GetRuleReq{ID: req.Id})
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}

// DeleteRule 删除规则
// @Summary 删除规则
// @Description 删除指定规则
// @Tags 规则引擎
// @ID DeleteRule
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "规则ID"
// @Success 200 {object} base_info.Success{} "删除成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/rule/{id} [delete]
func (rs *RuleService) DeleteRule(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := rs.rc.HandleDeleteRule(ctx, &command.DeleteRuleCommand{ID: req.Id})
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// UpdateRule 更新规则信息
// @Summary 更新规则信息
// @Description 更新规则基本信息
// @Tags 规则引擎
// @ID UpdateRule
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.UpdateRuleCommand true "规则信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/rule [put]
func (rs *RuleService) UpdateRule(ctx context.Context, req *command.UpdateRuleCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := rs.rc.HandleUpdateRule(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// UpdateRuleStatus 更新规则状态
// @Summary 更新规则状态
// @Description 更新规则状态（启用/禁用）
// @Tags 规则引擎
// @ID UpdateRuleStatus
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.UpdateRuleStatusCommand true "状态信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/rule/status [put]
func (rs *RuleService) UpdateRuleStatus(ctx context.Context, req *command.UpdateRuleStatusCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := rs.rc.HandleUpdateRuleStatus(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// GetEnabledRules 获取启用的规则列表
// @Summary 获取启用的规则列表
// @Description 获取所有启用的规则列表
// @Tags 规则引擎
// @ID GetEnabledRules
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetEnabledRulesReq true "查询参数"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.RuleDTO]{list=[]dto.RuleDTO}} "规则列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/rule/enabled [get]
func (rs *RuleService) GetEnabledRules(ctx context.Context, req *queries.GetEnabledRulesReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, total, err := rs.rh.HandleGetEnabledRules(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(models.NewPageRes(total, data))
}

// GetRulesByCategory 根据分类获取规则
// @Summary 根据分类获取规则
// @Description 根据分类获取规则列表
// @Tags 规则引擎
// @ID GetRulesByCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetRulesByCategoryReq true "查询参数"
// @Success 200 {object} base_info.Success{data=[]dto.RuleDTO} "规则列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/rule/by-category [get]
func (rs *RuleService) GetRulesByCategory(ctx context.Context, req *queries.GetRulesByCategoryReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := rs.rh.HandleGetRulesByCategory(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}

// GetRulesByType 根据类型获取规则
// @Summary 根据类型获取规则
// @Description 根据类型获取规则列表
// @Tags 规则引擎
// @ID GetRulesByType
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetRulesByTypeReq true "查询参数"
// @Success 200 {object} base_info.Success{data=[]dto.RuleDTO} "规则列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/rule/by-type [get]
func (rs *RuleService) GetRulesByType(ctx context.Context, req *queries.GetRulesByTypeReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := rs.rh.HandleGetRulesByType(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}
