package handler

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/template/application/dto"
	"github.com/flare-admin/flare-server-go/framework/support/template/application/queries"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/repository"
)

// CategoryQueryHandler 分类查询处理器
type CategoryQueryHandler struct {
	categoryRepo repository.ICategoryRepository
}

// NewCategoryQueryHandler 创建分类查询处理器
func NewCategoryQueryHandler(categoryRepo repository.ICategoryRepository) *CategoryQueryHandler {
	return &CategoryQueryHandler{
		categoryRepo: categoryRepo,
	}
}

// HandleGetCategoryList 处理获取分类列表查询
func (h *CategoryQueryHandler) HandleGetCategoryList(ctx context.Context, req *queries.GetCategoryListReq) ([]*dto.CategoryDTO, int64, *herrors.HError) {
	query := db_query.NewQueryBuilder()
	if req.Name != "" {
		query.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	if req.Code != "" {
		query.Where("code", db_query.Eq, req.Code)
	}
	if req.Status != 0 {
		query.Where("status", db_query.Eq, req.Status)
	}
	total, err2 := h.categoryRepo.Count(ctx, query)
	if err2 != nil {
		hlog.CtxErrorf(ctx, "Count categories error: %v", err2)
		return nil, 0, herrors.QueryFail(err2)
	}
	query.WithPage(&req.Page)
	query.OrderByDESC("sort")
	entities, err := h.categoryRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取分类列表失败:%v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	return dto.FromEntities(entities), total, nil
}

// HandleGetCategoryDetail 处理获取分类详情查询
func (h *CategoryQueryHandler) HandleGetCategoryDetail(ctx context.Context, id string) (*dto.CategoryDTO, *herrors.HError) {
	entity, err := h.categoryRepo.FindById(ctx, id)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取分类详情失败:%v", err)
		return nil, herrors.QueryFail(err)
	}

	return dto.FromEntity(entity), nil
}

// HandleGetCategoryByCode 处理根据编码获取分类查询
func (h *CategoryQueryHandler) HandleGetCategoryByCode(ctx context.Context, req *queries.GetCategoryByCodeReq) (*dto.CategoryDTO, *herrors.HError) {
	entity, err := h.categoryRepo.FindByCode(ctx, req.Code)
	if err != nil {
		hlog.CtxErrorf(ctx, "根据编码获取分类失败:%v", err)
		return nil, herrors.QueryFail(err)
	}

	return dto.FromEntity(entity), nil
}

// HandleGetAllCategories 处理获取所有分类查询
func (h *CategoryQueryHandler) HandleGetAllCategories(ctx context.Context) ([]*dto.CategoryDTO, *herrors.HError) {
	entities, err := h.categoryRepo.FindAll(ctx)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取所有分类失败:%v", err)
		return nil, herrors.QueryFail(err)
	}

	return dto.FromEntities(entities), nil
}

// HandleGetAllEnableCategories 处理获取所有启用分类查询
func (h *CategoryQueryHandler) HandleGetAllEnableCategories(ctx context.Context) ([]*dto.CategoryDTO, *herrors.HError) {
	query := db_query.NewQueryBuilder()
	query.Where("status", db_query.Eq, 1)
	query.OrderByDESC("sort")
	entities, err := h.categoryRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "获取所有分类失败:%v", err)
		return nil, herrors.QueryFail(err)
	}

	return dto.FromEntities(entities), nil
}
