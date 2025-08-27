package base_api

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/dto"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query"
)

type ITenantApi interface {
	// GetDefTenant 获取默认租户
	GetDefTenant(ctx context.Context) (*dto.TenantDto, error)
	// GetTenant 获取租户信息
	GetTenant(ctx context.Context, id string) (*dto.TenantDto, error)
}

type tenantApi struct {
	tenantQueryService query.ITenantQueryService
}

func NewTenantApi(tenantQueryService query.ITenantQueryService) ITenantApi {
	return &tenantApi{
		tenantQueryService: tenantQueryService,
	}
}

func (t tenantApi) GetDefTenant(ctx context.Context) (*dto.TenantDto, error) {
	return t.tenantQueryService.GetDefTenant(ctx)
}

func (t tenantApi) GetTenant(ctx context.Context, id string) (*dto.TenantDto, error) {
	return t.tenantQueryService.GetTenant(ctx, id)
}
