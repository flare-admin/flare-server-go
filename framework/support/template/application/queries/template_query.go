package queries

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
)

// GetTemplateListReq 获取模板列表请求
type GetTemplateListReq struct {
	db_query.Page
	Name       string `form:"name" query:"name" json:"name" comment:"模板名称"`
	CategoryID string `form:"category_id" query:"category_id" json:"category_id" comment:"分类ID"`
	Code       string `form:"code" query:"code" json:"code" comment:"模板编码"`
	Status     int32  `form:"status" query:"status" json:"status" comment:"状态"`
}

// GetTemplatesByCategoryReq 根据分类获取模板列表请求
type GetTemplatesByCategoryReq struct {
	CategoryID string `form:"category_id" query:"category_id" comment:"分类ID"`
}

// GetAllEnabledTemplatesByCategory 根据分类获取模板列表请求
type GetAllEnabledTemplatesByCategory struct {
	Code string `form:"code" query:"code" comment:"分类编码"`
}

// GetEnabledTemplateReq 获取启用的模板列表请求
type GetEnabledTemplateReq struct {
	db_query.Page
	Name         string `form:"name" query:"name" json:"name" comment:"模板名称"`
	CategoryCode string `form:"category_code" query:"category_code" json:"category_code" comment:"分类ID"`
	Code         string `form:"code" query:"code" json:"code" comment:"模板编码"`
}
