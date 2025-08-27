package valueobject

import "time"

// FileInfo 文件信息
type FileInfo struct {
	Key         string    // 文件唯一标识
	Bucket      string    // 存储桶
	Size        int64     // 文件大小
	ContentType string    // 文件类型
	ETag        string    // 文件ETag
	CreatedAt   time.Time // 创建时间
	UpdatedAt   time.Time // 更新时间
}
