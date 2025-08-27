package model

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"time"
)

// File 文件领域模型
type File struct {
	ID          string    // 文件ID
	Key         string    // 文件唯一标识
	Bucket      string    // 存储桶
	Size        int64     // 文件大小
	ContentType string    // 文件类型
	ETag        string    // 文件ETag
	StorageType string    // 存储类型
	URL         string    // 文件访问URL
	ExpiresAt   time.Time // URL过期时间
	IsDeleted   bool      // 是否删除
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
}

// NewFile 创建文件领域模型
func NewFile(
	id string,
	key string,
	bucket string,
	size int64,
	contentType string,
	etag string,
	storageType string,
) *File {
	now := utils.GetTimeNow()
	return &File{
		ID:          id,
		Key:         key,
		Bucket:      bucket,
		Size:        size,
		ContentType: contentType,
		ETag:        etag,
		StorageType: storageType,
		IsDeleted:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// SetURL 设置文件访问URL
func (f *File) SetURL(url string, expiresAt time.Time) {
	f.URL = url
	f.ExpiresAt = expiresAt
	f.UpdatedAt = utils.GetTimeNow()
}

// SoftDelete 软删除文件
func (f *File) SoftDelete() {
	f.IsDeleted = true
	f.UpdatedAt = utils.GetTimeNow()
}

// IsExpired 检查URL是否过期
func (f *File) IsExpired() bool {
	return !f.ExpiresAt.IsZero() && f.ExpiresAt.Before(utils.GetTimeNow())
}
