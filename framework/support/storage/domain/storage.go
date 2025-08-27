package domain

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/valueobject"
	"io"
)

// StorageService 存储服务接口
type StorageService interface {
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

	// UploadContent 上传内容
	UploadContent(ctx context.Context, objectName, content, contentType string) (string, error)

	// ReadFileToContent 读取文件内容
	ReadFileToContent(ctx context.Context, objectName string) (string, error)

	// ReadFileToContentByName 根据文件名读取文件内容
	ReadFileToContentByName(ctx context.Context, objectName string) (string, string, error)
}
