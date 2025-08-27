package query

import (
	"context"
	"time"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/dto"
)

// ILoginLogQuery 登录日志查询接口
type ILoginLogQuery interface {
	// Find 查询登录日志列表
	Find(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) ([]*dto.LoginLogDto, error)
	// Count 统计登录日志数量
	Count(ctx context.Context, tenantID string, month time.Time, qb *db_query.QueryBuilder) (int64, error)
}
