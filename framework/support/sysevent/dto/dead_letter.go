package dto

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/model"
)

type DeadLetterSubscribeModel struct {
	database.BaseIntTime
	Id     string `json:"id,omitempty"`     // ID
	Name   string `json:"name,omitempty"`   // 订阅名称
	Topic  string `json:"topic,omitempty"`  // 订阅的主题
	Group  string `json:"group,omitempty"`  // 订阅的组
	Status int32  `json:"status,omitempty"` // 状态
}

func DeadLetterToDto(vo *model.DeadLetterSubscribe) *DeadLetterSubscribeModel {
	return &DeadLetterSubscribeModel{
		Id:          vo.Id,
		Name:        vo.Name,
		Topic:       vo.Topic,
		Group:       vo.Channel,
		Status:      int32(vo.Status),
		BaseIntTime: vo.BaseIntTime,
	}
}

type GetDeadLetterSubscribeListReq struct {
	db_query.Page
	Topic string `json:"topic,omitempty" query:"topic"` // 主题
	Name  string `json:"name,omitempty" query:"name"`   // 订阅名称
	Group string `json:"group,omitempty" query:"group"` // 订阅的组
}
