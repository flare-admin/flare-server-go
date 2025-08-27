package model

// CacheType 缓存类型
type CacheType string

const (
	CacheTypeConfig      CacheType = "config"       // 配置缓存
	CacheTypeConfigGroup CacheType = "config_group" // 配置分组缓存
)

// CacheKey 生成缓存键
func CacheKey(cacheType CacheType, key string) string {
	return string(cacheType) + ":" + key
}
