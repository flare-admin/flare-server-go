package base

import (
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/base/casbin"
	"github.com/flare-admin/flare-server-go/framework/support/base/infrastructure/base/oplog"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	casbin.NewRepositoryImpl,
	oplog.NewDbOperationLogWriter,
)
