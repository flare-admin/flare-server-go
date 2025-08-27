package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/base/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence/entity"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence/repository"
	"gorm.io/gorm"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
)

type sysTenantRepo struct {
	*baserepo.BaseRepo[entity.Tenant, string]
}

func NewSysTenantRepo(data database.IDataBase) repository.ISysTenantRepo {
	model := new(entity.Tenant)
	// 同步表
	if err := data.AutoMigrate(model, &entity.TenantPermissions{}); err != nil {
		hlog.Fatalf("sync tenant tables to db error: %v", err)
	}
	return &sysTenantRepo{
		BaseRepo: baserepo.NewBaseRepo[entity.Tenant, string](data, entity.Tenant{}),
	}
}
func (r *sysTenantRepo) CommonGetByID(ctx context.Context, id string) (*entity.Tenant, error) {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	return r.BaseRepo.FindById(ctx, id)
}

// GetByCode 根据编码获取租户
func (r *sysTenantRepo) GetByCode(ctx context.Context, code string) (*entity.Tenant, error) {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	var tenant entity.Tenant
	err := r.Db(ctx).Where("code = ?", code).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// DelById 删除租户（包括关联关系）
func (r *sysTenantRepo) DelById(ctx context.Context, id string) error {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 检查是否为默认租户
		tenant, err := r.FindById(ctx, id)
		if err != nil {
			return fmt.Errorf("find tenant failed: %w", err)
		}
		if tenant.IsDefault == 1 {
			return fmt.Errorf("cannot delete default tenant")
		}

		// 删除租户下的所有用户
		if err := r.Db(ctx).Where("tenant_id = ?", id).Delete(&entity.SysUser{}).Error; err != nil {
			return err
		}

		// 删除租户下的所有角色
		if err := r.Db(ctx).Where("tenant_id = ?", id).Delete(&entity.Role{}).Error; err != nil {
			return err
		}

		// 删除租户
		return r.Db(ctx).Delete(&entity.Tenant{}, "id = ?", id).Error
	})
}

// Create 创建租户（重写基类方法，处理默认租户）
func (r *sysTenantRepo) Create(ctx context.Context, tenant *entity.Tenant) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		tenant.CreatedAt = utils.GetDateUnix()
		// 如果是默认租户，需要将其他租户设置为非默认
		if tenant.IsDefault == 1 {
			if err := r.Db(ctx).Model(&entity.Tenant{}).Where("is_default = ?", 1).
				Updates(map[string]interface{}{"is_default": 2}).Error; err != nil {
				return err
			}
		}

		// 创建租户
		return r.Db(ctx).Create(tenant).Error
	})
}

// Update 更新租户（重写基类方法，处理默认租户）
func (r *sysTenantRepo) Update(ctx context.Context, tenant *entity.Tenant) error {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		tenant.UpdatedAt = utils.GetDateUnix()
		// 如果是默认租户，需要将其他租户设置为非默认
		if tenant.IsDefault == 1 {
			if err := r.Db(ctx).Model(&entity.Tenant{}).
				Where("id != ? AND is_default = ?", tenant.ID, 1).
				Updates(map[string]interface{}{"is_default": 2}).Error; err != nil {
				return err
			}
		}

		// 更新租户
		return r.Db(ctx).Updates(tenant).Error
	})
}

