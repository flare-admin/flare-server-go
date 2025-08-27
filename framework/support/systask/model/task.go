package model

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
)

// Task 定时任务
type Task struct {
	database.BaseIntTime
	ID        string `gorm:"primarykey" json:"id"`
	GroupName string `gorm:"not null" json:"groupName"`        // 分组ID
	Name      string `gorm:"size:50;not null" json:"name"`     // 任务名称
	Handler   string `gorm:"size:100;not null" json:"handler"` // 任务处理器名称
	Args      string `gorm:"size:500" json:"args"`             // 任务参数，使用--key=value格式
	Cron      string `gorm:"size:50;not null" json:"cron"`     // cron表达式
	Status    int    `gorm:"default:2" json:"status"`          // 状态 2:停止 1:运行
	Remark    string `gorm:"size:200" json:"remark"`           // 备注
}

func (Task) TableName() string {
	return "sys_task"
}

func (Task) GetPrimaryKey() string {
	return "id"
}

// TaskLog 任务执行日志
type TaskLog struct {
	ID        string `gorm:"primarykey" json:"id"`
	TaskID    string `gorm:"not null" json:"task_id"` // 任务ID
	Output    string `gorm:"type:text" json:"output"` // 执行输出
	Error     string `gorm:"type:text" json:"error"`  // 错误信息
	Status    int    `gorm:"default:2" json:"status"` // 状态 2:失败 1:成功
	StartTime int64  `json:"start_time"`              // 开始时间
	EndTime   int64  `json:"end_time"`                // 结束时间
	Duration  int64  `json:"duration"`                // 执行时长
}

func (TaskLog) TableName() string {
	return "sys_task_logs"
}

func (TaskLog) GetPrimaryKey() string {
	return "id"
}
