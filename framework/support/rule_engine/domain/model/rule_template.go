package model

import (
	"encoding/json"
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
)

// RuleTemplate 规则模板领域模型
type RuleTemplate struct {
	// 基础信息
	ID          string `json:"id"`          // 模板ID
	Code        string `json:"code"`        // 模板编码
	Name        string `json:"name"`        // 模板名称
	Description string `json:"description"` // 模板描述
	CategoryID  string `json:"categoryId"`  // 分类ID

	// 模板配置
	Type    string `json:"type"`    // 模板类型：condition(条件模板) lua(lua脚本模板) formula(公式模板)
	Version string `json:"version"` // 模板版本
	Status  int32  `json:"status"`  // 状态：1-启用 2-禁用

	// 模板内容
	Conditions  string `json:"conditions"`  // 条件表达式(JSON格式)
	LuaScript   string `json:"luaScript"`   // Lua脚本代码
	Formula     string `json:"formula"`     // 计算公式
	FormulaVars string `json:"formulaVars"` // 公式变量映射(JSON格式)

	// 模板参数
	Parameters string `json:"parameters"` // 模板参数定义(JSON格式)

	// 优先级和排序
	Priority int32 `json:"priority"` // 优先级，数字越大优先级越高
	Sorting  int32 `json:"sorting"`  // 排序权重

	// 时间信息
	CreatedAt int64 `json:"createdAt"` // 创建时间
	UpdatedAt int64 `json:"updatedAt"` // 更新时间

	// 租户信息
	TenantID string `json:"tenantId"` // 租户ID
}

// NewRuleTemplate 创建规则模板
func NewRuleTemplate(code, name, description, categoryID, templateType string) *RuleTemplate {
	now := utils.GetDateUnix()
	return &RuleTemplate{
		ID:          "",
		Code:        code,
		Name:        name,
		Description: description,
		CategoryID:  categoryID,
		Type:        templateType,
		Version:     "1.0.0",
		Status:      1,
		Conditions:  "{}",
		LuaScript:   "",
		Formula:     "{}",
		FormulaVars: "{}",
		Parameters:  "{}",
		Priority:    0,
		Sorting:     0,
		CreatedAt:   now,
		UpdatedAt:   now,
		TenantID:    "",
	}
}

func (rt *RuleTemplate) Completion() {
	if rt.Conditions == "" {
		rt.Conditions = "{}"
	}
	if rt.LuaScript == "" {
		rt.LuaScript = "return true"
	}
	if rt.Formula == "" {
		rt.Formula = "{}"
	}
	if rt.FormulaVars == "" {
		rt.FormulaVars = "{}"
	}
	if rt.Parameters == "" {
		rt.Parameters = "{}"
	}
}

// SetConditions 设置条件表达式
func (rt *RuleTemplate) SetConditions(conditions map[string]interface{}) error {
	// 验证JSON格式
	conditionsJSON, err := json.Marshal(conditions)
	if err != nil {
		return fmt.Errorf("invalid conditions format: %v", err)
	}

	rt.Conditions = string(conditionsJSON)
	return nil
}

// GetConditions 获取条件表达式
func (rt *RuleTemplate) GetConditions() (map[string]interface{}, error) {
	var conditions map[string]interface{}
	if rt.Conditions == "" {
		return make(map[string]interface{}), nil
	}

	err := json.Unmarshal([]byte(rt.Conditions), &conditions)
	if err != nil {
		return nil, fmt.Errorf("invalid conditions format: %v", err)
	}

	return conditions, nil
}

// SetLuaScript 设置Lua脚本
func (rt *RuleTemplate) SetLuaScript(script string) error {
	if rt.Type != "lua" {
		return fmt.Errorf("template type is not lua")
	}

	rt.LuaScript = script
	return nil
}

// SetFormula 设置计算公式
func (rt *RuleTemplate) SetFormula(formula string, vars map[string]interface{}) error {
	if rt.Type != "formula" {
		return fmt.Errorf("template type is not formula")
	}

	rt.Formula = formula

	// 序列化变量映射
	varsJSON, err := json.Marshal(vars)
	if err != nil {
		return fmt.Errorf("invalid formula vars format: %v", err)
	}

	rt.FormulaVars = string(varsJSON)
	return nil
}

// GetFormulaVars 获取公式变量
func (rt *RuleTemplate) GetFormulaVars() (map[string]interface{}, error) {
	var vars map[string]interface{}
	if rt.FormulaVars == "" {
		return make(map[string]interface{}), nil
	}

	err := json.Unmarshal([]byte(rt.FormulaVars), &vars)
	if err != nil {
		return nil, fmt.Errorf("invalid formula vars format: %v", err)
	}

	return vars, nil
}

// SetParameters 设置模板参数
func (rt *RuleTemplate) SetParameters(parameters map[string]interface{}) error {
	// 验证JSON格式
	parametersJSON, err := json.Marshal(parameters)
	if err != nil {
		return fmt.Errorf("invalid parameters format: %v", err)
	}

	rt.Parameters = string(parametersJSON)
	return nil
}

