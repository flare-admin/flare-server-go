package repository

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/model"
	dr "github.com/flare-admin/flare-server-go/framework/support/storage/domain/repository"
	"github.com/flare-admin/flare-server-go/framework/support/storage/infrastructure/persistence/data"
	entity2 "github.com/flare-admin/flare-server-go/framework/support/storage/infrastructure/persistence/entity"
)

type StorageRepo struct {
	rp data.IFileRepository
}

func NewStorageRepo(rp data.IFileRepository) dr.StorageRepository {
	return &StorageRepo{rp: rp}
}

// Save 保存文件信息
func (s StorageRepo) Save(ctx context.Context, info *model.File) error {
	// 转换为实体
	entity := entity2.FromDomain(info)
	if entity.ID == "" {
		// 保存到数据库
		if _, err := s.rp.Add(ctx, entity); err != nil {
			hlog.CtxErrorf(ctx, "Save file error: %v", err)
			return err
		}
	} else {
		if err := s.rp.EditById(ctx, entity.ID, entity); err != nil {
			hlog.CtxErrorf(ctx, "Save file error: %v", err)
			return err
		}
	}
	return nil
}

// Find 查找文件信息
func (s StorageRepo) Find(ctx context.Context, key string) (*model.File, error) {
	// 从数据库查找
	entity, err := s.rp.FindByKey(ctx, key)
	if err != nil {
		hlog.CtxErrorf(ctx, "Find file error: %v", err)
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

	// 转换为领域模型
	return entity2.ToDomain(entity), nil
}

// Delete 删除文件信息
func (s StorageRepo) Delete(ctx context.Context, key string) error {
	// 软删除
	if err := s.rp.SoftDelete(ctx, key); err != nil {
		hlog.CtxErrorf(ctx, "Delete file error: %v", err)
		return err
	}
	return nil
}