// GetDefTenant 获取默认租户
func (r *sysTenantRepo) GetDefTenant(ctx context.Context) (*entity.Tenant, error) {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	var tenant entity.Tenant
	err := r.Db(ctx).Where("is_default = ?", 1).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// AssignPermissions 分配权限给租户
func (r *sysTenantRepo) AssignPermissions(ctx context.Context, tenantID string, permissionIDs []int64) error {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 删除原有权限关联
		if err := r.Db(ctx).Where("tenant_id = ?", tenantID).Delete(&entity.TenantPermissions{}).Error; err != nil {
			return err
		}

		// 创建新的权限关联
		pars := make([]*entity.TenantPermissions, 0)
		for _, permID := range permissionIDs {
			tp := &entity.TenantPermissions{
				TenantID:     tenantID,
				PermissionID: permID,
			}
			pars = append(pars, tp)
		}
		if err := r.Db(ctx).Create(&pars).Error; err != nil {
			return err
		}
		// 删除该租户下角色权限不在 permissionIDs 中的权限
		if err := r.Db(ctx).Where("tenant_id = ?", tenantID).
			Not("permission_id IN ?", permissionIDs).
			Delete(&entity.RolePermissions{}).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetPermissionsByTenantID 获取租户的权限列表
func (r *sysTenantRepo) GetPermissionsByTenantID(ctx context.Context, tenantID string) ([]*entity.Permissions, error) {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	var permissions []*entity.Permissions
	err := r.Db(ctx).
		Joins("JOIN sys_tenant_permissions tp ON tp.permission_id = sys_permissions.id").
		Where("tp.tenant_id = ?", tenantID).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// GetTenantPermissionsResource 根据租户ID获取权限列表
func (r *sysTenantRepo) GetTenantPermissionsResource(ctx context.Context, tenantID string) ([]*entity.PermissionsResource, error) {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	permissions := make([]*entity.PermissionsResource, 0)
	err := r.Db(ctx).
		Joins("JOIN sys_tenant_permissions ON sys_permissions_resource.permissions_id = sys_tenant_permissions.permission_id").
		Where("sys_tenant_permissions.tenant_id = ?", tenantID).
		Find(&permissions).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return permissions, nil
		}
		return nil, err
	}
	return permissions, nil
}

// GetTenantIDPermissionsByType 根据租户ID和权限类型获取权限列表
func (r *sysTenantRepo) GetTenantIDPermissionsByType(ctx context.Context, tenantID string, int8 int64) ([]*entity.Permissions, error) {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	var permissions []*entity.Permissions
	err := r.Db(ctx).
		Joins("JOIN sys_tenant_permissions tp ON tp.permission_id = sys_permissions.id").
		Where("tp.tenant_id = ? AND sys_permissions.type = ?", tenantID, int8).
		Find(&permissions).Error
	return permissions, err
}

// HasPermission 检查租户是否拥有指定权限
func (r *sysTenantRepo) HasPermission(ctx context.Context, tenantID string, permissionID int64) (bool, error) {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	var count int64
	err := r.Db(ctx).Model(&entity.TenantPermissions{}).
		Where("tenant_id = ? AND permission_id = ?", tenantID, permissionID).
		Count(&count).Error
	return count > 0, err
}

// DeleteWithRelations 删除租户及关联数据
func (r *sysTenantRepo) DeleteWithRelations(ctx context.Context, id string) error {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 检查是否为默认租户
		tenant, err := r.FindById(ctx, id)
		if err != nil {
			return fmt.Errorf("find tenant failed: %w", err)
		}
		if tenant.IsDefault == 1 {
			return fmt.Errorf("cannot delete default tenant")
		}

		// 删除租户下的所有用户
		if err := r.Db(ctx).Where("tenant_id = ?", id).Delete(&entity.SysUser{}).Error; err != nil {
			return err
		}

		// 删除租户下的所有角色
		if err := r.Db(ctx).Where("tenant_id = ?", id).Delete(&entity.Role{}).Error; err != nil {
			return err
		}

		// 删除租户
		return r.DelById(ctx, id)
	})
}
func (r *sysTenantRepo) GetAllEnabled(ctx context.Context) ([]*entity.Tenant, error) {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	res := make([]*entity.Tenant, 0)
	err := r.Db(ctx).Where("status = ?", 1).Find(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res, nil
		}
		return nil, err
	}
	return res, err
}

// Lock 锁定租户
func (r *sysTenantRepo) Lock(ctx context.Context, tenantID string, reason string) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 检查租户是否存在
		tenant, err := r.FindById(ctx, tenantID)
		if err != nil {
			return fmt.Errorf("find tenant failed: %w", err)
		}
		if tenant == nil {
			return fmt.Errorf("tenant not found: %s", tenantID)
		}

		// 更新租户状态和锁定原因
		return r.Db(ctx).Model(&entity.Tenant{}).
			Where("id = ?", tenantID).
			Updates(map[string]interface{}{
				"status":      model.StatusDisabled,
				"lock_reason": reason,
				"updated_at":  utils.GetDateUnix(),
			}).Error
	})
}

// Unlock 解锁租户
func (r *sysTenantRepo) Unlock(ctx context.Context, tenantID string) error {
	return r.GetDb().InTx(ctx, func(ctx context.Context) error {
		// 检查租户是否存在
		tenant, err := r.FindById(ctx, tenantID)
		if err != nil {
			return fmt.Errorf("find tenant failed: %w", err)
		}
		if tenant == nil {
			return fmt.Errorf("tenant not found: %s", tenantID)
		}

		// 更新租户状态和清空锁定原因
		return r.Db(ctx).Model(&entity.Tenant{}).
			Where("id = ?", tenantID).
			Updates(map[string]interface{}{
				"status":      model.StatusEnabled,
				"lock_reason": "",
				"updated_at":  utils.GetDateUnix(),
			}).Error
	})
}

// GetTenantRoles 获取租户角色
func (r *sysTenantRepo) GetTenantRoles(ctx context.Context, tenantID string) ([]*entity.Role, error) {
	ctx = actx.BuildIgnoreTenantCtx(ctx)
	var roles []*entity.Role
	err := r.Db(ctx).Model(&entity.Role{}).
		Joins("JOIN sys_tenant_role ON sys_tenant_role.role_id = sys_role.id").
		Where("sys_tenant_role.tenant_id = ? AND sys_role.status = ?", tenantID, 1).
		Find(&roles).Error
	return roles, err
}
