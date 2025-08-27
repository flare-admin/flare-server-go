package repository

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/infrastructure/entity"
)

// IConfigRepository 配置接口
type IConfigRepository interface {
	baserepo.IBaseRepo[entity.Config, string]
	FindByKey(ctx context.Context, key string) (*entity.Config, error)
	BatchUpdate(ctx context.Context, configs []*entity.Config) error
}

// configGroupRepository 配置分组仓储实现
type configRepository struct {
	*baserepo.BaseRepo[entity.Config, string]
}

// NewConfigRepository 创建配置仓储
func NewConfigRepository(db database.IDataBase) IConfigRepository {
	// 同步表
	tables := []interface{}{
		&entity.Config{},
	}
	if err := db.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync tables  error: %v", err)
	}
	return &configRepository{
		BaseRepo: baserepo.NewBaseRepo[entity.Config, string](db, entity.Config{}),
	}
}

func (c configRepository) FindByKey(ctx context.Context, key string) (*entity.Config, error) {
	var config entity.Config
	if err := c.Db(ctx).Where("key = ?", key).First(&config).Error; err != nil {
		if database.IfErrorNotFound(err) {
			return nil, nil
		}
		hlog.CtxErrorf(ctx, "Find config by key error: %v", err)
		return nil, err
	}
	return &config, nil
}

// BatchUpdate 批量更新配置
func (c configRepository) BatchUpdate(ctx context.Context, configs []*entity.Config) error {
	if len(configs) == 0 {
		return nil
	}

	// 使用事务进行批量更新
	tx := c.Db(ctx).Begin()
	if tx.Error != nil {
		hlog.CtxErrorf(ctx, "Begin transaction error: %v", tx.Error)
		return tx.Error
	}

	for _, config := range configs {
		if err := tx.Model(&entity.Config{}).Where("id = ?", config.ID).Updates(map[string]interface{}{
			"value":       config.Value,
			"type":        config.Type,
			"description": config.Description,
			"i18n_key":    config.I18nKey,
			"is_system":   config.IsSystem,
			"is_enabled":  config.IsEnabled,
			"sort":        config.Sort,
			"updated_at":  config.UpdatedAt,
		}).Error; err != nil {
			tx.Rollback()
			hlog.CtxErrorf(ctx, "Batch update config error: %v", err)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		hlog.CtxErrorf(ctx, "Commit transaction error: %v", err)
		return err
	}

	return nil
}
