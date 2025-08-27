package handler

import (
	"context"
	"encoding/json"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/repository"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/dto"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/queries"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/model"
)

// RuleQueryHandler 规则查询处理器
type RuleQueryHandler struct {
	ruleRepo repository.IRuleRepository
}

// NewRuleQueryHandler 创建规则查询处理器
func NewRuleQueryHandler(ruleRepo repository.IRuleRepository) *RuleQueryHandler {
	return &RuleQueryHandler{
		ruleRepo: ruleRepo,
	}
}

// HandleGetRuleList 处理获取规则列表查询
func (h *RuleQueryHandler) HandleGetRuleList(ctx context.Context, req *queries.GetRuleListReq) ([]*dto.RuleDTO, int64, *herrors.HError) {
	// 构建查询条件
	query := db_query.NewQueryBuilder()
	if req.Name != "" {
		query.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	if req.Code != "" {
		query.Where("code", db_query.Eq, req.Code)
	}
	if req.CategoryID != "" {
		query.Where("category_id", db_query.Eq, req.CategoryID)
	}
	if req.TemplateID != "" {
		query.Where("template_id", db_query.Eq, req.TemplateID)
	}
	if req.Type != "" {
		query.Where("type", db_query.Eq, req.Type)
	}
	if req.Trigger != "" {
		query.Where("trigger", db_query.Eq, req.Trigger)
	}
	if req.Scope != "" {
		query.Where("scope", db_query.Eq, req.Scope)
	}
	if req.ScopeID != "" {
		query.Where("scope_id", db_query.Eq, req.ScopeID)
	}
	if req.BusinessType != "" {
		query.Where("business_type", db_query.Eq, req.BusinessType)
	}
	if req.Status > 0 {
		query.Where("status", db_query.Eq, req.Status)
	}

	// 获取总数
	total, err := h.ruleRepo.Count(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Count rules error: %v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	// 分页查询
	query.OrderByASC("priority")
	query.OrderByASC("sorting")
	query.OrderByDESC("created_at")
	query.WithPage(&req.Page)
	rules, err := h.ruleRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Get rule list error: %v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = h.convertToDTO(rule)
	}

	return dtos, total, nil
}

// HandleGetRulesByCategory 处理根据分类获取规则列表查询
func (h *RuleQueryHandler) HandleGetRulesByCategory(ctx context.Context, req *queries.GetRulesByCategoryReq) ([]*dto.RuleDTO, *herrors.HError) {
	rules, err := h.ruleRepo.FindByCategoryID(ctx, req.CategoryID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = h.convertToDTO(rule)
	}

	return dtos, nil
}

// HandleGetRulesByTemplate 处理根据模板获取规则列表查询
func (h *RuleQueryHandler) HandleGetRulesByTemplate(ctx context.Context, req *queries.GetRulesByTemplateReq) ([]*dto.RuleDTO, *herrors.HError) {
	rules, err := h.ruleRepo.FindByTemplateID(ctx, req.TemplateID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = h.convertToDTO(rule)
	}

	return dtos, nil
}

// HandleGetRulesByType 处理根据类型获取规则列表查询
func (h *RuleQueryHandler) HandleGetRulesByType(ctx context.Context, req *queries.GetRulesByTypeReq) ([]*dto.RuleDTO, *herrors.HError) {
	rules, err := h.ruleRepo.FindByType(ctx, req.Type)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = h.convertToDTO(rule)
	}

	return dtos, nil
}

// HandleGetRulesByTrigger 处理根据触发条件获取规则列表查询
func (h *RuleQueryHandler) HandleGetRulesByTrigger(ctx context.Context, req *queries.GetRulesByTriggerReq) ([]*dto.RuleDTO, *herrors.HError) {
	rules, err := h.ruleRepo.FindByTrigger(ctx, req.Trigger)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = h.convertToDTO(rule)
	}

	return dtos, nil
}

// HandleGetRulesByScope 处理根据作用域获取规则列表查询
func (h *RuleQueryHandler) HandleGetRulesByScope(ctx context.Context, req *queries.GetRulesByScopeReq) ([]*dto.RuleDTO, *herrors.HError) {
	rules, err := h.ruleRepo.FindByScope(ctx, req.Scope)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = h.convertToDTO(rule)
	}

	return dtos, nil
}

// HandleGetRulesByBusinessType 处理根据业务类型获取规则列表查询
func (h *RuleQueryHandler) HandleGetRulesByBusinessType(ctx context.Context, req *queries.GetRulesByBusinessTypeReq) ([]*dto.RuleDTO, *herrors.HError) {
	rules, err := h.ruleRepo.FindByBusinessType(ctx, req.BusinessType)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = h.convertToDTO(rule)
	}

	return dtos, nil
}

