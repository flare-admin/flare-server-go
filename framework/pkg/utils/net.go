package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// ToQueryParams 将结构体转换为查询参数
func ToQueryParams(v interface{}) (string, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("input is not a struct")
	}

	values := url.Values{}
	t := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		value := val.Field(i)
		if value.CanInterface() {
			values.Set(field.Name, fmt.Sprintf("%v", value.Interface()))
		}
	}

	return values.Encode(), nil
}

// ToQueryParamsByJsonTag 根据JSON标签将结构体转换为查询参数
func ToQueryParamsByJsonTag(v interface{}) (string, error) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("input is not a struct")
	}

	values := url.Values{}
	t := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		// 获取标签的实际名称，去掉逗号后面的部分
		tagName := strings.Split(tag, ",")[0]
		value := val.Field(i)
		if value.CanInterface() {
			values.Set(tagName, fmt.Sprintf("%v", value.Interface()))
		}
	}

	return values.Encode(), nil
}

// MapToQueryParams 将 map 转换为查询参数字符串
// MapToQueryParams 将 map 转换为查询参数字符串
func MapToQueryParams(params map[string]interface{}) (string, error) {
	// 获取所有键并排序
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parms := make(map[string]string)
	for _, key := range keys {
		value := params[key]
		switch v := value.(type) {
		case string:
			parms[key] = v
		case int, int8, int16, int32, int64:
			parms[key] = fmt.Sprintf("%d", v)
		case uint, uint8, uint16, uint32, uint64:
			parms[key] = fmt.Sprintf("%d", v)
		case float32, float64:
			parms[key] = fmt.Sprintf("%f", v)
		case bool:
			parms[key] = fmt.Sprintf("%t", v)
		default:
			return "", fmt.Errorf("unsupported type for key %s", key)
		}
	}
	// 手动排序键值对，保持原有顺序
	var sortedParams []string
	for _, key := range keys {
		sortedParams = append(sortedParams, fmt.Sprintf("%s=%s", key, parms[key]))
	}
	return strings.Join(sortedParams, "&"), nil
}

// GenerateMD5Signature 生成MD5签名
func GenerateMD5Signature(data string) string {
	hasher := md5.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetSignToUpper 计算签名摘要
func GetSignToUpper(params map[string]interface{}, keyPar, key string) string {
	return strings.ToUpper(GetSign(params, keyPar, key))
}

// GetSign 计算签名摘要
func GetSign(params map[string]interface{}, keyPar, key string) string {
	var list []string
	for k, v := range params {
		if val, ok := v.(string); ok && val != "" {
			list = append(list, k+"="+val+"&")
		} else if val1, ok1 := v.(float64); ok1 {
			list = append(list, k+"="+strconv.FormatFloat(val1, 'f', -1, 64)+"&")
		} else if val2, ok2 := v.(int); ok2 {
			list = append(list, k+"="+strconv.Itoa(val2)+"&")
		} else if val3, ok3 := v.(int8); ok3 {
			list = append(list, k+"="+strconv.FormatInt(int64(val3), 10)+"&")
		} else if val4, ok4 := v.(int16); ok4 {
			list = append(list, k+"="+strconv.FormatInt(int64(val4), 10)+"&")
		} else if val5, ok5 := v.(int32); ok5 {
			list = append(list, k+"="+strconv.FormatInt(int64(val5), 10)+"&")
		} else if val6, ok6 := v.(int64); ok6 {
			list = append(list, k+"="+strconv.FormatInt(val6, 10)+"&")
		}
	}
	sort.Strings(list)
	var sb strings.Builder
	for _, str := range list {
		sb.WriteString(str)
	}
	sb.WriteString(fmt.Sprintf("%s=%s", keyPar, key))
	fmt.Println(sb.String())
	signature := md5Hash(sb.String())
	return signature
}

// GetSignByMap 计算签名摘要
func GetSignByMap(params map[string]string, key string) string {
	var list []string
	for k, v := range params {
		if v != "" {
			list = append(list, k+"="+v+"&")
		}
	}
	sort.Strings(list)
	var sb strings.Builder
	for _, str := range list {
		sb.WriteString(str)
	}
	sb.WriteString("key=" + key)
	fmt.Println(sb.String())
	signature := md5Hash(sb.String())
	return strings.ToUpper(signature)
}

// md5Hash 计算字符串的MD5值
func md5Hash(value string) string {
	hasher := md5.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}

func QueryStrToMap(queryString string) (map[string]interface{}, error) {
	// Parse the query string
	params, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, err
	}

	// Convert to a map and remove the "sign" key
	mapParams := make(map[string]interface{})
	for key, values := range params {
		// Take the first value only, as ParseQuery returns a map of string slices
		mapParams[key] = values[0]
	}
	return mapParams, nil
}

// StructToMap 将结构体根据json tag转换为map
func StructToMap(data interface{}) (map[string]interface{}, error) {
	// 首先将结构体转换为JSON字节
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// 将JSON字节解码为map
	var resultMap map[string]interface{}
	err = json.Unmarshal(jsonData, &resultMap)
	if err != nil {
		return nil, err
	}

	return resultMap, nil
}

// JoinURL 将多个字符串拼接成一个URL
// 参数:
//   - parts: 要拼接的URL部分
//
// 返回:
//   - string: 拼接后的URL
func JoinURL(parts ...string) string {
	var result strings.Builder
	for i, part := range parts {
		// 移除开头和结尾的斜杠
		part = strings.Trim(part, "/")

		// 如果不是第一个部分，添加斜杠
		if i > 0 {
			result.WriteString("/")
		}

		// 添加当前部分
		result.WriteString(part)
	}

	return result.String()
}
