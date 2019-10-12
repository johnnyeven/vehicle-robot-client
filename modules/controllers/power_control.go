package controllers

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/constants"
	bus2 "github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
)

const PowerControlTopic = "power.moving"
const MaxPower float64 = 255

type PowerController struct {
	motorLeft  *gpio.MotorDriver
	motorRight *gpio.MotorDriver
	message    *bus.MessageBus
}

func NewPowerController(motorLeft *gpio.MotorDriver, motorRight *gpio.MotorDriver, messageBus *bus.MessageBus) *PowerController {
	return &PowerController{
		motorLeft:  motorLeft,
		motorRight: motorRight,
		message:    messageBus,
	}
}

func (c *PowerController) Forward(speed uint8) error {
	err := c.motorLeft.Forward(speed)
	if err != nil {
		return err
	}
	err = c.motorRight.Forward(speed)
	return err
}

func (c *PowerController) Backward(speed uint8) error {
	err := c.motorLeft.Backward(speed)
	if err != nil {
		return err
	}
	err = c.motorRight.Backward(speed)
	return err
}

func (c *PowerController) TurnLeft(speed uint8) error {
	err := c.motorLeft.Off()
	if err != nil {
		return err
	}
	err = c.motorRight.Forward(speed)
	return err
}

func (c *PowerController) TurnRight(speed uint8) error {
	err := c.motorRight.Off()
	if err != nil {
		return err
	}
	err = c.motorLeft.Forward(speed)
	return err
}

func (c *PowerController) Stop() error {
	err := c.motorRight.Off()
	if err != nil {
		return err
	}
	err = c.motorLeft.Off()
	return err
}

func (c *PowerController) Start() {
	c.message.RegisterTopic(PowerControlTopic)
	c.message.RegisterHandler("camera-moving-handler", PowerControlTopic, func(e *bus2.Event) {
		var err error
		if evt, ok := e.Data.(*client.PowerMovingRequest); ok {
			switch evt.Direction {
			case constants.MOVING_DIRECTION__FORWARD:
				err = c.Forward(uint8(evt.Speed * MaxPower))
			case constants.MOVING_DIRECTION__BACKWARD:
				err = c.Backward(uint8(evt.Speed * MaxPower))
			case constants.MOVING_DIRECTION__TURN_LEFT:
				err = c.TurnLeft(uint8(evt.Speed * MaxPower))
			case constants.MOVING_DIRECTION__TURN_RIGHT:
				err = c.TurnRight(uint8(evt.Speed * MaxPower))
			case constants.MOVING_DIRECTION__STOP:
				err = c.Stop()
			}

			if err != nil {
				logrus.Errorf("[PowerController] camera-moving-handler moving err: %v, event: %+v", err, evt)
			}
		} else {
			logrus.Errorf("[PowerController] camera-moving-handler Data type err: %s", "not PowerMovingRequest struct")
		}
	})
}
