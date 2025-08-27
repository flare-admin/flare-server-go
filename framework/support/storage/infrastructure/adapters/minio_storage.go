package adapters

import (
	"context"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/valueobject"
	"io"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioStorageAdapter MinIO存储适配器
type MinioStorageAdapter struct {
	config  *configs.MinioStorage
	client  *minio.Client
	expires int64
}

// NewMinioStorageAdapter 创建MinIO存储适配器
func NewMinioStorageAdapter(config *configs.MinioStorage, expires int64) (domain.StorageAdapter, error) {
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("创建MinIO客户端失败: %w", err)
	}
	adapter := &MinioStorageAdapter{
		config:  config,
		client:  client,
		expires: expires,
	}
	if err := adapter.EnsureBucket(context.Background()); err != nil {
		return nil, err
	}
	return adapter, nil
}
func (a *MinioStorageAdapter) GetStorageType() string {
	return "minio"
}

// PutObject 上传文件
func (a *MinioStorageAdapter) PutObject(ctx context.Context, key string, reader io.Reader, contentType string) (*valueobject.FileInfo, error) {
	// 确保存储桶存在
	if err := a.EnsureBucket(ctx); err != nil {
		return nil, err
	}

	// 上传文件
	info, err := a.client.PutObject(ctx, a.config.Bucket, key, reader, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 返回文件信息
	return &valueobject.FileInfo{
		Key:         key,
		Bucket:      a.config.Bucket,
		Size:        info.Size,
		ContentType: contentType,
		ETag:        info.ETag,
		CreatedAt:   utils.GetTimeNow(),
		UpdatedAt:   utils.GetTimeNow(),
	}, nil
}

// GetObject 获取文件
func (a *MinioStorageAdapter) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	object, err := a.client.GetObject(ctx, a.config.Bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取文件失败: %w", err)
	}
	return object, nil
}

// DeleteObject 删除文件
func (a *MinioStorageAdapter) DeleteObject(ctx context.Context, key string) error {
	err := a.client.RemoveObject(ctx, a.config.Bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// GetObjectURL 获取文件访问URL
func (a *MinioStorageAdapter) GetObjectURL(ctx context.Context, key string) (string, error) {
	// 如果配置了公共URL，直接拼接
	if a.config.PublicURL != "" {
		return fmt.Sprintf("%s/%s", strings.TrimRight(a.config.PublicURL, "/"), key), nil
	}

	// 否则生成预签名URL
	url, err := a.client.PresignedGetObject(ctx, a.config.Bucket, key, time.Duration(a.expires)*time.Hour*24, nil)
	if err != nil {
		return "", fmt.Errorf("获取文件URL失败: %w", err)
	}
	return url.String(), nil
}

// ListObjects 列出文件
func (a *MinioStorageAdapter) ListObjects(ctx context.Context, prefix string) ([]*valueobject.FileInfo, error) {
	var files []*valueobject.FileInfo

	// 列出对象
	objectCh := a.client.ListObjects(ctx, a.config.Bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("列出文件失败: %w", object.Err)
		}

		files = append(files, &valueobject.FileInfo{
			Key:         object.Key,
			Bucket:      a.config.Bucket,
			Size:        object.Size,
			ContentType: object.ContentType,
			ETag:        object.ETag,
			CreatedAt:   object.LastModified,
			UpdatedAt:   object.LastModified,
		})
	}

	return files, nil
}

// CopyObject 复制文件
func (a *MinioStorageAdapter) CopyObject(ctx context.Context, srcKey, destKey string) error {
	_, err := a.client.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: a.config.Bucket,
		Object: destKey,
	}, minio.CopySrcOptions{
		Bucket: a.config.Bucket,
		Object: srcKey,
	})
	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}
	return nil
}

// MoveObject 移动文件
func (a *MinioStorageAdapter) MoveObject(ctx context.Context, srcKey, destKey string) error {
	if err := a.CopyObject(ctx, srcKey, destKey); err != nil {
		return err
	}
	return a.DeleteObject(ctx, srcKey)
}

// EnsureBucket 确保存储桶存在
func (a *MinioStorageAdapter) EnsureBucket(ctx context.Context) error {
	//exists, err := a.client.BucketExists(ctx, a.config.Bucket)
	//if err != nil {
	//	return fmt.Errorf("检查存储桶失败: %w", err)
	//}
	//
	//if !exists {
	//	err = a.client.MakeBucket(ctx, a.config.Bucket, minio.MakeBucketOptions{})
	//	if err != nil {
	//		return fmt.Errorf("创建存储桶失败: %w", err)
	//	}
	//}

	return nil
}
