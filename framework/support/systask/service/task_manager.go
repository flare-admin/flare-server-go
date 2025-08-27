package service

import (
	"github.com/flare-admin/flare-server-go/framework/support/systask/model"
)

// ITaskManager 任务管理器接口
type ITaskManager interface {
	// Initialize 初始化并加载任务
	Initialize() error
	// RegisterHandler 注册任务处理器
	RegisterHandler(name string, handler TaskHandler)
	// ExecuteOnce 执行一次任务
	ExecuteOnce(task *model.Task) error
	// StartTask 开启任务
	StartTask(task *model.Task) error
	// StopTask 停止任务
	StopTask(taskID string) error
	// ReloadTask 重新加载任务
	ReloadTask(task *model.Task) error
	// Start 启动任务管理器
	Start()
	// Stop 停止任务管理器
	Stop()
}

// TaskHandler 任务处理器
type TaskHandler func(args map[string]string) error
