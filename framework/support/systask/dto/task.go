package dto

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
)

// TaskListReq 任务列表请求参数
type TaskListReq struct {
	db_query.Page
	GroupName string `query:"groupName"` // 分组名称
	Name      string `query:"name"`      // 任务名称
	Status    int    `query:"status"`    // 状态 -1:全部 0:停止 1:运行
}

// GetLogsReq 获取日志请求参数
type GetLogsReq struct {
	db_query.Page
	TaskID    string `query:"taskId"`    // 任务ID
	Status    int    `query:"status"`    // 执行状态 -1:全部 0:失败 1:成功
	StartTime int64  `query:"startTime"` // 开始时间
	EndTime   int64  `query:"endTime"`   // 结束时间
}

// Task 定时任务
type Task struct {
	database.BaseIntTime
	ID        string `json:"id"`
	GroupName string `json:"groupName"` // 分组ID
	Name      string `json:"name"`      // 任务名称
	Handler   string `json:"handler"`   // 任务处理器名称
	Args      string `json:"args"`      // 任务参数，使用--key=value格式
	Cron      string `json:"cron"`      // cron表达式
	Status    int    `json:"status"`    // 状态 0:停止 1:运行
	Remark    string `json:"remark"`    // 备注
}

// CreateTask 创建定时任务
type CreateTask struct {
	GroupName string `json:"groupName"` // 分组ID
	Name      string `json:"name"`      // 任务名称
	Handler   string `json:"handler"`   // 任务处理器名称
	Args      string `json:"args"`      // 任务参数，使用--key=value格式
	Cron      string `json:"cron"`      // cron表达式
	Status    int    `json:"status"`    // 状态 2:停止 1:运行
	Remark    string `json:"remark"`    // 备注
}

// UpdateTask 修改定时任务
type UpdateTask struct {
	ID        string `json:"id"`
	GroupName string `json:"groupName"` // 分组ID
	Name      string `json:"name"`      // 任务名称
	Handler   string `json:"handler"`   // 任务处理器名称
	Args      string `json:"args"`      // 任务参数，使用--key=value格式
	Cron      string `json:"cron"`      // cron表达式
	Status    int    `json:"status"`    // 状态 2:停止 1:运行
	Remark    string `json:"remark"`    // 备注
}

type TaskLog struct {
	database.BaseIntTime
	ID        string `json:"id"`
	TaskID    string `json:"task_id"`
	StartTime int64  `json:"start_time"`
	EndTime   int64  `json:"end_time"`
	Duration  int64  `json:"duration"`
	Status    int    `json:"status"`
	Output    string `json:"output"`
	Error     string `json:"error"`
}
