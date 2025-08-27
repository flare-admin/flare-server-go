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

type subscribeRepo struct {
	*baserepo.BaseRepo[model.Subscribe, string]
}

func NewSubscribeRepo(data database.IDataBase) biz.ISubscribeRepo {
	// 同步表
	tables := []interface{}{
		&model.Subscribe{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync tables  error: %v", err)
	}
	return &subscribeRepo{
		BaseRepo: baserepo.NewBaseRepo[model.Subscribe, string](data, model.Subscribe{}),
	}
}
func (s subscribeRepo) GetByStatus(ctx context.Context, status int8) ([]*model.Subscribe, error) {
	var result []*model.Subscribe
	err := s.Db(ctx).Model(&model.Subscribe{}).Where("status = ?", status).Find(&result).Error
	return result, err
}
func (s subscribeRepo) GetAll(ctx context.Context) ([]*model.Subscribe, error) {
	var result []*model.Subscribe
	err := s.Db(ctx).Model(&model.Subscribe{}).Find(&result).Error
	return result, err
}
func (s subscribeRepo) GetByTopicAndGroup(ctx context.Context, topic, group string) (*model.Subscribe, error) {
	var result model.Subscribe
	err := s.Db(ctx).Model(&model.Subscribe{}).Where("deleted_at = 0 and topic = ? and group_name = ?", topic, group).First(&result).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 或返回自定义 NotFound 错误
		}
		return nil, err
	}
	return &result, nil
}

type subscribeParameterRepo struct {
	*baserepo.BaseRepo[model.SubscribeParameter, string]
}

func NewSubscribeParameterRepo(data database.IDataBase) biz.ISubscribeParameterRepo {
	// 同步表
	tables := []interface{}{
		&model.SubscribeParameter{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync Parameter tables to mysql error: %v", err)
	}
	return &subscribeParameterRepo{
		BaseRepo: baserepo.NewBaseRepo[model.SubscribeParameter, string](data, model.SubscribeParameter{}),
	}
}
func (s subscribeParameterRepo) GetBySubscribe(ctx context.Context, subscribe string) ([]*model.SubscribeParameter, error) {
	var result []*model.SubscribeParameter
	err := s.Db(ctx).Model(&model.SubscribeParameter{}).Where("subscribe_id = ?", subscribe).Find(&result).Error
	return result, err
}

func (s subscribeParameterRepo) DeleteBySubscribe(ctx context.Context, subscribe string) error {
	return s.Db(ctx).Unscoped().Where("subscribe_id = ?", subscribe).Delete(&model.SubscribeParameter{}).Error
}

func (s subscribeParameterRepo) DeleteByIds(ctx context.Context, ids ...string) error {
	return s.Db(ctx).Unscoped().Where("id in (?)", ids).Delete(&model.SubscribeParameter{}).Error
}
