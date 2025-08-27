package infrastructure

import (
	"github.com/flare-admin/flare-server-go/framework/infrastructure/database"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/events"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/idempotence"
	"github.com/flare-admin/flare-server-go/framework/infrastructure/mq"
	"github.com/flare-admin/flare-server-go/framework/pkg/hserver/middleware/repeated_submit"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	database.ProviderSet,
	idempotence.NewIdempotencyTool,
	repeated_submit.NewDefRepeatedSubmitLock,
	mq.ProviderSet,
	events.ProviderSet,
)
