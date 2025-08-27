package handlers

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/dto"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/base/application/queries"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query"
)

type DataPermissionQueryHandler struct {
	permQuery query.IDataPermissionQuery
}

func NewDataPermissionQueryHandler(permQuery query.IDataPermissionQuery) *DataPermissionQueryHandler {
	return &DataPermissionQueryHandler{
		permQuery: permQuery,
	}
}

// HandleGetByRoleID 处理获取角色数据权限
func (h *DataPermissionQueryHandler) HandleGetByRoleID(ctx context.Context, query queries.GetDataPermissionQuery) (*dto.DataPermissionDto, herrors.Herr) {
	return h.permQuery.GetByRoleID(ctx, query.RoleID)
}
