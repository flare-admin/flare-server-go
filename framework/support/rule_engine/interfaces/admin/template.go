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

// TemplateService 规则模板管理服务
type TemplateService struct {
	th *queryhandler.TemplateQueryHandler
	tc *commandhandler.TemplateCommandHandler
	ef *casbin.Enforcer
}

// NewTemplateService 创建规则模板管理服务
func NewTemplateService(
	th *queryhandler.TemplateQueryHandler,
	tc *commandhandler.TemplateCommandHandler,
	ef *casbin.Enforcer,
) *TemplateService {
	return &TemplateService{
		th: th,
		tc: tc,
		ef: ef,
	}
}

// RegisterRouter 注册路由
func (ts *TemplateService) RegisterRouter(rg *route.RouterGroup, t token.IToken) {
	g := rg.Group("/v1/rule-engine/template", jwt.Handler(t))
	{
		// 模板管理
		g.POST("", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则模板管理",
			Action:      "新建",
		}), hserver.NewHandlerFu[command.CreateTemplateCommand](ts.CreateTemplate)) // 新建模板

		g.GET("", hserver.NewHandlerFu[queries.GetTemplateListReq](ts.GetTemplateList)) // 获取模板列表

		g.GET("/:id", hserver.NewHandlerFu[models.StringIdReq](ts.GetTemplateInfo)) // 获取模板信息

		g.DELETE("/:id", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则模板管理",
			Action:      "删除",
		}), hserver.NewHandlerFu[models.StringIdReq](ts.DeleteTemplate)) // 删除模板

		g.PUT("", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则模板管理",
			Action:      "修改",
		}), hserver.NewHandlerFu[command.UpdateTemplateCommand](ts.UpdateTemplate)) // 修改模板信息

		g.PUT("/status", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则模板管理",
			Action:      "更新状态",
		}), hserver.NewHandlerFu[command.UpdateTemplateStatusCommand](ts.UpdateTemplateStatus)) // 更新模板状态

		g.GET("/enabled", hserver.NewHandlerFu[queries.GetEnabledTemplatesReq](ts.GetEnabledTemplates)) // 获取启用的模板列表

		g.GET("/enabled/all", hserver.NewHandlerFu[queries.GetTemplatesByCategoryReq](ts.GetAllEnabledTemplatesByCategory)) // 获取分类下所有启用模板

		g.GET("/by-category", hserver.NewHandlerFu[queries.GetTemplatesByCategoryReq](ts.GetTemplatesByCategory)) // 根据分类获取模板

		g.GET("/by-type", hserver.NewHandlerFu[queries.GetTemplatesByTypeReq](ts.GetTemplatesByType)) // 根据类型获取模板
	}
}

