package application

import (
	"context"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/valueobject"
	"io"
	"mime/multipart"
	"path"
	"path/filepath"
)

// StorageService 存储服务
type StorageService struct {
	storage domain.StorageService
}

// NewStorageService 创建存储服务
func NewStorageService(storage domain.StorageService) *StorageService {
	return &StorageService{
		storage: storage,
	}
}

// SingleFile 上传单个文件
func (s *StorageService) SingleFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	// 打开文件
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("打开文件失败: %w", err)
	}
	defer src.Close()

	// 生成文件标识
	ext := filepath.Ext(file.Filename)
	key := s.GetFileName(file.Filename)

	// 获取文件类型
	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = utils.GetContentType(ext)
	}

	// 上传文件
	info, err := s.storage.PutObject(ctx, key, src, contentType)
	if err != nil {
		return "", fmt.Errorf("上传文件失败: %w", err)
	}

	return info.Key, nil
}

// MultiFile 上传多个文件
func (s *StorageService) MultiFile(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {
	if len(files) == 0 {
		return nil, nil
	}

	keys := make([]string, 0, len(files))
	for _, file := range files {
		key, err := s.SingleFile(ctx, file)
		if err != nil {
			return nil, fmt.Errorf("上传文件 %s 失败: %w", file.Filename, err)
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// UploadFile 上传文件
func (s *StorageService) UploadFile(ctx context.Context, key string, reader io.Reader, contentType string) (*valueobject.FileInfo, error) {
	return s.storage.PutObject(ctx, key, reader, contentType)
}

// DownloadFile 下载文件
func (s *StorageService) DownloadFile(ctx context.Context, key string) (io.ReadCloser, error) {
	return s.storage.GetObject(ctx, key)
}

// UploadContent 上传内容
func (s *StorageService) UploadContent(ctx context.Context, objectName, content, contentType string) (string, error) {
	return s.storage.UploadContent(ctx, objectName, content, contentType)
}

// ReadFileToContent 读取文件内容
func (s *StorageService) ReadFileToContent(ctx context.Context, objectName string) (string, error) {
	return s.storage.ReadFileToContent(ctx, objectName)
}

// ReadFileToContentByName 读取文件内容
func (s *StorageService) ReadFileToContentByName(ctx context.Context, objectName string) (string, string, error) {
	return s.storage.ReadFileToContentByName(ctx, objectName)
}

// DeleteFile 删除文件
func (s *StorageService) DeleteFile(ctx context.Context, key string) error {
	return s.storage.DeleteObject(ctx, key)
}

// GetFileURL 获取文件URL
func (s *StorageService) GetFileURL(ctx context.Context, key string) (string, error) {
	return s.storage.GetObjectURL(ctx, key)
}

// ListFiles 列出文件
func (s *StorageService) ListFiles(ctx context.Context, prefix string) ([]*valueobject.FileInfo, error) {
	return s.storage.ListObjects(ctx, prefix)
}

// CopyFile 复制文件
func (s *StorageService) CopyFile(ctx context.Context, srcKey, destKey string) error {
	return s.storage.CopyObject(ctx, srcKey, destKey)
}

// MoveFile 移动文件
func (s *StorageService) MoveFile(ctx context.Context, srcKey, destKey string) error {
	return s.storage.MoveObject(ctx, srcKey, destKey)
}

func (s *StorageService) GetFileName(fullFilename string) string {
	//获取文件名带后缀
	suffix := s.GetFileSuffix(fullFilename)
	now := utils.GetTimeNow()
	format := now.Format("20060102150405")
	return fmt.Sprintf("%s_%d%s", format, now.Unix(), suffix)
}
func (s *StorageService) GetFileSuffix(fullFilename string) string {
	//获取文件名带后缀
	filenameWithSuffix := path.Base(fullFilename)
	fileSuffix := path.Ext(filenameWithSuffix)
	return fileSuffix
}
