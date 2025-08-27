package translator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/support/dictionary/data"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Redis缓存key前缀
	dictCacheKeyPrefix = "dict:translator:"
	// 缓存过期时间
	dictCacheExpiration = 24 * time.Hour
)

// translator 字典翻译器实现
type translator struct {
	repo  data.IDictionaryRepo
	rdb   *redis.Client
	cache sync.Map // 缓存字典数据，key为categoryID，value为map[string]TranslateInfo
}

// NewTranslator 创建翻译器实例
func NewTranslator(repo data.IDictionaryRepo, rdb *redis.Client) ITranslator {
	return &translator{
		repo: repo,
		rdb:  rdb,
	}
}

// getCacheKey 生成缓存key
func getCacheKey(categoryID string) string {
	return dictCacheKeyPrefix + categoryID
}

// loadDictData 加载字典数据到缓存
func (t *translator) loadDictData(ctx context.Context, categoryID string) (map[string]TranslateInfo, error) {
	cacheKey := getCacheKey(categoryID)

	// 先从Redis缓存获取
	cacheData, err := t.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var dictMap map[string]TranslateInfo
		if err := json.Unmarshal([]byte(cacheData), &dictMap); err == nil {
			return dictMap, nil
		}
	}
	var myInt int = 1
	// 查询字典选项
	options, err := t.repo.GetOptions(ctx, categoryID, "", &myInt)
	if err != nil {
		return nil, err
	}

	// 构建翻译映射
	dictMap := make(map[string]TranslateInfo)
	for _, opt := range options {
		dictMap[opt.Value] = TranslateInfo{
			Label:   opt.Label,
			I18nKey: opt.I18nKey,
		}
	}

	// 存入Redis缓存
	if cacheData, err := json.Marshal(dictMap); err == nil {
		t.rdb.Set(ctx, cacheKey, cacheData, dictCacheExpiration)
	}

	return dictMap, nil
}

// GetTranslateInfo 根据分类和值获取翻译信息
func (t *translator) GetTranslateInfo(ctx context.Context, categoryID string, value interface{}) (*TranslateInfo, error) {
	// 转换值为字符串
	strValue := toString(value)
	if strValue == "" {
		return nil, fmt.Errorf("无效的值类型")
	}

	dictData, err := t.loadDictData(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	if info, ok := dictData[strValue]; ok {
		return &info, nil
	}
	return nil, fmt.Errorf("未找到对应的字典项: category=%s, value=%s", categoryID, strValue)
}

// GetLabel 根据分类和值获取标签
func (t *translator) GetLabel(ctx context.Context, categoryID string, value interface{}) (string, error) {
	info, err := t.GetTranslateInfo(ctx, categoryID, value)
	if err != nil {
		return "", err
	}
	return info.Label, nil
}

// GetI18nKey 根据分类和值获取国际化key
func (t *translator) GetI18nKey(ctx context.Context, categoryID string, value interface{}) (string, error) {
	info, err := t.GetTranslateInfo(ctx, categoryID, value)
	if err != nil {
		return "", err
	}
	return info.I18nKey, nil
}

// ClearCache 清除指定分类的缓存
func (t *translator) ClearCache(ctx context.Context, categoryID string) error {
	if categoryID == "" {
		// 清除所有字典缓存
		pattern := dictCacheKeyPrefix + "*"
		keys, err := t.rdb.Keys(ctx, pattern).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			return t.rdb.Del(ctx, keys...).Err()
		}
		return nil
	}
	// 清除指定分类的缓存
	return t.rdb.Del(ctx, getCacheKey(categoryID)).Err()
}

// Translate 翻译结构体
func (t *translator) Translate(ctx context.Context, obj interface{}) error {
	if obj == nil {
		return nil
	}
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// 获取dict标签
		dictTag := fieldType.Tag.Get("dict")
		if dictTag == "" {
			continue
		}

		// 解析标签配置
		config := parseDictTag(dictTag)
		if config == nil {
			continue
		}

		// 获取字段值
		fieldValue := field.Interface()
		strValue := toString(fieldValue)
		if strValue == "" {
			continue
		}

		// 加载字典数据
		dictData, err := t.loadDictData(ctx, config.Category)
		if err != nil {
			return err
		}

		// 获取翻译信息
		if info, ok := dictData[strValue]; ok {
			// 设置翻译后的值
			if config.Field != "" {
				targetField := val.FieldByName(config.Field)
				if targetField.IsValid() && targetField.CanSet() {
					switch targetField.Kind() {
					case reflect.String:
						targetField.SetString(info.Label)
					default:
						hlog.Errorf("无效的字段类型: %v", targetField.Kind())
					}
				}
			}
			// 设置i18nKey
			if config.I18nKey != "" {
				targetField := val.FieldByName(config.I18nKey)
				if targetField.IsValid() && targetField.CanSet() {
					switch targetField.Kind() {
					case reflect.String:
						targetField.SetString(info.I18nKey)
					default:
						hlog.Errorf("无效的字段类型: %v", targetField.Kind())
					}
				}
			}
		}
	}

	return nil
}

// TranslateSlice 翻译切片
func (t *translator) TranslateSlice(ctx context.Context, slice interface{}) error {
	val := reflect.ValueOf(slice)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		return nil
	}

	for i := 0; i < val.Len(); i++ {
		if err := t.Translate(ctx, val.Index(i).Interface()); err != nil {
			return err
		}
	}

	return nil
}

// parseDictTag 解析字典标签
// 参数 tag：字典标签，格式为 "category=gender,field=GenderLabel,i18=GenderI18Key"
// 返回值：返回一个包含解析信息的 DictTag 结构体指针
func parseDictTag(tag string) *DictTag {
	// 初始化一个空的 DictTag 结构体
	dictTag := &DictTag{}

	// 解析标签字符串，使用逗号分割每个键值对
	pairs := strings.Split(tag, ",")
	for _, pair := range pairs {
		// 进一步解析每个键值对
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			// 如果格式不正确，跳过该项
			continue
		}

		// 根据键名填充对应的 DictTag 字段
		switch kv[0] {
		case "category":
			dictTag.Category = kv[1]
		case "field":
			dictTag.Field = kv[1]
		case "i18":
			dictTag.I18nKey = kv[1]
		}
	}

	// 返回填充后的 DictTag 结构体指针
	return dictTag
}

// ValueType 支持的值类型
type ValueType interface {
	~string | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

// 修改 toString 方法，增加对数字类型的处理
func toString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case int8:
		return fmt.Sprintf("%d", v)
	case int16:
		return fmt.Sprintf("%d", v)
	case int32:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case uint:
		return fmt.Sprintf("%d", v)
	case uint8:
		return fmt.Sprintf("%d", v)
	case uint16:
		return fmt.Sprintf("%d", v)
	case uint32:
		return fmt.Sprintf("%d", v)
	case uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return fmt.Sprintf("%.0f", v) // 去掉小数部分
	case float64:
		return fmt.Sprintf("%.0f", v) // 去掉小数部分
	default:
		return fmt.Sprintf("%v", v)
	}
}
