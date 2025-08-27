package sysevent_service

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/dto"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/service"

	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	_ "github.com/flare-admin/flare-server-go/framework/pkg/hserver/base_info"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/casbin"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/jwt"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/oplog"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
)

type EventService struct {
	as      service.IEventServerApi
	sbs     service.ISubscribeServerApi
	ds      service.IDeadLetterServiceApi
	ef      *casbin.Enforcer
	modeNma string
}

func NewEventService(as service.IEventServerApi, sbs service.ISubscribeServerApi, ds service.IDeadLetterServiceApi, ef *casbin.Enforcer) *EventService {
	return &EventService{
		as:      as,
		sbs:     sbs,
		ds:      ds,
		ef:      ef,
		modeNma: "事件管理",
	}
}

func (a *EventService) RegisterRouter(rg *route.RouterGroup, t token.IToken) {
	g := rg.Group("/v1/event", jwt.Handler(t))
	{
		g.POST("", casbin.Handler(a.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      a.modeNma,
			Action:      "新增",
		}), hserver.NewHandlerFu[dto.AddEventReq](a.AddEvent)) // 新增事件

		g.GET("/:id", casbin.Handler(a.ef),
			hserver.NewHandlerFu[models.StringIdReq](a.GetEvent)) // 获取事件详情

		g.PUT("", casbin.Handler(a.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      a.modeNma,
			Action:      "修改",
		}), hserver.NewHandlerFu[dto.UpdateEventReq](a.UpdateEvent)) // 更新事件

		g.GET("", casbin.Handler(a.ef),
			hserver.NewHandlerFu[dto.GetEventListReq](a.GetEventList)) // 获取事件列表

		g.PUT("/update_status", casbin.Handler(a.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      a.modeNma,
			Action:      "更新状态",
		}), hserver.NewHandlerFu[dto.UpdateEventStatusReq](a.UpdateEventStatus)) // 更新事件状态

		g.GET("/all", hserver.NewNotParHandlerFu(a.GetAllEventList)) // 获取所有事件

		// 事件订阅
		sg := g.Group("/subscribe")
		{
			sg.POST("", casbin.Handler(a.ef), oplog.Record(oplog.LogOption{
				IncludeBody: true,
				Module:      a.modeNma,
				Action:      "新增订阅",
			}), hserver.NewHandlerFu[dto.AddSubscribeReq](a.AddSubscribe)) // 新增订阅

			sg.GET("/:id", casbin.Handler(a.ef),
				hserver.NewHandlerFu[models.StringIdReq](a.GetSubscribe)) // 获取订阅详情

			sg.PUT("", casbin.Handler(a.ef), oplog.Record(oplog.LogOption{
				IncludeBody: true,
				Module:      a.modeNma,
				Action:      "修改订阅",
			}), hserver.NewHandlerFu[dto.UpdateSubscribeReq](a.UpdateSubscribe)) // 更新订阅

			sg.GET("", casbin.Handler(a.ef),
				hserver.NewHandlerFu[dto.GetSubscribeListReq](a.GetSubscribeList)) // 获取订阅列表

			sg.PUT("/enable", casbin.Handler(a.ef), oplog.Record(oplog.LogOption{
				IncludeBody: true,
				Module:      a.modeNma,
				Action:      "启用订阅",
			}), hserver.NewHandlerFu[dto.EnableReq](a.EnableSubscribe)) // 启用订阅

			sg.PUT("/disable/:id", casbin.Handler(a.ef), oplog.Record(oplog.LogOption{
				IncludeBody: true,
				Module:      a.modeNma,
				Action:      "停用订阅",
			}), hserver.NewHandlerFu[models.StringIdReq](a.DisableSubscribe)) // 停用订阅
		}

		// 死信队列
		dg := g.Group("/dead_letter")
		{
			dg.GET("", casbin.Handler(a.ef),
				hserver.NewHandlerFu[dto.GetDeadLetterSubscribeListReq](a.GetDeadLetterSubscribeList)) // 获取死信队列列表

			dg.PUT("/retry/:id", casbin.Handler(a.ef), oplog.Record(oplog.LogOption{
				IncludeBody: true,
				Module:      a.modeNma,
				Action:      "重试",
			}), hserver.NewHandlerFu[models.StringIdReq](a.DeadLetterQueueRetry)) // 死信队列重试
		}
	}
}

// AddEvent 新增事件
// @Summary 新增事件
// @Description 新增事件
// @Tags 事件
// @ID AddEvent
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req body dto.AddEventReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event [post]
func (a *EventService) AddEvent(ctx context.Context, req *dto.AddEventReq) *hserver.ResponseResult {
	err := a.as.Add(ctx, req)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res
}

// UpdateEvent 修改事件
// @Summary 修改事件
// @Description 修改事件
// @Tags 事件
// @ID UpdateEvent
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req body dto.UpdateEventReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event [put]
func (a *EventService) UpdateEvent(ctx context.Context, req *dto.UpdateEventReq) *hserver.ResponseResult {
	err := a.as.Update(ctx, req)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res
}

// UpdateEventStatus 修改事件
// @Summary 修改事件
// @Description 修改事件
// @Tags 事件
// @ID UpdateEventStatus
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req body dto.UpdateEventStatusReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event/update_status [put]
func (a *EventService) UpdateEventStatus(ctx context.Context, req *dto.UpdateEventStatusReq) *hserver.ResponseResult {
	err := a.as.UpdateStatus(ctx, req.Id, req.Status)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res
}

