package commands

import "github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"

// DeleteCacheCommand 删除缓存命令
type DeleteCacheCommand struct {
	GroupID string `json:"group_id"` // 分组ID
	Key     string `json:"key"`      // 缓存键
}

// Validate 验证命令
func (c *DeleteCacheCommand) Validate() error {
	if c.GroupID == "" {
		return herrors.NewBadReqError("分组ID不能为空")
	}
	if c.Key == "" {
		return herrors.NewBadReqError("缓存键不能为空")
	}
	return nil
}

// DeleteGroupCommand 删除分组命令
type DeleteGroupCommand struct {
	GroupID string `json:"group_id"` // 分组ID
}

// Validate 验证命令
func (c *DeleteGroupCommand) Validate() error {
	if c.GroupID == "" {
		return herrors.NewBadReqError("分组ID不能为空")
	}
	return nil
}
