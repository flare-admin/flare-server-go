package queries

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
)

// ==================== 模板相关查询 ====================

// GetTemplateListReq 获取模板列表请求
type GetTemplateListReq struct {
	db_query.Page
	Name       string `form:"name" query:"name" json:"name"`                   // 模板名称
	CategoryID string `form:"categoryId" query:"categoryId" json:"categoryId"` // 分类ID
	Code       string `form:"code" query:"code" json:"code"`                   // 模板编码
	Type       string `form:"type" query:"type" json:"type"`                   // 模板类型
	Status     int32  `form:"status" query:"status" json:"status"`             // 状态
}

// GetTemplatesByCategoryReq 根据分类获取模板列表请求
type GetTemplatesByCategoryReq struct {
	CategoryID string `form:"categoryId" query:"categoryId"` // 分类ID
}

// GetTemplatesByTypeReq 根据类型获取模板列表请求
type GetTemplatesByTypeReq struct {
	Type string `form:"type" query:"type"` // 模板类型
}

// GetEnabledTemplatesReq 获取启用的模板列表请求
type GetEnabledTemplatesReq struct {
	db_query.Page
	Name       string `form:"name" query:"name" json:"name"`                   // 模板名称
	CategoryID string `form:"categoryId" query:"categoryId" json:"categoryId"` // 分类ID
	Code       string `form:"code" query:"code" json:"code"`                   // 模板编码
	Type       string `form:"type" query:"type" json:"type"`                   // 模板类型
}

// GetTemplateReq 获取模板详情请求
type GetTemplateReq struct {
	ID string `form:"id" query:"id"` // 模板ID
}

// GetTemplateByCodeReq 根据编码获取模板请求
type GetTemplateByCodeReq struct {
	Code string `form:"code" query:"code"` // 模板编码
}

// ==================== 分类相关查询 ====================

// GetCategoryListReq 获取分类列表请求
type GetCategoryListReq struct {
	db_query.Page
	Name         string `form:"name" query:"name" json:"name"`                         // 分类名称
	Code         string `form:"code" query:"code" json:"code"`                         // 分类编码
	ParentID     string `form:"parentId" query:"parentId" json:"parentId"`             // 父分类ID
	Type         string `form:"type" query:"type" json:"type"`                         // 分类类型
	BusinessType string `form:"businessType" query:"businessType" json:"businessType"` // 业务类型
	Status       int32  `form:"status" query:"status" json:"status"`                   // 状态
}

// GetCategoriesByParentReq 根据父分类获取子分类列表请求
type GetCategoriesByParentReq struct {
	ParentID string `form:"parentId" query:"parentId"` // 父分类ID
}

// GetCategoriesByBusinessTypeReq 根据业务类型获取分类列表请求
type GetCategoriesByBusinessTypeReq struct {
	BusinessType string `form:"businessType" query:"businessType"` // 业务类型
}

// GetCategoriesByTypeReq 根据分类类型获取分类列表请求
type GetCategoriesByTypeReq struct {
	Type string `form:"type" query:"type"` // 分类类型
}

// GetRootCategoriesReq 获取根分类列表请求
type GetRootCategoriesReq struct {
	BusinessType string `form:"businessType" query:"businessType"` // 业务类型
}

// GetCategoryTreeReq 获取分类树请求
type GetCategoryTreeReq struct {
	CategoryID string `form:"categoryId" query:"categoryId"` // 分类ID（为空时获取所有根分类）
}

// GetCategoryReq 获取分类详情请求
type GetCategoryReq struct {
	ID string `form:"id" query:"id"` // 分类ID
}

// GetCategoryByCodeReq 根据编码获取分类请求
type GetCategoryByCodeReq struct {
	Code string `form:"code" query:"code"` // 分类编码
}

// ==================== 规则相关查询 ====================

// GetRuleListReq 获取规则列表请求
type GetRuleListReq struct {
	db_query.Page
	Name         string `form:"name" query:"name" json:"name"`                         // 规则名称
	Code         string `form:"code" query:"code" json:"code"`                         // 规则编码
	CategoryID   string `form:"categoryId" query:"categoryId" json:"categoryId"`       // 分类ID
	TemplateID   string `form:"templateId" query:"templateId" json:"templateId"`       // 模板ID
	Type         string `form:"type" query:"type" json:"type"`                         // 规则类型
	Trigger      string `form:"trigger" query:"trigger" json:"trigger"`                // 触发条件
	Scope        string `form:"scope" query:"scope" json:"scope"`                      // 作用域
	ScopeID      string `form:"scopeId" query:"scopeId" json:"scopeId"`                // 作用域ID
	BusinessType string `form:"businessType" query:"businessType" json:"businessType"` // 业务类型
	Status       int32  `form:"status" query:"status" json:"status"`                   // 状态
}

// GetRulesByCategoryReq 根据分类获取规则列表请求
type GetRulesByCategoryReq struct {
	CategoryID string `form:"categoryId" query:"categoryId"` // 分类ID
}

// GetRulesByTemplateReq 根据模板获取规则列表请求
type GetRulesByTemplateReq struct {
	TemplateID string `form:"templateId" query:"templateId"` // 模板ID
}

// GetRulesByTypeReq 根据类型获取规则列表请求
type GetRulesByTypeReq struct {
	Type string `form:"type" query:"type"` // 规则类型
}

// GetRulesByTriggerReq 根据触发条件获取规则列表请求
type GetRulesByTriggerReq struct {
	Trigger string `form:"trigger" query:"trigger"` // 触发条件
}

// GetRulesByScopeReq 根据作用域获取规则列表请求
type GetRulesByScopeReq struct {
	Scope string `form:"scope" query:"scope"` // 作用域
}

// GetRulesByBusinessTypeReq 根据业务类型获取规则列表请求
type GetRulesByBusinessTypeReq struct {
	BusinessType string `form:"businessType" query:"businessType"` // 业务类型
}

// GetEnabledRulesReq 获取启用的规则列表请求
type GetEnabledRulesReq struct {
	db_query.Page
	Name         string `form:"name" query:"name" json:"name"`                         // 规则名称
	CategoryID   string `form:"categoryId" query:"categoryId" json:"categoryId"`       // 分类ID
	Type         string `form:"type" query:"type" json:"type"`                         // 规则类型
	Trigger      string `form:"trigger" query:"trigger" json:"trigger"`                // 触发条件
	BusinessType string `form:"businessType" query:"businessType" json:"businessType"` // 业务类型
}

// GetRuleReq 获取规则详情请求
type GetRuleReq struct {
	ID string `form:"id" query:"id"` // 规则ID
}

// GetRuleByCodeReq 根据编码获取规则请求
type GetRuleByCodeReq struct {
	Code string `form:"code" query:"code"` // 规则编码
}