// GetParameters 获取模板参数
func (rt *RuleTemplate) GetParameters() (map[string]interface{}, error) {
	var parameters map[string]interface{}
	if rt.Parameters == "" {
		return make(map[string]interface{}), nil
	}

	err := json.Unmarshal([]byte(rt.Parameters), &parameters)
	if err != nil {
		return nil, fmt.Errorf("invalid parameters format: %v", err)
	}

	return parameters, nil
}

// Enable 启用模板
func (rt *RuleTemplate) Enable() {
	rt.Status = 1
	rt.UpdatedAt = utils.GetDateUnix()
}

// Disable 禁用模板
func (rt *RuleTemplate) Disable() {
	rt.Status = 2
	rt.UpdatedAt = utils.GetDateUnix()
}

// IsEnabled 是否启用
func (rt *RuleTemplate) IsEnabled() bool {
	return rt.Status == 1
}

// Validate 验证模板
func (rt *RuleTemplate) Validate() error {
	if rt.Code == "" {
		return fmt.Errorf("template code cannot be empty")
	}

	if rt.Name == "" {
		return fmt.Errorf("template name cannot be empty")
	}

	//if rt.CategoryID == "" {
	//	return fmt.Errorf("category ID cannot be empty")
	//}

	// 验证模板类型
	validTypes := []string{"condition", "lua", "formula"}
	isValidType := false
	for _, validType := range validTypes {
		if rt.Type == validType {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("invalid template type: %s", rt.Type)
	}

	// 根据类型验证内容
	switch rt.Type {
	case "condition":
		if _, err := rt.GetConditions(); err != nil {
			return fmt.Errorf("invalid conditions: %v", err)
		}
	case "lua":
		if rt.LuaScript == "" {
			return fmt.Errorf("lua script cannot be empty for lua template")
		}
	case "formula":
		if rt.Formula == "" {
			return fmt.Errorf("formula cannot be empty for formula template")
		}
		if _, err := rt.GetFormulaVars(); err != nil {
			return fmt.Errorf("invalid formula vars: %v", err)
		}
	}

	// 验证参数格式
	if _, err := rt.GetParameters(); err != nil {
		return fmt.Errorf("invalid parameters: %v", err)
	}

	return nil
}

// Update 更新模板
func (rt *RuleTemplate) Update(name, description string) {
	rt.Name = name
	rt.Description = description
	rt.UpdatedAt = utils.GetDateUnix()
}

// SetPriority 设置优先级
func (rt *RuleTemplate) SetPriority(priority int32) {
	rt.Priority = priority
	rt.UpdatedAt = utils.GetDateUnix()
}

// SetSorting 设置排序
func (rt *RuleTemplate) SetSorting(sorting int32) {
	rt.Sorting = sorting
	rt.UpdatedAt = utils.GetDateUnix()
}

// GetTemplateContent 获取模板内容
func (rt *RuleTemplate) GetTemplateContent() map[string]interface{} {
	content := make(map[string]interface{})

	switch rt.Type {
	case "condition":
		if conditions, err := rt.GetConditions(); err == nil {
			content["conditions"] = conditions
		}
	case "lua":
		content["luaScript"] = rt.LuaScript
	case "formula":
		content["formula"] = rt.Formula
		if vars, err := rt.GetFormulaVars(); err == nil {
			content["formulaVars"] = vars
		}
	}

	// 添加参数
	if parameters, err := rt.GetParameters(); err == nil {
		content["parameters"] = parameters
	}

	return content
}

// ApplyParameters 应用参数到模板内容
func (rt *RuleTemplate) ApplyParameters(params map[string]interface{}) (map[string]interface{}, error) {
	content := rt.GetTemplateContent()

	// 获取模板参数定义
	templateParams, err := rt.GetParameters()
	if err != nil {
		return nil, fmt.Errorf("failed to get template parameters: %v", err)
	}

	// 应用参数到内容中
	for paramName, paramValue := range params {
		if _, exists := templateParams[paramName]; exists {
			// 递归替换内容中的参数占位符
			replacedContent := replaceParameterInContent(content, paramName, paramValue)
			// 类型断言，确保返回正确的类型
			if contentMap, ok := replacedContent.(map[string]interface{}); ok {
				content = contentMap
			} else {
				return nil, fmt.Errorf("failed to apply parameter %s: invalid content type", paramName)
			}
		}
	}

	return content, nil
}

// replaceParameterInContent 递归替换内容中的参数占位符
func replaceParameterInContent(content interface{}, paramName string, paramValue interface{}) interface{} {
	switch v := content.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			if key == paramName {
				result[key] = paramValue
			} else {
				result[key] = replaceParameterInContent(value, paramName, paramValue)
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, value := range v {
			result[i] = replaceParameterInContent(value, paramName, paramValue)
		}
		return result
	case string:
		// 简单的字符串替换，可以根据需要扩展为更复杂的模板引擎
		if v == fmt.Sprintf("{{%s}}", paramName) {
			return paramValue
		}
		return v
	default:
		return v
	}
}