// GetEvent 获取事件
// @Summary 获取事件
// @Description 获取事件
// @Tags 事件
// @ID GetEvent
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req path models.StringIdReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{data=dto.EventModel} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event/:id [get]
func (a *EventService) GetEvent(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	re, err := a.as.GetDetails(ctx, req.Id)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res.WithData(re)
}

// GetEventList 获取事件
// @Summary 获取事件
// @Description 获取事件
// @Tags 事件
// @ID GetEventList
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req query dto.GetEventListReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{data=[]dto.EventModel} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event [get]
func (a *EventService) GetEventList(ctx context.Context, req *dto.GetEventListReq) *hserver.ResponseResult {
	re, err := a.as.GetList(ctx, req)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res.WithData(re)
}

// GetAllEventList 获取所有模版列表
// @Summary 获取所有模版列表
// @Description 获取所有模版列表
// @Tags 事件
// @ID GetAllEventList
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Success 200 {object} base_info.Success{data=[]dto.EventModel} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event/all [get]
func (a *EventService) GetAllEventList(ctx context.Context) *hserver.ResponseResult {
	rply, err := a.as.GetAllByStatusList(ctx, 2)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(herrors.TohError(err))
	}
	return res.WithData(rply)
}

// AddSubscribe 新增事件订阅
// @Summary 新增事件订阅
// @Description 新增事件订阅
// @Tags 事件订阅
// @ID AddSubscribe
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req body dto.AddSubscribeReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event/subscribe [post]
func (a *EventService) AddSubscribe(ctx context.Context, req *dto.AddSubscribeReq) *hserver.ResponseResult {
	err := a.sbs.Add(ctx, req)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res
}

// UpdateSubscribe 修改事件订阅
// @Summary 修改事件订阅
// @Description 修改事件订阅
// @Tags 事件订阅
// @ID UpdateSubscribe
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req body dto.UpdateSubscribeReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event/subscribe [put]
func (a *EventService) UpdateSubscribe(ctx context.Context, req *dto.UpdateSubscribeReq) *hserver.ResponseResult {
	err := a.sbs.Update(ctx, req)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res
}

// EnableSubscribe 启用订阅
// @Summary 启用订阅
// @Description 启用订阅
// @Tags 事件订阅
// @ID EnableSubscribe
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req body dto.EnableReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event/subscribe/enable [put]
func (a *EventService) EnableSubscribe(ctx context.Context, req *dto.EnableReq) *hserver.ResponseResult {
	err := a.sbs.Enable(ctx, req.Id, req.IgnoringHistory)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res
}

// DisableSubscribe 停止订阅
// @Summary 停止订阅
// @Description 停止订阅
// @Tags 事件订阅
// @ID DisableSubscribe
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req path models.StringIdReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event/subscribe/disable [put]
func (a *EventService) DisableSubscribe(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	err := a.sbs.Disable(ctx, req.Id)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res
}

// GetSubscribe 获取事件订阅
// @Summary 获取事件订阅
// @Description 获取事件订阅
// @Tags 事件订阅
// @ID GetSubscribe
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req path models.StringIdReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{data=dto.SubscribeModel} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event/subscribe/:id [get]
func (a *EventService) GetSubscribe(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	re, err := a.sbs.GetDetails(ctx, req.Id)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res.WithData(re)
}

// GetSubscribeList 获取事件订阅
// @Summary 获取事件订阅
// @Description 获取事件订阅
// @Tags 事件
// @ID GetSubscribeList
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req query dto.GetSubscribeListReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{data=[]dto.SubscribeModel} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event/subscribe [get]
func (a *EventService) GetSubscribeList(ctx context.Context, req *dto.GetSubscribeListReq) *hserver.ResponseResult {
	re, err := a.sbs.GetList(ctx, req)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res.WithData(re)
}

// GetDeadLetterSubscribeList 获取死信队列列表
// @Summary 获取死信队列列表
// @Description 获取死信队列列表
// @Tags 事件
// @ID GetDeadLetterSubscribeList
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req query dto.GetDeadLetterSubscribeListReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{data=[]dto.DeadLetterSubscribeModel} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event/dead_letter [get]
func (a *EventService) GetDeadLetterSubscribeList(ctx context.Context, req *dto.GetDeadLetterSubscribeListReq) *hserver.ResponseResult {
	re, err := a.ds.GetList(ctx, req)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res.WithData(re)
}

// DeadLetterQueueRetry 死信队列重试
// @Summary 死信队列重试
// @Description 死信队列重试
// @Tags 事件
// @ID DeadLetterQueueRetry
// @Accept application/json
// @Produce application/json
// @Param        Authorization  header  string                     true  "Bearer token"
// @Param req path models.StringIdReq true "属性说明请在对应model中查看"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/event/dead_letter/retry/:id [put]
func (a *EventService) DeadLetterQueueRetry(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	err := a.ds.Retry(ctx, req.Id)
	res := hserver.DefaultResponseResult()
	if herrors.HaveError(err) {
		res = res.WithError(herrors.TohError(err))
	}
	return res
}
