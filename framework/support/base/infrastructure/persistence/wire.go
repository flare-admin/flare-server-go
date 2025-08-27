package persistence

import (
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence/data"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence/mapper"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/persistence/repository"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	data.ProviderSet,
	mapper.ProviderSet,
	repository.ProviderSet,
)
