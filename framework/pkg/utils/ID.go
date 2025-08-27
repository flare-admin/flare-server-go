package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 验证身份证号码的正则表达式
var idCardRegex = regexp.MustCompile(`^\d{17}[\dXx]$`)

// 身份证号码中每一位的加权因子
var weightFactors = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}

// 校验码对照表
var checkCodeMap = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}

// IsValidIDCard 验证身份证号码的正确性
func IsValidIDCard(id string) bool {
	// 长度和格式校验
	if !idCardRegex.MatchString(id) {
		return false
	}

	// 将身份证号码的前17位提取出来
	id17 := id[:17]

	// 计算校验码
	sum := 0
	for i := 0; i < 17; i++ {
		num, err := strconv.Atoi(string(id17[i]))
		if err != nil {
			return false
		}
		sum += num * weightFactors[i]
	}

	// 计算得到的校验码
	mod := sum % 11
	checkCode := checkCodeMap[mod]

	// 比较校验码
	return checkCode == id[17] || (checkCode == 'X' && strings.ToUpper(string(id[17])) == "X")
}

// ExtractBirthdayAndAge 从身份证号中提取生日和年龄
func ExtractBirthdayAndAge(id string) (birthday string, age int, err error) {
	if len(id) != 18 {
		return "", 0, fmt.Errorf("invalid ID length")
	}
	// 提取生日
	year, err := strconv.Atoi(id[6:10])
	if err != nil {
		return "", 0, fmt.Errorf("invalid year in ID")
	}
	month, err := strconv.Atoi(id[10:12])
	if err != nil {
		return "", 0, fmt.Errorf("invalid month in ID")
	}
	day, err := strconv.Atoi(id[12:14])
	if err != nil {
		return "", 0, fmt.Errorf("invalid day in ID")
	}
	birthday = fmt.Sprintf("%04d-%02d-%02d", year, month, day)

	// 计算年龄
	now := GetTimeNow()
	birthdate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	age = now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		age--
	}

	return birthday, age, nil
}

// ExtractBirthdayAgeAndGender 从身份证号中提取生日、年龄和性别
// return gender 1 女 2男
func ExtractBirthdayAgeAndGender(id string) (birthday string, age int, gender int, err error) {
	if len(id) != 18 {
		return "", 0, 1, fmt.Errorf("invalid ID length")
	}
	// 提取生日
	year, err := strconv.Atoi(id[6:10])
	if err != nil {
		return "", 0, 1, fmt.Errorf("invalid year in ID")
	}
	month, err := strconv.Atoi(id[10:12])
	if err != nil {
		return "", 0, 1, fmt.Errorf("invalid month in ID")
	}
	day, err := strconv.Atoi(id[12:14])
	if err != nil {
		return "", 0, 1, fmt.Errorf("invalid day in ID")
	}
	birthday = fmt.Sprintf("%04d-%02d-%02d", year, month, day)

	// 计算年龄
	now := GetTimeNow()
	birthdate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	age = now.Year() - birthdate.Year()
	if now.YearDay() < birthdate.YearDay() {
		age--
	}

	// 提取性别
	genderDigit, err := strconv.Atoi(string(id[16]))
	if err != nil {
		return "", 0, 1, fmt.Errorf("invalid gender digit in ID")
	}
	if genderDigit%2 == 0 {
		gender = 1
	} else {
		gender = 2
	}

	return birthday, age, gender, nil
}
