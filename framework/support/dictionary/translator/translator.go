package translator

import (
	"context"
)

// ITranslator 字典翻译器接口
type ITranslator interface {
	// GetTranslateInfo 根据分类和值获取翻译信息
	GetTranslateInfo(ctx context.Context, categoryID string, value interface{}) (*TranslateInfo, error)
	// GetLabel 根据分类和值获取标签
	GetLabel(ctx context.Context, categoryID string, value interface{}) (string, error)
	// GetI18nKey 根据分类和值获取国际化key
	GetI18nKey(ctx context.Context, categoryID string, value interface{}) (string, error)
	// Translate 翻译结构体
	Translate(ctx context.Context, obj interface{}) error
	// TranslateSlice 翻译切片
	TranslateSlice(ctx context.Context, slice interface{}) error
	// ClearCache 清除指定分类的缓存
	ClearCache(ctx context.Context, categoryID string) error
}
