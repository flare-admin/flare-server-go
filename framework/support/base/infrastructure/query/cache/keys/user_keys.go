package keys

import "fmt"

const (
	UserPrefix = "sysUser"
)

// UserKey 用户缓存key
func UserKey(tenantID, userID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:detail:%s", UserPrefix, userID)
	}
	return fmt.Sprintf("%s:%s:detail:%s", tenantID, UserPrefix, userID)
}

// UserPermissionsKey 用户权限缓存key
func UserPermissionsKey(tenantID, userID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:permissions:%s", UserPrefix, userID)
	}
	return fmt.Sprintf("%s:%s:permissions:%s", tenantID, UserPrefix, userID)
}

// UserRolesKey 用户角色缓存key
func UserRolesKey(tenantID, userID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:roles:%s", UserPrefix, userID)
	}
	return fmt.Sprintf("%s:%s:roles:%s", tenantID, UserPrefix, userID)
}

// UserMenusKey 用户菜单缓存key
func UserMenusKey(tenantID, userID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:menus:%s", UserPrefix, userID)
	}
	return fmt.Sprintf("%s:%s:menus:%s", tenantID, UserPrefix, userID)
}

// UserRoleCodesKey 用户角色编码缓存key
func UserRoleCodesKey(tenantID, userID string) string {
	return fmt.Sprintf("%s:role:codes:%s", UserPrefix, userID)
}

// UserDepartmentKey 用户部门缓存key
func UserDepartmentKey(tenantID, userID string) string {
	return fmt.Sprintf("user:department:%s", userID)
}

// UserKeys 生成用户相关的所有缓存key
func UserKeys(tenantID, userID string) []string {
	return []string{
		UserKey(tenantID, userID),
		UserPermissionsKey(tenantID, userID),
		UserRolesKey(tenantID, userID),
		UserMenusKey(tenantID, userID),
		UserRoleCodesKey(tenantID, userID),
	}
}

// UserListKey 用户列表缓存key
func UserListKey(tenantID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:list", UserPrefix)
	}
	return fmt.Sprintf("%s:%s:list", tenantID, UserPrefix)
}
