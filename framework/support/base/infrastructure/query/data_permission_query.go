package query

import (
	"context"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/dto"
)

type IDataPermissionQuery interface {
	// GetByRoleID 获取角色的数据权限
	GetByRoleID(ctx context.Context, roleID int64) (*dto.DataPermissionDto, herrors.Herr)
}
