package repository

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/infrastructure/entity"
)

// IConfigGroupRepository 配置分类接口
type IConfigGroupRepository interface {
	baserepo.IBaseRepo[entity.ConfigGroup, string]
	FindByCode(ctx context.Context, code string) (*entity.ConfigGroup, error)
}

// configGroupRepository 配置分组仓储实现
type configGroupRepository struct {
	*baserepo.BaseRepo[entity.ConfigGroup, string]
}

// NewConfigGroupRepository 创建配置分组仓储
func NewConfigGroupRepository(db database.IDataBase) IConfigGroupRepository {
	// 同步表
	tables := []interface{}{
		&entity.ConfigGroup{},
	}
	if err := db.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync tables  error: %v", err)
	}
	return &configGroupRepository{
		BaseRepo: baserepo.NewBaseRepo[entity.ConfigGroup, string](db),
	}
}
func (r *configGroupRepository) FindByCode(ctx context.Context, code string) (*entity.ConfigGroup, error) {
	var group entity.ConfigGroup
	if err := r.Db(ctx).Where("code = ?", code).First(&group).Error; err != nil {
		if database.IfErrorNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &group, nil
}
