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

// CategoryService 规则引擎分类管理服务
type CategoryService struct {
	ch *queryhandler.CategoryQueryHandler
	cc *commandhandler.CategoryCommandHandler
	ef *casbin.Enforcer
}

// NewCategoryService 创建规则引擎分类管理服务
func NewCategoryService(
	ch *queryhandler.CategoryQueryHandler,
	cc *commandhandler.CategoryCommandHandler,
	ef *casbin.Enforcer,
) *CategoryService {
	return &CategoryService{
		ch: ch,
		cc: cc,
		ef: ef,
	}
}

// RegisterRouter 注册路由
func (cs *CategoryService) RegisterRouter(rg *route.RouterGroup, t token.IToken) {
	g := rg.Group("/v1/rule-engine/category", jwt.Handler(t))
	{
		// 分类管理
		g.POST("", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则引擎分类管理",
			Action:      "新建",
		}), hserver.NewHandlerFu[command.CreateCategoryCommand](cs.CreateCategory)) // 新建分类

		g.GET("", hserver.NewHandlerFu[queries.GetCategoryListReq](cs.GetCategoryList)) // 获取分类列表

		g.GET("/:id", hserver.NewHandlerFu[models.StringIdReq](cs.GetCategoryInfo)) // 获取分类信息

		g.DELETE("/:id", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则引擎分类管理",
			Action:      "删除",
		}), hserver.NewHandlerFu[models.StringIdReq](cs.DeleteCategory)) // 删除分类

		g.PUT("", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则引擎分类管理",
			Action:      "修改",
		}), hserver.NewHandlerFu[command.UpdateCategoryCommand](cs.UpdateCategory)) // 修改分类信息

		g.PUT("/status", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      "规则引擎分类管理",
			Action:      "更新状态",
		}), hserver.NewHandlerFu[command.UpdateCategoryStatusCommand](cs.UpdateCategoryStatus)) // 更新分类状态

		g.GET("/by-business-type", hserver.NewHandlerFu[queries.GetCategoriesByBusinessTypeReq](cs.GetCategoriesByBusinessType)) // 根据业务类型获取分类

		g.GET("/by-type", hserver.NewHandlerFu[queries.GetCategoriesByTypeReq](cs.GetCategoriesByType)) // 根据分类类型获取分类
	}
}

// CreateCategory 创建分类
// @Summary 创建规则引擎分类
// @Description 创建新规则引擎分类
// @Tags 规则引擎
// @ID CreateRuleCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.CreateCategoryCommand true "分类信息"
// @Success 200 {object} base_info.Success{} "创建成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/category [post]
func (cs *CategoryService) CreateCategory(ctx context.Context, req *command.CreateCategoryCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := cs.cc.HandleCreateCategory(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// GetCategoryList 获取分类列表
// @Summary 获取规则引擎分类列表
// @Description 分页获取规则引擎分类列表
// @Tags 规则引擎
// @ID GetRuleCategoryList
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetCategoryListReq true "查询参数"
// @Success 200 {object} base_info.Success{data=models.PageRes[dto.CategoryDTO]{list=[]dto.CategoryDTO}} "分类列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/category [get]
func (cs *CategoryService) GetCategoryList(ctx context.Context, req *queries.GetCategoryListReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, total, err := cs.ch.HandleGetCategoryList(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(models.NewPageRes(total, data))
}

// GetCategoryInfo 获取分类信息
// @Summary 获取规则引擎分类信息
// @Description 获取指定规则引擎分类的详细信息
// @Tags 规则引擎
// @ID GetRuleCategoryInfo
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "分类ID"
// @Success 200 {object} base_info.Success{data=dto.CategoryDTO} "分类信息"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/category/{id} [get]
func (cs *CategoryService) GetCategoryInfo(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := cs.ch.HandleGetCategory(ctx, &queries.GetCategoryReq{ID: req.Id})
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}

// DeleteCategory 删除分类
// @Summary 删除规则引擎分类
// @Description 删除指定规则引擎分类
// @Tags 规则引擎
// @ID DeleteRuleCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "分类ID"
// @Success 200 {object} base_info.Success{} "删除成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/category/{id} [delete]
func (cs *CategoryService) DeleteCategory(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := cs.cc.HandleDeleteCategory(ctx, &command.DeleteCategoryCommand{ID: req.Id})
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// UpdateCategory 更新分类信息
// @Summary 更新规则引擎分类信息
// @Description 更新规则引擎分类基本信息
// @Tags 规则引擎
// @ID UpdateRuleCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.UpdateCategoryCommand true "分类信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/category [put]
func (cs *CategoryService) UpdateCategory(ctx context.Context, req *command.UpdateCategoryCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := cs.cc.HandleUpdateCategory(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// UpdateCategoryStatus 更新分类状态
// @Summary 更新规则引擎分类状态
// @Description 更新规则引擎分类状态（启用/禁用）
// @Tags 规则引擎
// @ID UpdateRuleCategoryStatus
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body command.UpdateCategoryStatusCommand true "状态信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/category/status [put]
func (cs *CategoryService) UpdateCategoryStatus(ctx context.Context, req *command.UpdateCategoryStatusCommand) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	err := cs.cc.HandleUpdateCategoryStatus(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res
}

// GetCategoriesByBusinessType 根据业务类型获取分类
// @Summary 根据业务类型获取分类
// @Description 根据业务类型获取分类列表
// @Tags 规则引擎
// @ID GetRuleCategoriesByBusinessType
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetCategoriesByBusinessTypeReq true "查询参数"
// @Success 200 {object} base_info.Success{data=[]dto.CategoryDTO} "分类列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/category/by-business-type [get]
func (cs *CategoryService) GetCategoriesByBusinessType(ctx context.Context, req *queries.GetCategoriesByBusinessTypeReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := cs.ch.HandleGetCategoriesByBusinessType(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}

// GetCategoriesByType 根据分类类型获取分类
// @Summary 根据分类类型获取分类
// @Description 根据分类类型获取分类列表
// @Tags 规则引擎
// @ID GetRuleCategoriesByType
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req query queries.GetCategoriesByTypeReq true "查询参数"
// @Success 200 {object} base_info.Success{data=[]dto.CategoryDTO} "分类列表"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /v1/rule-engine/category/by-type [get]
func (cs *CategoryService) GetCategoriesByType(ctx context.Context, req *queries.GetCategoriesByTypeReq) *hserver.ResponseResult {
	res := hserver.DefaultResponseResult()
	data, err := cs.ch.HandleGetCategoriesByType(ctx, req)
	if herrors.HaveError(err) {
		return res.WithError(err)
	}
	return res.WithData(data)
}
