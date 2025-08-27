package utils

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

// IsValidJSON 检查字符串是否是有效的 JSON
func IsValidJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

// ConvertToJSONString 将格式化字符串转换为 JSON 字符串，所有值为字符串
func ConvertToJSONString(rawString string) (string, error) {
	// 如果是有效的 JSON，直接返回
	if IsValidJSON(rawString) {
		return rawString, nil
	}

	// 去掉大括号并分割
	rawString = strings.Trim(rawString, "{}")
	pairs := strings.Split(rawString, ",")

	// 用于存储转换后的键值对
	jsonMap := make(map[string]interface{})

	// 正则表达式匹配键值对
	re := regexp.MustCompile(`(?P<key>[^:]+):(?P<value>[^,]*)`)

	// 处理每一对
	for _, pair := range pairs {
		matches := re.FindStringSubmatch(pair)
		if len(matches) == 3 {
			key := strings.TrimSpace(matches[1])
			value := strings.TrimSpace(matches[2])

			// 转换值为字符串，处理空值
			if value == "" {
				value = ""
				// 存储到 map 中
				jsonMap[key] = value
			} else if key == "tradeType" || key == "status" {
				num, err := strconv.Atoi(value)
				if err != nil {
					jsonMap[key] = value
				}
				jsonMap[key] = num
			} else {
				// 不再添加额外的引号
				value = "" + value + "" // 只在前后添加引号
				// 存储到 map 中
				jsonMap[key] = value
			}
		}
	}
	// 转换为 JSON 字符串
	jsonData, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
