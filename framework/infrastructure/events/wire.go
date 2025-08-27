package events

import (
	sysevents "github.com/flare-admin/flare-server-go/framework/pkg/events"
	"github.com/flare-admin/flare-server-go/framework/pkg/mq"
	"github.com/flare-admin/flare-server-go/framework/pkg/mqevent"
	"github.com/flare-admin/flare-server-go/framework/pkg/mqevent/manager"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	manager.NewEventBusManager,
	NewNatsEventBus,
	//系统事件
	sysevents.NewEventBus,
)

// NewNatsEventBus 创建 NATS 事件总线
func NewNatsEventBus(sr mq.Server) mqevent.IMQEventBus {
	return mqevent.NewMQEventBus(sr)
}
