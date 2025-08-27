package casbin

import (
	"context"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/constant"
	psb "github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/casbin"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence/repository"
)

type RepositoryImpl struct {
	rr repository.ISysRoleRepo
	pr repository.IPermissionsRepo
	tr repository.ISysTenantRepo
}

func NewRepositoryImpl(rr repository.ISysRoleRepo, pr repository.IPermissionsRepo, tr repository.ISysTenantRepo) psb.IPermissionsRepository {
	return &RepositoryImpl{
		rr: rr,
		pr: pr,
		tr: tr,
	}
}

// FindAllEnabled 获取所有启用的角色及其权限
func (r *RepositoryImpl) FindAllEnabled(ctx context.Context) ([]*psb.Role, error) {
	ctx = context.Background()
	// 获取所有启用的角色
	roles, err := r.rr.FindAllEnabled(ctx)
	if err != nil {
		hlog.CtxErrorf(ctx, "casbin [FindAllEnabled] error: %v", err)
		return nil, err
	}
	if len(roles) == 0 {
		return []*psb.Role{}, nil
	}

	// 获取角色ID列表
	roleIds := make([]int64, len(roles))
	for i, role := range roles {
		roleIds[i] = role.ID
	}
	roleResourcesMap, err := r.pr.GetResourcesByRolesGrouped(ctx, roleIds)
	// 转换为 casbin 角色格式
	var casbinRoles []*psb.Role
	for _, role := range roles {
		casbinRole := &psb.Role{
			Id:       fmt.Sprintf("%d", role.ID),
			Code:     role.Code,
			TenantID: role.TenantID,
		}
		// 添加权限
		if resources, ok := roleResourcesMap[role.ID]; ok {
			for _, resource := range resources {
				casbinRole.Permissions = append(casbinRole.Permissions, psb.ApiPermissions{
					Id:     fmt.Sprintf("%d", resource.PermissionsID),
					Method: resource.Method,
					Path:   resource.Path,
				})
			}
		}

		casbinRoles = append(casbinRoles, casbinRole)
	}
	// 获取租户对应的权限列表
	tenants, err := r.tr.GetAllEnabled(ctx)
	if err != nil {
		hlog.CtxErrorf(ctx, "casbin [FindAllEnabled] error: %v", err)
		return nil, err
	}
	if len(tenants) > 0 {
		for _, tenant := range tenants {
			tenantResources, err := r.tr.GetTenantPermissionsResource(ctx, tenant.ID)
			if err != nil {
				hlog.CtxErrorf(ctx, "casbin [FindAllEnabled] error: %v", err)
				return nil, err
			}
			casbinRole := &psb.Role{
				Id:       constant.RoleTenantAdmin,
				Code:     constant.RoleTenantAdmin,
				TenantID: tenant.ID,
			}
			// 添加权限
			for _, resource := range tenantResources {
				casbinRole.Permissions = append(casbinRole.Permissions, psb.ApiPermissions{
					Id:     fmt.Sprintf("%d", resource.ID),
					Method: resource.Method,
					Path:   resource.Path,
				})
			}

			casbinRoles = append(casbinRoles, casbinRole)
		}
	}
	return casbinRoles, nil
}
