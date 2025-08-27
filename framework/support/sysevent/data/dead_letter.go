package data

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/biz"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/model"
)

type deadLetterRepo struct {
	*baserepo.BaseRepo[model.DeadLetterSubscribe, string]
}

func NewDeadLetterSubscribeRepo(data database.IDataBase) biz.IDeadLetterSubscribeRepo {
	// 同步表
	tables := []interface{}{
		&model.DeadLetterSubscribe{},
	}
	if err := data.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync tables  error: %v", err)
	}
	return &deadLetterRepo{
		BaseRepo: baserepo.NewBaseRepo[model.DeadLetterSubscribe, string](data, model.DeadLetterSubscribe{}),
	}
}
func (d deadLetterRepo) GetBy(ctx context.Context, topic, channel, messageId string) (*model.DeadLetterSubscribe, error) {
	var data model.DeadLetterSubscribe
	err := d.Db(ctx).Model(&model.DeadLetterSubscribe{}).Where("topic = ? AND channel = ? AND msg_id = ?", topic, channel, messageId).First(&data).Error
	return &data, err
}
func (d deadLetterRepo) UpdateStatus(ctx context.Context, id string, status int8) error {
	return d.Db(ctx).Model(&model.DeadLetterSubscribe{}).Where("id = ?", id).Update("status", status).Update("updated_at", utils.GetDateUnix()).Error
}
