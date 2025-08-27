package query

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/dto"
)

// ITenantQueryService 租户查询服务接口
type ITenantQueryService interface {
	GetTenant(ctx context.Context, id string) (*dto.TenantDto, error)
	FindTenants(ctx context.Context, qb *db_query.QueryBuilder) ([]*dto.TenantDto, error)
	CountTenants(ctx context.Context, qb *db_query.QueryBuilder) (int64, error)
	GetDefTenant(ctx context.Context) (*dto.TenantDto, error)
	GetTenantPermissions(ctx context.Context, tenantID string) ([]*dto.PermissionsDto, error)
}
