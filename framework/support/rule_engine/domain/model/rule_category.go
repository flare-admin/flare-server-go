package model

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
)

// RuleCategory 规则分类领域模型
type RuleCategory struct {
	// 基础信息
	ID          string `json:"id"`          // 分类ID
	Code        string `json:"code"`        // 分类编码
	Name        string `json:"name"`        // 分类名称
	Description string `json:"description"` // 分类描述

	// 分类配置
	Type     string `json:"type"`     // 分类类型：business(业务分类) system(系统分类) custom(自定义分类)
	ParentID string `json:"parentId"` // 父分类ID
	Level    int32  `json:"level"`    // 分类层级
	Path     string `json:"path"`     // 分类路径，如：/1/2/3
	Sorting  int32  `json:"sorting"`  // 排序权重

	// 状态信息
	Status int32 `json:"status"` // 状态：1-启用 2-禁用
	IsLeaf bool  `json:"isLeaf"` // 是否为叶子节点

	// 业务配置
	BusinessType string `json:"businessType"` // 业务类型：order(订单) user(用户) product(商品) payment(支付) withdrawal(提现) declaration(申报)

	// 时间信息
	CreatedAt int64 `json:"createdAt"` // 创建时间
	UpdatedAt int64 `json:"updatedAt"` // 更新时间

	// 租户信息
	TenantID string `json:"tenantId"` // 租户ID
}

// NewRuleCategory 创建规则分类
func NewRuleCategory(code, name, description, categoryType, businessType string) *RuleCategory {
	now := utils.GetDateUnix()
	return &RuleCategory{
		ID:           "",
		Code:         code,
		Name:         name,
		Description:  description,
		Type:         categoryType,
		ParentID:     "",
		Level:        1,
		Path:         "",
		Sorting:      0,
		Status:       1,
		IsLeaf:       true,
		BusinessType: businessType,
		CreatedAt:    now,
		UpdatedAt:    now,
		TenantID:     "",
	}
}

// SetParent 设置父分类
func (rc *RuleCategory) SetParent(parentID string, parentPath string) {
	rc.ParentID = parentID
	if parentPath == "" {
		rc.Path = "/" + rc.ID
		rc.Level = 1
	} else {
		rc.Path = parentPath + "/" + rc.ID
		rc.Level = int32(len(rc.Path) / 2) // 简单计算层级
	}
}

// Enable 启用分类
func (rc *RuleCategory) Enable() {
	rc.Status = 1
	rc.UpdatedAt = utils.GetDateUnix()
}

// Disable 禁用分类
func (rc *RuleCategory) Disable() {
	rc.Status = 2
	rc.UpdatedAt = utils.GetDateUnix()
}

// IsEnabled 是否启用
func (rc *RuleCategory) IsEnabled() bool {
	return rc.Status == 1
}

// SetAsLeaf 设置为叶子节点
func (rc *RuleCategory) SetAsLeaf(isLeaf bool) {
	rc.IsLeaf = isLeaf
	rc.UpdatedAt = utils.GetDateUnix()
}

// Update 更新分类
func (rc *RuleCategory) Update(name, description string) {
	rc.Name = name
	rc.Description = description
	rc.UpdatedAt = utils.GetDateUnix()
}

// SetSorting 设置排序
func (rc *RuleCategory) SetSorting(sorting int32) {
	rc.Sorting = sorting
	rc.UpdatedAt = utils.GetDateUnix()
}

// Validate 验证分类
func (rc *RuleCategory) Validate() error {
	if rc.Code == "" {
		return fmt.Errorf("category code cannot be empty")
	}

	if rc.Name == "" {
		return fmt.Errorf("category name cannot be empty")
	}

	// 验证分类类型
	validTypes := []string{"business", "system", "custom"}
	isValidType := false
	for _, validType := range validTypes {
		if rc.Type == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("invalid category type: %s", rc.Type)
	}

	// 验证业务类型
	validBusinessTypes := []string{"order", "user", "product", "payment", "withdrawal", "declaration", "lottery", "recharge", "exchange"}
	isValidBusinessType := false
	for _, validBusinessType := range validBusinessTypes {
		if rc.BusinessType == validBusinessType {
			isValidBusinessType = true
			break
		}
	}

	if !isValidBusinessType {
		return fmt.Errorf("invalid business type: %s", rc.BusinessType)
	}

	return nil
}

// IsRoot 是否为根分类
func (rc *RuleCategory) IsRoot() bool {
	return rc.ParentID == ""
}

// GetPathLevel 获取路径层级
func (rc *RuleCategory) GetPathLevel() int32 {
	if rc.Path == "" {
		return 0
	}

	// 计算路径中的层级数
	count := 0
	for _, char := range rc.Path {
		if char == '/' {
			count++
		}
	}
	return int32(count - 1) // 减去开头的斜杠
}

// IsDescendantOf 是否为指定分类的后代
func (rc *RuleCategory) IsDescendantOf(categoryPath string) bool {
	if categoryPath == "" || rc.Path == "" {
		return false
	}

	// 检查当前路径是否以指定路径开头
	return len(rc.Path) > len(categoryPath) && rc.Path[:len(categoryPath)] == categoryPath
}

// IsAncestorOf 是否为指定分类的祖先
func (rc *RuleCategory) IsAncestorOf(categoryPath string) bool {
	if rc.Path == "" || categoryPath == "" {
		return false
	}

	// 检查指定路径是否以当前路径开头
	return len(categoryPath) > len(rc.Path) && categoryPath[:len(rc.Path)] == rc.Path
}

// GetAncestors 获取所有祖先分类ID
func (rc *RuleCategory) GetAncestors() []string {
	if rc.Path == "" {
		return []string{}
	}

	// 解析路径获取所有祖先ID
	pathParts := rc.Path[1:] // 去掉开头的斜杠
	if pathParts == "" {
		return []string{}
	}

	parts := []string{}
	for i := 0; i < len(pathParts); i += 2 {
		if i+1 < len(pathParts) {
			parts = append(parts, string(pathParts[i:i+2]))
		}
	}

	return parts
}

// GetDescendants 获取所有后代分类ID（需要从数据库查询）
func (rc *RuleCategory) GetDescendants() []string {
	// 这个方法需要在仓储层实现，通过查询数据库获取所有以当前路径开头的分类
	// 这里只是占位符
	return []string{}
}

// IsLeafCategory 是否为叶子分类
func (rc *RuleCategory) IsLeafCategory() bool {
	return rc.IsLeaf
}

// GetFullPath 获取完整路径
func (rc *RuleCategory) GetFullPath() string {
	return rc.Path
}

// GetLevel 获取层级
func (rc *RuleCategory) GetLevel() int32 {
	return rc.Level
}
