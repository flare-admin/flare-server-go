package keys

import "fmt"

const (
	PermissionKeyPrefix = "permission"
)

// PermissionKey 权限详情缓存key
func PermissionKey(id int64) string {
	return fmt.Sprintf("%s:detail:%d", PermissionKeyPrefix, id)
}

// PermissionTreeKey 权限树缓存key
func PermissionTreeKey(tenantID string, permType interface{}) string {
	if permType == nil {
		return fmt.Sprintf("%s:%s:tree", tenantID, PermissionKeyPrefix)
	}
	return fmt.Sprintf("%s:%s:tree:%v", tenantID, PermissionKeyPrefix, permType)
}

// PermissionListKey 权限列表缓存key
func PermissionListKey(tenantID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:list", PermissionKeyPrefix)
	}
	return fmt.Sprintf("%s:%s:list", tenantID, PermissionKeyPrefix)
}

// PermissionEnabledKey 启用权限列表缓存key
func PermissionEnabledKey(tenantID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:enabled", PermissionKeyPrefix)
	}
	return fmt.Sprintf("%s:%s:enabled", tenantID, PermissionKeyPrefix)
}

// PermissionSimpleTreeKey 简化权限树缓存key
func PermissionSimpleTreeKey(tenantID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:simple:tree", PermissionKeyPrefix)
	}
	return fmt.Sprintf("%s:%s:simple:tree", tenantID, PermissionKeyPrefix)
}

// 权限相关的缓存键
func PermissionChildrenKey(tenantID string, parentID int64) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:children:%d", PermissionKeyPrefix, parentID)
	}
	return fmt.Sprintf("%s:children:%s:%d", tenantID, PermissionKeyPrefix, parentID)
}

func PermissionResourceKey(tenantID string, permID int64) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:resource:%d", PermissionKeyPrefix, permID)
	}
	return fmt.Sprintf("%s:resource:%s:%d", tenantID, PermissionKeyPrefix, permID)
}
