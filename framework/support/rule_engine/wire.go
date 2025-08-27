package rule_engine

import (
	comhandler "github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/command/handler"
	queryhandler "github.com/flare-admin/flare-server-go/framework/support/rule_engine/application/queries/handler"
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/infrastructure/persistence/data"
	"github.com/google/wire"
	// 领域层
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/domain/service"
	// 基础设施层
	"github.com/flare-admin/flare-server-go/framework/support/rule_engine/infrastructure/repository"
	// 接口层
	admin "github.com/flare-admin/flare-server-go/framework/support/rule_engine/interfaces/admin"
	api "github.com/flare-admin/flare-server-go/framework/support/rule_engine/interfaces/api"
)

var BaseProviderSet = wire.NewSet(
	// 基础设施层
	data.NewRuleCategoryRepository,
	data.NewRuleTemplateRepository,
	data.NewRuleRepository,

	repository.NewRuleTemplateRepository,
	repository.NewRuleCategoryRepository,
	repository.NewRuleRepository,

	// 领域层
	service.NewRuleTemplateService,
	service.NewRuleCategoryService,
	service.NewRuleService,
	service.NewRuleExecutionService,

	// 应用层
	comhandler.NewTemplateCommandHandler,
	comhandler.NewCategoryCommandHandler,
	comhandler.NewRuleCommandHandler,

	queryhandler.NewTemplateQueryHandler,
	queryhandler.NewCategoryQueryHandler,
	queryhandler.NewRuleQueryHandler,

	// 接口层
	api.NewRuleEngineService,
)
var ProviderSet = wire.NewSet(
	BaseProviderSet,
	// 接口层
	admin.NewTemplateService,
	admin.NewCategoryService,
	admin.NewRuleService,
	NewServer,
)
