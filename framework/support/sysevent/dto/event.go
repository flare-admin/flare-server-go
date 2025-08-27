package dto

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/model"
)

// EventModel 事件模型
type EventModel struct {
	Id        string `json:"id"`        // ID
	Name      string `json:"name"`      // 名称
	Topic     string `json:"topic"`     // 主题
	Dis       string `json:"dis"`       // 描述
	Status    int8   `json:"status"`    // 状态
	CreatedAt int64  `json:"createdAt"` // 创建时间
	UpdatedAt int64  `json:"updatedAt"` // 更新时间
}

// AddEventReq 添加事件请求
type AddEventReq struct {
	Name  string `json:"name" query:"name"`   // 名称
	Topic string `json:"topic" query:"topic"` // 主题
	Dis   string `json:"dis" query:"dis"`     // 描述
}

// UpdateEventReq 更新事件请求
type UpdateEventReq struct {
	Id    string `json:"id" query:"id"`       // ID
	Name  string `json:"name" query:"name"`   // 名称
	Topic string `json:"topic" query:"topic"` // 主题
	Dis   string `json:"dis" query:"dis"`     // 描述
}

// GetEventListReq 获取事件列表请求
type GetEventListReq struct {
	db_query.Page
	Topic  string `json:"topic" query:"topic"`   // 主题
	Status int32  `json:"status" query:"status"` // 状态
}

// 更新事件状态
type UpdateEventStatusReq struct {
	Id     string `json:"id,omitempty"`
	Status int32  `json:"status,omitempty"` // 事件状态 1->未达标, 2->完成，3->奖励发放完成
}

// 更新事件状态
type EnableReq struct {
	Id              string `json:"id,omitempty"`
	IgnoringHistory int32  `json:"ignoring_history,omitempty"` // 是否忽略历史 1->或略, 2->不忽略
}

// EventToDto 将model转换为dto
func EventToDto(event *model.Event) *EventModel {
	if event == nil {
		return nil
	}
	return &EventModel{
		Id:     event.Id,
		Name:   event.Name,
		Topic:  event.Topic,
		Dis:    event.Dis,
		Status: event.Status,
	}
}
