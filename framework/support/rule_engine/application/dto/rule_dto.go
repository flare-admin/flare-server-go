package dto

// ==================== 模板相关DTO ====================

// TemplateDTO 模板数据传输对象
type TemplateDTO struct {
	ID          string `json:"id"`          // 模板ID
	Code        string `json:"code"`        // 模板编码
	Name        string `json:"name"`        // 模板名称
	Description string `json:"description"` // 模板描述
	CategoryID  string `json:"categoryId"`  // 分类ID
	Type        string `json:"type"`        // 模板类型
	Version     string `json:"version"`     // 模板版本
	Status      int    `json:"status"`      // 状态
	Conditions  string `json:"conditions"`  // 条件表达式
	LuaScript   string `json:"luaScript"`   // Lua脚本代码
	Formula     string `json:"formula"`     // 计算公式
	FormulaVars string `json:"formulaVars"` // 公式变量映射
	Parameters  string `json:"parameters"`  // 模板参数定义
	Priority    int32  `json:"priority"`    // 优先级
	Sorting     int32  `json:"sorting"`     // 排序权重
	CreatedAt   int64  `json:"createdAt"`   // 创建时间
	UpdatedAt   int64  `json:"updatedAt"`   // 更新时间
	TenantID    string `json:"tenantId"`    // 租户ID
}

// ==================== 分类相关DTO ====================

// CategoryDTO 分类数据传输对象
type CategoryDTO struct {
	ID           string `json:"id"`           // 分类ID
	Code         string `json:"code"`         // 分类编码
	Name         string `json:"name"`         // 分类名称
	Description  string `json:"description"`  // 分类描述
	ParentID     string `json:"parentId"`     // 父分类ID
	Type         string `json:"type"`         // 分类类型
	BusinessType string `json:"businessType"` // 业务类型
	Level        int32  `json:"level"`        // 层级
	Path         string `json:"path"`         // 路径
	IsLeaf       bool   `json:"isLeaf"`       // 是否叶子节点
	Sorting      int32  `json:"sorting"`      // 排序权重
	Status       int    `json:"status"`       // 状态
	CreatedAt    int64  `json:"createdAt"`    // 创建时间
	UpdatedAt    int64  `json:"updatedAt"`    // 更新时间
	TenantID     string `json:"tenantId"`     // 租户ID
}

// CategoryTreeDTO 分类树数据传输对象
type CategoryTreeDTO struct {
	CategoryDTO
	Children []*CategoryTreeDTO `json:"children"` // 子分类列表
}

// ==================== 规则相关DTO ====================

// RuleDTO 规则数据传输对象
type RuleDTO struct {
	ID              string        `json:"id"`              // 规则ID
	Code            string        `json:"code"`            // 规则编码
	Name            string        `json:"name"`            // 规则名称
	Description     string        `json:"description"`     // 规则描述
	CategoryID      string        `json:"categoryId"`      // 分类ID
	TemplateID      string        `json:"templateId"`      // 模板ID
	Type            string        `json:"type"`            // 规则类型
	Triggers        []string      `json:"triggers"`        // 触发条件
	ExecutionTiming string        `json:"executionTiming"` // 执行时机：before(前置) after(后置) both(前后都执行)
	Scope           string        `json:"scope"`           // 作用域
	ScopeID         string        `json:"scopeId"`         // 作用域ID（商品ID、用户ID、订单ID等）
	Condition       *ConditionDTO `json:"condition"`       // 条件配置
	LuaScript       string        `json:"luaScript"`       // Lua脚本
	Formula         string        `json:"formula"`         // 计算公式
	Action          string        `json:"action"`          // 触发动作：allow(允许) deny(拒绝) modify(修改) notify(通知) redirect(重定向)
	Priority        int32         `json:"priority"`        // 优先级
	Sorting         int32         `json:"sorting"`         // 排序权重
	Status          int32         `json:"status"`          // 状态
	CreatedAt       int64         `json:"createdAt"`       // 创建时间
	UpdatedAt       int64         `json:"updatedAt"`       // 更新时间
	TenantID        string        `json:"tenantId"`        // 租户ID
}

// ConditionDTO 条件配置数据传输对象
type ConditionDTO struct {
	Type       string                 `json:"type"`       // 条件类型
	Expression string                 `json:"expression"` // 条件表达式
	Parameters map[string]interface{} `json:"parameters"` // 条件参数
}

// ==================== 规则执行相关DTO ====================

