package keys

import "fmt"

const (
	DataPermissionPrefix = "data_permission"
)

// DataPermissionKey 数据权限缓存key
func DataPermissionKey(tenantID string, roleID int64) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:role:%d", DataPermissionPrefix, roleID)
	}
	return fmt.Sprintf("%s:%s:role:%d", tenantID, DataPermissionPrefix, roleID)
}

// DataPermissionListKey 数据权限列表缓存key
func DataPermissionListKey(tenantID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:list", DataPermissionPrefix)
	}
	return fmt.Sprintf("%s:%s:list", tenantID, DataPermissionPrefix)
}
