package domain

import (
	"context"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/model"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/repository"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/valueobject"
	"io"
	"strings"
)

// storageService 存储服务领域实现
type storageService struct {
	adapter StorageAdapter
	repo    repository.StorageRepository
}

// NewStorageService 创建存储服务
func NewStorageService(adapter StorageAdapter, repo repository.StorageRepository) StorageService {
	return &storageService{
		adapter: adapter,
		repo:    repo,
	}
}

// PutObject 上传文件
func (s *storageService) PutObject(ctx context.Context, key string, reader io.Reader, contentType string) (*valueobject.FileInfo, error) {
	// 使用适配器上传文件
	info, err := s.adapter.PutObject(ctx, key, reader, contentType)
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}
	destInfo := model.NewFile("", key, info.Bucket, info.Size, info.ContentType, info.ETag, s.adapter.GetStorageType())
	// 保存文件信息到仓储
	if err := s.repo.Save(ctx, destInfo); err != nil {
		return nil, fmt.Errorf("保存文件信息失败: %w", err)
	}

	return info, nil
}

// GetObject 获取文件
func (s *storageService) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	// 从仓储获取文件信息
	info, err := s.repo.Find(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 使用适配器获取文件内容
	return s.adapter.GetObject(ctx, info.Key)
}

// DeleteObject 删除文件
func (s *storageService) DeleteObject(ctx context.Context, key string) error {
	// 使用适配器删除文件
	if err := s.adapter.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	// 从仓储删除文件信息
	if err := s.repo.Delete(ctx, key); err != nil {
		return fmt.Errorf("删除文件信息失败: %w", err)
	}

	return nil
}

// GetObjectURL 获取文件访问URL
func (s *storageService) GetObjectURL(ctx context.Context, key string) (string, error) {
	// 从仓储获取文件信息
	info, err := s.repo.Find(ctx, key)
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 使用适配器获取文件URL
	return s.adapter.GetObjectURL(ctx, info.Key)
}

// ListObjects 列出文件
func (s *storageService) ListObjects(ctx context.Context, prefix string) ([]*valueobject.FileInfo, error) {
	// 从仓储获取文件列表
	return s.adapter.ListObjects(ctx, prefix)
}

// CopyObject 复制文件
func (s *storageService) CopyObject(ctx context.Context, srcKey, destKey string) error {
	// 从仓储获取源文件信息
	srcInfo, err := s.repo.Find(ctx, srcKey)
	if err != nil {
		return fmt.Errorf("获取源文件信息失败: %w", err)
	}

	// 使用适配器复制文件
	if err := s.adapter.CopyObject(ctx, srcKey, destKey); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	destInfo := model.NewFile("", destKey, srcInfo.Bucket, srcInfo.Size, srcInfo.ContentType, srcInfo.ETag, srcInfo.StorageType)
	// 保存新文件信息到仓储
	if err := s.repo.Save(ctx, destInfo); err != nil {
		return fmt.Errorf("保存新文件信息失败: %w", err)
	}

	return nil
}

// MoveObject 移动文件
func (s *storageService) MoveObject(ctx context.Context, srcKey, destKey string) error {
	// 先复制文件
	if err := s.CopyObject(ctx, srcKey, destKey); err != nil {
		return err
	}

	// 再删除源文件
	return s.DeleteObject(ctx, srcKey)
}

// UploadContent 上传内容
func (s *storageService) UploadContent(ctx context.Context, objectName, content, contentType string) (string, error) {
	// 使用 PutObject 上传内容
	info, err := s.PutObject(ctx, objectName, strings.NewReader(content), contentType)
	if err != nil {
		return "", err
	}

	// 获取文件URL
	url, err := s.GetObjectURL(ctx, info.Key)
	if err != nil {
		return "", err
	}

	return url, nil
}

// ReadFileToContent 读取文件内容
func (s *storageService) ReadFileToContent(ctx context.Context, objectName string) (string, error) {
	// 获取文件内容
	reader, err := s.GetObject(ctx, objectName)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	// 读取内容
	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("读取文件内容失败: %w", err)
	}

	return string(content), nil
}

// ReadFileToContentByName 根据文件名读取文件内容
func (s *storageService) ReadFileToContentByName(ctx context.Context, objectName string) (string, string, error) {
	// 读取文件内容
	content, err := s.ReadFileToContent(ctx, objectName)
	if err != nil {
		return "", "", err
	}

	// 获取文件信息
	info, err := s.repo.Find(ctx, objectName)
	if err != nil {
		return content, "", nil
	}

	return content, info.ContentType, nil
}
