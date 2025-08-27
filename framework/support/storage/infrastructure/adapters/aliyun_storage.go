package adapters

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/valueobject"
	"io"
	"strconv"
	"strings"
)

// AliyunStorageAdapter 阿里云存储适配器
type AliyunStorageAdapter struct {
	config  *configs.AliyunStorage
	expires int64
	client  *oss.Client
}

// NewAliyunStorageAdapter 创建阿里云存储适配器
func NewAliyunStorageAdapter(config *configs.AliyunStorage, expires int64) (domain.StorageAdapter, error) {
	client, err := oss.New(
		config.Endpoint,
		config.AccessKeyID,
		config.AccessKeySecret,
	)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云OSS客户端失败: %w", err)
	}
	adapter := &AliyunStorageAdapter{
		config:  config,
		client:  client,
		expires: expires,
	}
	if err := adapter.EnsureBucket(context.Background()); err != nil {
		return nil, err
	}
	return adapter, nil
}
func (a *AliyunStorageAdapter) GetStorageType() string {
	return "aliyun"
}

// PutObject 上传文件
func (a *AliyunStorageAdapter) PutObject(_ context.Context, key string, reader io.Reader, contentType string) (*valueobject.FileInfo, error) {
	bucket := a.config.BucketName
	// 获取存储桶
	bucketClient, err := a.client.Bucket(bucket)
	if err != nil {
		return nil, fmt.Errorf("获取存储桶失败: %w", err)
	}

	// 上传文件
	err = bucketClient.PutObject(key, reader, oss.ContentType(contentType))
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 获取文件信息
	meta, err := bucketClient.GetObjectDetailedMeta(key)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}
	get := meta.Get("Content-Length")
	size := 0
	atoi, err := strconv.Atoi(get)
	if err == nil {
		size = atoi
	}
	// 返回文件信息
	return &valueobject.FileInfo{
		Key:         key,
		Bucket:      bucket,
		Size:        int64(size),
		ContentType: contentType,
		ETag:        meta.Get("ETag"),
		CreatedAt:   utils.GetTimeNow(),
		UpdatedAt:   utils.GetTimeNow(),
	}, nil
}

// GetObject 获取文件
func (a *AliyunStorageAdapter) GetObject(_ context.Context, key string) (io.ReadCloser, error) {
	bucketClient, err := a.client.Bucket(a.config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("获取存储桶失败: %w", err)
	}

	object, err := bucketClient.GetObject(key)
	if err != nil {
		return nil, fmt.Errorf("获取文件失败: %w", err)
	}

	return object, nil
}

// DeleteObject 删除文件
func (a *AliyunStorageAdapter) DeleteObject(_ context.Context, key string) error {
	bucketClient, err := a.client.Bucket(a.config.BucketName)
	if err != nil {
		return fmt.Errorf("获取存储桶失败: %w", err)
	}

	err = bucketClient.DeleteObject(key)
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// GetObjectURL 获取文件访问URL
func (a *AliyunStorageAdapter) GetObjectURL(_ context.Context, key string) (string, error) {
	// 如果配置了公共URL，直接拼接
	if a.config.PublicURL != "" {
		return fmt.Sprintf("%s/%s", strings.TrimRight(a.config.PublicURL, "/"), key), nil
	}

	// 否则生成签名URL
	bucketClient, err := a.client.Bucket(a.config.BucketName)
	if err != nil {
		return "", fmt.Errorf("获取存储桶失败: %w", err)
	}

	url, err := bucketClient.SignURL(key, oss.HTTPGet, a.expires)
	if err != nil {
		return "", fmt.Errorf("生成签名URL失败: %w", err)
	}

	return url, nil
}

// ListObjects 列出文件
func (a *AliyunStorageAdapter) ListObjects(_ context.Context, prefix string) ([]*valueobject.FileInfo, error) {
	bucketClient, err := a.client.Bucket(a.config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("获取存储桶失败: %w", err)
	}

	var files []*valueobject.FileInfo
	marker := ""
	for {
		lsRes, err := bucketClient.ListObjects(oss.Marker(marker), oss.Prefix(prefix))
		if err != nil {
			return nil, fmt.Errorf("列出文件失败: %w", err)
		}

		for _, object := range lsRes.Objects {
			files = append(files, &valueobject.FileInfo{
				Key:         object.Key,
				Bucket:      a.config.BucketName,
				Size:        object.Size,
				ContentType: object.Type,
				ETag:        object.ETag,
				CreatedAt:   object.LastModified,
				UpdatedAt:   object.LastModified,
			})
		}

		if !lsRes.IsTruncated {
			break
		}
		marker = lsRes.NextMarker
	}

	return files, nil
}

// CopyObject 复制文件
func (a *AliyunStorageAdapter) CopyObject(_ context.Context, srcKey, destKey string) error {
	client, err := a.client.Bucket(a.config.BucketName)
	if err != nil {
		return fmt.Errorf("获取源存储桶失败: %w", err)
	}

	_, err = client.CopyObject(srcKey, destKey)
	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	return nil
}

// MoveObject 移动文件
func (a *AliyunStorageAdapter) MoveObject(ctx context.Context, srcKey, destKey string) error {
	if err := a.CopyObject(ctx, srcKey, destKey); err != nil {
		return err
	}
	return a.DeleteObject(ctx, srcKey)
}

// EnsureBucket 确保存储桶存在
func (a *AliyunStorageAdapter) EnsureBucket(_ context.Context) error {
	exists, err := a.client.IsBucketExist(a.config.BucketName)
	if err != nil {
		return fmt.Errorf("检查存储桶失败: %w", err)
	}

	if !exists {
		err = a.client.CreateBucket(a.config.BucketName)
		if err != nil {
			return fmt.Errorf("创建存储桶失败: %w", err)
		}
	}

	return nil
}
