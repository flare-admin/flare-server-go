package utils

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/repository"
)

func IsTenantAdmin(ctx context.Context, user *model.User, tenantRepo repository.ITenantRepository) (bool, *model.Tenant, error) {
	tenantId := actx.GetTenantId(ctx)
	if tenantId != "" {
		tenant, err := tenantRepo.FindByID(context.Background(), tenantId)
		if err != nil {
			return false, nil, err
		}
		//租户管理处理
		if user != nil && tenant.AdminUser.ID == user.ID {
			return true, tenant, nil
		}
	}
	return false, nil, nil
}
