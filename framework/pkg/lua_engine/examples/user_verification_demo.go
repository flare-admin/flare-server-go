package examples

import (
	"fmt"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"time"

	"github.com/flare-admin/flare-server-go/framework/pkg/lua_engine"
)

// UserVerificationDemo 用户实名验证演示
func UserVerificationDemo() {
	executor := lua_engine.NewRuleExecutor()

	// 测试用例1：正常用户
	normalUser := map[string]interface{}{
		"user": map[string]interface{}{
			"id":          "12345",
			"nickname":    "张三",
			"phone":       "13800138000",
			"account":     "zhangsan",
			"email":       "zhangsan@example.com",
			"age":         25,
			"gender":      2,
			"personId":    "110101199001011234",
			"isReal":      1,
			"realAt":      utils.GetDateUnix() - 86400, // 1天前实名
			"level":       2,
			"testAccount": 2,
			"status":      1,
		},
	}

	// 测试用例2：未实名用户
	unverifiedUser := map[string]interface{}{
		"user": map[string]interface{}{
			"id":          "12346",
			"nickname":    "李四",
			"phone":       "13800138001",
			"account":     "lisi",
			"email":       "lisi@example.com",
			"age":         20,
			"gender":      1,
			"personId":    "",
			"isReal":      0,
			"realAt":      0,
			"level":       0,
			"testAccount": 2,
			"status":      1,
		},
	}

	// 测试用例3：年龄不足用户
	youngUser := map[string]interface{}{
		"user": map[string]interface{}{
			"id":          "12347",
			"nickname":    "王五",
			"phone":       "13800138002",
			"account":     "wangwu",
			"email":       "wangwu@example.com",
			"age":         16,
			"gender":      2,
			"personId":    "110101200801011234",
			"isReal":      1,
			"realAt":      utils.GetDateUnix() - 86400,
			"level":       1,
			"testAccount": 2,
			"status":      1,
		},
	}

	// 测试用例4：测试账号
	testUser := map[string]interface{}{
		"user": map[string]interface{}{
			"id":          "12348",
			"nickname":    "测试用户",
			"phone":       "13800138003",
			"account":     "testuser",
			"email":       "test@example.com",
			"age":         25,
			"gender":      2,
			"personId":    "110101199001011235",
			"isReal":      1,
			"realAt":      utils.GetDateUnix() - 86400,
			"level":       1,
			"testAccount": 1, // 测试账号
			"status":      1,
		},
	}

	// 读取Lua脚本
	script := `-- 这里应该是完整的Lua脚本内容，为了演示简化处理`

	// 执行测试用例
	testCases := []struct {
		name    string
		context map[string]interface{}
	}{
		{"正常用户", normalUser},
		{"未实名用户", unverifiedUser},
		{"年龄不足用户", youngUser},
		{"测试账号", testUser},
	}

	for _, tc := range testCases {
		fmt.Printf("\n=== 测试用例: %s ===\n", tc.name)

		opts := &lua_engine.ExecuteOptions{
			Timeout: 5 * time.Second,
			Context: tc.context,
		}

		result, err := executor.Execute(script, opts)
		if err != nil {
			fmt.Printf("执行错误: %v\n", err)
			continue
		}

		// 输出结果
		fmt.Printf("验证结果: %t\n", result.Valid)
		fmt.Printf("执行动作: %s\n", result.Action)
		fmt.Printf("执行时间: %dms\n", result.ExecuteTime)

		if result.Error != "" {
			fmt.Printf("错误信息: %s\n", result.Error)
			fmt.Printf("错误原因: %s\n", result.ErrorReason)
		}

		if result.Variables != nil {
			fmt.Printf("输出变量: %+v\n", result.Variables)
		}
	}
}

// 国际化错误信息映射示例
var ErrorMessages = map[string]map[string]string{
	"zh-CN": {
		"user.not_found":                "用户信息不存在",
		"user.real_name_required":       "用户未完成实名认证",
		"user.real_name_expired":        "实名认证时间超过7天限制",
		"user.status_invalid":           "用户状态异常",
		"user.age_not_qualified":        "用户年龄未满18岁",
		"user.id_card_required":         "用户未提供身份证号",
		"user.id_card_format_invalid":   "身份证号格式不正确",
		"user.phone_required":           "用户未提供手机号",
		"user.phone_format_invalid":     "手机号格式不正确",
		"user.test_account_not_allowed": "测试账号不能进行实名操作",
		"rule.missing_valid_field":      "规则必须返回布尔类型的valid变量",
	},
	"en-US": {
		"user.not_found":                "User information not found",
		"user.real_name_required":       "User has not completed real-name verification",
		"user.real_name_expired":        "Real-name verification time exceeds 7-day limit",
		"user.status_invalid":           "User status is abnormal",
		"user.age_not_qualified":        "User age is under 18",
		"user.id_card_required":         "User has not provided ID card number",
		"user.id_card_format_invalid":   "ID card number format is incorrect",
		"user.phone_required":           "User has not provided phone number",
		"user.phone_format_invalid":     "Phone number format is incorrect",
		"user.test_account_not_allowed": "Test accounts cannot perform real-name operations",
		"rule.missing_valid_field":      "Rule must return a boolean valid variable",
	},
}

// GetLocalizedMessage 获取本地化错误信息
func GetLocalizedMessage(locale, errorReason string) string {
	if messages, ok := ErrorMessages[locale]; ok {
		if message, ok := messages[errorReason]; ok {
			return message
		}
	}
	return errorReason // 如果找不到对应的消息，返回错误原因
}

// ProcessExecuteResult 处理执行结果并支持国际化
func ProcessExecuteResult(result *lua_engine.ExecuteResult, locale string) {
	fmt.Printf("验证结果: %t\n", result.Valid)
	fmt.Printf("执行动作: %s\n", result.Action)
	fmt.Printf("执行时间: %dms\n", result.ExecuteTime)

	if result.Error != "" {
		// 使用本地化错误信息
		localizedError := GetLocalizedMessage(locale, result.ErrorReason)
		fmt.Printf("错误信息: %s\n", localizedError)
		fmt.Printf("错误原因: %s\n", result.ErrorReason)
	}

	if result.Variables != nil {
		fmt.Printf("输出变量: %+v\n", result.Variables)
	}
}
