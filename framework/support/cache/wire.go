package cache

import (
	as "github.com/flare-admin/flare-server-go/framework/support/cache/application/service"
	"github.com/flare-admin/flare-server-go/framework/support/cache/domain/service"
	is "github.com/flare-admin/flare-server-go/framework/support/cache/infrastructure/service"
	cacheapi "github.com/flare-admin/flare-server-go/framework/support/cache/interfaces/api"
	"github.com/flare-admin/flare-server-go/framework/support/cache/interfaces/rest"
	"github.com/google/wire"
)

// ProviderSet 缓存模块的依赖注入集合
var ProviderSet = wire.NewSet(
	// 基础设施层
	is.NewRedisCacheService,

	// 领域服务层
	service.NewInternalCacheService,
	service.NewCacheService,
	wire.Bind(new(service.InternalDomainCacheService), new(*service.InternalCacheServiceImpl)),
	wire.Bind(new(cacheapi.InternalCacheService), new(*service.InternalCacheServiceImpl)),
	// 应用服务层
	as.NewCacheService,

	// 接口层
	rest.NewCacheHandler,
)

// BaseProviderSet 基础依赖
var BaseProviderSet = wire.NewSet(
	// 基础设施层
	is.NewRedisCacheService,
	service.NewInternalCacheService,
	wire.Bind(new(cacheapi.InternalCacheService), new(*service.InternalCacheServiceImpl)),
)
