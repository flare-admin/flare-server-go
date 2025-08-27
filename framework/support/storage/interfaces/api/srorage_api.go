package storage_api

import (
	"context"
	"github.com/flare-admin/flare-server-go/framework/support/storage/application"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/valueobject"
	"io"
)

type IStorageApi interface {
	// PutObject 上传文件
	PutObject(ctx context.Context, key string, reader io.Reader, contentType string) (*valueobject.FileInfo, error)

	// GetObject 获取文件
	GetObject(ctx context.Context, key string) (io.ReadCloser, error)

	// DeleteObject 删除文件
	DeleteObject(ctx context.Context, key string) error

	// GetObjectURL 获取文件访问URL
	GetObjectURL(ctx context.Context, key string) (string, error)

	// UploadContent 上传内容
	UploadContent(ctx context.Context, objectName, content, contentType string) (string, error)

	// ReadFileToContent 读取文件内容
	ReadFileToContent(ctx context.Context, objectName string) (string, error)

	// ReadFileToContentByName 根据文件名读取文件内容
	ReadFileToContentByName(ctx context.Context, objectName string) (string, string, error)
}

type StorageApi struct {
	service *application.StorageService
}

func NewStorageApi(service *application.StorageService) IStorageApi {
	return &StorageApi{
		service: service,
	}
}

// PutObject 上传文件
func (s *StorageApi) PutObject(ctx context.Context, key string, reader io.Reader, contentType string) (*valueobject.FileInfo, error) {
	return s.service.UploadFile(ctx, key, reader, contentType)
}

// GetObject 获取文件
func (s *StorageApi) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	return s.service.DownloadFile(ctx, key)
}

// DeleteObject 删除文件
func (s *StorageApi) DeleteObject(ctx context.Context, key string) error {
	return s.service.DeleteFile(ctx, key)
}

// GetObjectURL 获取文件访问URL
func (s *StorageApi) GetObjectURL(ctx context.Context, key string) (string, error) {
	return s.service.GetFileURL(ctx, key)
}

// UploadContent 上传内容
func (s *StorageApi) UploadContent(ctx context.Context, objectName, content, contentType string) (string, error) {
	return s.service.UploadContent(ctx, objectName, content, contentType)
}

// ReadFileToContent 读取文件内容
func (s *StorageApi) ReadFileToContent(ctx context.Context, objectName string) (string, error) {
	return s.service.ReadFileToContent(ctx, objectName)
}

// ReadFileToContentByName 根据文件名读取文件内容
func (s *StorageApi) ReadFileToContentByName(ctx context.Context, objectName string) (string, string, error) {
	return s.service.ReadFileToContentByName(ctx, objectName)
}
