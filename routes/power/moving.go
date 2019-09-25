package power

import (
	tp "github.com/henrylee2cn/teleport"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/modules/controllers"
)

func (p *Power) Moving(req *client.PowerMovingRequest) *tp.Status {
	_, err := p.messageBus.Emit(controllers.PowerControlTopic, req, "")
	if err != nil {
		return tp.NewStatus(99, "", err)
	}
	return nil
}
