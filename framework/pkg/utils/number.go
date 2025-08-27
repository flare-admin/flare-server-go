package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func IsValidAmount(amount string) bool {
	// 定义金额的正则表达式
	regexPattern := `^-?\d{1,3}(,\d{3})*(\.\d{1,2})?$`
	re := regexp.MustCompile(regexPattern)
	return re.MatchString(amount)
}

// IsValidAmountAndConvert 校验字符串是否为合法金额并返回float64值
func IsValidAmountAndConvert(amount string) (float64, error) {
	// 定义金额的正则表达式
	regexPattern := `^-?\d{1,3}(,\d{3})*(\.\d{1,2})?$|^-?\d+(\.\d{1,2})?$`
	re := regexp.MustCompile(regexPattern)

	// 检查是否匹配正则表达式
	if !re.MatchString(amount) {
		return 0, fmt.Errorf("invalid amount format")
	}

	// 去除金额中的千位分隔符
	cleanAmount := strings.ReplaceAll(amount, ",", "")

	// 转换为float64
	value, err := strconv.ParseFloat(cleanAmount, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse amount: %v", err)
	}

	return value, nil
}

// IsValidInteger 校验字符串是否是有效的整数
func IsValidInteger(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64) // 尝试将字符串解析为整数
	return err == nil                     // 如果没有错误，则返回 true
}
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

// ChunkSlice 将切片分成大小为 batchSize 的多个子切片
func ChunkSlice[T any](slice []T, batchSize int) [][]T {
	if batchSize <= 0 {
		return nil
	}

	var chunks [][]T
	for i := 0; i < len(slice); i += batchSize {
		end := i + batchSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

// ContainsString checks if a slice contains a specific string.
func ContainsString(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func StringToFloatSlice(str string) ([]float64, error) {
	if len(str) == 0 {
		return make([]float64, 0), nil
	}
	str = strings.Trim(str, " ")
	// 使用 strings.Split 将字符串按逗号分割成字符串切片
	strSlice := strings.Split(str, ",")

	// 创建一个float64切片来存储转换后的值
	floatSlice := make([]float64, len(strSlice))

	// 遍历字符串切片，并将每个字符串转换为float64
	for i, v := range strSlice {
		// 将字符串转换为float64
		num, err := strconv.ParseFloat(v, 64)
		if err != nil {
			fmt.Println("转换错误:", err)
			return nil, err
		}
		// 存储转换后的值
		floatSlice[i] = num
	}
	return floatSlice, nil
}
func StringToInt8Slice(str string) ([]int8, error) {
	if len(str) == 0 {
		return make([]int8, 0), nil
	}
	str = strings.Trim(str, " ")
	strSlice := strings.Split(str, ",")
	int8Slice := make([]int8, len(strSlice))
	for i, v := range strSlice {
		num, err := strconv.ParseInt(v, 10, 8)
		if err != nil {
			fmt.Println("转换错误:", err)
			return nil, err
		}
		int8Slice[i] = int8(num)
	}
	return int8Slice, nil
}

// StringToAmountsAndWallets 将字符串转换为金额和钱包类型
func StringToAmountsAndWallets(amounts, wallets string, sizeEqual bool) ([]float64, []int8, error) {
	amountsSlice, err := StringToFloatSlice(amounts)
	if err != nil {
		return nil, nil, err
	}
	walletsSlice, err := StringToInt8Slice(wallets)
	if err != nil {
		return nil, nil, err
	}
	if len(amountsSlice) != len(walletsSlice) && !sizeEqual {
		return nil, nil, fmt.Errorf("amounts and wallets size not match")
	}
	return amountsSlice, walletsSlice, nil
}

// StringToAmountsAndChickWallets 将字符串转换为金额并校验钱包
func StringToAmountsAndChickWallets(amounts string, wallets []int8, sizeEqual bool) ([]float64, []int8, error) {
	amountsSlice, err := StringToFloatSlice(amounts)
	if err != nil {
		return nil, nil, err
	}
	if len(amountsSlice) != len(wallets) && !sizeEqual {
		return nil, nil, fmt.Errorf("amounts and wallets size not match")
	}
	return amountsSlice, wallets, nil
}

func StringToInt64(str string) (int64, error) {
	if len(str) == 0 {
		return 0, nil
	}
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func StringToFloat64(str string) (float64, error) {
	if len(str) == 0 {
		return 0, nil
	}
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}
	return num, nil
}

// IntToBool 将int转换为bool，1对应true，2对应false
func IntToBool(i int) bool {
	return i == 1
}

// BoolToInt 将bool转换为int，true对应1，false对应2
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 2
}
