package service

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/support/systask/dto"
)

type ITaskService interface {
	// CreateTask 创建任务
	CreateTask(ctx context.Context, req *dto.CreateTask) herrors.Herr
	// UpdateTask 更新任务
	UpdateTask(ctx context.Context, req *dto.UpdateTask) herrors.Herr
	// DeleteTask 删除任务
	DeleteTask(ctx context.Context, taskID string) herrors.Herr
	// EnableTask 启用任务
	EnableTask(ctx context.Context, taskID string) herrors.Herr
	// DisableTask 停用任务
	DisableTask(ctx context.Context, taskID string) herrors.Herr
	// GetTaskList 分页查询任务列表
	GetTaskList(ctx context.Context, req *dto.TaskListReq) (*models.PageRes[dto.Task], herrors.Herr)
	// GetTaskLogs 查询任务日志
	GetTaskLogs(ctx context.Context, req *dto.GetLogsReq) (*models.PageRes[dto.TaskLog], herrors.Herr)
	// ExecuteOnce 执行一次任务
	ExecuteOnce(ctx context.Context, taskID string) herrors.Herr
}
