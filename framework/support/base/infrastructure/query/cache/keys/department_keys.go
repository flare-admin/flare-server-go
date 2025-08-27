package keys

import "fmt"

const (
	DeptPrefix = "dept"
)

// DepartmentKey 部门缓存key
func DepartmentKey(tenantID string, deptID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:detail:%s", DeptPrefix, deptID)
	}
	return fmt.Sprintf("%s:%s:detail:%s", tenantID, DeptPrefix, deptID)
}

// DepartmentTreeKey 部门树缓存key
func DepartmentTreeKey(tenantID string, parentID string) string {
	if tenantID == "" {
		if parentID == "" {
			return fmt.Sprintf("%s:tree", DeptPrefix)
		}
		return fmt.Sprintf("%s:tree:%s", DeptPrefix, parentID)
	}
	if parentID == "" {
		return fmt.Sprintf("%s:%s:tree", tenantID, DeptPrefix)
	}
	return fmt.Sprintf("%s:%s:tree:%s", tenantID, DeptPrefix, parentID)
}

// UserDepartmentsKey 用户部门缓存key
func UserDepartmentsKey(tenantID string, userID string) string {
	return fmt.Sprintf("%s:%s:user:%s:depts", tenantID, DeptPrefix, userID)
}

// DepartmentKeys 生成部门相关的所有缓存key
func DepartmentKeys(tenantID string, deptID string) []string {
	return []string{
		DepartmentKey(tenantID, deptID),
		DepartmentTreeKey(tenantID, ""),
	}
}

// DepartmentChildrenKey 部门子节点列表缓存key
func DepartmentChildrenKey(tenantID string, parentID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:children:%s", DeptPrefix, parentID)
	}
	return fmt.Sprintf("%s:%s:children:%s", tenantID, DeptPrefix, parentID)
}

// DepartmentUsersKey 部门用户列表缓存key
func DepartmentUsersKey(tenantID string, deptID string) string {
	if tenantID == "" {
		return fmt.Sprintf("%s:users:%s", DeptPrefix, deptID)
	}
	return fmt.Sprintf("%s:%s:users:%s", tenantID, DeptPrefix, deptID)
}
