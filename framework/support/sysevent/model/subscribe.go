package model

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"strconv"
)

// Subscribe ， 订阅事件
type Subscribe struct {
	database.BaseIntTime
	Id        string `gorm:"column:id;primary_key" json:"id"`
	Name      string `json:"name" gorm:"column:name;size:255;comment:订阅名称"`
	Topic     string `json:"topic" gorm:"column:topic;size:150;not null;comment:事件主题"`
	Group     string `json:"group" gorm:"column:group_name;size:150;not null;comment:分组名称"`
	Dis       string `json:"dis" gorm:"column:dis;type:text;comment:内容"`
	Constants string `json:"constants" gorm:"column:constants;type:text;comment:自定义的订阅常量json"`
	Start     int64  `json:"start" gorm:"column:start;not null;default:0;comment:开始时间"`
	End       int64  `json:"end" gorm:"column:end;not null;default:0;comment:结束时间"`
	Status    int8   `gorm:"column:status;default:1;comment:公告状态 事件状态 1->新建, 2->启用，3->停用 ;NOT NULL" json:"status"`
	TenantID  string `json:"tenantId" gorm:"column:tenant_id;size:255;comment:租户ID"` // 租户ID
}

func (Subscribe) TableName() string {
	return "sys_events_subscribes"
}

func (Subscribe) GetPrimaryKey() string {
	return "id"
}

// SubscribeParameter ， 订阅过程参数
type SubscribeParameter struct {
	database.BaseIntTime
	Id          string `gorm:"column:id;primary_key" json:"id"`
	SubscribeId string `gorm:"column:subscribe_id;comment:订阅的id" json:"subscribe_id"`
	Key         string `json:"name" gorm:"column:name;size:255;not null;comment:参数的key"`
	Value       string `json:"value" gorm:"column:value;size:150;comment:参数的值"`
	Dis         string `json:"dis" gorm:"column:dis;size:255;comment:参数描述"`
	DataType    string `gorm:"size:50;column:data_type;comment:数据类型" json:"data_type"` // string, int, float, boolean, date, select, multiselect, switch
}

func (SubscribeParameter) TableName() string {
	return "sys_events_subscribes_parameters"
}

func (SubscribeParameter) GetPrimaryKey() string {
	return "id"
}

// GetSubordinatesParameterCacheKey //缓存key
var GetSubordinatesParameterCacheKey = func(topic, group, tenantID string) string {
	return fmt.Sprintf("SubordinatesParameterCache:%s_%s_%s", topic, group, tenantID)
}

func ConvertParametersToMap(params []*SubscribeParameter) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for _, param := range params {
		var convertedValue interface{}
		var err error
		// 根据 DataType 转换 Value
		switch param.DataType {
		case "integer":
			convertedValue, err = strconv.Atoi(param.Value)
		case "number", "money":
			convertedValue, err = strconv.ParseFloat(param.Value, 64)
		case "boolean":
			convertedValue, err = strconv.ParseBool(param.Value)
		case "datetime", "date", "time":
			convertedValue, err = strconv.Atoi(param.Value)
		case "switch":
			convertedValue, err = strconv.Atoi(param.Value)
		case "wallet":
			convertedValue, err = utils.StringToInt8Slice(param.Value)
		default:
			// 默认为 string
			convertedValue = param.Value
		}
		// 如果有转换错误，返回错误
		if err != nil {
			return nil, fmt.Errorf("failed to convert value for key %s: %v", param.Key, err)
		}
		// 将转换后的值放入结果 map 中
		result[param.Key] = convertedValue
	}

	return result, nil
}
