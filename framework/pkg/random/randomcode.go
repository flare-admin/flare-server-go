package random

import (
	"crypto/rand"
	"math/big"
)

// GenerateRandomNumericString 生成指定长度的纯数字字符串
func GenerateRandomNumericString(length int) (string, error) {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := cryptoRandInt(int64(len(digits)))
		if err != nil {
			return "", err
		}
		result[i] = digits[num]
	}
	return string(result), nil
}

// GenerateRandomAlphaNumericString 生成指定长度的数字字母字符串
// 包含数字(0-9)和大写字母(A-Z)
func GenerateRandomAlphaNumericString(length int) (string, error) {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := cryptoRandInt(int64(len(chars)))
		if err != nil {
			return "", err
		}
		result[i] = chars[num]
	}
	return string(result), nil
}

// 使用crypto/rand生成随机整数
func cryptoRandInt(max int64) (int64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}
