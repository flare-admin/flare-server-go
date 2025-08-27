package dto

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
)

// CategoryCreateReq 创建分类请求
type CategoryCreateReq struct {
	ID          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	I18nKey     string `json:"i18nKey"` // 国际化key
	Description string `json:"description"`
}

// CategoryUpdateReq 更新分类请求
type CategoryUpdateReq struct {
	ID          string `json:"id" path:"id"`
	Name        string `json:"name"`
	I18nKey     string `json:"i18nKey"` // 国际化key
	Description string `json:"description"`
}

// CategoryQueryReq 分页查询分类列表
type CategoryQueryReq struct {
	db_query.Page
	Id      string `json:"id" query:"id"`
	Name    string `json:"name" query:"name"`
	I18nKey string `json:"i18nKey" query:"i18nKey"`
}

type Category struct {
	database.BaseModel
	ID          string `json:"id"`
	Name        string `json:"name"`
	I18nKey     string `json:"i18nKey"` // 国际化key
	Description string `json:"description"`
}

// OptionCreateReq 创建选项请求
type OptionCreateReq struct {
	CategoryID string `json:"categoryId" binding:"required"`
	Label      string `json:"label"` //默认名称
	Value      string `json:"value" binding:"required"`
	I18nKey    string `json:"i18nKey"` // 国际化key
	Sort       int    `json:"sort"`
	Status     int    `json:"status"`
	Remark     string `json:"remark"`
}

// OptionUpdateReq 更新选项请求
type OptionUpdateReq struct {
	ID      string `json:"id" path:"id"`
	Label   string `json:"label"` //默认名称
	Value   string `json:"value"`
	I18nKey string `json:"i18nKey"` // 国际化key
	Sort    int    `json:"sort"`
	Status  int    `json:"status"`
	Remark  string `json:"remark"`
}

// OptionQueryReq 查询选项请求
type OptionQueryReq struct {
	CategoryID string `json:"categoryId" query:"categoryId"`
	Keyword    string `json:"keyword" query:"keyword"`
	Status     *int   `json:"status" query:"status"`
}

type Option struct {
	database.BaseModel
	ID         string `json:"id"`         // 选项ID
	CategoryID string `json:"categoryId"` //分类ID
	Label      string `json:"label"`      //默认名称
	I18nKey    string `json:"i18nKey"`    //国际化key
	Value      string `json:"value"`      //选项值
	Sort       int    `json:"sort"`       // 排序号
	Status     int    `json:"status"`     // 状态:1-启用,0-禁用
	Remark     string `json:"remark"`     // 备注
}
