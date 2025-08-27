package manager

import (
	"context"
	"errors"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/systask/data"
	"github.com/flare-admin/flare-server-go/framework/support/systask/model"
	"github.com/flare-admin/flare-server-go/framework/support/systask/service"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/robfig/cron/v3"
)

type TaskManager struct {
	cron     *cron.Cron
	tasks    map[string]*model.Task
	taskJobs map[string]cron.EntryID
	handlers map[string]service.TaskHandler
	mutex    sync.RWMutex
	repo     data.ITaskRepo
}

func NewTaskManager(repo data.ITaskRepo) service.ITaskManager {
	// 使用带秒的 cron 表达式，并设置时区
	c := cron.New(cron.WithSeconds(), cron.WithLocation(time.Local))
	return &TaskManager{
		cron:     c,
		tasks:    make(map[string]*model.Task),
		taskJobs: make(map[string]cron.EntryID),
		handlers: make(map[string]service.TaskHandler),
		repo:     repo,
	}
}

// Initialize 初始化并加载任务
func (tm *TaskManager) Initialize() error {
	// 从数据库加载已启用的任务
	tasks, err := tm.repo.FindEnabledTasks(context.Background())
	if err != nil {
		return fmt.Errorf("load enabled tasks error: %v", err)
	}

	hlog.Infof("Loading %d enabled tasks", len(tasks))

	for _, task := range tasks {
		if err := tm.StartTask(&task); err != nil {
			hlog.Errorf("start task error: %v, task: %+v", err, task)
			continue
		}
	}

	// 启动 cron
	tm.Start()
	hlog.Info("Task manager initialized and started")

	// 打印所有已加载的任务
	for _, task := range tm.tasks {
		entryID := tm.taskJobs[task.ID]
		entry := tm.cron.Entry(entryID)
		hlog.Infof("Loaded task: ID=%s, Name=%s, Cron=%s, Next Run=%s",
			task.ID, task.Name, task.Cron, entry.Next.Format("2006-01-02 15:04:05"))
	}

	return nil
}

// RegisterHandler 注册任务处理器
func (tm *TaskManager) RegisterHandler(name string, handler service.TaskHandler) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	tm.handlers[name] = handler
}

// executeTask 执行任务
func (tm *TaskManager) executeTask(task *model.Task, args map[string]string) error {
	handler, exists := tm.handlers[task.Handler]
	if !exists {
		return errors.New("handler not found: " + task.Handler)
	}

	// 打印执行时间，帮助调试
	hlog.Infof("Executing task: ID=%s, Name=%s, Time=%s",
		task.ID, task.Name, utils.GetTimeNow().Format("2006-01-02 15:04:05"))

	// 创建任务日志
	taskLog := &model.TaskLog{
		TaskID:    task.ID,
		StartTime: utils.GetDateUnix(),
	}

	// 使用defer确保无论如何都会记录日志
	defer func() {
		taskLog.EndTime = utils.GetDateUnix()
		taskLog.Duration = taskLog.EndTime - taskLog.StartTime
		if err := tm.repo.SaveLog(taskLog); err != nil {
			hlog.Errorf("save task log error: %v, log: %+v", err, taskLog)
		}
	}()

	// 捕获panic
	defer func() {
		if r := recover(); r != nil {
			stackTrace := string(debug.Stack())
			errMsg := fmt.Sprintf("panic: %v\nstack: %s", r, stackTrace)
			taskLog.Status = 2
			taskLog.Error = errMsg
			hlog.Errorf("task panic: %s", errMsg)
		}
	}()

	// 执行任务
	err := handler(args)
	if err != nil {
		taskLog.Status = 2
		taskLog.Error = err.Error()
		hlog.Errorf("task execute error: %v, task: %+v", err, task)
		return err
	}

	taskLog.Status = 1
	return nil
}

// ExecuteOnce 执行一次任务
func (tm *TaskManager) ExecuteOnce(task *model.Task) error {
	args := tm.parseArgs(task.Args)
	return tm.executeTask(task, args)
}

// StartTask 开启任务
func (tm *TaskManager) StartTask(task *model.Task) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	// 先检查任务是否已经存在，如果存在则先停止
	if _, exists := tm.taskJobs[task.ID]; exists {
		tm.cron.Remove(tm.taskJobs[task.ID])
		delete(tm.taskJobs, task.ID)
	}

	args := tm.parseArgs(task.Args)

	// 打印日志，帮助调试
	hlog.Infof("Adding task: ID=%s, Name=%s, Cron=%s", task.ID, task.Name, task.Cron)

	entryID, err := tm.cron.AddFunc(task.Cron, func() {
		_ = tm.executeTask(task, args)
	})

	if err != nil {
		return fmt.Errorf("add cron job error: %v", err)
	}

	// 打印已添加的任务信息
	entry := tm.cron.Entry(entryID)
	hlog.Infof("Task added successfully: ID=%s, Name=%s, Next Run=%s",
		task.ID, task.Name, entry.Next.Format("2006-01-02 15:04:05"))

	tm.taskJobs[task.ID] = entryID
	tm.tasks[task.ID] = task
	return nil
}

// StopTask 停止任务
func (tm *TaskManager) StopTask(taskID string) error {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()

	if entryID, exists := tm.taskJobs[taskID]; exists {
		tm.cron.Remove(entryID)
		delete(tm.taskJobs, taskID)
		delete(tm.tasks, taskID)
	}
	return nil
}

// ReloadTask 重新加载任务
func (tm *TaskManager) ReloadTask(task *model.Task) error {
	if err := tm.StopTask(task.ID); err != nil {
		return err
	}

	if task.Status == 1 {
		return tm.StartTask(task)
	}
	return nil
}

// Start 启动任务管理器
func (tm *TaskManager) Start() {
	tm.cron.Start()
}

// Stop 停止任务管理器
func (tm *TaskManager) Stop() {
	tm.cron.Stop()
}

// parseArgs 解析任务参数
func (tm *TaskManager) parseArgs(argsStr string) map[string]string {
	args := make(map[string]string)
	if argsStr == "" {
		return args
	}

	pairs := strings.Split(argsStr, " ")
	for _, pair := range pairs {
		if strings.HasPrefix(pair, "--") {
			kv := strings.SplitN(pair[2:], "=", 2)
			if len(kv) == 2 {
				args[kv[0]] = kv[1]
			}
		}
	}
	return args
}
