package adapters

import (
	"context"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/configs"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/valueobject"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TencentStorageAdapter 腾讯云存储适配器
type TencentStorageAdapter struct {
	config  *configs.TencentStorage
	client  *cos.Client
	expires int64
}

// NewTencentStorageAdapter 创建腾讯云存储适配器
func NewTencentStorageAdapter(config *configs.TencentStorage, expires int64) (domain.StorageAdapter, error) {
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", config.Bucket, config.Region))
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  config.SecretID,
			SecretKey: config.SecretKey,
		},
	})
	adapter := &TencentStorageAdapter{
		config:  config,
		client:  client,
		expires: expires,
	}
	if err := adapter.EnsureBucket(context.Background()); err != nil {
		return nil, fmt.Errorf("创建存储桶失败: %w", err)
	}
	return adapter, nil
}
func (a *TencentStorageAdapter) GetStorageType() string {
	return "tencent"
}

// PutObject 上传文件
func (a *TencentStorageAdapter) PutObject(ctx context.Context, key string, reader io.Reader, contentType string) (*valueobject.FileInfo, error) {
	// 上传文件
	_, err := a.client.Object.Put(ctx, key, reader, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: contentType,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	// 获取文件信息
	meta, err := a.client.Object.Head(ctx, key, nil)
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %w", err)
	}

	// 返回文件信息
	return &valueobject.FileInfo{
		Key:         key,
		Bucket:      a.config.Bucket,
		Size:        meta.ContentLength,
		ContentType: contentType,
		ETag:        "",
		CreatedAt:   utils.GetTimeNow(),
		UpdatedAt:   utils.GetTimeNow(),
	}, nil
}

// GetObject 获取文件
func (a *TencentStorageAdapter) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	resp, err := a.client.Object.Get(ctx, key, nil)
	if err != nil {
		return nil, fmt.Errorf("获取文件失败: %w", err)
	}
	return resp.Body, nil
}

// DeleteObject 删除文件
func (a *TencentStorageAdapter) DeleteObject(ctx context.Context, key string) error {
	_, err := a.client.Object.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}
	return nil
}

// GetObjectURL 获取文件访问URL
func (a *TencentStorageAdapter) GetObjectURL(ctx context.Context, key string) (string, error) {
	// 如果配置了公共URL，直接拼接
	if a.config.PublicURL != "" {
		return fmt.Sprintf("%s/%s", strings.TrimRight(a.config.PublicURL, "/"), key), nil
	}

	// 否则生成预签名URL
	url, err := a.client.Object.GetPresignedURL(ctx, http.MethodGet, key, a.config.SecretID, a.config.SecretKey, time.Duration(a.expires)*time.Hour*24, nil)
	if err != nil {
		return "", fmt.Errorf("生成预签名URL失败: %w", err)
	}

	return url.String(), nil
}

// ListObjects 列出文件
func (a *TencentStorageAdapter) ListObjects(ctx context.Context, prefix string) ([]*valueobject.FileInfo, error) {
	var files []*valueobject.FileInfo
	marker := ""

	for {
		opt := &cos.BucketGetOptions{
			Prefix:  prefix,
			Marker:  marker,
			MaxKeys: 1000,
		}

		result, _, err := a.client.Bucket.Get(ctx, opt)
		if err != nil {
			return nil, fmt.Errorf("列出文件失败: %w", err)
		}

		for _, object := range result.Contents {
			resp, err := a.client.Object.Head(context.Background(), object.Key, nil)
			if err != nil {
				fmt.Printf("无法获取对象 %s 的信息: %v\n", object.Key, err)
				continue
			}

			contentType := resp.Header.Get("Content-Type")
			files = append(files, &valueobject.FileInfo{
				Key:         object.Key,
				Bucket:      a.config.Bucket,
				Size:        object.Size,
				ContentType: contentType,
				ETag:        object.ETag,
			})
		}

		if !result.IsTruncated {
			break
		}
		marker = result.NextMarker
	}

	return files, nil
}

// CopyObject 复制文件
func (a *TencentStorageAdapter) CopyObject(ctx context.Context, srcKey, destKey string) error {
	source := fmt.Sprintf("%s.cos.%s.myqcloud.com/%s", a.config.Bucket, a.config.Region, srcKey)
	_, _, err := a.client.Object.Copy(ctx, destKey, source, nil)
	if err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}
	return nil
}

// MoveObject 移动文件
func (a *TencentStorageAdapter) MoveObject(ctx context.Context, srcKey, destKey string) error {
	if err := a.CopyObject(ctx, srcKey, destKey); err != nil {
		return err
	}
	return a.DeleteObject(ctx, srcKey)
}

// EnsureBucket 确保存储桶存在
func (a *TencentStorageAdapter) EnsureBucket(ctx context.Context) error {
	_, err := a.client.Bucket.Head(ctx)
	if err != nil {
		return fmt.Errorf("检查存储桶失败: %w", err)
	}
	return nil
}
