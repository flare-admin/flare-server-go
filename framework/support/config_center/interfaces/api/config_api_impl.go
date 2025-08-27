package config_api

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/herrors"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/handlers"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/queries"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/domain/errors"
)

// ConfigApi 配置API实现
type ConfigApi struct {
	queryHandler *handlers.ConfigQueryHandler
}

// NewConfigApi 创建配置API
func NewConfigApi(queryHandler *handlers.ConfigQueryHandler) IConfigApi {
	return &ConfigApi{
		queryHandler: queryHandler,
	}
}

// GetValue 根据配置键获取配置值
func (a *ConfigApi) GetValue(ctx context.Context, key string) (interface{}, herrors.Herr) {
	value, err := a.queryHandler.HandleGetValue(ctx, queries.GetConfigValueQuery{
		Key: key,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "Get config value error: %v", err)
		return nil, errors.GetConfigFail(err)
	}
	return value, nil
}

// GetValueMap 根据配置键列表获取配置值映射
func (a *ConfigApi) GetValueMap(ctx context.Context, keys []string) (map[string]interface{}, herrors.Herr) {
	valueMap, err := a.queryHandler.HandleGetValueMap(ctx, queries.GetConfigValueMapQuery{
		Keys: keys,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "Get config value map error: %v", err)
		return nil, errors.GetConfigFail(err)
	}
	return valueMap, nil
}

// GetValueByGroup 根据分组编码获取配置值映射
func (a *ConfigApi) GetValueByGroup(ctx context.Context, groupCode string) (map[string]interface{}, herrors.Herr) {
	// 获取分组下的所有配置
	configs, err := a.queryHandler.HandleGetValueByGroupCode(ctx, groupCode)
	if err != nil {
		hlog.CtxErrorf(ctx, "Get configs by group error: %v", err)
		return nil, errors.GetConfigFail(err)
	}

	// 构建配置键列表
	keys := make([]string, 0, len(configs))
	for _, config := range configs {
		keys = append(keys, config.Key)
	}

	// 获取配置值映射
	return a.GetValueMap(ctx, keys)
}

// GetValueByGroupWithType 根据分组编码获取配置值映射，支持类型映射
func (a *ConfigApi) GetValueByGroupWithType(ctx context.Context, groupCode string, data interface{}) herrors.Herr {
	// 获取分组下的所有配置
	configs, her := a.queryHandler.HandleGetValueByGroupCode(ctx, groupCode)
	if herrors.HaveError(her) {
		hlog.CtxErrorf(ctx, "Get configs by group error: %v", her)
		return errors.GetConfigFail(her)
	}

	// 构建配置值映射
	valueMap := make(map[string]interface{})
	for _, config := range configs {
		// 根据配置类型解析值
		var value interface{}
		switch config.Type {
		case "string", "regex":
			value = config.Value
		case "int":
			var v int64
			if err := json.Unmarshal([]byte(config.Value), &v); err != nil {
				hlog.CtxErrorf(ctx, "Parse config value error: %v", err)
				return errors.ConfigTypeInvalidFail
			}
			value = v
		case "float":
			var v float64
			if err := json.Unmarshal([]byte(config.Value), &v); err != nil {
				hlog.CtxErrorf(ctx, "Parse config value error: %v", err)
				return errors.ConfigTypeInvalidFail
			}
			value = v
		case "bool":
			var v bool
			if err := json.Unmarshal([]byte(config.Value), &v); err != nil {
				hlog.CtxErrorf(ctx, "Parse config value error: %v", err)
				return errors.ConfigTypeInvalidFail
			}
			value = v
		case "json":
			var v interface{}
			if err := json.Unmarshal([]byte(config.Value), &v); err != nil {
				hlog.CtxErrorf(ctx, "Parse config value error: %v", err)
				return errors.ConfigTypeInvalidFail
			}
			value = v
		case "array":
			var v []interface{}
			if err := json.Unmarshal([]byte(config.Value), &v); err != nil {
				hlog.CtxErrorf(ctx, "Parse config value error: %v", err)
				return errors.ConfigTypeInvalidFail
			}
			value = v
		case "object":
			var v map[string]interface{}
			if err := json.Unmarshal([]byte(config.Value), &v); err != nil {
				hlog.CtxErrorf(ctx, "Parse config value error: %v", err)
				return errors.ConfigTypeInvalidFail
			}
			value = v
		case "time":
			//var v time.Time
			//if err := json.Unmarshal([]byte(config.Value), &v); err != nil {
			//	hlog.CtxErrorf(ctx, "Parse config value error: %v", err)
			//	return errors.ConfigTypeInvalidFail
			//}
			//value = v
			value = config.Value
		case "date":
			//var v time.Time
			//if err := json.Unmarshal([]byte(config.Value), &v); err != nil {
			//	hlog.CtxErrorf(ctx, "Parse config value error: %v", err)
			//	return errors.ConfigTypeInvalidFail
			//}
			//value = v.Format("2006-01-02")
			value = config.Value
		case "datetime":
			//var v time.Time
			//if err := json.Unmarshal([]byte(config.Value), &v); err != nil {
			//	hlog.CtxErrorf(ctx, "Parse config value error: %v", err)
			//	return errors.ConfigTypeInvalidFail
			//}
			//value = v.Format("2006-01-02 15:04:05")
			value = config.Value
		default:
			hlog.CtxErrorf(ctx, "Unsupported config type: %s", config.Type)
			return errors.ConfigTypeInvalidFail
		}
		valueMap[config.Key] = value
	}

	// 将配置值映射转换为目标结构体
	jsonData, err1 := json.Marshal(valueMap)
	if err1 != nil {
		hlog.CtxErrorf(ctx, "Marshal config value map error: %v", err1)
		return errors.ConfigTypeInvalidFail
	}

	if err := json.Unmarshal(jsonData, data); err != nil {
		hlog.CtxErrorf(ctx, "Unmarshal config value map error: %v", err)
		return errors.ConfigTypeInvalidFail
	}

	return nil
}
