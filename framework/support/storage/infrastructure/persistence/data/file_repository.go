package data

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/pkg/database/baserepo"
	"github.com/flare-admin/flare-server-go/framework/support/storage/infrastructure/persistence/entity"
)

// IFileRepository 文件对象仓储接口
type IFileRepository interface {
	baserepo.IBaseRepo[entity.File, string]
	FindByKey(ctx context.Context, key string) (*entity.File, error)
	FindByKeys(ctx context.Context, keys []string) ([]*entity.File, error)
	SoftDelete(ctx context.Context, key string) error
	BatchSoftDelete(ctx context.Context, keys []string) error
}

// fileRepository 文件对象仓储实现
type fileRepository struct {
	*baserepo.BaseRepo[entity.File, string]
}

// NewFileRepository 创建文件对象仓储
func NewFileRepository(db database.IDataBase) IFileRepository {
	// 同步表
	tables := []interface{}{
		&entity.File{},
	}
	if err := db.AutoMigrate(tables...); err != nil {
		hlog.Fatalf("sync tables error: %v", err)
	}
	return &fileRepository{
		BaseRepo: baserepo.NewBaseRepo[entity.File, string](db),
	}
}

// FindByKey 根据文件标识查找文件
func (r *fileRepository) FindByKey(ctx context.Context, key string) (*entity.File, error) {
	var file entity.File
	if err := r.Db(ctx).Where("key = ? AND is_deleted = ?", key, false).First(&file).Error; err != nil {
		if database.IfErrorNotFound(err) {
			return nil, nil
		}
		hlog.CtxErrorf(ctx, "Find file by key error: %v", err)
		return nil, err
	}
	return &file, nil
}

// FindByKeys 根据文件标识列表查找文件
func (r *fileRepository) FindByKeys(ctx context.Context, keys []string) ([]*entity.File, error) {
	var files []*entity.File
	if err := r.Db(ctx).Where("key IN ? AND is_deleted = ?", keys, false).Find(&files).Error; err != nil {
		hlog.CtxErrorf(ctx, "Find files by keys error: %v", err)
		return nil, err
	}
	return files, nil
}

// SoftDelete 软删除文件
func (r *fileRepository) SoftDelete(ctx context.Context, key string) error {
	if err := r.Db(ctx).Model(&entity.File{}).Where("key = ?", key).Update("is_deleted", true).Error; err != nil {
		hlog.CtxErrorf(ctx, "Soft delete file error: %v", err)
		return err
	}
	return nil
}

// BatchSoftDelete 批量软删除文件
func (r *fileRepository) BatchSoftDelete(ctx context.Context, keys []string) error {
	if err := r.Db(ctx).Model(&entity.File{}).Where("key IN ?", keys).Update("is_deleted", true).Error; err != nil {
		hlog.CtxErrorf(ctx, "Batch soft delete files error: %v", err)
		return err
	}
	return nil
}
