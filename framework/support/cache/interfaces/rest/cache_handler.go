package rest

import (
	"context"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/jwt"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/oplog"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/support/cache/application/commands"
	"github.com/flare-admin/flare-server-go/framework/support/cache/application/queries"
	"github.com/flare-admin/flare-server-go/framework/support/cache/application/service"
)

// CacheHandler 缓存处理器
type CacheHandler struct {
	cacheService *service.CacheService
	moduleName   string
}

// NewCacheHandler 创建缓存处理器
func NewCacheHandler(cacheService *service.CacheService) *CacheHandler {
	return &CacheHandler{
		cacheService: cacheService,
		moduleName:   "缓存管理",
	}
}

// RegisterRouter 注册路由
func (h *CacheHandler) RegisterRouter(g *route.RouterGroup, t token.IToken) {
	v1 := g.Group("/v1")
	cr := v1.Group("/cache", jwt.Handler(t))
	{
		// 获取所有缓存分组
		cr.GET("/groups", hserver.NewNotParHandlerFu(h.ListGroups))

		// 清除指定分组的缓存
		cr.DELETE("/group/:groupID", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      h.moduleName,
			Action:      "清除分组缓存",
		}), hserver.NewHandlerFu[commands.DeleteGroupCommand](h.DeleteGroup))

		// 获取指定缓存的值
		cr.GET("/:groupID/:key", hserver.NewHandlerFu[queries.GetCacheQuery](h.GetCache))

		// 清除指定的缓存
		cr.DELETE("/:groupID/:key", oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      h.moduleName,
			Action:      "清除缓存",
		}), hserver.NewHandlerFu[commands.DeleteCacheCommand](h.DeleteCache))

		// 获取分组下的所有键
		cr.GET("/group/:groupID/keys", hserver.NewHandlerFu[queries.ListGroupKeysQuery](h.ListGroupKeys))
	}
}

// ListGroups 获取所有缓存分组
// @Summary 获取所有缓存分组
// @Description 获取所有缓存分组及其统计信息
// @Tags 缓存管理
// @Accept json
// @Produce json
// @Success 200 {object} base_info.Success{data=[]model.CacheGroupInfo}
// @Router /v1/cache/groups [get]
func (h *CacheHandler) ListGroups(ctx context.Context) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	groups, err := h.cacheService.ListGroups(ctx)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(groups)
}

// DeleteGroup 清除指定分组的缓存
// @Summary 清除指定分组的缓存
// @Description 清除指定分组的缓存
// @Tags 缓存管理
// @Accept json
// @Produce json
// @Param groupID path string true "分组ID"
// @Success 200 {object} base_info.Success
// @Router /v1/cache/group/{groupID} [delete]
func (h *CacheHandler) DeleteGroup(ctx context.Context, cmd *commands.DeleteGroupCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := h.cacheService.DeleteGroup(ctx, cmd)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// GetCache 获取指定缓存的值
// @Summary 获取指定缓存的值
// @Description 获取指定缓存的值
// @Tags 缓存管理
// @Accept json
// @Produce json
// @Param groupID path string true "分组ID"
// @Param key path string true "缓存键"
// @Success 200 {object} base_info.Success{data=interface{}}
// @Router /v1/cache/{groupID}/{key} [get]
func (h *CacheHandler) GetCache(ctx context.Context, query *queries.GetCacheQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	value, err := h.cacheService.GetCache(ctx, query)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(value)
}

// DeleteCache 清除指定的缓存
// @Summary 清除指定的缓存
// @Description 清除指定的缓存
// @Tags 缓存管理
// @Accept json
// @Produce json
// @Param groupID path string true "分组ID"
// @Param key path string true "缓存键"
// @Success 200 {object} base_info.Success
// @Router /v1/cache/{groupID}/{key} [delete]
func (h *CacheHandler) DeleteCache(ctx context.Context, cmd *commands.DeleteCacheCommand) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	err := h.cacheService.DeleteCache(ctx, cmd)
	if err != nil {
		return result.WithError(err)
	}
	return result
}

// ListGroupKeys 获取分组下的所有键
// @Summary 获取分组下的所有键
// @Description 获取指定分组下的所有缓存键
// @Tags 缓存管理
// @Accept json
// @Produce json
// @Param groupID path string true "分组ID"
// @Success 200 {object} base_info.Success{data=[]string}
// @Router /v1/cache/group/{groupID}/keys [get]
func (h *CacheHandler) ListGroupKeys(ctx context.Context, query *queries.ListGroupKeysQuery) *hserver.ResponseResult {
	result := hserver.DefaultResponseResult()
	keys, err := h.cacheService.ListGroupKeys(ctx, query)
	if err != nil {
		return result.WithError(err)
	}
	return result.WithData(keys)
}
