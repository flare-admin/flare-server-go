package sysevent

import (
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/base"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/biz"
	"github.com/flare-admin/flare-server-go/framework/support/sysevent/data"
	syseventservice "github.com/flare-admin/flare-server-go/framework/support/sysevent/interfaces"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	data.NewEventRepo,
	data.NewSubscribeRepo,
	data.NewSubscribeParameterRepo,
	data.NewDeadLetterSubscribeRepo,
	biz.NewEventUseCase,
	biz.NewSubscribeUseCase,
	biz.NewDeadLetterSubscribeUseCase,

	base.NewSubscribeManagerUseCase,

	syseventservice.NewEventService,
)
