package dictionaryinterfaces

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/dto"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/service"

	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	_ "github.com/flare-admin/flare-server-go/framework/pkg/hserver/base_info"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/casbin"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/jwt"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/oplog"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
)

type DictionaryService struct {
	ds      service.IDictionaryService
	ef      *casbin.Enforcer
	modeNma string
}

func NewDictionaryService(ds service.IDictionaryService, ef *casbin.Enforcer) *DictionaryService {
	return &DictionaryService{
		ds:      ds,
		ef:      ef,
		modeNma: "数据字典",
	}
}

func (s *DictionaryService) RegisterRouter(rg *route.RouterGroup, t token.IToken) {
	g := rg.Group("/v1/dictionary", jwt.Handler(t))
	{
		// 分类管理
		g.POST("/category", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.modeNma,
			Action:      "新增分类",
		}), hserver.NewHandlerFu[dto.CategoryCreateReq](s.CreateCategory))

		g.PUT("/category", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.modeNma,
			Action:      "更新分类",
		}), hserver.NewHandlerFu[dto.CategoryUpdateReq](s.UpdateCategory))

		g.GET("/category/:id", casbin.Handler(s.ef), hserver.NewHandlerFu[models.StringIdReq](s.GetCategory))
		g.GET("/categories", casbin.Handler(s.ef), hserver.NewHandlerFu[dto.CategoryQueryReq](s.ListCategories))
		g.DELETE("/category/:id", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.modeNma,
			Action:      "删除分类",
		}), hserver.NewHandlerFu[models.StringIdReq](s.DelCategory))
		// 选项管理
		g.POST("/option", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.modeNma,
			Action:      "新增选项",
		}), hserver.NewHandlerFu[dto.OptionCreateReq](s.CreateOption))

		g.PUT("/option/:id", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.modeNma,
			Action:      "更新选项",
		}), hserver.NewHandlerFu[dto.OptionUpdateReq](s.UpdateOption))

		g.DELETE("/option/:id", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.modeNma,
			Action:      "删除选项",
		}), hserver.NewHandlerFu[models.StringIdReq](s.DeleteOption))

		g.GET("/options", hserver.NewHandlerFu[dto.OptionQueryReq](s.GetOptions))
	}
}

// CreateCategory 创建分类
// @Summary 创建字典分类
// @Description 创建新的字典分类
// @Tags 数据字典
// @ID CreateDictionaryCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body dto.CategoryCreateReq true "分类信息"
// @Success 200 {object} base_info.Success{} "创建成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /api/admin/v1/dictionary/category [post]
func (s *DictionaryService) CreateCategory(ctx context.Context, req *dto.CategoryCreateReq) *hserver.ResponseResult {
	err := s.ds.CreateCategory(ctx, req)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}

// UpdateCategory 更新分类
// @Summary 更新字典分类
// @Description 更新已存在的字典分类
// @Tags 数据字典
// @ID UpdateDictionaryCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "分类ID"
// @Param req body dto.CategoryUpdateReq true "更新信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /api/admin/v1/dictionary/category/{id} [put]
func (s *DictionaryService) UpdateCategory(ctx context.Context, req *dto.CategoryUpdateReq) *hserver.ResponseResult {
	err := s.ds.UpdateCategory(ctx, req.ID, req)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}

// GetCategory 获取分类详情
// @Summary 获取字典分类详情
// @Description 根据ID获取字典分类详情
// @Tags 数据字典
// @ID GetDictionaryCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "分类ID"
// @Success 200 {object} base_info.Success{data=dto.Category} "获取成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /api/admin/v1/dictionary/category/{id} [get]
func (s *DictionaryService) GetCategory(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	category, err := s.ds.GetCategory(ctx, req.Id)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res.WithData(category)
}

// DelCategory 删除分类
// @Summary 删除分类
// @Description 删除分类
// @Tags 数据字典
// @ID DelDictionaryCategory
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "分类ID"
// @Success 200 {object} base_info.Success{data=dto.Category} "获取成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /api/admin/v1/dictionary/category/{id} [delete]
func (s *DictionaryService) DelCategory(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	err := s.ds.DelCategory(ctx, req.Id)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}

// ListCategories 获取分类列表
// @Summary 获取字典分类列表
// @Description 分页获取字典分类列表
// @Tags 数据字典
// @ID ListDictionaryCategories
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param id query string false "分类ID"
// @Param name query string false "分类名称"
// @Param i18n_key query string false "国际化key"
// @Success 200 {object} base_info.Success{data=[]dto.Category} "获取成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /api/admin/v1/dictionary/categories [get]
func (s *DictionaryService) ListCategories(ctx context.Context, req *dto.CategoryQueryReq) *hserver.ResponseResult {
	categories, total, err := s.ds.ListCategories(ctx, req)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res.WithData(&models.PageRes[dto.Category]{
		Total: total,
		List:  categories,
	})
}

// CreateOption 创建选项
// @Summary 创建字典选项
// @Description 创建新的字典选项
// @Tags 数据字典
// @ID CreateDictionaryOption
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param req body dto.OptionCreateReq true "选项信息"
// @Success 200 {object} base_info.Success{} "创建成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /api/admin/v1/dictionary/option [post]
func (s *DictionaryService) CreateOption(ctx context.Context, req *dto.OptionCreateReq) *hserver.ResponseResult {
	err := s.ds.CreateOption(ctx, req)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}

// UpdateOption 更新选项
// @Summary 更新字典选项
// @Description 更新已存在的字典选项
// @Tags 数据字典
// @ID UpdateDictionaryOption
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "选项ID"
// @Param req body dto.OptionUpdateReq true "更新信息"
// @Success 200 {object} base_info.Success{} "更新成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /api/admin/v1/dictionary/option/{id} [put]
func (s *DictionaryService) UpdateOption(ctx context.Context, req *dto.OptionUpdateReq) *hserver.ResponseResult {
	err := s.ds.UpdateOption(ctx, req.ID, req)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}

// DeleteOption 删除选项
// @Summary 删除字典选项
// @Description 删除字典选项
// @Tags 数据字典
// @ID DeleteDictionaryOption
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "选项ID"
// @Success 200 {object} base_info.Success{} "删除成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /api/admin/v1/dictionary/option/{id} [delete]
func (s *DictionaryService) DeleteOption(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	err := s.ds.DeleteOption(ctx, req.Id)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}

// GetOptions 获取选项列表
// @Summary 获取字典选项列表
// @Description 获取字典选项列表
// @Tags 数据字典
// @ID GetDictionaryOptions
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer token"
// @Param category_id query string false "分类ID"
// @Param keyword query string false "关键字"
// @Param status query int false "状态"
// @Success 200 {object} base_info.Success{data=[]dto.Option} "获取成功"
// @Failure 400 {object} base_info.Swagger400Resp "参数错误"
// @Failure 401 {object} base_info.Swagger401Resp "未授权"
// @Failure 500 {object} base_info.Swagger500Resp "服务器内部错误"
// @Router /api/admin/v1/dictionary/options [get]
func (s *DictionaryService) GetOptions(ctx context.Context, req *dto.OptionQueryReq) *hserver.ResponseResult {
	options, err := s.ds.GetOptions(ctx, req)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res.WithData(options)
}