// RuleExecutionResult 规则执行结果
type RuleExecutionResult struct {
	RuleID      string                 `json:"ruleId"`      // 规则ID
	RuleName    string                 `json:"ruleName"`    // 规则名称
	Matched     bool                   `json:"matched"`     // 是否匹配
	Result      map[string]interface{} `json:"result"`      // 执行结果
	ExecuteTime int64                  `json:"executeTime"` // 执行时间（毫秒）
	Error       string                 `json:"error"`       // 错误信息
	ExecuteAt   int64                  `json:"executeAt"`   // 执行时间戳
}

// ExecuteRulesResponse 执行规则响应
type ExecuteRulesResponse struct {
	Results          []RuleExecutionResult `json:"results"`          // 执行结果列表
	TotalExecuteTime int64                 `json:"totalExecuteTime"` // 总执行时间（毫秒）
	MatchedCount     int32                 `json:"matchedCount"`     // 匹配规则数量
	ExecuteAt        int64                 `json:"executeAt"`        // 执行时间戳
}

// ==================== 统计相关DTO ====================

// RuleStatisticsDTO 规则统计信息
type RuleStatisticsDTO struct {
	RuleID              string  `json:"ruleId"`              // 规则ID
	RuleName            string  `json:"ruleName"`            // 规则名称
	TotalExecuteCount   int64   `json:"totalExecuteCount"`   // 总执行次数
	SuccessExecuteCount int64   `json:"successExecuteCount"` // 成功执行次数
	FailedExecuteCount  int64   `json:"failedExecuteCount"`  // 失败执行次数
	SuccessRate         float64 `json:"successRate"`         // 成功率
	AverageExecuteTime  float64 `json:"averageExecuteTime"`  // 平均执行时间（毫秒）
	LastExecuteAt       int64   `json:"lastExecuteAt"`       // 最后执行时间
	TodayExecuteCount   int64   `json:"todayExecuteCount"`   // 今日执行次数
	WeekExecuteCount    int64   `json:"weekExecuteCount"`    // 本周执行次数
	MonthExecuteCount   int64   `json:"monthExecuteCount"`   // 本月执行次数
}

// GetRuleStatisticsReq 获取规则统计请求
type GetRuleStatisticsReq struct {
	RuleID         string `json:"ruleId" form:"ruleId" query:"ruleId"`                         // 规则ID
	StartTime      int64  `json:"startTime" form:"startTime" query:"startTime"`                // 开始时间
	EndTime        int64  `json:"endTime" form:"endTime" query:"endTime"`                      // 结束时间
	StatisticsType string `json:"statisticsType" form:"statisticsType" query:"statisticsType"` // 统计类型：daily(日统计) weekly(周统计) monthly(月统计)
}

// RuleStatisticsResponse 规则统计响应
type RuleStatisticsResponse struct {
	Statistics RuleStatisticsDTO `json:"statistics"` // 统计信息
	TimeSeries []TimeSeriesData  `json:"timeSeries"` // 时间序列数据
}

// TimeSeriesData 时间序列数据
type TimeSeriesData struct {
	Time               int64   `json:"time"`               // 时间点
	ExecuteCount       int64   `json:"executeCount"`       // 执行次数
	SuccessCount       int64   `json:"successCount"`       // 成功次数
	FailedCount        int64   `json:"failedCount"`        // 失败次数
	AverageExecuteTime float64 `json:"averageExecuteTime"` // 平均执行时间
}

// RuleResultDTO 规则执行结果数据传输对象
type RuleResultDTO struct {
	RuleID      string                 `json:"ruleId"`      // 规则ID
	RuleCode    string                 `json:"ruleCode"`    // 规则编码
	RuleName    string                 `json:"ruleName"`    // 规则名称
	Success     bool                   `json:"success"`     // 执行是否成功
	Result      map[string]interface{} `json:"result"`      // 执行结果
	Message     string                 `json:"message"`     // 执行消息
	ExecuteTime int64                  `json:"executeTime"` // 执行时间
	Duration    int64                  `json:"duration"`    // 执行耗时（毫秒）
}

// ExecuteRuleRequestDTO 执行规则请求数据传输对象
type ExecuteRuleRequestDTO struct {
	RuleID  string                 `json:"ruleId"`  // 规则ID
	Context map[string]interface{} `json:"context"` // 执行上下文
}

// ExecuteRuleByCodeRequestDTO 根据编码执行规则请求数据传输对象
type ExecuteRuleByCodeRequestDTO struct {
	Code    string                 `json:"code"`    // 规则编码
	Context map[string]interface{} `json:"context"` // 执行上下文
}

// ExecuteRulesByTriggerRequestDTO 根据触发条件执行规则请求数据传输对象
type ExecuteRulesByTriggerRequestDTO struct {
	Trigger string                 `json:"trigger"` // 触发条件
	Context map[string]interface{} `json:"context"` // 执行上下文
}
