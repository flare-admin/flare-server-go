package data

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/support/systask/dto"
	"github.com/flare-admin/flare-server-go/framework/support/systask/model"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type ITaskRepo interface {
	baserepo.IBaseRepo[model.Task, string]
	// FindEnabledTasks 查询所有开启的任务
	FindEnabledTasks(ctx context.Context) ([]model.Task, error)
	SaveLog(log *model.TaskLog) error
	GetLogs(ctx context.Context, req dto.GetLogsReq) ([]model.TaskLog, int64, error)
}

type taskRepo struct {
	*baserepo.BaseRepo[model.Task, string]
}

func NewTaskRepo(data database.IDataBase) ITaskRepo {
	// 同步表
	tables := []interface{}{
		&model.Task{},
		&model.TaskLog{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync tables  error: %v", err)
	}
	return &taskRepo{
		BaseRepo: baserepo.NewBaseRepo[model.Task, string](data),
	}
}

// FindEnabledTasks 查询所有开启的任务
func (t taskRepo) FindEnabledTasks(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task
	err := t.Db(ctx).Where("status = ?", 1).Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// SaveLog 保存任务执行日志
func (t taskRepo) SaveLog(log *model.TaskLog) error {
	// 生成ID
	log.ID = t.GenStringId()
	// 保存日志
	return t.Db(context.Background()).Create(log).Error
}

// GetLogs 获取任务执行日志
func (t taskRepo) GetLogs(ctx context.Context, req dto.GetLogsReq) ([]model.TaskLog, int64, error) {
	var logs []model.TaskLog
	var total int64
	tx := t.Db(ctx).Model(&model.TaskLog{})

	// 构建查询条件
	if req.TaskID != "" {
		tx = tx.Where("task_id = ?", req.TaskID)
	}
	if req.Status > 0 {
		tx = tx.Where("status = ?", req.Status)
	}
	if req.StartTime > 0 {
		tx = tx.Where("start_time >= ?", req.StartTime)
	}
	if req.EndTime > 0 {
		tx = tx.Where("end_time <= ?", req.EndTime)
	}

	// 按时间倒序
	tx = tx.Order("start_time DESC")

	if err := tx.Count(&total).Error; err != nil {
		return nil, total, err
	}
	// 执行查询
	err := tx.Scopes(database.Operation(req.Current, req.Size)).Scan(&logs).Error
	if err != nil {
		return nil, total, err
	}

	return logs, total, nil
}
