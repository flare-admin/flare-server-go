package utils

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// ToDBMoney 转成数据库的金额计算
func ToDBMoney(money float64) int64 {
	if money == 0 {
		return 0
	}
	return int64(money * 100)
}

// ToDtoMoney 转成数据中的金额
func ToDtoMoney(money int64) string {
	if money == 0 {
		return "0.00"
	}
	return fmt.Sprintf("%.2f", PointsToYuan(money))
}
func ReplaceDecimal(money string, replacement string, size int) string {
	if len(replacement) > size {
		replacement = replacement[:size]
	}
	// 检查字符串中是否包含小数点
	if strings.Contains(money, ".") {
		// 找到小数点的位置
		dotIndex := strings.Index(money, ".")
		// 替换小数部分
		return money[:dotIndex+1] + replacement
	}
	return money // 如果没有小数点，返回原字符串
}

// FlToDtoMoney 转成数据中的金额
func FlToDtoMoney(money float64) string {
	if money == 0 {
		return "0.00"
	}
	return fmt.Sprintf("%.2f", money/100)
}

func PointsToYuan(money int64) float64 {
	return float64(money) / 100
}
func YuanToWanYuan(money float64) float64 {
	return float64(money) / 10000
}

func PointsToYuanF32(money int64) float32 {
	return float32(money) / 100
}
func PointsToIntegerYuan(money int64) int64 {
	return int64(math.Ceil(float64(money) / 100))
}

// IntegerToDBMoney 转成数据库的金额计算
func IntegerToDBMoney(money int64) int64 {
	return money * 100
}
func FlPointsToYuan(money float64) float64 {
	return math.Floor(money) / 100
}

func StringToPoint(money string) (int64, error) {
	float, err := strconv.ParseFloat(money, 64)
	if err != nil {
		return 0, err
	}
	return ToDBMoney(float), nil
}
func StringToYuan(money string) (float64, error) {
	float, err := strconv.ParseFloat(money, 64)
	if err != nil {
		return 0, err
	}
	return float, nil
}

// CalculateActualAmount 计算实际金额
func CalculateActualAmount(amountStr string, decimals int) (*big.Rat, error) {
	// 创建一个 big.Int 来存储原始金额
	amount := new(big.Int)

	// 将字符串转换为 big.Int
	_, success := amount.SetString(amountStr, 10)
	if !success {
		return nil, fmt.Errorf("invalid amount: %s", amountStr)
	}

	// 计算 10 的 decimals 次方
	base := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// 计算实际金额
	actualAmount := new(big.Rat).SetFrac(amount, base)

	return actualAmount, nil
}

// ConvertAndTrim 将实际金额转换为字符串，并截取小数点后四位，不包含小数点
func ConvertAndTrim(amount float64, decimals int, size int) string {
	// 将浮点数转换为字符串
	amountStr := strconv.FormatFloat(amount, 'f', decimals, 64)

	// 查找小数点的位置
	dotIndex := strings.Index(amountStr, ".")
	if dotIndex == -1 {
		// 如果没有小数点，返回原字符串
		return amountStr
	}
	// 截取小数点后四位，不包含小数点
	endIndex := dotIndex + size // 仅取小数点后四位
	if endIndex > len(amountStr) {
		endIndex = len(amountStr) // 确保不越界
	}
	return amountStr[dotIndex+1 : endIndex] // 从小数点后开始截取
}

// SplitAmount 拆解整数部分为 int64，并返回小数点后面的前四位作为 int64
func SplitAmount(amountStr string, decimals, size int) (int64, int64, error) {
	amount, err := CalculateActualAmount(amountStr, decimals)
	if err != nil {
		return 0, 0, err
	}
	amountFormatted := amount.FloatString(4)
	// 按小数点分割整数和小数部分
	parts := strings.Split(amountFormatted, ".")
	var integerPart int64
	var decimalPart int64

	// 处理整数部分
	if len(parts) > 0 {
		integerPart, _ = strconv.ParseInt(parts[0], 10, 64) // 转换为 int64
	}

	// 处理小数部分
	if len(parts) > 1 {
		// 只取前四位，补零
		decimalPartStr := parts[1]
		if len(decimalPartStr) >= 4 {
			decimalPartStr = decimalPartStr[:4]
		} else {
			decimalPartStr = fmt.Sprintf("%-4s", decimalPartStr)
		}
		decimalPart, _ = strconv.ParseInt(decimalPartStr, 10, 64) // 转换为 int64
	} else {
		// 如果没有小数部分，设置小数部分为 0
		decimalPart = 0
	}

	return integerPart, decimalPart, nil
}

// ConvertAmountToCents 将金额字符串转换为整数（单位：分）
func ConvertAmountToCents(amountStr string, decimals int) (int64, error) {
	// 假设 CalculateActualAmount 已处理精度
	amount, err := CalculateActualAmount(amountStr, decimals)
	if err != nil {
		return 0, err
	}

	// 格式化金额为足够的位数，确保小数点后有至少两位
	amountFormatted := amount.FloatString(decimals + 2)
	parts := strings.Split(amountFormatted, ".")

	// 获取整数部分
	integerPart := parts[0]

	// 获取小数点后两位（直接截取，不四舍五入）
	var decimalPart string
	if len(parts) > 1 && len(parts[1]) >= 2 {
		decimalPart = parts[1][:2]
	} else if len(parts) > 1 {
		decimalPart = fmt.Sprintf("%-2s", parts[1]) // 不足两位补零
	} else {
		decimalPart = "00"
	}

	// 合并整数部分和小数部分，去掉小数点表示成分
	amountInCentsStr := integerPart + decimalPart
	amountInCents, err := strconv.ParseInt(amountInCentsStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return amountInCents, nil
}

// SumWithoutMax 计算总和去掉最大的
func SumWithoutMax(amounts []float64) float64 {
	if len(amounts) == 0 {
		return 0 // 如果数组为空，返回 0
	}

	// 1. 找到最大值
	maxValue := amounts[0]
	for _, amount := range amounts {
		if amount > maxValue {
			maxValue = amount
		}
	}

	// 2. 计算总和（去掉最大值）
	var total float64
	var maxRemoved bool // 用于标记是否已经移除了最大值
	for _, amount := range amounts {
		if !maxRemoved && amount == maxValue {
			maxRemoved = true // 只去掉第一个最大值
			continue
		}
		total += amount
	}

	return total
}
