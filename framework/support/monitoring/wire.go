package monitoring

import (
	"github.com/flare-admin/flare-server-go/framework/support/monitoring/application/handlers"
	"github.com/flare-admin/flare-server-go/framework/support/monitoring/domain/service"
	"github.com/flare-admin/flare-server-go/framework/support/monitoring/interfaces/rest"
	"github.com/google/wire"
)

// ProviderSet is monitoring providers.
var ProviderSet = wire.NewSet(
	service.NewMetricsService,
	handlers.NewMetricsQueryHandler,
	rest.NewMetricsController,
	NewServer,
)
