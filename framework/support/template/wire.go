package template

import (
	commondhandler "github.com/flare-admin/flare-server-go/framework/support/template/application/command/handler"
	queryhandler "github.com/flare-admin/flare-server-go/framework/support/template/application/queries/handler"
	"github.com/flare-admin/flare-server-go/framework/support/template/domain/service"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/persistence/data"
	"github.com/flare-admin/flare-server-go/framework/support/template/infrastructure/repository"
	"github.com/flare-admin/flare-server-go/framework/support/template/interfaces/admin"
	templateapi "github.com/flare-admin/flare-server-go/framework/support/template/interfaces/api"
	"github.com/google/wire"
)

// ProviderSet 提供模板模块的所有依赖
var ProviderSet = wire.NewSet(
	// 公共依赖
	BaseProviderSet,
	// 命令处理器
	commondhandler.NewTemplateCommandHandler,
	commondhandler.NewCategoryCommandHandler,

	// 查询处理器
	queryhandler.NewTemplateQueryHandler,
	queryhandler.NewCategoryQueryHandler,

	// 接口层
	admin.NewTemplateService,
	admin.NewCategoryService,

	// 服务
	NewTempServer,
)

// BaseProviderSet 公共依赖
var BaseProviderSet = wire.NewSet(
	//数据层
	data.NewTemplateRepository,
	data.NewCategoryRepository,

	// 仓储层
	repository.NewTemplateRepository,
	repository.NewCategoryRepository,

	// 领域服务
	service.NewTemplateService,
	service.NewCategoryService,

	// 内部接口
	templateapi.NewTemplateService,
)
