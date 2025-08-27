package model

import (
	"encoding/json"
	"github.com/flare-admin/flare-server-go/framework/pkg/utils"
	"strings"
)

// Cache 缓存模型
type Cache struct {
	TenantID  string      `json:"tenant_id"`  // 租户ID
	Key       string      `json:"key"`        // 缓存键
	Value     interface{} `json:"value"`      // 缓存值
	GroupID   string      `json:"group_id"`   // 分组ID
	ExpireAt  int64       `json:"expire_at"`  // 过期时间
	CreatedAt int64       `json:"created_at"` // 创建时间
	UpdatedAt int64       `json:"updated_at"` // 更新时间
}

// IsExpired 是否过期
func (c *Cache) IsExpired() bool {
	return c.ExpireAt > 0 && c.ExpireAt < utils.GetDateUnix()
}

// GetValue 获取缓存值
func (c *Cache) GetValue(target interface{}) error {
	if c.Value == nil {
		return nil
	}

	// 如果值是字符串，尝试反序列化
	if str, ok := c.Value.(string); ok {
		// 先尝试反序列化
		err := json.Unmarshal([]byte(str), target)
		if err == nil {
			return nil
		}

		// 反序列化失败，根据目标类型处理
		switch t := target.(type) {
		case *string:
			*t = str
		case *int:
			var v float64
			if err := json.Unmarshal([]byte(str), &v); err == nil {
				*t = int(v)
			}
		case *int64:
			var v float64
			if err := json.Unmarshal([]byte(str), &v); err == nil {
				*t = int64(v)
			}
		case *float64:
			if err := json.Unmarshal([]byte(str), t); err == nil {
				return nil
			}
		case *bool:
			if str == "true" {
				*t = true
			} else if str == "false" {
				*t = false
			}
		case *[]interface{}:
			if err := json.Unmarshal([]byte(str), t); err == nil {
				return nil
			}
		case *map[string]interface{}:
			if err := json.Unmarshal([]byte(str), t); err == nil {
				return nil
			}
		}
		return nil
	}

	// 如果值是基础类型，直接赋值
	switch v := c.Value.(type) {
	case string:
		if t, ok := target.(*string); ok {
			*t = v
		}
	case int:
		if t, ok := target.(*int); ok {
			*t = v
		} else if t, ok := target.(*int64); ok {
			*t = int64(v)
		} else if t, ok := target.(*float64); ok {
			*t = float64(v)
		}
	case int64:
		if t, ok := target.(*int64); ok {
			*t = v
		} else if t, ok := target.(*int); ok {
			*t = int(v)
		} else if t, ok := target.(*float64); ok {
			*t = float64(v)
		}
	case float64:
		if t, ok := target.(*float64); ok {
			*t = v
		} else if t, ok := target.(*int); ok {
			*t = int(v)
		} else if t, ok := target.(*int64); ok {
			*t = int64(v)
		}
	case bool:
		if t, ok := target.(*bool); ok {
			*t = v
		}
	default:
		// 其他类型尝试序列化后反序列化
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return json.Unmarshal(data, target)
	}

	return nil
}

// SetValue 设置缓存值
func (c *Cache) SetValue(value interface{}) error {
	if value == nil {
		c.Value = nil
		return nil
	}

	// 基础类型直接保存
	switch v := value.(type) {
	case string, int, int64, float64, bool:
		c.Value = v
		return nil
	}

	// 其他类型序列化为字符串
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	c.Value = string(data)
	return nil
}

// BuildKey 构建缓存键
func (c *Cache) BuildKey() string {
	parts := []string{c.TenantID}
	if c.GroupID != "" {
		parts = append(parts, c.GroupID)
	}
	parts = append(parts, c.Key)
	return strings.Join(parts, ":")
}

// BuildGroupPattern 构建分组匹配模式
func (c *Cache) BuildGroupPattern() string {
	c.Key = "*"
	return c.BuildKey()
}
