package model

import (
	"strings"
)

// CacheKey 生成缓存键
// 格式: tenant:group:key
// tenant: 租户ID
// group: 分组ID（可选）
// key: 具体的键名
func CacheKey(tenantID string, groupID string, key string) string {
	parts := []string{tenantID}
	if groupID != "" {
		parts = append(parts, groupID)
	}
	parts = append(parts, key)
	return strings.Join(parts, ":")
}

// ParseCacheKey 解析缓存键
// 返回: tenantID, groupID, key
func ParseCacheKey(cacheKey string) (string, string, string) {
	parts := strings.Split(cacheKey, ":")
	if len(parts) < 2 {
		return "", "", cacheKey
	}
	if len(parts) == 2 {
		return parts[0], "", parts[1]
	}
	return parts[0], parts[1], parts[2]
}
