package handlers

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/dto"
	"time"

	"github.com/flare-admin/flare-server-go/framework/pkg/actx"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/support/base/application/queries"

	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query"
)

type OperationLogQueryHandler struct {
	query query.IOperationLogQuery
}

func NewOperationLogQueryHandler(query query.IOperationLogQuery) *OperationLogQueryHandler {
	return &OperationLogQueryHandler{
		query: query,
	}
}

// HandleList 处理查询操作日志列表
func (h *OperationLogQueryHandler) HandleList(ctx context.Context, q *queries.ListOperationLogQuery) (*models.PageRes[dto.OperationLogDto], herrors.Herr) {
	tm := utils.GetTimeNow()
	if q.Month != "" {
		// 解析查询月份
		month, err := time.Parse("200601", q.Month)
		if err != nil {
			return nil, herrors.NewBadReqError("invalid month format")
		}
		tm = month
	}

	// 构建查询条件
	qb := db_query.NewQueryBuilder()
	if q.Username != "" {
		qb.Where("username", db_query.Like, "%"+q.Username+"%")
	}
	if q.Module != "" {
		qb.Where("ip", db_query.Like, "%"+q.Module+"%")
	}
	if q.Action != "" {
		qb.Where("status", db_query.Eq, q.Action)
	}
	if q.StartTime > 0 {
		qb.Where("login_time", db_query.Gte, time.Unix(q.StartTime, 0))
	}
	if q.EndTime > 0 {
		qb.Where("login_time", db_query.Lte, time.Unix(q.EndTime, 0))
	}

	// 设置排序
	qb.OrderBy("created_at", false)

	// 设置分页
	qb.WithPage(&q.Page)

	tenant := actx.GetTenantId(ctx)
	// 获取总数
	total, err := h.query.Count(ctx, tenant, tm, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	// 查询数据
	logs, err := h.query.Find(ctx, tenant, tm, qb)
	if err != nil {
		return nil, herrors.NewErr(err)
	}

	return &models.PageRes[dto.OperationLogDto]{
		List:  logs,
		Total: total,
	}, nil
}
