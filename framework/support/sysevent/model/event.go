package model

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
)

// Event ， 事件
type Event struct {
	database.BaseIntTime
	Id     string `gorm:"column:id;primary_key" json:"id"`
	Name   string `json:"name" gorm:"column:name;size:255;comment:事件名称"`
	Topic  string `json:"topic" gorm:"column:topic;size:150;primary_key;not null;comment:事件主题"`
	Dis    string `json:"dis" gorm:"column:dis;type:text;comment:内容"`
	Status int8   `gorm:"column:status;default:1;comment:公告状态 事件状态 1->新建, 2->启用，3->停用 ;NOT NULL" json:"status"`
}

func (Event) TableName() string {
	return "sys_events"
}

func (Event) GetPrimaryKey() string {
	return "id"
}

// GetEventCacheKey //缓存key
var GetEventCacheKey = func(topic string) string {
	return fmt.Sprintf("SysEventCache:%s", topic)
}
