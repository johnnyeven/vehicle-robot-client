package power

import (
	tp "github.com/henrylee2cn/teleport"
	"github.com/johnnyeven/libtools/bus"
)

type Power struct {
	tp.PushCtx

	messageBus *bus.MessageBus
}

func NewPowerRouter(messageBus *bus.MessageBus) *Power {
	return &Power{
		messageBus: messageBus,
	}
}
