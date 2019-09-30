package camera

import (
	tp "github.com/henrylee2cn/teleport"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/modules/controllers"
	"github.com/sirupsen/logrus"
)

func (c *Camera) Holder(req *client.CameraHolderRequest) *tp.Status {
	logrus.Debug(req.Direction, req.Angle)
	_, err := c.messageBus.Emit(controllers.PowerControlTopic, req, "")
	if err != nil {
		return tp.NewStatus(99, "", err)
	}
	return nil
}
