package application

import (
	"github.com/flare-admin/flare-server-go/framework/support/base/application/handlers"
	"github.com/google/wire"
)

var HandlerSet = wire.NewSet(
	handlers.NewUserCommandHandler,
	handlers.NewUserQueryHandler,
)
