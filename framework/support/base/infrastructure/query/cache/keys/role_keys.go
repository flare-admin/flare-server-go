package keys

import (
	"fmt"
)

const (
	// 角色缓存key前缀
	RolePrefix = "role"
)

// RoleKey 生成角色缓存key
func RoleKey(tenantID string, id int64) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:detail:%d", RolePrefix, id)
	}
	return fmt.Sprintf("%s:%s:detail:%d", tenantID, RolePrefix, id)
}

// RoleCodeKey 生成角色编码缓存key
func RoleCodeKey(tenantID string, code string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:code:%s", RolePrefix, code)
	}
	return fmt.Sprintf("%s:%s:code:%s", tenantID, RolePrefix, code)
}

// RolePermKey 生成角色权限缓存key
func RolePermKey(tenantID string, roleID int64) string {
	return fmt.Sprintf("%s:%s:perm:%d", tenantID, RolePrefix, roleID)
}

// RolePermissionsKey 角色权限缓存key
func RolePermissionsKey(tenantID string, roleID int64) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:permissions:%d", RolePrefix, roleID)
	}
	return fmt.Sprintf("%s:%s:permissions:%d", tenantID, RolePrefix, roleID)
}

// RoleListKey 角色列表缓存key
func RoleListKey(tenantID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:list", RolePrefix)
	}
	return fmt.Sprintf("%s:%s:list", tenantID, RolePrefix)
}

// RoleKeys 生成角色相关的所有缓存key
func RoleKeys(tenantID string, roleID int64) []string {
	return []string{
		RoleKey(tenantID, roleID),
		RolePermKey(tenantID, roleID),
		RolePermissionsKey(tenantID, roleID),
		RoleListKey(tenantID),
	}
}

// RoleUsersKey 角色用户缓存key
func RoleUsersKey(tenantID string, roleID int64) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:users:%d", RolePrefix, roleID)
	}
	return fmt.Sprintf("%s:%s:users:%d", tenantID, RolePrefix, roleID)
}
