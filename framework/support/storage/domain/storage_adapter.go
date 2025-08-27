package domain

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/valueobject"
	"io"
)

// StorageAdapter 存储适配器接口
type StorageAdapter interface {
	// GetStorageType 获取存储类型
	GetStorageType() string
	// PutObject 上传文件
	PutObject(ctx context.Context, key string, reader io.Reader, contentType string) (*valueobject.FileInfo, error)

	// GetObject 获取文件
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)

	// DeleteObject 删除文件
	DeleteObject(ctx context.Context, key string) error

	// GetObjectURL 获取文件访问URL
	GetObjectURL(ctx context.Context, key string) (string, error)

	// ListObjects 列出文件
	ListObjects(ctx context.Context, prefix string) ([]*valueobject.FileInfo, error)

	// CopyObject 复制文件（仅在桶内）
	CopyObject(ctx context.Context, srcKey, destKey string) error

	// MoveObject 移动文件（仅在桶内）
	MoveObject(ctx context.Context, srcKey, destKey string) error

	// EnsureBucket 确保存储桶存在
	EnsureBucket(ctx context.Context) error
}
