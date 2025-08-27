package config_center

import (
	"github.com/flare-admin/flare-server-go/framework/support/config_center/application/handlers"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/infrastructure/repository"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/interfaces/api"
	"github.com/flare-admin/flare-server-go/framework/support/config_center/interfaces/rest"
	"github.com/google/wire"
)

// ProviderSet 配置中心依赖注入
var ProviderSet = wire.NewSet(
	repository.NewConfigRepository,
	repository.NewConfigGroupRepository,
	handlers.NewConfigCommandHandler,
	handlers.NewConfigQueryHandler,
	handlers.NewConfigGroupCommandHandler,
	handlers.NewConfigGroupQueryHandler,
	rest.NewConfigHandler,
	config_api.NewConfigApi,
)

// BaseProviderSet 基础依赖
var BaseProviderSet = wire.NewSet(
	repository.NewConfigRepository,
	repository.NewConfigGroupRepository,
	handlers.NewConfigQueryHandler,
	config_api.NewConfigApi,
)
