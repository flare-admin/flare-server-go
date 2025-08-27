package queries

import "github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"

// GetCacheQuery 获取缓存查询
type GetCacheQuery struct {
	GroupID string `json:"group_id"` // 分组ID
	Key     string `json:"key"`      // 缓存键
}

// Validate 验证查询
func (q *GetCacheQuery) Validate() error {
	if q.GroupID == "" {
		return herrors.NewBadReqError("分组ID不能为空")
	}
	if q.Key == "" {
		return herrors.NewBadReqError("缓存键不能为空")
	}
	return nil
}

// GetGroupQuery 获取分组查询
type GetGroupQuery struct {
	GroupID string `json:"group_id"` // 分组ID
}

// Validate 验证查询
func (q *GetGroupQuery) Validate() error {
	if q.GroupID == "" {
		return herrors.NewBadReqError("分组ID不能为空")
	}
	return nil
}

// ListGroupKeysQuery 获取分组下所有键的查询
type ListGroupKeysQuery struct {
	GroupID string `json:"group_id"` // 分组ID
}

// Validate 验证查询
func (q *ListGroupKeysQuery) Validate() error {
	if q.GroupID == "" {
		return herrors.NewBadReqError("分组ID不能为空")
	}
	return nil
}
