package translator

// TranslateInfo 翻译信息
type TranslateInfo struct {
	Label   string `json:"label"`   // 显示标签
	I18nKey string `json:"i18nKey"` // 国际化key
}

// DictTag 字典标签配置
type DictTag struct {
	Category string `json:"category"` // 字典分类ID
	Field    string `json:"field"`    // 翻译后赋值的字段
	I18nKey  string `json:"i18nKey"`  // 国际化key
}
