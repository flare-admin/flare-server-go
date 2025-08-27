package database

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"gorm.io/gorm"
)

type BaseIntTime struct {
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;not null;default:0;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;not null;default:0;comment:更新时间"`
	DeletedAt int64 `json:"deletedAt" gorm:"column:deleted_at;not null;default:0;comment:删除时间"`
}

// 插入前自动设置
func (m *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	now := utils.GetDateUnix()
	m.CreatedAt = now
	m.UpdatedAt = now
	return
}

// 更新前自动设置
func (m *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	m.UpdatedAt = utils.GetDateUnix()
	return
}

type BaseModel struct {
	BaseIntTime
	Creator  string `json:"creator"  gorm:"column:creator;not null;default:'';comment:创建者"`
	Updater  string `json:"updater"  gorm:"column:updater;not null;default:'';comment:更新人"`
	TenantID string `json:"tenantId" gorm:"index:tenantIndex,column:tenant_id;default:'';comment:租户ID"`
}
