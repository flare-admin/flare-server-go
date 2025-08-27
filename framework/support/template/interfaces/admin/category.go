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
	_ "github.com/flare-admin/flare-server-go/framework/support/template/application/dto"
	"github.com/flare-admin/flare-server-go/framework/support/template/application/queries"
	queryhandler "github.com/flare-admin/flare-server-go/framework/support/template/application/queries/handler"
)

// CategoryService 模板管理服务
type CategoryService struct {
	th *queryhandler.CategoryQueryHandler
	tc *commandhandler.CategoryCommandHandler
	ef *casbin.Enforcer
}

// NewCategoryService 创建模板管理服务
func NewCategoryService(
	th *queryhandler.CategoryQueryHandler,
	tc *commandhandler.CategoryCommandHandler,
	ef *casbin.Enforcer,
) *CategoryService {
	return &CategoryService{
		th: th,
		tc: tc,
		ef: ef,
	}
}

// RegisterRouter 注册路由
func (cs *CategoryService) RegisterRouter(rg *route.RouterGroup, t token.IToken) {
	g := rg.Group("/v1/template/center/category", jwt.Handler(t))
	{
		// 模板管理
		g.POST("", casbin.Handler(cs.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "模板管理",
			Action:      "新建",
		}), hserver.NewHandlerFu[command.CreateCategoryCommand](cs.CreateCategory)) // 新建分类

		g.GET("", casbin.Handler(cs.ef), hserver.NewHandlerFu[queries.GetCategoryListReq](cs.GetCategoryList)) // 获取分类列表

		g.GET("/:id", casbin.Handler(cs.ef), hserver.NewHandlerFu[models.StringIdReq](cs.GetCategoryInfo)) // 获取分类信息

		g.DELETE("/:id", casbin.Handler(cs.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "模板管理",
			Action:      "删除",
		}), hserver.NewHandlerFu[models.StringIdReq](cs.DeleteCategory)) // 删除分类

		g.PUT("", casbin.Handler(cs.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "模板管理",
			Action:      "修改",
		}), hserver.NewHandlerFu[command.UpdateCategoryCommand](cs.UpdateCategory)) // 修改分类信息

		g.PUT("/status", casbin.Handler(cs.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "模板管理",
			Action:      "更新状态",
		}), hserver.NewHandlerFu[command.UpdateCategoryStatusCommand](cs.UpdateCategoryStatus)) // 更新分类状态

		g.GET("/enable/all", casbin.Handler(cs.ef), hserver.NewNotParHandlerFu(cs.GetAllCategories)) // 获取所有启用分类
	}
}

// CreateCategory 创建分类
// @Summary 创建分类
// @Description 创建新分类
// @Tags 模板管理
// @ID CreateCenterTemplateCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.CreateCategoryCommand true "分类信息"
// @Success 200 {object} base_info.Success{} "创建成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center/category [post]
func (cs *CategoryService) CreateCategory(ctx context.Context, req *command.CreateCategoryCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := cs.tc.HandleCreateCategory(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// GetCategoryList 获取分类列表
// @Summary 获取分类列表
// @Description 分页获取分类列表
// @Tags 模板管理
// @ID GetCenterTemplateCategoryList
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetCategoryListReq true "查询参数"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.CategoryDTO]} "分类列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center/category [get]
func (cs *CategoryService) GetCategoryList(ctx context.Context, req *queries.GetCategoryListReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, total, err := cs.th.HandleGetCategoryList(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(models.NewPageRes(total, data))
}

// GetCategoryInfo 获取分类信息
// @Summary 获取分类信息
// @Description 获取指定分类的详细信息
// @Tags 模板管理
// @ID GetCenterTemplateCategoryInfo
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "分类ID"
// @Success 200 {object} base_info.Success{data=dto.CategoryDTO} "分类信息"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center/category/{id} [get]
func (cs *CategoryService) GetCategoryInfo(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := cs.th.HandleGetCategoryDetail(ctx, req.Id)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}

// DeleteCategory 删除分类
// @Summary 删除分类
// @Description 删除指定分类
// @Tags 模板管理
// @ID DeleteCenterTemplateCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "分类ID"
// @Success 200 {object} base_info.Success{} "删除成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center/category/{id} [delete]
func (cs *CategoryService) DeleteCategory(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := cs.tc.HandleDeleteCategory(ctx, &command.DeleteCategoryCommand{ID: req.Id})
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// UpdateCategory 更新分类信息
// @Summary 更新分类信息
// @Description 更新分类基本信息
// @Tags 模板管理
// @ID UpdateCenterTemplateCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.UpdateCategoryCommand true "分类信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center/category [put]
func (cs *CategoryService) UpdateCategory(ctx context.Context, req *command.UpdateCategoryCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := cs.tc.HandleUpdateCategory(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// UpdateCategoryStatus 更新分类状态
// @Summary 更新分类状态
// @Description 更新分类状态（启用/禁用）
// @Tags 模板管理
// @ID UpdateCenterTemplateCategoryStatus
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.UpdateCategoryStatusCommand true "状态信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router/v1/template/center/category/status [put]
func (cs *CategoryService) UpdateCategoryStatus(ctx context.Context, req *command.UpdateCategoryStatusCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := cs.tc.HandleUpdateCategoryStatus(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// GetAllCategories 获取所有启用分类
// @Summary 获取所有启用分类
// @Description 获取所有启用分类列表
// @Tags 模板管理
// @ID GetAllCenterTemplateCategories
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} base_info.Success{data=[]dto.CategoryDTO} "分类列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/template/center/category/enable/all [get]
func (cs *CategoryService) GetAllCategories(ctx context.Context) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := cs.th.HandleGetAllEnableCategories(ctx)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}
