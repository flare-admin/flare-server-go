package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/template/application/dto"
	"github.com/flare-admin/flare-server-go/framework/support/template/application/queries"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/repository"
)

// TemplateQueryHandler 模板查询处理器
type TemplateQueryHandler struct {
	templateRepo repository.ITemplateRepository
	categoryRepo repository.ICategoryRepository
}

// NewTemplateQueryHandler 创建模板查询处理器
func NewTemplateQueryHandler(templateRepo repository.ITemplateRepository, categoryRepo repository.ICategoryRepository) *TemplateQueryHandler {
	return &TemplateQueryHandler{
		templateRepo: templateRepo,
		categoryRepo: categoryRepo,
	}
}

// HandleGetTemplateList 处理获取模板列表查询
func (h *TemplateQueryHandler) HandleGetTemplateList(ctx context.Context, req *queries.GetTemplateListReq) ([]*dto.TemplateDTO, int64, *herrors.HError) {
	query := db_query.NewQueryBuilder()
	if req.Name != "" {
		query.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	if req.Status != 0 {
		query.Where("status", db_query.Eq, req.Status)
	}
	if req.CategoryID != "" {
		query.Where("category_id", db_query.Eq, req.CategoryID)
	}
	if req.Code != "" {
		query.Where("code", db_query.Eq, req.Code)
	}
	query.OrderByASC("status")
	total, err2 := h.templateRepo.Count(ctx, query)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "Count templates error: %v", err2)
		return nil, 0, herrors.QueryFail(err2)
	}
	query.WithPage(&req.Page)
	entities, err := h.templateRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取模版列表失败:%v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	return dto.TemplateFromEntities(entities), total, nil
}

// HandleGetTemplateDetail 处理获取模板详情查询
func (h *TemplateQueryHandler) HandleGetTemplateDetail(ctx context.Context, id string) (*dto.TemplateDTO, *herrors.HError) {
	entity, err := h.templateRepo.FindById(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取模板详情失败:%v", err)
		return nil, herrors.QueryFail(err)
	}

	return dto.TemplateFromEntity(entity), nil
}

// HandleGetTemplatesByCategory 处理获取分类下的模板列表查询
func (h *TemplateQueryHandler) HandleGetTemplatesByCategory(ctx context.Context, req *queries.GetTemplatesByCategoryReq) ([]*dto.TemplateDTO, *herrors.HError) {
	entities, err := h.templateRepo.FindByCategoryID(ctx, req.CategoryID)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取分类下的模板列表失败:%v", err)
		return nil, herrors.QueryFail(err)
	}

	return dto.TemplateFromEntities(entities), nil
}

// HandleGetAllEnabledTemplatesByCategory 处理获取分类下所有启用模版
func (h *TemplateQueryHandler) HandleGetAllEnabledTemplatesByCategory(ctx context.Context, req *queries.GetTemplatesByCategoryReq) ([]*dto.TemplateDTO, *herrors.HError) {
	query := db_query.NewQueryBuilder()
	if req.CategoryID != "" {
		query.Where("category_id", db_query.Eq, req.CategoryID)
	}
	query.Where("status", db_query.Eq, 1)
	entities, err := h.templateRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取模版列表失败:%v", err)
		return nil, herrors.QueryFail(err)
	}

	return dto.TemplateFromEntities(entities), nil
}

// HandleEnabledGetTemplatesByCategory 处理启用的模板列表查询
func (h *TemplateQueryHandler) HandleEnabledGetTemplatesByCategory(ctx context.Context, req *queries.GetEnabledTemplateReq) ([]*dto.TemplateDTO, int64, *herrors.HError) {
	var categoryID string
	if req.CategoryCode != "" {
		category, err := h.categoryRepo.FindByCode(ctx, req.CategoryCode)
		if err != nil {
			if database.IfErrorNotFound(err) {
				return nil, 0, nil
			}
			return nil, 0, herrors.QueryFail(err)
		}
		categoryID = category.ID
	}

	query := db_query.NewQueryBuilder()
	if req.Name != "" {
		query.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	if categoryID != "" {
		query.Where("category_id", db_query.Eq, categoryID)
	}
	if req.Code != "" {
		query.Where("code", db_query.Eq, req.Code)
	}
	query.Where("status", db_query.Eq, 1)

	total, err2 := h.templateRepo.Count(ctx, query)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "Count templates error: %v", err2)
		return nil, 0, herrors.QueryFail(err2)
	}
	query.OrderByASC("status")
	query.WithPage(&req.Page)
	entities, err := h.templateRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取模版列表失败:%v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	return dto.TemplateFromEntities(entities), total, nil
}
