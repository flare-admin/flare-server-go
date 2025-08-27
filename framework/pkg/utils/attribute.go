package utils

import (
	"fmt"
	"strconv"
	"time"
)

// AttributeType 属性类型
const (
	AttributeTypeString   = "string"   // 字符串
	AttributeTypeNumber   = "number"   // 数字
	AttributeTypeBoolean  = "boolean"  // 布尔值
	AttributeTypeDate     = "date"     // 日期
	AttributeTypeDateTime = "datetime" // 日期时间
	AttributeTypeTime     = "time"     // 时间
)

// 错误常量
const (
	ErrTypeMismatchString  = "TYPE_MISMATCH_STRING"   // 类型不匹配：期望string
	ErrTypeMismatchNumber  = "TYPE_MISMATCH_NUMBER"   // 类型不匹配：期望number
	ErrTypeMismatchBoolean = "TYPE_MISMATCH_BOOLEAN"  // 类型不匹配：期望boolean
	ErrTypeNotSupport      = "TYPE_NOT_SUPPORT"       // 不支持的属性类型
	ErrTimestampInvalid    = "TIMESTAMP_INVALID"      // 时间戳无效
	ErrTimestampOutOfRange = "TIMESTAMP_OUT_OF_RANGE" // 时间戳超出范围
)

// ValidateAttributeType 验证属性类型
func ValidateAttributeType(attrType string, value interface{}) error {
	switch attrType {
	case AttributeTypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf(ErrTypeMismatchString)
		}
	case AttributeTypeNumber:
		if _, ok := value.(float64); !ok {
			return fmt.Errorf(ErrTypeMismatchNumber)
		}
	case AttributeTypeBoolean:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf(ErrTypeMismatchBoolean)
		}
	case AttributeTypeDate:
		// 验证是否为时间戳
		timestamp, err := toTimestamp(value)
		if err != nil {
			return fmt.Errorf(ErrTimestampInvalid)
		}
		// 验证时间戳是否在合理范围内
		if timestamp < 0 || timestamp > GetDateUnix()*2 {
			return fmt.Errorf(ErrTimestampOutOfRange)
		}
	case AttributeTypeDateTime:
		// 验证是否为时间戳
		timestamp, err := toTimestamp(value)
		if err != nil {
			return fmt.Errorf(ErrTimestampInvalid)
		}
		// 验证时间戳是否在合理范围内
		if timestamp < 0 || timestamp > GetDateUnix()*2 {
			return fmt.Errorf(ErrTimestampOutOfRange)
		}
	case AttributeTypeTime:
		// 验证是否为时间戳
		timestamp, err := toTimestamp(value)
		if err != nil {
			return fmt.Errorf(ErrTimestampInvalid)
		}
		// 验证时间戳是否在合理范围内
		if timestamp < 0 || timestamp > GetDateUnix()*2 {
			return fmt.Errorf(ErrTimestampOutOfRange)
		}
	default:
		return fmt.Errorf(ErrTypeNotSupport)
	}
	return nil
}

// toTimestamp 转换为时间戳
func toTimestamp(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		// 尝试解析为时间戳
		timestamp, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return timestamp, nil
		}
		// 尝试解析为日期时间字符串
		t, err := time.Parse(time.RFC3339, v)
		if err == nil {
			return t.Unix(), nil
		}
		// 尝试解析为日期字符串
		t, err = time.Parse("2006-01-02", v)
		if err == nil {
			return t.Unix(), nil
		}
		return 0, fmt.Errorf(ErrTimestampInvalid)
	default:
		return 0, fmt.Errorf(ErrTimestampInvalid)
	}
}
