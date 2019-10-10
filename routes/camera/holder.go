package camera

import (
	tp "github.com/henrylee2cn/teleport"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/modules/controllers"
)

func (c *Camera) Holder(req *client.CameraHolderRequest) *tp.Status {
	_, err := c.messageBus.Emit(controllers.CameraHolderTopic, req, "")
	if err != nil {
		return tp.NewStatus(99, "", err)
	}
	return nil
}
