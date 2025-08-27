package systask

import (
	"github.com/flare-admin/flare-server-go/framework/support/systask/biz"
	"github.com/flare-admin/flare-server-go/framework/support/systask/data"
	"github.com/flare-admin/flare-server-go/framework/support/systask/interfaces"
	"github.com/flare-admin/flare-server-go/framework/support/systask/manager"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(data.NewTaskRepo, biz.NewTaskBiz, manager.NewTaskManager, interfaces.NewTaskService)
