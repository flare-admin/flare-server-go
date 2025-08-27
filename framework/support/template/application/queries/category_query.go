package queries

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
)

// GetCategoryListReq 获取分类列表请求
type GetCategoryListReq struct {
	db_query.Page
	Name   string `form:"name" query:"name" json:"name" comment:"分类名称"`
	Code   string `form:"code" query:"code" json:"code" comment:"分类编码"`
	Status int32  `form:"status" query:"status" json:"status" comment:"状态"`
}

// GetCategoryByCodeReq 根据编码获取分类请求
type GetCategoryByCodeReq struct {
	Code string `form:"code" query:"code" binding:"required" comment:"分类编码"`
}
