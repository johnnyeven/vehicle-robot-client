package camera

import (
	tp "github.com/henrylee2cn/teleport"
	"github.com/johnnyeven/libtools/bus"
)

type Camera struct {
	tp.PushCtx

	messageBus *bus.MessageBus
}

func NewCameraRouter(messageBus *bus.MessageBus) *Camera {
	return &Camera{
		messageBus: messageBus,
	}
}
