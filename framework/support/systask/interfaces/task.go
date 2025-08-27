package interfaces

import (
	"context"

	"github.com/cloudwego/hertz/pkg/route"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver"
	_ "github.com/flare-admin/flare-server-go/framework/pkg/hserver/base_info"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/casbin"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/jwt"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/oplog"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/pkg/token"
	"github.com/flare-admin/flare-server-go/framework/support/systask/dto"
	"github.com/flare-admin/flare-server-go/framework/support/systask/service"
)

type TaskService struct {
	ts         service.ITaskService
	ef         *casbin.Enforcer
	moduleName string
}

func NewTaskService(ts service.ITaskService, ef *casbin.Enforcer) *TaskService {
	return &TaskService{
		ts:         ts,
		ef:         ef,
		moduleName: "定时任务",
	}
}

func (s *TaskService) RegisterRouter(r *route.RouterGroup, t token.IToken) {
	g := r.Group("/v1/tasks", jwt.Handler(t))
	{
		g.POST("", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.moduleName,
			Action:      "创建任务",
		}), hserver.NewHandlerFu[dto.CreateTask](s.CreateTask))

		g.PUT("", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.moduleName,
			Action:      "更新任务",
		}), hserver.NewHandlerFu[dto.UpdateTask](s.UpdateTask))

		g.DELETE("/:id", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.moduleName,
			Action:      "删除任务",
		}), hserver.NewHandlerFu[models.StringIdReq](s.DeleteTask))

		g.PUT("/enable/:id", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.moduleName,
			Action:      "启用任务",
		}), hserver.NewHandlerFu[models.StringIdReq](s.EnableTask))

		g.PUT("/disable/:id", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.moduleName,
			Action:      "停用任务",
		}), hserver.NewHandlerFu[models.StringIdReq](s.DisableTask))

		g.GET("", casbin.Handler(s.ef), hserver.NewHandlerFu[dto.TaskListReq](s.GetTaskList))
		g.GET("/logs", casbin.Handler(s.ef), hserver.NewHandlerFu[dto.GetLogsReq](s.GetTaskLogs))

		g.POST("/execute/:id", casbin.Handler(s.ef), oplog.Record(oplog.LogOption{
			IncludeBody: true,
			Module:      s.moduleName,
			Action:      "执行任务",
		}), hserver.NewHandlerFu[models.StringIdReq](s.ExecuteOnce))
	}
}

// CreateTask 创建任务
// @Summary 创建任务
// @Description 创建定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param req body dto.CreateTask true "创建任务参数"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/tasks [post]
func (s *TaskService) CreateTask(ctx context.Context, req *dto.CreateTask) *hserver.ResponseResult {
	err := s.ts.CreateTask(ctx, req)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}

// UpdateTask 更新任务
// @Summary 更新任务
// @Description 更新定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param req body dto.UpdateTask true "更新任务参数"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/tasks [put]
func (s *TaskService) UpdateTask(ctx context.Context, req *dto.UpdateTask) *hserver.ResponseResult {
	err := s.ts.UpdateTask(ctx, req)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}

// DeleteTask 删除任务
// @Summary 删除任务
// @Description 删除定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "任务ID"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/tasks/{id} [delete]
func (s *TaskService) DeleteTask(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	err := s.ts.DeleteTask(ctx, req.Id)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}

// EnableTask 启用任务
// @Summary 启用任务
// @Description 启用定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "任务ID"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/tasks/enable/{id} [put]
func (s *TaskService) EnableTask(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	err := s.ts.EnableTask(ctx, req.Id)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}

// DisableTask 停用任务
// @Summary 停用任务
// @Description 停用定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "任务ID"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/tasks/disable/{id} [put]
func (s *TaskService) DisableTask(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	err := s.ts.DisableTask(ctx, req.Id)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}

// GetTaskList 获取任务列表
// @Summary 获取任务列表
// @Description 分页获取定时任务列表
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Param groupName query string false "分组名称"
// @Param name query string false "任务名称"
// @Param status query int false "状态"
// @Success 200 {object} base_info.Success{data=[]dto.Task} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"{data=db.Response[dto.Task]}
// @Router /api/admin/v1/tasks [get]
func (s *TaskService) GetTaskList(ctx context.Context, req *dto.TaskListReq) *hserver.ResponseResult {
	data, err := s.ts.GetTaskList(ctx, req)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res.WithData(data)
}

// GetTaskLogs 获取任务日志
// @Summary 获取任务日志
// @Description 分页获取定时任务执行日志
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param page query int false "页码"
// @Param size query int false "每页数量"
// @Param taskId query string false "任务ID"
// @Param status query int false "状态"
// @Param startTime query int64 false "开始时间"
// @Param endTime query int64 false "结束时间"
// @Success 200 {object} base_info.Success{data=[]dto.TaskLog} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/tasks/logs [get]
func (s *TaskService) GetTaskLogs(ctx context.Context, req *dto.GetLogsReq) *hserver.ResponseResult {
	data, err := s.ts.GetTaskLogs(ctx, req)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res.WithData(data)
}

// ExecuteOnce 执行一次任务
// @Summary 执行一次任务
// @Description 手动执行一次定时任务
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path string true "任务ID"
// @Success 200 {object} base_info.Success{} "data的属性说明请在对应model中查看"
// @Failure 400 {object} base_info.Swagger400Resp "code为400 参数输入错误"
// @Failure 401 {object} base_info.Swagger401Resp "code为401 token未带上"
// @Failure 500 {object} base_info.Swagger500Resp "code为500 服务端内部错误"
// @Router /api/admin/v1/tasks/execute/{id} [post]
func (s *TaskService) ExecuteOnce(ctx context.Context, req *models.StringIdReq) *hserver.ResponseResult {
	err := s.ts.ExecuteOnce(ctx, req.Id)
	res := hserver.DefaultResponseResult()
	if err != nil {
		res = res.WithError(err)
	}
	return res
}
