package dto

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/db_query"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/model"
)

type SubscribeParameterModel struct {
	Id          string `json:"id,omitempty"`           // ID
	Key         string `json:"key,omitempty"`          // 属性的key
	Value       string `json:"value,omitempty"`        // 参数值
	Dis         string `json:"dis,omitempty"`          // 参数描述
	DataType    string `json:"data_type,omitempty"`    // 属性类型（如string, int, float, boolean, date, select, multiselect, switch等）
	SubscribeId string `json:"subscribe_id,omitempty"` // 订阅事件的id
}

type AddSubscribeReq struct {
	Name      string                     `json:"name,omitempty" query:"name"`           // 订阅名称
	Topic     string                     `json:"topic,omitempty" query:"topic"`         // 订阅的主题
	Group     string                     `json:"group,omitempty" query:"group"`         // 订阅的组
	Dis       string                     `json:"dis,omitempty" query:"dis"`             // 描述
	Parameter []*SubscribeParameterModel `json:"parameter,omitempty" query:"parameter"` // 参数
}

// 修改事件
type UpdateSubscribeReq struct {
	Id        string                     `json:"id,omitempty" query:"id"`               // ID
	Name      string                     `json:"name,omitempty" query:"name"`           // 订阅名称
	Topic     string                     `json:"topic,omitempty" query:"topic"`         // 订阅的主题
	Group     string                     `json:"group,omitempty" query:"group"`         // 订阅的组
	Dis       string                     `json:"dis,omitempty" query:"dis"`             // 描述
	Parameter []*SubscribeParameterModel `json:"parameter,omitempty" query:"parameter"` // 参数
}

// 获取事件
// -------------------------------------------
type GetSubscribeListReq struct {
	db_query.Page
	Topic  string `json:"topic,omitempty" query:"topic"`   // 主题
	Status int32  `json:"status,omitempty" query:"status"` // 状态
	Name   string `json:"name,omitempty" query:"name"`     // 订阅名称
	Group  string `json:"group,omitempty" query:"group"`   // 订阅的组
}

type SubscribeModel struct {
	database.BaseIntTime
	Id        string                     `json:"id,omitempty"`        // ID
	Name      string                     `json:"name,omitempty"`      // 订阅名称
	Topic     string                     `json:"topic,omitempty"`     // 订阅的主题
	Group     string                     `json:"group,omitempty"`     // 订阅的组
	Dis       string                     `json:"dis,omitempty"`       // 描述
	Status    int32                      `json:"status,omitempty"`    // 状态 1->新建, 2->启用，3->停用
	Start     int64                      `json:"start,omitempty"`     // 订阅开始时间
	End       int64                      `json:"end,omitempty"`       // 订阅结束时间
	Parameter []*SubscribeParameterModel `json:"parameter,omitempty"` // 参数
}

// SubscribeToDto 将Subscribe模型转换为DTO
func SubscribeToDto(subscribe *model.Subscribe) *SubscribeModel {
	if subscribe == nil {
		return nil
	}
	return &SubscribeModel{
		BaseIntTime: subscribe.BaseIntTime,
		Id:          subscribe.Id,
		Name:        subscribe.Name,
		Topic:       subscribe.Topic,
		Group:       subscribe.Group,
		Dis:         subscribe.Dis,
		Status:      int32(subscribe.Status),
		Start:       subscribe.Start,
		End:         subscribe.End,
		Parameter:   nil, // 需要单独处理Parameter
	}
}

// SubscribeParameterToDto 将SubscribeParameter模型转换为DTO
func SubscribeParameterToDto(param *model.SubscribeParameter) *SubscribeParameterModel {
	if param == nil {
		return nil
	}
	return &SubscribeParameterModel{
		Id:          param.Id,
		Key:         param.Key,
		Value:       param.Value,
		Dis:         param.Dis,
		DataType:    param.DataType,
		SubscribeId: param.SubscribeId,
	}
}

// SubscribeParametersToDtos 将SubscribeParameter列表转换为DTO列表
func SubscribeParametersToDtos(params []*model.SubscribeParameter) []*SubscribeParameterModel {
	if params == nil {
		return nil
	}
	dtos := make([]*SubscribeParameterModel, 0, len(params))
	for _, param := range params {
		if dto := SubscribeParameterToDto(param); dto != nil {
			dtos = append(dtos, dto)
		}
	}
	return dtos
}

// SubscribeToDtoWithParams 将Subscribe模型及其参数转换为完整的DTO
func SubscribeToDtoWithParams(subscribe *model.Subscribe, params []*model.SubscribeParameter) *SubscribeModel {
	if subscribe == nil {
		return nil
	}
	dto := SubscribeToDto(subscribe)
	if dto != nil {
		dto.Parameter = SubscribeParametersToDtos(params)
	}
	return dto
}

// EnableSubscribeReq 启用订阅请求
type EnableSubscribeReq struct {
	Id              string `json:"id" query:"id"`                           // ID
	IgnoringHistory int32  `json:"ignoringHistory" query:"ignoringHistory"` // 是否忽略历史消息 0->否, 1->是
}
