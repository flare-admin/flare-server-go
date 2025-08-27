package infrastructure

import (
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/base"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/converter"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/handlers"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/query"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	persistence.ProviderSet,
	query.ProviderSet,
	converter.ProviderSet,
	base.ProviderSet,
	handlers.ProviderSet,
)
