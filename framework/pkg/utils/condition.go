package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ConditionOperator 条件运算符
const (
	ConditionOperatorEQ      = "eq"      // 等于
	ConditionOperatorNEQ     = "neq"     // 不等于
	ConditionOperatorGT      = "gt"      // 大于
	ConditionOperatorGTE     = "gte"     // 大于等于
	ConditionOperatorLT      = "lt"      // 小于
	ConditionOperatorLTE     = "lte"     // 小于等于
	ConditionOperatorBETWEEN = "between" // 区间
	ConditionOperatorIN      = "in"      // 包含
)

// ConditionEvaluator 条件评估器
type ConditionEvaluator struct {
	operator string
	value    string
}

// NewConditionEvaluator 创建条件评估器
func NewConditionEvaluator(operator, value string) *ConditionEvaluator {
	return &ConditionEvaluator{
		operator: operator,
		value:    value,
	}
}

// Evaluate 评估条件
func (e *ConditionEvaluator) Evaluate(actual interface{}) bool {
	switch e.operator {
	case ConditionOperatorEQ:
		return e.evaluateEqual(actual)
	case ConditionOperatorNEQ:
		return !e.evaluateEqual(actual)
	case ConditionOperatorGT:
		return e.compareValues(actual, e.value) > 0
	case ConditionOperatorGTE:
		return e.compareValues(actual, e.value) >= 0
	case ConditionOperatorLT:
		return e.compareValues(actual, e.value) < 0
	case ConditionOperatorLTE:
		return e.compareValues(actual, e.value) <= 0
	case ConditionOperatorBETWEEN:
		return e.evaluateBetween(actual)
	case ConditionOperatorIN:
		return e.evaluateIn(actual)
	default:
		return false
	}
}

// evaluateEqual 评估等于
func (e *ConditionEvaluator) evaluateEqual(actual interface{}) bool {
	return fmt.Sprint(actual) == e.value
}

// evaluateBetween 评估区间
func (e *ConditionEvaluator) evaluateBetween(actual interface{}) bool {
	values := strings.Split(e.value, ",")
	if len(values) != 2 {
		return false
	}

	min, err1 := strconv.ParseFloat(strings.TrimSpace(values[0]), 64)
	max, err2 := strconv.ParseFloat(strings.TrimSpace(values[1]), 64)
	if err1 != nil || err2 != nil {
		return false
	}

	actualFloat, err := e.toFloat64(actual)
	if err != nil {
		return false
	}

	return actualFloat >= min && actualFloat <= max
}

// evaluateIn 评估包含
func (e *ConditionEvaluator) evaluateIn(actual interface{}) bool {
	values := strings.Split(e.value, ",")
	actualStr := fmt.Sprint(actual)
	for _, v := range values {
		if strings.TrimSpace(v) == actualStr {
			return true
		}
	}
	return false
}

// compareValues 比较值
func (e *ConditionEvaluator) compareValues(actual interface{}, target string) int {
	actualFloat, err1 := e.toFloat64(actual)
	if err1 != nil {
		return 0
	}

	targetFloat, err2 := strconv.ParseFloat(target, 64)
	if err2 != nil {
		return 0
	}

	if actualFloat < targetFloat {
		return -1
	} else if actualFloat > targetFloat {
		return 1
	}
	return 0
}

// toFloat64 转换为float64
func (e *ConditionEvaluator) toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("unsupported type: %T", value)
	}
}

// EvaluateCondition 评估条件的便捷方法
func EvaluateCondition(operator, targetValue string, actualValue interface{}) bool {
	evaluator := NewConditionEvaluator(operator, targetValue)
	return evaluator.Evaluate(actualValue)
}