// CreateTemplate 创建模板
// @Summary 创建规则模板
// @Description 创建新规则模板
// @Tags 规则引擎
// @ID CreateRuleTemplate
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.CreateTemplateCommand true "模板信息"
// @Success 200 {object} base_info.Success{} "创建成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/template [post]
func (ts *TemplateService) CreateTemplate(ctx context.Context, req *command.CreateTemplateCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := ts.tc.HandleCreateTemplate(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// GetTemplateList 获取模板列表
// @Summary 获取规则模板列表
// @Description 分页获取规则模板列表
// @Tags 规则引擎
// @ID GetRuleTemplateList
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetTemplateListReq true "查询参数"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.TemplateDTO]{list=[]dto.TemplateDTO}} "模板列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/template [get]
func (ts *TemplateService) GetTemplateList(ctx context.Context, req *queries.GetTemplateListReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, total, err := ts.th.HandleGetTemplateList(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(models.NewPageRes(total, data))
}

// GetTemplateInfo 获取模板信息
// @Summary 获取规则模板信息
// @Description 获取指定规则模板的详细信息
// @Tags 规则引擎
// @ID GetRuleTemplateInfo
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "模板ID"
// @Success 200 {object} base_info.Success{data=dto.TemplateDTO} "模板信息"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/template/{id} [get]
func (ts *TemplateService) GetTemplateInfo(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := ts.th.HandleGetTemplate(ctx, &queries.GetTemplateReq{ID: req.Id})
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}

// DeleteTemplate 删除模板
// @Summary 删除规则模板
// @Description 删除指定规则模板
// @Tags 规则引擎
// @ID DeleteRuleTemplate1
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "模板ID"
// @Success 200 {object} base_info.Success{} "删除成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/template/{id} [delete]
func (ts *TemplateService) DeleteTemplate(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := ts.tc.HandleDeleteTemplate(ctx, &command.DeleteTemplateCommand{ID: req.Id})
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// UpdateTemplate 更新模板信息
// @Summary 更新规则模板信息
// @Description 更新规则模板基本信息
// @Tags 规则引擎
// @ID UpdateRuleTemplate
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.UpdateTemplateCommand true "模板信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/template [put]
func (ts *TemplateService) UpdateTemplate(ctx context.Context, req *command.UpdateTemplateCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := ts.tc.HandleUpdateTemplate(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// UpdateTemplateStatus 更新模板状态
// @Summary 更新规则模板状态
// @Description 更新规则模板状态（启用/禁用）
// @Tags 规则引擎
// @ID UpdateRuleTemplateStatus
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.UpdateTemplateStatusCommand true "状态信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/template/status [put]
func (ts *TemplateService) UpdateTemplateStatus(ctx context.Context, req *command.UpdateTemplateStatusCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := ts.tc.HandleUpdateTemplateStatus(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// GetEnabledTemplates 获取启用的模板列表
// @Summary 获取启用的规则模板列表
// @Description 获取所有启用的规则模板列表
// @Tags 规则引擎
// @ID GetEnabledRuleTemplates
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetEnabledTemplatesReq true "查询参数"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.TemplateDTO]{list=[]dto.TemplateDTO}} "模板列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/template/enabled [get]
func (ts *TemplateService) GetEnabledTemplates(ctx context.Context, req *queries.GetEnabledTemplatesReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, total, err := ts.th.HandleGetEnabledTemplates(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(models.NewPageRes(total, data))
}

// GetAllEnabledTemplatesByCategory 获取分类下所有启用模板
// @Summary 获取分类下所有启用模板
// @Description 获取分类下所有启用模板
// @Tags 规则引擎
// @ID GetAllEnabledRuleTemplatesByCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetTemplatesByCategoryReq true "查询参数"
// @Success 200 {object} base_info.Success{data=[]dto.TemplateDTO} "模板列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/template/enabled/all [get]
func (ts *TemplateService) GetAllEnabledTemplatesByCategory(ctx context.Context, req *queries.GetTemplatesByCategoryReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := ts.th.HandleGetTemplatesByCategory(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}

// GetTemplatesByCategory 根据分类获取模板
// @Summary 根据分类获取模板
// @Description 根据分类获取模板列表
// @Tags 规则引擎
// @ID GetRuleTemplatesByCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetTemplatesByCategoryReq true "查询参数"
// @Success 200 {object} base_info.Success{data=[]dto.TemplateDTO} "模板列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/template/by-category [get]
func (ts *TemplateService) GetTemplatesByCategory(ctx context.Context, req *queries.GetTemplatesByCategoryReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := ts.th.HandleGetTemplatesByCategory(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}

// GetTemplatesByType 根据类型获取模板
// @Summary 根据类型获取模板
// @Description 根据类型获取模板列表
// @Tags 规则引擎
// @ID GetRuleTemplatesByType
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetTemplatesByTypeReq true "查询参数"
// @Success 200 {object} base_info.Success{data=[]dto.TemplateDTO} "模板列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/template/by-type [get]
func (ts *TemplateService) GetTemplatesByType(ctx context.Context, req *queries.GetTemplatesByTypeReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := ts.th.HandleGetTemplatesByType(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}
