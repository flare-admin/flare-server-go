package command

// ==================== 模板相关命令 ====================

// CreateTemplateCommand 创建模板命令
type CreateTemplateCommand struct {
	Name        string `json:"name" form:"name" query:"name"`                      // 模板名称
	Code        string `json:"code" form:"code" query:"code"`                      // 模板编码
	Description string `json:"description" form:"description" query:"description"` // 模板描述
	CategoryID  string `json:"categoryId" form:"categoryId" query:"categoryId"`    // 分类ID
	Type        string `json:"type" form:"type" query:"type"`                      // 模板类型
	Conditions  string `json:"conditions" form:"conditions" query:"conditions"`    // 条件表达式
	LuaScript   string `json:"luaScript" form:"luaScript" query:"luaScript"`       // Lua脚本代码
	Formula     string `json:"formula" form:"formula" query:"formula"`             // 计算公式
	FormulaVars string `json:"formulaVars" form:"formulaVars" query:"formulaVars"` // 公式变量映射
	Parameters  string `json:"parameters" form:"parameters" query:"parameters"`    // 模板参数定义
	Priority    int32  `json:"priority" form:"priority" query:"priority"`          // 优先级
	Sorting     int32  `json:"sorting" form:"sorting" query:"sorting"`             // 排序权重
}

// UpdateTemplateCommand 更新模板命令
type UpdateTemplateCommand struct {
	ID          string `json:"id" form:"id" query:"id"`                            // 模板ID
	Name        string `json:"name" form:"name" query:"name"`                      // 模板名称
	Code        string `json:"code" form:"code" query:"code"`                      // 模板编码
	Description string `json:"description" form:"description" query:"description"` // 模板描述
	CategoryID  string `json:"categoryId" form:"categoryId" query:"categoryId"`    // 分类ID
	Type        string `json:"type" form:"type" query:"type"`                      // 模板类型
	Conditions  string `json:"conditions" form:"conditions" query:"conditions"`    // 条件表达式
	LuaScript   string `json:"luaScript" form:"luaScript" query:"luaScript"`       // Lua脚本代码
	Formula     string `json:"formula" form:"formula" query:"formula"`             // 计算公式
	FormulaVars string `json:"formulaVars" form:"formulaVars" query:"formulaVars"` // 公式变量映射
	Parameters  string `json:"parameters" form:"parameters" query:"parameters"`    // 模板参数定义
	Priority    int32  `json:"priority" form:"priority" query:"priority"`          // 优先级
	Sorting     int32  `json:"sorting" form:"sorting" query:"sorting"`             // 排序权重
}

// UpdateTemplateStatusCommand 更新模板状态命令
type UpdateTemplateStatusCommand struct {
	ID     string `json:"id" form:"id" query:"id"`             // 模板ID
	Status int32  `json:"status" form:"status" query:"status"` // 状态
}

// DeleteTemplateCommand 删除模板命令
type DeleteTemplateCommand struct {
	ID string `json:"id" form:"id" query:"id"` // 模板ID
}

// ==================== 分类相关命令 ====================

// CreateCategoryCommand 创建分类命令
type CreateCategoryCommand struct {
	Name         string `json:"name" form:"name" query:"name"`                         // 分类名称
	Code         string `json:"code" form:"code" query:"code"`                         // 分类编码
	Description  string `json:"description" form:"description" query:"description"`    // 分类描述
	ParentID     string `json:"parentId" form:"parentId" query:"parentId"`             // 父分类ID
	Type         string `json:"type" form:"type" query:"type"`                         // 分类类型
	BusinessType string `json:"businessType" form:"businessType" query:"businessType"` // 业务类型
	Sorting      int32  `json:"sorting" form:"sorting" query:"sorting"`                // 排序权重
}

// UpdateCategoryCommand 更新分类命令
type UpdateCategoryCommand struct {
	ID           string `json:"id" form:"id" query:"id"`                               // 分类ID
	Name         string `json:"name" form:"name" query:"name"`                         // 分类名称
	Code         string `json:"code" form:"code" query:"code"`                         // 分类编码
	Description  string `json:"description" form:"description" query:"description"`    // 分类描述
	ParentID     string `json:"parentId" form:"parentId" query:"parentId"`             // 父分类ID
	Type         string `json:"type" form:"type" query:"type"`                         // 分类类型
	BusinessType string `json:"businessType" form:"businessType" query:"businessType"` // 业务类型
	Sorting      int32  `json:"sorting" form:"sorting" query:"sorting"`                // 排序权重
}

// UpdateCategoryStatusCommand 更新分类状态命令
type UpdateCategoryStatusCommand struct {
	ID     string `json:"id" form:"id" query:"id"`             // 分类ID
	Status int32  `json:"status" form:"status" query:"status"` // 状态
}

// DeleteCategoryCommand 删除分类命令
type DeleteCategoryCommand struct {
	ID string `json:"id" form:"id" query:"id"` // 分类ID
}

// ==================== 规则相关命令 ====================

