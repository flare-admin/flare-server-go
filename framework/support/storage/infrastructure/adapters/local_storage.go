package adapters

import (
	"context"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/valueobject"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorageAdapter 本地存储适配器
type LocalStorageAdapter struct {
	config *configs.LocalStorage
}

// NewLocalStorageAdapter 创建本地存储适配器
func NewLocalStorageAdapter(config *configs.LocalStorage) domain.StorageAdapter {
	return &LocalStorageAdapter{
		config: config,
	}
}
func (a *LocalStorageAdapter) GetStorageType() string {
	return "local"
}

// PutObject 上传文件
func (a *LocalStorageAdapter) PutObject(ctx context.Context, key string, reader io.Reader, contentType string) (*valueobject.FileInfo, error) {
	// 确保目录存在
	dir := filepath.Join(a.config.RootPath, filepath.Dir(key))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %w", err)
	}

	// 创建文件
	filePath := filepath.Join(a.config.RootPath, key)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建文件失败: %w", err)
	}
	defer file.Close()

	// 写入文件
	size, err := io.Copy(file, reader)
	if err != nil {
		return nil, fmt.Errorf("写入文件失败: %w", err)
	}

	// 返回文件信息
	return &valueobject.FileInfo{
		Key:         key,
		Bucket:      a.config.RootPath,
		Size:        size,
		ContentType: contentType,
		ETag:        fmt.Sprintf("%d", utils.GetTimeNow().UnixNano()),
		CreatedAt:   utils.GetTimeNow(),
		UpdatedAt:   utils.GetTimeNow(),
	}, nil
}

// GetObject 获取文件
func (a *LocalStorageAdapter) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	filePath := filepath.Join(a.config.RootPath, key)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	return file, nil
}

// DeleteObject 删除文件
func (a *LocalStorageAdapter) DeleteObject(ctx context.Context, key string) error {
	filePath := filepath.Join(a.config.RootPath, key)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// GetObjectURL 获取文件访问URL
func (a *LocalStorageAdapter) GetObjectURL(ctx context.Context, key string) (string, error) {
	return fmt.Sprintf("%s/%s", strings.TrimRight(a.config.PublicPath, "/"), key), nil
}

// ListObjects 列出文件
func (a *LocalStorageAdapter) ListObjects(ctx context.Context, prefix string) ([]*valueobject.FileInfo, error) {
	dir := a.config.RootPath
	var files []*valueobject.FileInfo

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}

			if strings.HasPrefix(relPath, prefix) {
				files = append(files, &valueobject.FileInfo{
					Key:         relPath,
					Bucket:      a.config.RootPath,
					Size:        info.Size(),
					ContentType: utils.GetContentType(filepath.Ext(relPath)),
					ETag:        fmt.Sprintf("%d", info.ModTime().UnixNano()),
					CreatedAt:   info.ModTime(),
					UpdatedAt:   info.ModTime(),
				})
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("列出文件失败: %w", err)
	}

	return files, nil
}

// CopyObject 复制文件
func (a *LocalStorageAdapter) CopyObject(ctx context.Context, srcKey, destKey string) error {
	srcPath := filepath.Join(a.config.RootPath, srcKey)
	destPath := filepath.Join(a.config.RootPath, destKey)

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	// 复制文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	return nil
}

// MoveObject 移动文件
func (a *LocalStorageAdapter) MoveObject(ctx context.Context, srcKey, destKey string) error {
	if err := a.CopyObject(ctx, srcKey, destKey); err != nil {
		return err
	}
	return a.DeleteObject(ctx, srcKey)
}

// EnsureBucket 确保存储桶存在
func (a *LocalStorageAdapter) EnsureBucket(_ context.Context) error {
	return os.MkdirAll(a.config.RootPath, 0755)
}
