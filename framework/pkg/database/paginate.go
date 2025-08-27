package database

import "gorm.io/gorm"

// Page 分页参数
type Page struct {
	Size    int `json:"size" query:"size"`       // 页码大小，最大100
	Current int `json:"current" query:"current"` // 页码，从1开始
}

func NewPage(current, size int) *Page {
	return &Page{
		Current: current,
		Size:    size,
	}
}

func (r *Page) Fix() {
	if r.Current <= 0 {
		r.Current = 1
	}

	if r.Size <= 0 {
		r.Size = 10
	} else if r.Size > 500 {
		r.Size = 500
	}
}

func Operation(current, size int) func(db *gorm.DB) *gorm.DB {
	r := NewPage(current, size)
	return func(db *gorm.DB) *gorm.DB {
		r.Fix()
		offset := (r.Current - 1) * r.Size
		return db.Offset(offset).Limit(int(r.Size))
	}
}
