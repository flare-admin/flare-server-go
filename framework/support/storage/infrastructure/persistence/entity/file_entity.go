package entity

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/database"
	"github.com/flare-admin/flare-server-go/framework/support/storage/domain/model"
	"time"
)

// File 文件对象实体
type File struct {
	database.BaseModel
	ID          string `gorm:"column:id;primaryKey;pk:true"`                          // 文件ID
	Key         string `gorm:"column:key;type:varchar(255);not null;uniqueIndex"`     // 文件唯一标识
	Bucket      string `gorm:"column:bucket;type:varchar(100);not null"`              // 存储桶
	Size        int64  `gorm:"column:size;type:bigint;not null"`                      // 文件大小
	ContentType string `gorm:"column:content_type;type:varchar(100);not null"`        // 文件类型
	ETag        string `gorm:"column:etag;type:varchar(100)"`                         // 文件ETag
	StorageType string `gorm:"column:storage_type;type:varchar(20);not null"`         // 存储类型
	URL         string `gorm:"column:url;type:varchar(500)"`                          // 文件访问URL
	ExpiresAt   int64  `gorm:"column:expires_at"`                                     // URL过期时间
	IsDeleted   bool   `gorm:"column:is_deleted;type:boolean;not null;default:false"` // 是否删除
}

// GetPrimaryKey 获取主键
func (f File) GetPrimaryKey() string {
	return "id"
}

// TableName 表名
func (File) TableName() string {
	return "files"
}

func FromDomain(domain *model.File) *File {
	return &File{
		ID:          domain.ID,
		Key:         domain.Key,
		Bucket:      domain.Bucket,
		Size:        domain.Size,
		ContentType: domain.ContentType,
		ETag:        domain.ETag,
		StorageType: domain.StorageType,
		URL:         domain.URL,
		ExpiresAt:   domain.ExpiresAt.Unix(),
		IsDeleted:   domain.IsDeleted,
	}
}

func ToDomain(e *File) *model.File {
	return &model.File{
		ID:          e.ID,
		Key:         e.Key,
		Bucket:      e.Bucket,
		Size:        e.Size,
		ContentType: e.ContentType,
		ETag:        e.ETag,
		StorageType: e.StorageType,
		URL:         e.URL,
		ExpiresAt:   time.Unix(e.ExpiresAt, 0),
		IsDeleted:   e.IsDeleted,
	}
}
