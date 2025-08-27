package rest

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/jwt"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/support/monitoring/application/handlers"
	"github.com/flare-admin/flare-server-go/framework/support/monitoring/application/queries"

	"github.com/cloudwego/hertz/pkg/route"
)

type MetricsController struct {
	queryHandler *handlers.MetricsQueryHandler
}

func NewMetricsController(queryHandler *handlers.MetricsQueryHandler) *MetricsController {
	return &MetricsController{
		queryHandler: queryHandler,
	}
}

func (c *MetricsController) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	metrics := g.Group("/v1/metrics", jwt.Handler(t))
	{
		metrics.GET("/system", hserver.NewHandlerFu[queries.GetSystemMetricsQuery](c.GetSystemMetrics))
		metrics.GET("/runtime", hserver.NewNotParHandlerFu(c.GetRuntimeMetrics))
	}
}

// GetSystemMetrics 获取系统指标
// @Summary 获取系统指标
// @Description 获取系统CPU、内存、磁盘使用率
// @Tags 监控指标
// @Accept json
// @Produce json
// @Success 200 {object} base_info.Success{data=dto.SystemMetricsDto}
// @Router /v1/metrics/system [get]
func (c *MetricsController) GetSystemMetrics(ctx context.Context, q *queries.GetSystemMetricsQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleGetSystemMetrics(ctx, q)
	if herrors.HaveError(err) {
		return result.WithError(err)
	}
	return result.WithData(data)
}

// GetRuntimeMetrics 获取运行时指标
// @Summary 获取运行时指标
// @Description 获取Go运行时内存、GC等指标
// @Tags 监控指标
// @Accept json
// @Produce json
// @Success 200 {object} base_info.Success{data=dto.RuntimeMetricsDto}
// @Router /v1/metrics/system [get]
func (c *MetricsController) GetRuntimeMetrics(ctx context.Context) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	data, err := c.queryHandler.HandleGetRuntimeMetrics(ctx)
	if herrors.HaveError(err) {
		return result.WithError(err)
	}
	return result.WithData(data)
}
