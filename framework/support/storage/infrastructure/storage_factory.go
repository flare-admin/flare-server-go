package infrastructure

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain"
	"github.com/flare-admin/flare-server-go/framework/support/storage/infrastructure/adapters"
)

// StorageFactory 存储工厂
type StorageFactory struct {
	config *configs.StorageConfig
}

// NewStorageFactory 创建存储工厂
func NewStorageFactory(config *configs.Bootstrap) *StorageFactory {
	return &StorageFactory{
		config: config.Storage,
	}
}

// CreateAdapter 创建存储适配器
func (f *StorageFactory) CreateAdapter() (domain.StorageAdapter, error) {
	switch f.config.Type {
	case "local":
		return adapters.NewLocalStorageAdapter(f.config.Local), nil
	case "minio":
		return adapters.NewMinioStorageAdapter(f.config.Minio, f.config.UrlExpires)
	case "aliyun":
		return adapters.NewAliyunStorageAdapter(f.config.Aliyun, f.config.UrlExpires)
	case "tencent":
		return adapters.NewTencentStorageAdapter(f.config.Tencent, f.config.UrlExpires)
	default:
		return nil, fmt.Errorf("不支持的存储类型: %s", f.config.Type)
	}
}

func NewStorageAdapter(factory *StorageFactory) (domain.StorageAdapter, error) {
	return factory.CreateAdapter()
}
