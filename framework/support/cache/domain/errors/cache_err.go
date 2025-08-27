package cache_errors

import (
	"errors"
)

// 基础错误定义
var (
	// ErrNil 表示缓存键不存在
	ErrNil = errors.New("cache: nil")

	// ErrKeyNotFound 表示缓存键未找到
	ErrKeyNotFound = errors.New("cache: key not found")

	// ErrKeyExpired 表示缓存键已过期
	ErrKeyExpired = errors.New("cache: key expired")

	// ErrInvalidKey 表示无效的缓存键
	ErrInvalidKey = errors.New("cache: invalid key")

	// ErrInvalidValue 表示无效的缓存值
	ErrInvalidValue = errors.New("cache: invalid value")

	// ErrConnectionFailed 表示缓存连接失败
	ErrConnectionFailed = errors.New("cache: connection failed")

	// ErrOperationTimeout 表示缓存操作超时
	ErrOperationTimeout = errors.New("cache: operation timeout")
)
