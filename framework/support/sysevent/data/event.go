package data

import (
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/biz"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/model"
	"gorm.io/gorm"
)

type eventRepo struct {
	*baserepo.BaseRepo[model.Event, string]
}

func NewEventRepo(data database.IDataBase) biz.IEventRepo {
	// 同步表
	tables := []interface{}{
		&model.Event{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync tables  error: %v", err)
	}
	return &eventRepo{
		BaseRepo: baserepo.NewBaseRepo[model.Event, string](data, model.Event{}),
	}
}

func (e eventRepo) GetByTopic(ctx context.Context, topic string) (*model.Event, error) {
	var event model.Event
	err := e.Db(ctx).Where("topic = ?", topic).First(&event).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}
func (e eventRepo) GetAllByStatusList(ctx context.Context, status int32) ([]*model.Event, error) {
	var events []*model.Event
	err := e.Db(ctx).Where("status = ? and deleted_at = 0", status).Find(&events).Error
	return events, err
}
