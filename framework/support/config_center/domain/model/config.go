package model

// ConfigType 配置类型
type ConfigType string

const (
	ConfigTypeString   ConfigType = "string"
	ConfigTypeInt      ConfigType = "int"
	ConfigTypeFloat    ConfigType = "float"
	ConfigTypeBool     ConfigType = "bool"
	ConfigTypeJSON     ConfigType = "json"
	ConfigTypeArray    ConfigType = "array"
	ConfigTypeObject   ConfigType = "object"
	ConfigTypeTime     ConfigType = "time"
	ConfigTypeDate     ConfigType = "date"
	ConfigTypeDateTime ConfigType = "datetime"
	ConfigTypeRegex    ConfigType = "regex" // 正则表达式类型，实际存储为字符串，通过正则标识符判断
)
