package dto

type PermissionsTreeResult struct {
	Tree []*PermissionsTreeDto `json:"tree"`
	Ids  []int64               `json:"ids"`
}

// PermissionsTreeDto 简化的权限树数据传输对象
type PermissionsTreeDto struct {
	ID          int64                 `json:"id"`          // ID
	Code        string                `json:"code"`        // 编码
	Name        string                `json:"name"`        // 名称
	Localize    string                `json:"localize"`    // 本地化
	Icon        string                `json:"icon"`        // 图标
	Description string                `json:"description"` // 描述
	Sequence    int                   `json:"sequence"`    // 排序
	Type        int8                  `json:"type"`        // 类型
	Component   string                `json:"component"`   // 组件
	Path        string                `json:"path"`        // 路径
	Properties  string                `json:"properties"`  // 属性
	Status      int8                  `json:"status"`      // 状态
	ParentID    int64                 `json:"parentId"`    // 父级ID
	ParentPath  string                `json:"parent_path"` // 父级路径
	Children    []*PermissionsTreeDto `json:"children,omitempty"`
}
