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
const MaxAngle uint8 = 180
const CentreAngle uint8 = 90

func CameraHolderController(servoHorizon *gpio.ServoDriver, servoVertical *gpio.ServoDriver, messageBus *bus.MessageBus) {
	var currentHorizonAngle = CentreAngle
	var currentVerticalAngle = CentreAngle
	logrus.Infof("[HolderController] Init servos to center angle: %d", CentreAngle)
	err := servoHorizon.Move(currentHorizonAngle)
	if err != nil {
		logrus.Errorf("[HolderController] horizon servo move failed with err: %v", err)
		return
	}
	err = servoVertical.Move(currentVerticalAngle)
	if err != nil {
		logrus.Errorf("[HolderController] vertical servo move failed with err: %v", err)
		return
	}
	messageBus.RegisterHandler("camera-holder-handler", CameraHolderTopic, func(e *bus2.Event) {
		var err error
		if evt, ok := e.Data.(*client.CameraHolderRequest); ok {
			if evt.Direction == constants.HOLDER_DIRECTION__HORIZEN {
				currentHorizonAngle = servoAngleChange(currentHorizonAngle, evt.Offset)
				err = servoHorizon.Move(currentHorizonAngle)
			} else {
				currentVerticalAngle = servoAngleChange(currentVerticalAngle, evt.Offset)
				err = servoVertical.Move(currentVerticalAngle)
			}

			if err != nil {
				logrus.Errorf("[HolderController] camera-holder-handler moving err: %v, event: %+v", err, evt)
			}
		} else {
			logrus.Errorf("[HolderController] camera-holder-handler Data type err: %s", "not CameraHolderRequest struct")
		}
	})
}

func servoAngleChange(current uint8, offset float64) uint8 {
	current = uint8(float64(current) + offset)
	if current < 0 {
		current = 0
	} else if current > MaxAngle {
		current = MaxAngle
	}
	return current
}
