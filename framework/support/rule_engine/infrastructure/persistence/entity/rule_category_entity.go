package entity

import "github.com/flare-admin/flare-server-go/framework/pkg/database"

// RuleCategory 规则分类实体
type RuleCategory struct {
	database.BaseModel
	ID           string `gorm:"primarykey"`
	Code         string `gorm:"size:100;not null;uniqueIndex;comment:分类编码"`
	Name         string `gorm:"size:100;not null;comment:分类名称"`
	Description  string `gorm:"size:500;comment:分类描述"`
	Type         string `gorm:"size:50;not null;comment:分类类型：business(业务分类) system(系统分类) custom(自定义分类)"`
	ParentID     string `gorm:"index;comment:父分类ID"`
	Level        int32  `gorm:"not null;default:1;comment:分类层级"`
	Path         string `gorm:"size:500;comment:分类路径，如：/1/2/3"`
	Sorting      int32  `gorm:"not null;default:0;comment:排序权重"`
	Status       int    `gorm:"not null;default:1;comment:状态：1-启用 2-禁用"`
	IsLeaf       bool   `gorm:"not null;default:true;comment:是否为叶子节点"`
	BusinessType string `gorm:"size:50;not null;comment:业务类型：order(订单) user(用户) product(商品) payment(支付) withdrawal(提现) declaration(申报)"`
	TenantID     string `gorm:"size:50;comment:租户ID"`
}

func (RuleCategory) TableName() string {
	return "rule_categories"
}

// GetPrimaryKey 获取主键
func (RuleCategory) GetPrimaryKey() string {
	return "id"
}