// HandleGetEnabledRules 处理获取启用的规则列表查询
func (h *RuleQueryHandler) HandleGetEnabledRules(ctx context.Context, req *queries.GetEnabledRulesReq) ([]*dto.RuleDTO, int64, *herrors.HError) {
	// 构建查询条件
	query := db_query.NewQueryBuilder()
	if req.Name != "" {
		query.Where("name", db_query.Like, "%"+req.Name+"%")
	}
	if req.CategoryID != "" {
		query.Where("category_id", db_query.Eq, req.CategoryID)
	}
	if req.Type != "" {
		query.Where("type", db_query.Eq, req.Type)
	}
	if req.Trigger != "" {
		query.Where("trigger", db_query.Eq, req.Trigger)
	}
	if req.BusinessType != "" {
		query.Where("business_type", db_query.Eq, req.BusinessType)
	}
	query.Where("status", db_query.Eq, 1) // 只查询启用的规则

	// 获取总数
	total, err := h.ruleRepo.Count(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Count enabled rules error: %v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	// 分页查询
	query.OrderByASC("priority")
	query.OrderByASC("sorting")
	query.OrderByDESC("created_at")
	query.WithPage(&req.Page)
	rules, err := h.ruleRepo.Find(ctx, query)
	if err != nil {
		hlog.CtxErrorf(ctx, "Get enabled rules error: %v", err)
		return nil, 0, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = h.convertToDTO(rule)
	}

	return dtos, total, nil
}

// HandleGetRule 处理获取规则详情查询
func (h *RuleQueryHandler) HandleGetRule(ctx context.Context, req *queries.GetRuleReq) (*dto.RuleDTO, *herrors.HError) {
	rule, err := h.ruleRepo.FindByID(ctx, req.ID)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	return h.convertToDTO(rule), nil
}

// HandleGetRuleByCode 处理根据编码获取规则查询
func (h *RuleQueryHandler) HandleGetRuleByCode(ctx context.Context, req *queries.GetRuleByCodeReq) (*dto.RuleDTO, *herrors.HError) {
	rule, err := h.ruleRepo.FindByCode(ctx, req.Code)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	return h.convertToDTO(rule), nil
}

// HandleGetAllRules 处理获取所有规则查询
func (h *RuleQueryHandler) HandleGetAllRules(ctx context.Context) ([]*dto.RuleDTO, *herrors.HError) {
	rules, err := h.ruleRepo.FindAll(ctx)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	dtos := make([]*dto.RuleDTO, len(rules))
	for i, rule := range rules {
		dtos[i] = h.convertToDTO(rule)
	}

	return dtos, nil
}

// convertToDTO 转换为DTO
func (h *RuleQueryHandler) convertToDTO(rule *model.Rule) *dto.RuleDTO {
	return &dto.RuleDTO{
		ID:              rule.ID,
		Code:            rule.Code,
		Name:            rule.Name,
		Description:     rule.Description,
		CategoryID:      rule.CategoryID,
		TemplateID:      rule.TemplateID,
		Type:            rule.Type,
		Action:          rule.Action,
		Triggers:        rule.Triggers,
		Scope:           rule.Scope,
		ScopeID:         rule.ScopeID,
		ExecutionTiming: rule.ExecutionTiming,
		Condition:       h.convertConditionToDTO(rule.Conditions),
		LuaScript:       rule.LuaScript,
		Formula:         rule.Formula,
		Priority:        rule.Priority,
		Sorting:         rule.Sorting,
		Status:          rule.Status,
		CreatedAt:       rule.CreatedAt,
		UpdatedAt:       rule.UpdatedAt,
		TenantID:        rule.TenantID,
	}
}

// convertConditionToDTO 转换条件配置为DTO
func (h *RuleQueryHandler) convertConditionToDTO(conditions string) *dto.ConditionDTO {
	if conditions == "" {
		return nil
	}

	// 解析条件JSON
	var conditionMap map[string]interface{}
	if err := json.Unmarshal([]byte(conditions), &conditionMap); err != nil {
		return &dto.ConditionDTO{
			Type:       "expression",
			Expression: conditions,
			Parameters: make(map[string]interface{}),
		}
	}

	// 提取条件类型和表达式
	conditionType, _ := conditionMap["type"].(string)
	if conditionType == "" {
		conditionType = "expression"
	}

	expression, _ := conditionMap["expression"].(string)
	if expression == "" {
		expression = conditions
	}

	parameters, _ := conditionMap["parameters"].(map[string]interface{})
	if parameters == nil {
		parameters = make(map[string]interface{})
	}

	return &dto.ConditionDTO{
		Type:       conditionType,
		Expression: expression,
		Parameters: parameters,
	}
}
