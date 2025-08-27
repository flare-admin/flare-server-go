package handlers

import (
	"context"

	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/base/application/commands"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/service"
)

type DataPermissionCommandHandler struct {
	permService *service.DataPermissionService
}

func NewDataPermissionCommandHandler(permService *service.DataPermissionService) *DataPermissionCommandHandler {
	return &DataPermissionCommandHandler{
		permService: permService,
	}
}

// HandleAssign 处理分配数据权限
func (h *DataPermissionCommandHandler) HandleAssign(ctx context.Context, cmd *commands.AssignDataPermissionCommand) herrors.Herr {
	// 1. 构建数据权限领域模型
	perm := &model.DataPermission{
		RoleID:   cmd.RoleID,
		Scope:    model.DataScope(cmd.Scope),
		DeptIDs:  cmd.DeptIDs,
		TenantID: actx.GetTenantId(ctx),
	}

	// 2. 调用领域服务分配数据权限
	return h.permService.AssignDataPermission(ctx, perm)
}

// HandleRemove 处理移除数据权限
func (h *DataPermissionCommandHandler) HandleRemove(ctx context.Context, cmd *commands.RemoveDataPermissionCommand) herrors.Herr {
	return h.permService.RemoveDataPermission(ctx, cmd.RoleID)
}
