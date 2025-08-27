package model

import "fmt"

// CacheGroup 缓存分组常量
const (
	// CacheGroupUser 用户信息缓存分组
	CacheGroupUser = "user"
	// CacheGroupDict 数据字典缓存分组
	CacheGroupDict = "dict"
	// CacheGroupConfig 系统配置缓存分组
	CacheGroupConfig = "config"
	// CacheGroupDefault 默认缓存分组
	CacheGroupDefault = "default"
)

// CacheGroupInfo 缓存分组信息
type CacheGroupInfo struct {
	GroupID     string `json:"group_id"`     // 分组ID
	Name        string `json:"name"`         // 分组名称
	Description string `json:"description"`  // 分组描述
	KeyCount    int64  `json:"key_count"`    // 键数量
	MemoryUsage int64  `json:"memory_usage"` // 内存使用量（字节）
}

// String 返回缓存分组信息的字符串表示
func (info *CacheGroupInfo) String() string {
	return fmt.Sprintf("分组: %s, 键数量: %d, 内存使用: %d 字节",
		info.GroupID, info.KeyCount, info.MemoryUsage)
}

// GetDefaultCacheGroups 获取默认缓存分组列表
func GetDefaultCacheGroups() []*CacheGroupInfo {
	return []*CacheGroupInfo{
		{
			GroupID:     CacheGroupUser,
			Name:        "用户信息",
			Description: "用户相关的缓存信息，如用户基本信息、权限等",
		},
		{
			GroupID:     CacheGroupDict,
			Name:        "数据字典",
			Description: "系统数据字典缓存，如枚举值、配置项等",
		},
		{
			GroupID:     CacheGroupConfig,
			Name:        "系统配置",
			Description: "系统配置信息缓存，如系统参数、业务配置等",
		},
		{
			GroupID:     CacheGroupDefault,
			Name:        "默认分组",
			Description: "系统默认缓存分组",
		},
	}
}

// IsValidCacheGroup 检查是否是有效的缓存分组
func IsValidCacheGroup(groupID string) bool {
	switch groupID {
	case CacheGroupUser, CacheGroupDict, CacheGroupConfig, CacheGroupDefault:
		return true
	default:
		return false
	}
}
