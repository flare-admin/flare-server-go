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

// TemplateQueryHandler 模板查询处理器
type TemplateQueryHandler struct {
	templateRepo repository.ITemplateRepository
}

// NewTemplateQueryHandler 创建模板查询处理器
func NewTemplateQueryHandler(templateRepo repository.ITemplateRepository) *TemplateQueryHandler {
	return &TemplateQueryHandler{
		templateRepo: templateRepo,
	}
}

// HandleGetTemplateList 处理获取模板列表查询
func (h *TemplateQueryHandler) HandleGetTemplateList(ctx context.Context, req *queries.GetTemplateListReq) ([]*dto.TemplateDTO, int64, *herrors.HError) {
	// 构建查询条件
	query := db_query.NewQueryBuilder()
	if req.Name != "" {
		query.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	if req.CategoryID != "" {
		query.Where("category_id", db_query.Eq, req.CategoryID)
	}
	if req.Code != "" {
		query.Where("code", db_query.Eq, req.Code)
	}
	if req.Type != "" {
		query.Where("type", db_query.Eq, req.Type)
	}
	if req.Status > 0 {
		query.Where("status", db_query.Eq, req.Status)
	}

	// 获取总数
	total, err := h.templateRepo.Count(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Count templates error: %v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	// 分页查询
	query.OrderByDESC("created_at")
	query.WithPage(&req.Page)
	templates, err := h.templateRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Get template list error: %v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.TemplateDTO, len(templates))
	for i, template := range templates {
		dtos[i] = h.convertToDTO(template)
	}

	return dtos, total, nil
}

// HandleGetTemplatesByCategory 处理根据分类获取模板列表查询
func (h *TemplateQueryHandler) HandleGetTemplatesByCategory(ctx context.Context, req *queries.GetTemplatesByCategoryReq) ([]*dto.TemplateDTO, *herrors.HError) {
	templates, err := h.templateRepo.FindByCategoryID(ctx, req.CategoryID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.TemplateDTO, len(templates))
	for i, template := range templates {
		dtos[i] = h.convertToDTO(template)
	}

	return dtos, nil
}

// HandleGetTemplatesByType 处理根据类型获取模板列表查询
func (h *TemplateQueryHandler) HandleGetTemplatesByType(ctx context.Context, req *queries.GetTemplatesByTypeReq) ([]*dto.TemplateDTO, *herrors.HError) {
	templates, err := h.templateRepo.FindByType(ctx, req.Type)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.TemplateDTO, len(templates))
	for i, template := range templates {
		dtos[i] = h.convertToDTO(template)
	}

	return dtos, nil
}

// HandleGetEnabledTemplates 处理获取启用的模板列表查询
func (h *TemplateQueryHandler) HandleGetEnabledTemplates(ctx context.Context, req *queries.GetEnabledTemplatesReq) ([]*dto.TemplateDTO, int64, *herrors.HError) {
	// 构建查询条件
	query := db_query.NewQueryBuilder()
	if req.Name != "" {
		query.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	if req.CategoryID != "" {
		query.Where("category_id", db_query.Eq, req.CategoryID)
	}
	if req.Code != "" {
		query.Where("code", db_query.Eq, req.Code)
	}
	if req.Type != "" {
		query.Where("type", db_query.Eq, req.Type)
	}
	query.Where("status", db_query.Eq, 1) // 只查询启用的模板

	// 获取总数
	total, err := h.templateRepo.Count(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Count enabled templates error: %v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	// 分页查询
	query.OrderByASC("sorting")
	query.OrderByDESC("created_at")
	query.WithPage(&req.Page)
	templates, err := h.templateRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Get enabled templates error: %v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.TemplateDTO, len(templates))
	for i, template := range templates {
		dtos[i] = h.convertToDTO(template)
	}

	return dtos, total, nil
}

// HandleGetTemplate 处理获取模板详情查询
func (h *TemplateQueryHandler) HandleGetTemplate(ctx context.Context, req *queries.GetTemplateReq) (*dto.TemplateDTO, *herrors.HError) {
	template, err := h.templateRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	return h.convertToDTO(template), nil
}

// HandleGetTemplateByCode 处理根据编码获取模板查询
func (h *TemplateQueryHandler) HandleGetTemplateByCode(ctx context.Context, req *queries.GetTemplateByCodeReq) (*dto.TemplateDTO, *herrors.HError) {
	template, err := h.templateRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	return h.convertToDTO(template), nil
}

// HandleGetAllTemplates 处理获取所有模板查询
func (h *TemplateQueryHandler) HandleGetAllTemplates(ctx context.Context) ([]*dto.TemplateDTO, *herrors.HError) {
	templates, err := h.templateRepo.FindAll(ctx)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.TemplateDTO, len(templates))
	for i, template := range templates {
		dtos[i] = h.convertToDTO(template)
	}

	return dtos, nil
}

// convertToDTO 转换为DTO
func (h *TemplateQueryHandler) convertToDTO(template *model.RuleTemplate) *dto.TemplateDTO {
	return &dto.TemplateDTO{
		ID:          template.ID,
		Code:        template.Code,
		Name:        template.Name,
		Description: template.Description,
		CategoryID:  template.CategoryID,
		Type:        template.Type,
		Version:     template.Version,
		Status:      int(template.Status),
		Conditions:  template.Conditions,
		LuaScript:   template.LuaScript,
		Formula:     template.Formula,
		FormulaVars: template.FormulaVars,
		Parameters:  template.Parameters,
		Priority:    template.Priority,
		Sorting:     template.Sorting,
		CreatedAt:   template.CreatedAt,
		UpdatedAt:   template.UpdatedAt,
		TenantID:    template.TenantID,
	}
}
