package commands

// BatchUpdateConfigCommand 批量更新配置命令
type BatchUpdateConfigCommand struct {
	Configs []UpdateConfigCommand `json:"configs" binding:"required"` // 配置列表
}
