package command

// CreateCategoryCommand 创建分类命令
type CreateCategoryCommand struct {
	Name        string `json:"name" binding:"required" comment:"分类名称"`
	Code        string `json:"code" binding:"required" comment:"分类编码"`
	Description string `json:"description" comment:"分类描述"`
	Sort        int    `json:"sort" comment:"排序"`
}

// UpdateCategoryCommand 更新分类命令
type UpdateCategoryCommand struct {
	ID          string `json:"id" binding:"required" comment:"分类ID"`
	Name        string `json:"name" binding:"required" comment:"分类名称"`
	Code        string `json:"code" binding:"required" comment:"分类编码"`
	Description string `json:"description" comment:"分类描述"`
	Sort        int    `json:"sort" comment:"排序"`
}

// UpdateCategoryStatusCommand 更新分类状态命令
type UpdateCategoryStatusCommand struct {
	ID     string `json:"id" binding:"required" comment:"分类ID"`
	Status int    `json:"status" binding:"required" comment:"状态"`
}

// DeleteCategoryCommand 删除分类命令
type DeleteCategoryCommand struct {
	ID string `json:"id" binding:"required" comment:"分类ID"`
}