// CreateRuleCommand 创建规则命令
type CreateRuleCommand struct {
	Name            string           `json:"name" form:"name" query:"name"`                                  // 规则名称
	Code            string           `json:"code" form:"code" query:"code"`                                  // 规则编码
	Description     string           `json:"description" form:"description" query:"description"`             // 规则描述
	CategoryID      string           `json:"categoryId" form:"categoryId" query:"categoryId"`                // 分类ID
	TemplateID      string           `json:"templateId" form:"templateId" query:"templateId"`                // 模板ID
	Type            string           `json:"type" form:"type" query:"type"`                                  // 规则类型
	Triggers        []string         `json:"triggers" form:"triggers" query:"triggers"`                      // 触发条件
	Scope           string           `json:"scope" form:"scope" query:"scope"`                               // 作用域
	ScopeID         string           `json:"scopeId" form:"scopeId" query:"scopeId"`                         // 作用域ID
	ExecutionTiming string           `json:"executionTiming" form:"executionTiming" query:"executionTiming"` // 执行时机：before(前置) after(后置) both(前后都执行)
	Condition       *ConditionConfig `json:"condition" form:"condition" query:"condition"`                   // 条件配置
	LuaScript       string           `json:"luaScript" form:"luaScript" query:"luaScript"`                   // Lua脚本
	Formula         string           `json:"formula" form:"formula" query:"formula"`                         // 计算公式
	Action          string           `json:"action" form:"action" query:"action"`                            // 触发动作：allow(允许) deny(拒绝) modify(修改) notify(通知) redirect(重定向)
	Priority        int32            `json:"priority" form:"priority" query:"priority"`                      // 优先级
	Sorting         int32            `json:"sorting" form:"sorting" query:"sorting"`                         // 排序权重
}

// UpdateRuleCommand 更新规则命令
type UpdateRuleCommand struct {
	ID              string           `json:"id" form:"id" query:"id"`                                        // 规则ID
	Name            string           `json:"name" form:"name" query:"name"`                                  // 规则名称
	Code            string           `json:"code" form:"code" query:"code"`                                  // 规则编码
	Description     string           `json:"description" form:"description" query:"description"`             // 规则描述
	CategoryID      string           `json:"categoryId" form:"categoryId" query:"categoryId"`                // 分类ID
	TemplateID      string           `json:"templateId" form:"templateId" query:"templateId"`                // 模板ID
	Type            string           `json:"type" form:"type" query:"type"`                                  // 规则类型
	Triggers        []string         `json:"triggers" form:"triggers" query:"triggers"`                      // 触发条件
	Scope           string           `json:"scope" form:"scope" query:"scope"`                               // 作用域
	ScopeID         string           `json:"scopeId" form:"scopeId" query:"scopeId"`                         // 作用域ID
	ExecutionTiming string           `json:"executionTiming" form:"executionTiming" query:"executionTiming"` // 执行时机：before(前置) after(后置) both(前后都执行)
	Condition       *ConditionConfig `json:"condition" form:"condition" query:"condition"`                   // 条件配置
	LuaScript       string           `json:"luaScript" form:"luaScript" query:"luaScript"`                   // Lua脚本
	Formula         string           `json:"formula" form:"formula" query:"formula"`                         // 计算公式
	Action          string           `json:"action" form:"action" query:"action"`                            // 触发动作：allow(允许) deny(拒绝) modify(修改) notify(通知) redirect(重定向)
	Priority        int32            `json:"priority" form:"priority" query:"priority"`                      // 优先级
	Sorting         int32            `json:"sorting" form:"sorting" query:"sorting"`                         // 排序权重
}

// UpdateRuleStatusCommand 更新规则状态命令
type UpdateRuleStatusCommand struct {
	ID     string `json:"id" form:"id" query:"id"`             // 规则ID
	Status int32  `json:"status" form:"status" query:"status"` // 状态
}

// DeleteRuleCommand 删除规则命令
type DeleteRuleCommand struct {
	ID string `json:"id" form:"id" query:"id"` // 规则ID
}

// ExecuteRuleCommand 执行规则命令
type ExecuteRuleCommand struct {
	RuleID  string                 `json:"ruleId" form:"ruleId" query:"ruleId"`    // 规则ID
	Context map[string]interface{} `json:"context" form:"context" query:"context"` // 执行上下文
}

// ExecuteRuleByCodeCommand 根据编码执行规则命令
type ExecuteRuleByCodeCommand struct {
	Code    string                 `json:"code" form:"code" query:"code"`          // 规则编码
	Context map[string]interface{} `json:"context" form:"context" query:"context"` // 执行上下文
}

// ExecuteRulesByTriggerCommand 根据触发条件执行规则命令
type ExecuteRulesByTriggerCommand struct {
	Trigger string                 `json:"trigger" form:"trigger" query:"trigger"` // 触发条件
	Context map[string]interface{} `json:"context" form:"context" query:"context"` // 执行上下文
}

// ==================== 配置对象 ====================

// ConditionConfig 条件配置
type ConditionConfig struct {
	Type       string                 `json:"type" form:"type" query:"type"`                   // 条件类型
	Expression string                 `json:"expression" form:"expression" query:"expression"` // 条件表达式
	Parameters map[string]interface{} `json:"parameters" form:"parameters" query:"parameters"` // 条件参数
}
