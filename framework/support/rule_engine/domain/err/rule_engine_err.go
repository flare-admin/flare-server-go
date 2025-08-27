package err

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
)

var (
	// ==================== 模板相关错误 ====================

	// RuleTemplateCreateFailed 创建规则模板失败
	RuleTemplateCreateFailed = herrors.NewServerError("RuleTemplateCreateFailed")
	// RuleTemplateUpdateFailed 更新规则模板失败
	RuleTemplateUpdateFailed = herrors.NewServerError("RuleTemplateUpdateFailed")
	// RuleTemplateDeleteFailed 删除规则模板失败
	RuleTemplateDeleteFailed = herrors.NewServerError("RuleTemplateDeleteFailed")
	// RuleTemplateGetFailed 获取规则模板失败
	RuleTemplateGetFailed = herrors.NewServerError("RuleTemplateGetFailed")
	// RuleTemplateDisabled 规则模板已禁用
	RuleTemplateDisabled = herrors.NewBusinessServerError("RuleTemplateDisabled")
	// RuleTemplateNotExist 规则模板不存在
	RuleTemplateNotExist = herrors.NewBusinessServerError("RuleTemplateNotExist")
	// RuleTemplateCodeExists 规则模板编码已存在
	RuleTemplateCodeExists = herrors.NewBusinessServerError("RuleTemplateCodeExists")
	// RuleTemplateValidationFailed 规则模板数据验证失败
	RuleTemplateValidationFailed = herrors.NewServerError("RuleTemplateValidationFailed")
	// RuleTemplateTypeInvalid 规则模板类型无效
	RuleTemplateTypeInvalid = herrors.NewBusinessServerError("RuleTemplateTypeInvalid")
	// RuleTemplateContentInvalid 规则模板内容无效
	RuleTemplateContentInvalid = herrors.NewBusinessServerError("RuleTemplateContentInvalid")

	// ==================== 规则相关错误 ====================

	// RuleCreateFailed 创建规则失败
	RuleCreateFailed = herrors.NewServerError("RuleCreateFailed")
	// RuleUpdateFailed 更新规则失败
	RuleUpdateFailed = herrors.NewServerError("RuleUpdateFailed")
	// RuleDeleteFailed 删除规则失败
	RuleDeleteFailed = herrors.NewServerError("RuleDeleteFailed")
	// RuleGetFailed 获取规则失败
	RuleGetFailed = herrors.NewServerError("RuleGetFailed")
	// RuleDisabled 规则已禁用
	RuleDisabled = herrors.NewBusinessServerError("RuleDisabled")
	// RuleNotExist 规则不存在
	RuleNotExist = herrors.NewBusinessServerError("RuleNotExist")
	// RuleCodeExists 规则编码已存在
	RuleCodeExists = herrors.NewBusinessServerError("RuleCodeExists")
	// RuleValidationFailed 规则数据验证失败
	RuleValidationFailed = herrors.NewServerError("RuleValidationFailed")
	// RuleTypeInvalid 规则类型无效
	RuleTypeInvalid = herrors.NewBusinessServerError("RuleTypeInvalid")
	// RuleScopeInvalid 规则作用域无效
	RuleScopeInvalid = herrors.NewBusinessServerError("RuleScopeInvalid")
	// RuleActionInvalid 规则动作无效
	RuleActionInvalid = herrors.NewBusinessServerError("RuleActionInvalid")
	// RuleTriggerInvalid 规则触发条件无效
	RuleTriggerInvalid = herrors.NewBusinessServerError("RuleTriggerInvalid")
	// RuleContentInvalid 规则内容无效
	RuleContentInvalid = herrors.NewBusinessServerError("RuleContentInvalid")
	// RuleExecutionFailed 规则执行失败
	RuleExecutionFailed = herrors.NewServerError("RuleExecutionFailed")
	// RuleConditionNotMatch 规则条件不匹配
	RuleConditionNotMatch = herrors.NewBusinessServerError("RuleConditionNotMatch")
	// RuleLuaScriptError Lua脚本执行错误
	RuleLuaScriptError = herrors.NewServerError("RuleLuaScriptError")
	// RuleFormulaError 公式计算错误
	RuleFormulaError = herrors.NewServerError("RuleFormulaError")

	// ==================== 分类相关错误 ====================

	// RuleCategoryCreateFailed 创建规则分类失败
	RuleCategoryCreateFailed = herrors.NewServerError("RuleCategoryCreateFailed")
	// RuleCategoryUpdateFailed 更新规则分类失败
	RuleCategoryUpdateFailed = herrors.NewServerError("RuleCategoryUpdateFailed")
	// RuleCategoryDeleteFailed 删除规则分类失败
	RuleCategoryDeleteFailed = herrors.NewServerError("RuleCategoryDeleteFailed")
	// RuleCategoryGetFailed 获取规则分类失败
	RuleCategoryGetFailed = herrors.NewServerError("RuleCategoryGetFailed")
	// RuleCategoryDisabled 规则分类已禁用
	RuleCategoryDisabled = herrors.NewBusinessServerError("RuleCategoryDisabled")
	// RuleCategoryNotExist 规则分类不存在
	RuleCategoryNotExist = herrors.NewBusinessServerError("RuleCategoryNotExist")
	// RuleCategoryCodeExists 规则分类编码已存在
	RuleCategoryCodeExists = herrors.NewBusinessServerError("RuleCategoryCodeExists")
	// RuleCategoryValidationFailed 规则分类数据验证失败
	RuleCategoryValidationFailed = herrors.NewServerError("RuleCategoryValidationFailed")
	// RuleCategoryTypeInvalid 规则分类类型无效
	RuleCategoryTypeInvalid = herrors.NewBusinessServerError("RuleCategoryTypeInvalid")
	// RuleCategoryBusinessTypeInvalid 规则分类业务类型无效
	RuleCategoryBusinessTypeInvalid = herrors.NewBusinessServerError("RuleCategoryBusinessTypeInvalid")
	// RuleCategoryHasChildren 规则分类存在子分类
	RuleCategoryHasChildren = herrors.NewBusinessServerError("RuleCategoryHasChildren")
	// RuleCategoryHasRules 规则分类存在规则
	RuleCategoryHasRules = herrors.NewBusinessServerError("RuleCategoryHasRules")
	// RuleCategoryHasTemplates 规则分类存在模板
	RuleCategoryHasTemplates = herrors.NewBusinessServerError("RuleCategoryHasTemplates")

	// ==================== 规则引擎相关错误 ====================

	// RuleEngineExecuteFailed 规则引擎执行失败
	RuleEngineExecuteFailed = herrors.NewServerError("RuleEngineExecuteFailed")
	// RuleEngineContextInvalid 规则引擎上下文无效
	RuleEngineContextInvalid = herrors.NewBusinessServerError("RuleEngineContextInvalid")
	// RuleContextInvalid 规则上下文无效
	RuleContextInvalid = herrors.NewServerError("RuleContextInvalid")
	// RuleEngineNoRulesFound 未找到匹配的规则
	RuleEngineNoRulesFound = herrors.NewBusinessServerError("RuleEngineNoRulesFound")
	// RuleEngineTimeout 规则引擎执行超时
	RuleEngineTimeout = herrors.NewServerError("RuleEngineTimeout")
	// RuleEngineCacheError 规则引擎缓存错误
	RuleEngineCacheError = herrors.NewServerError("RuleEngineCacheError")
	// RuleEngineStatisticsError 规则引擎统计错误
	RuleEngineStatisticsError = herrors.NewServerError("RuleEngineStatisticsError")
)
