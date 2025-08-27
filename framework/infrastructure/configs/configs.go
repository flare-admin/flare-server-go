package configs

import (
	"github.com/flare-admin/flare-server-go/framework/pkg/constant"
	"github.com/spf13/viper"
)

var (
	LocalizePath = ""
	Mode         constant.EnvMode // 开发环境
	configViper  *viper.Viper
)

// InitConfig 初始化配置
func InitConfig(path string, mode constant.EnvMode) error {
	configViper = viper.New()
	configViper.SetConfigFile(filePathByMode(mode, path))
	return configViper.ReadInConfig()
}

// LoadConfig 加载配置到指定的结构体
func LoadConfig(config interface{}) error {
	return configViper.Unmarshal(config)
}

// LoadConfigByKey 根据指定的key加载配置到结构体
func LoadConfigByKey(key string, config interface{}) error {
	return configViper.UnmarshalKey(key, config)
}

// LoadModuleConfig 加载模块配置
// module: 模块名称，如 "pay.new_pay"
// config: 配置结构体指针
func LoadModuleConfig(module string, config interface{}) error {
	return configViper.UnmarshalKey(module, config)
}

// GetString 获取字符串配置
func GetString(key string) string {
	return configViper.GetString(key)
}

// GetInt 获取整数配置
func GetInt(key string) int {
	return configViper.GetInt(key)
}

// GetBool 获取布尔配置
func GetBool(key string) bool {
	return configViper.GetBool(key)
}

// GetStringMap 获取map配置
func GetStringMap(key string) map[string]interface{} {
	return configViper.GetStringMap(key)
}

// Load 原有的Load函数，现在使用新的配置加载方式
func Load(path, lcp, env, logPath string) (*Bootstrap, error) {
	// 读取环境配置
	if mode := env; mode != "" {
		Mode = constant.EnvMode(mode)
	} else { // 默认「生产环境」
		Mode = constant.Development
	}

	// 初始化配置
	if err := InitConfig(path, Mode); err != nil {
		return nil, err
	}

	var bc Bootstrap
	if err := LoadConfig(&bc); err != nil {
		return nil, err
	}

	if lcp == "" {
		lcp = path + "/localize"
	}
	LocalizePath = lcp
	if logPath != "" {
		bc.Log.OutPath = logPath
	}
	// 检查一下数据库的日志级别有没有单独配置
	if bc.Data.DataBase.LogLevel == 0 {
		bc.Data.DataBase.LogLevel = bc.Log.Level
	}
	return &bc, nil
}

func filePathByMode(mode constant.EnvMode, path string) string {
	switch mode {
	case constant.Development:
		path = path + "/config_dev.yaml"
	case constant.Prerelease:
		path = path + "/config_pre.yaml"
	case constant.Production:
		path = path + "/config_pro.yaml"
	}
	return path
}

func GwtI18RootPath() string {
	return LocalizePath
}
