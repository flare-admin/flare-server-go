package config_api

import (
	"context"

	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
)

// IConfigApi 配置API接口
type IConfigApi interface {
	// GetValue 根据配置键获取配置值
	// key: 配置键
	// 返回: 配置值
	GetValue(ctx context.Context, key string) (interface{}, herrors.Herr)

	// GetValueMap 根据配置键列表获取配置值映射
	// keys: 配置键列表
	// 返回: 配置值映射
	GetValueMap(ctx context.Context, keys []string) (map[string]interface{}, herrors.Herr)

	// GetValueByGroup 根据分组编码获取配置值映射
	// groupCode: 分组编码
	// 返回: 配置值映射
	GetValueByGroup(ctx context.Context, groupCode string) (map[string]interface{}, herrors.Herr)

	// GetValueByGroupWithType 根据分组编码获取配置值映射，支持类型映射
	// groupCode: 分组编码
	// typeMap: 类型映射，key为配置键，value为配置类型
	// 返回: 配置值映射
	GetValueByGroupWithType(ctx context.Context, groupCode string, data interface{}) herrors.Herr
}
