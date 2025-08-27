package biz

import (
	"context"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/snowflake_id"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/models"
	"github.com/flare-admin/flare-server-go/framework/support/systask/data"
	"github.com/flare-admin/flare-server-go/framework/support/systask/dto"
	"github.com/flare-admin/flare-server-go/framework/support/systask/model"
	service2 "github.com/flare-admin/flare-server-go/framework/support/systask/service"

	"github.com/robfig/cron/v3"
)

type TaskUseCase struct {
	repo        data.ITaskRepo
	taskManager service2.ITaskManager
	ig          snowflake_id.IIdGenerate
}

func NewTaskBiz(repo data.ITaskRepo, taskManager service2.ITaskManager, ig snowflake_id.IIdGenerate) service2.ITaskService {
	return &TaskUseCase{
		repo:        repo,
		taskManager: taskManager,
		ig:          ig,
	}
}

// validateCron 校验cron表达式
func (uc *TaskUseCase) validateCron(cronExpr string) error {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	_, err := parser.Parse(cronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %v", err)
	}
	return nil
}

// CreateTask 创建任务
func (uc *TaskUseCase) CreateTask(ctx context.Context, req *dto.CreateTask) herrors.Herr {
	// 校验cron表达式
	if err := uc.validateCron(req.Cron); err != nil {
		return herrors.NewServerHError(err)
	}

	task := &model.Task{
		ID:        uc.ig.GenStringId(),
		GroupName: req.GroupName,
		Name:      req.Name,
		Handler:   req.Handler,
		Args:      req.Args,
		Cron:      req.Cron,
		Status:    req.Status,
		Remark:    req.Remark,
	}

	if _, err := uc.repo.Add(ctx, task); err != nil {
		return herrors.CreateFail(fmt.Errorf("create task error: %v", err))
	}

	if task.Status == 1 {
		if err := uc.taskManager.StartTask(task); err != nil {
			return herrors.CreateFail(fmt.Errorf("start task error: %v", err))
		}
	}
	return nil
}

// UpdateTask 更新任务
func (uc *TaskUseCase) UpdateTask(ctx context.Context, req *dto.UpdateTask) herrors.Herr {
	// 校验cron表达式
	if err := uc.validateCron(req.Cron); err != nil {
		return herrors.NewParameterHError(err)
	}

	task := &model.Task{
		GroupName: req.GroupName,
		Name:      req.Name,
		Handler:   req.Handler,
		Args:      req.Args,
		Cron:      req.Cron,
		Status:    req.Status,
		Remark:    req.Remark,
	}

	if err := uc.repo.EditById(ctx, req.ID, task); err != nil {
		return herrors.UpdateFail(fmt.Errorf("update task error: %v", err))
	}

	if err := uc.taskManager.ReloadTask(task); err != nil {
		return herrors.UpdateFail(fmt.Errorf("reload task error: %v", err))
	}
	return nil
}

// DeleteTask 删除任务
func (uc *TaskUseCase) DeleteTask(ctx context.Context, taskID string) herrors.Herr {
	// 先停止任务
	if err := uc.taskManager.StopTask(taskID); err != nil {
		return herrors.DeleteFail(fmt.Errorf("stop task error: %v", err))
	}
	// 删除任务
	if err := uc.repo.DelById(ctx, taskID); err != nil {
		return herrors.DeleteFail(fmt.Errorf("delete task error: %v", err))
	}
	return nil
}

// EnableTask 启用任务
func (uc *TaskUseCase) EnableTask(ctx context.Context, taskID string) herrors.Herr {
	task, err := uc.repo.FindById(ctx, taskID)
	if err != nil {
		return herrors.NotFound(fmt.Errorf("task not found: %v", err))
	}

	task.Status = 1
	if err := uc.repo.EditById(ctx, taskID, task); err != nil {
		return herrors.UpdateFail(fmt.Errorf("update task status error: %v", err))
	}

	if err := uc.taskManager.StartTask(task); err != nil {
		return herrors.UpdateFail(fmt.Errorf("start task error: %v", err))
	}
	return nil
}

// DisableTask 停用任务
func (uc *TaskUseCase) DisableTask(ctx context.Context, taskID string) herrors.Herr {
	task, err := uc.repo.FindById(ctx, taskID)
	if err != nil {
		return herrors.NotFound(fmt.Errorf("task not found: %v", err))
	}

	task.Status = 2
	if err := uc.repo.EditById(ctx, taskID, task); err != nil {
		return herrors.UpdateFail(fmt.Errorf("update task status error: %v", err))
	}

	if err := uc.taskManager.StopTask(taskID); err != nil {
		return herrors.UpdateFail(fmt.Errorf("stop task error: %v", err))
	}
	return nil
}

// GetTaskList 分页查询任务列表
func (uc *TaskUseCase) GetTaskList(ctx context.Context, req *dto.TaskListReq) (*models.PageRes[dto.Task], herrors.Herr) {
	// 构建查询条件
	qb := db_query.NewQueryBuilder()
	if req.GroupName != "" {
		qb.Where("group_name", db_query.Like, "%"+req.GroupName+"%")
	}
	if req.Status != 0 {
		qb.Where("status", db_query.Eq, req.Status)
	}
	if req.Name != "" {
		qb.Where("name", db_query.Like, "%"+req.Name+"%")
	}

	qb.WithPage(&req.Page)

	// 查询总数
	total, err := uc.repo.Count(ctx, qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 查询数据
	tasks, err := uc.repo.Find(ctx, qb)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	var taskDTOs []*dto.Task
	for _, task := range tasks {
		taskDTOs = append(taskDTOs, &dto.Task{
			BaseIntTime: task.BaseIntTime,
			ID:          task.ID,
			GroupName:   task.GroupName,
			Name:        task.Name,
			Handler:     task.Handler,
			Args:        task.Args,
			Cron:        task.Cron,
			Status:      task.Status,
			Remark:      task.Remark,
		})
	}

	return &models.PageRes[dto.Task]{
		List:  taskDTOs,
		Total: total,
	}, nil
}

// GetTaskLogs 查询任务日志
func (uc *TaskUseCase) GetTaskLogs(ctx context.Context, req *dto.GetLogsReq) (*models.PageRes[dto.TaskLog], herrors.Herr) {
	logs, total, err := uc.repo.GetLogs(ctx, *req)
	if err != nil {
		return nil, herrors.QueryFail(err)
	}

	// 转换为DTO
	var logDTOs []*dto.TaskLog
	for _, log := range logs {
		logDTOs = append(logDTOs, &dto.TaskLog{
			ID:        log.ID,
			TaskID:    log.TaskID,
			StartTime: log.StartTime,
			EndTime:   log.EndTime,
			Duration:  log.Duration,
			Status:    log.Status,
			Output:    log.Output,
			Error:     log.Error,
		})
	}

	return &models.PageRes[dto.TaskLog]{
		List:  logDTOs,
		Total: total,
	}, nil
}

// ExecuteOnce 执行一次任务
func (uc *TaskUseCase) ExecuteOnce(ctx context.Context, taskID string) herrors.Herr {
	task, err := uc.repo.FindById(ctx, taskID)
	if err != nil {
		return herrors.NotFound(fmt.Errorf("task not found: %v", err))
	}

	if err := uc.taskManager.ExecuteOnce(task); err != nil {
		return herrors.OperateFail(fmt.Errorf("execute task error: %v", err))
	}
	return nil
}
