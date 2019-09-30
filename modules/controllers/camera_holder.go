package controllers

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/constants"
	bus2 "github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
)

const CameraHolderTopic = "camera.holder"
const MaxAngle float64 = 180

func CameraHolderController(servoHorizon *gpio.ServoDriver, servoVertical *gpio.ServoDriver, messageBus *bus.MessageBus) {
	messageBus.RegisterHandler("camera-holder-handler", CameraHolderTopic, func(e *bus2.Event) {
		var err error
		if evt, ok := e.Data.(*client.CameraHolderRequest); ok {
			if evt.Direction == constants.HOLDER_DIRECTION__HORIZEN {
				err = servoHorizon.Move(evt.Angle)
			} else {
				err = servoVertical.Move(evt.Angle)
			}

			if err != nil {
				logrus.Errorf("[HolderController] camera-holder-handler moving err: %v, event: %+v", err, evt)
			}
		} else {
			logrus.Errorf("[HolderController] camera-holder-handler Data type err: %s", "not CameraHolderRequest struct")
		}
	})
}
