package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/repository"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/dto"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/queries"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
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
	// 构建查询条件
	query := db_query.NewQueryBuilder()
	if req.Name != "" {
		query.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	if req.Code != "" {
		query.Where("code", db_query.Eq, req.Code)
	}
	if req.ParentID != "" {
		query.Where("parent_id", db_query.Eq, req.ParentID)
	}
	if req.Type != "" {
		query.Where("type", db_query.Eq, req.Type)
	}
	if req.BusinessType != "" {
		query.Where("business_type", db_query.Eq, req.BusinessType)
	}
	if req.Status > 0 {
		query.Where("status", db_query.Eq, req.Status)
	}

	// 获取总数
	total, err := h.categoryRepo.Count(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Count categories error: %v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	// 分页查询
	query.OrderByASC("sorting")
	query.OrderByDESC("created_at")
	query.WithPage(&req.Page)
	categories, err := h.categoryRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Get category list error: %v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = h.convertToDTO(category)
	}

	return dtos, total, nil
}

// HandleGetCategoriesByParent 处理根据父分类获取子分类列表查询
func (h *CategoryQueryHandler) HandleGetCategoriesByParent(ctx context.Context, req *queries.GetCategoriesByParentReq) ([]*dto.CategoryDTO, *herrors.HError) {
	categories, err := h.categoryRepo.FindByParentID(ctx, req.ParentID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = h.convertToDTO(category)
	}

	return dtos, nil
}

// HandleGetCategoriesByBusinessType 处理根据业务类型获取分类列表查询
func (h *CategoryQueryHandler) HandleGetCategoriesByBusinessType(ctx context.Context, req *queries.GetCategoriesByBusinessTypeReq) ([]*dto.CategoryDTO, *herrors.HError) {
	categories, err := h.categoryRepo.FindByBusinessType(ctx, req.BusinessType)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = h.convertToDTO(category)
	}

	return dtos, nil
}

// HandleGetCategoriesByType 处理根据分类类型获取分类列表查询
func (h *CategoryQueryHandler) HandleGetCategoriesByType(ctx context.Context, req *queries.GetCategoriesByTypeReq) ([]*dto.CategoryDTO, *herrors.HError) {
	categories, err := h.categoryRepo.FindByType(ctx, req.Type)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = h.convertToDTO(category)
	}

	return dtos, nil
}

// HandleGetRootCategories 处理获取根分类列表查询
func (h *CategoryQueryHandler) HandleGetRootCategories(ctx context.Context, req *queries.GetRootCategoriesReq) ([]*dto.CategoryDTO, *herrors.HError) {
	categories, err := h.categoryRepo.FindRootCategories(ctx)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = h.convertToDTO(category)
	}

	return dtos, nil
}

// HandleGetCategoryTree 处理获取分类树查询
func (h *CategoryQueryHandler) HandleGetCategoryTree(ctx context.Context, req *queries.GetCategoryTreeReq) ([]*dto.CategoryTreeDTO, *herrors.HError) {
	var categories []*model.RuleCategory
	var err error

	if req.CategoryID == "" {
		// 获取所有根分类
		categories, err = h.categoryRepo.FindRootCategories(ctx)
	} else {
		// 获取指定分类及其所有后代
		category, err := h.categoryRepo.FindByID(ctx, req.CategoryID)
		if err != nil {
			return nil, herrors.QueryFail(err)
		}

		descendants, err := h.categoryRepo.FindDescendants(ctx, category.Path)
		if err != nil {
			return nil, herrors.QueryFail(err)
		}

		categories = append([]*model.RuleCategory{category}, descendants...)
	}

	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 构建树结构
	treeMap := make(map[string]*dto.CategoryTreeDTO)
	var rootNodes []*dto.CategoryTreeDTO

	for _, category := range categories {
		treeDTO := &dto.CategoryTreeDTO{
			CategoryDTO: *h.convertToDTO(category),
			Children:    make([]*dto.CategoryTreeDTO, 0),
		}
		treeMap[category.ID] = treeDTO

		if category.ParentID == "" {
			rootNodes = append(rootNodes, treeDTO)
		} else {
			if parent, exists := treeMap[category.ParentID]; exists {
				parent.Children = append(parent.Children, treeDTO)
			}
		}
	}

	return rootNodes, nil
}

// HandleGetCategory 处理获取分类详情查询
func (h *CategoryQueryHandler) HandleGetCategory(ctx context.Context, req *queries.GetCategoryReq) (*dto.CategoryDTO, *herrors.HError) {
	category, err := h.categoryRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	return h.convertToDTO(category), nil
}

// HandleGetCategoryByCode 处理根据编码获取分类查询
func (h *CategoryQueryHandler) HandleGetCategoryByCode(ctx context.Context, req *queries.GetCategoryByCodeReq) (*dto.CategoryDTO, *herrors.HError) {
	category, err := h.categoryRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	return h.convertToDTO(category), nil
}

// HandleGetAllCategories 处理获取所有分类查询
func (h *CategoryQueryHandler) HandleGetAllCategories(ctx context.Context) ([]*dto.CategoryDTO, *herrors.HError) {
	categories, err := h.categoryRepo.FindAll(ctx)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.CategoryDTO, len(categories))
	for i, category := range categories {
		dtos[i] = h.convertToDTO(category)
	}

	return dtos, nil
}

// convertToDTO 转换为DTO
func (h *CategoryQueryHandler) convertToDTO(category *model.RuleCategory) *dto.CategoryDTO {
	return &dto.CategoryDTO{
		ID:           category.ID,
		Code:         category.Code,
		Name:         category.Name,
		Description:  category.Description,
		ParentID:     category.ParentID,
		Type:         category.Type,
		BusinessType: category.BusinessType,
		Level:        category.Level,
		Path:         category.Path,
		IsLeaf:       category.IsLeaf,
		Sorting:      category.Sorting,
		Status:       int(category.Status),
		CreatedAt:    category.CreatedAt,
		UpdatedAt:    category.UpdatedAt,
		TenantID:     category.TenantID,
	}
}
