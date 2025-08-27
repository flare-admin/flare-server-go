package repository

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/model"
)

// StorageRepository 存储仓储接口
type StorageRepository interface {
	// Save 保存文件信息
	Save(ctx context.Context, info *model.File) error

	// Find 查找文件信息
	Find(ctx context.Context, key string) (*model.File, error)

	// Delete 删除文件信息
	Delete(ctx context.Context, key string) error
}
