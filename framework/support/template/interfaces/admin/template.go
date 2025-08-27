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
	"github.com/flare-admin/flare-server-go/framework/support/template/application/command"
	commandhandler "github.com/flare-admin/flare-server-go/framework/support/template/application/command/handler"
	"github.com/flare-admin/flare-server-go/framework/support/template/application/dto"
	"github.com/flare-admin/flare-server-go/framework/support/template/application/queries"
	queryhandler "github.com/flare-admin/flare-server-go/framework/support/template/application/queries/handler"
)

// TemplateService 模板管理服务
type TemplateService struct {
	th *queryhandler.TemplateQueryHandler
	tc *commandhandler.TemplateCommandHandler
	ef *casbin.Enforcer
}

// NewTemplateService 创建模板管理服务
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
	g := rg.Group("/v1/template/center", jwt.Handler(t))
	{
		// 模板管理
		g.POST("", casbin.Handler(ts.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "模板管理",
			Action:      "新建",
		}), hserver.NewHandlerFu[command.CreateTemplateCommand](ts.CreateTemplate)) // 新建模板

		g.GET("", casbin.Handler(ts.ef), hserver.NewHandlerFu[queries.GetTemplateListReq](ts.GetTemplateList)) // 获取模板列表

		g.GET("/:id", casbin.Handler(ts.ef), hserver.NewHandlerFu[models.StringIdReq](ts.GetTemplateInfo)) // 获取模板信息

		g.DELETE("/:id", casbin.Handler(ts.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "模板管理",
			Action:      "删除",
		}), hserver.NewHandlerFu[models.StringIdReq](ts.DeleteTemplate)) // 删除模板

		g.PUT("", casbin.Handler(ts.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "模板管理",
			Action:      "修改",
		}), hserver.NewHandlerFu[command.UpdateTemplateCommand](ts.UpdateTemplate)) // 修改模板信息

		g.PUT("/status", casbin.Handler(ts.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "模板管理",
			Action:      "更新状态",
		}), hserver.NewHandlerFu[command.UpdateTemplateStatusCommand](ts.UpdateTemplateStatus)) // 更新模板状态

		g.GET("/enabled", casbin.Handler(ts.ef), hserver.NewHandlerFu[queries.GetEnabledTemplateReq](ts.GetEnabledTemplates)) // 获取启用的模板列表

		g.GET("/enabled/all", casbin.Handler(ts.ef), hserver.NewHandlerFu[queries.GetTemplatesByCategoryReq](ts.GetAllEnabledTemplatesByCategory)) // 获取分类下所有启用模版
	}
}

// CreateTemplate 创建模板
// @Summary 创建模板
// @Description 创建新模板
// @Tags 模板管理
// @ID CreateCenterTemplate
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.CreateTemplateCommand true "模板信息"
// @Success 200 {object} base_info.Success{} "创建成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center [post]
func (ts *TemplateService) CreateTemplate(ctx context.Context, req *command.CreateTemplateCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := ts.tc.HandleCreateTemplate(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// GetTemplateList 获取模板列表
// @Summary 获取模板列表
// @Description 分页获取模板列表
// @Tags 模板管理
// @ID GetCenterTemplateList
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetTemplateListReq true "查询参数"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.TemplateDTO]} "模板列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center [get]
func (ts *TemplateService) GetTemplateList(ctx context.Context, req *queries.GetTemplateListReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, total, err := ts.th.HandleGetTemplateList(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(models.NewPageRes(total, data))
}

// GetTemplateInfo 获取模板信息
// @Summary 获取模板信息
// @Description 获取指定模板的详细信息
// @Tags 模板管理
// @ID GetCenterTemplateInfo
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "模板ID"
// @Success 200 {object} base_info.Success{data=dto.TemplateDTO} "模板信息"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center/{id} [get]
func (ts *TemplateService) GetTemplateInfo(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := ts.th.HandleGetTemplateDetail(ctx, req.Id)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	_ = dto.TemplateDTO{}
	return res.WithData(data)
}

// DeleteTemplate 删除模板
// @Summary 删除模板
// @Description 删除指定模板
// @Tags 模板管理
// @ID DeleteCenterTemplate
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "模板ID"
// @Success 200 {object} base_info.Success{} "删除成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center/{id} [delete]
func (ts *TemplateService) DeleteTemplate(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := ts.tc.HandleDeleteTemplate(ctx, &command.DeleteTemplateCommand{ID: req.Id})
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// UpdateTemplate 更新模板信息
// @Summary 更新模板信息
// @Description 更新模板基本信息
// @Tags 模板管理
// @ID UpdateCenterTemplate
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.UpdateTemplateCommand true "模板信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center [put]
func (ts *TemplateService) UpdateTemplate(ctx context.Context, req *command.UpdateTemplateCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := ts.tc.HandleUpdateTemplate(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// UpdateTemplateStatus 更新模板状态
// @Summary 更新模板状态
// @Description 更新模板状态（启用/禁用）
// @Tags 模板管理
// @ID UpdateCenterTemplateStatus
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.UpdateTemplateStatusCommand true "状态信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center/status [put]
func (ts *TemplateService) UpdateTemplateStatus(ctx context.Context, req *command.UpdateTemplateStatusCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := ts.tc.HandleUpdateTemplateStatus(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// GetEnabledTemplates 获取启用的模板列表
// @Summary 获取启用的模板列表
// @Description 获取所有启用的模板列表
// @Tags 模板管理
// @ID GetEnabledCenterTemplates
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetEnabledTemplateReq true "查询参数"
// @Success 200 {object} base_info.Success{data=[]dto.TemplateDTO} "模板列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center/enabled [get]
func (ts *TemplateService) GetEnabledTemplates(ctx context.Context, req *queries.GetEnabledTemplateReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, total, err := ts.th.HandleEnabledGetTemplatesByCategory(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(models.NewPageRes(total, data))
}

// GetAllEnabledTemplatesByCategory 	获取分类下所有启用模版
// @Summary 获取分类下所有启用模版
// @Description 获取分类下所有启用模版
// @Tags 模板管理
// @ID GetAllEnabledCenterTemplatesByCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetTemplatesByCategoryReq true "查询参数"
// @Success 200 {object} base_info.Success{data=[]dto.TemplateDTO} "模板列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center/enabled/all [get]
func (ts *TemplateService) GetAllEnabledTemplatesByCategory(ctx context.Context, req *queries.GetTemplatesByCategoryReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := ts.th.HandleGetAllEnabledTemplatesByCategory(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}
