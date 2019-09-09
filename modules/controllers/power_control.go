package controllers

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/constants"
	bus2 "github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
)

type PowerController struct {
	motorLeft  *gpio.MotorDriver
	motorRight *gpio.MotorDriver
	message    *bus.MessageBus
}

type PowerControlEvent struct {
	Direction constants.MovingDirection
	Speed     uint8
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

func (c *PowerController) Start() {
	c.message.RegisterHandler("moving-control-handler", "control.moving", func(e *bus2.Event) {
		var err error
		if evt, ok := e.Data.(PowerControlEvent); ok {
			switch evt.Direction {
			case constants.MOVING_DIRECTION__FORWARD:
				err = c.Forward(evt.Speed)
			case constants.MOVING_DIRECTION__BACKWARD:
				err = c.Backward(evt.Speed)
			case constants.MOVING_DIRECTION__TURN_LEFT:
				err = c.TurnLeft(evt.Speed)
			case constants.MOVING_DIRECTION__TURN_RIGHT:
				err = c.TurnRight(evt.Speed)
			}

			if err != nil {
				logrus.Errorf("[PowerController] moving-control-handler moving err: %v, event: %+v", err, evt)
			}
		} else {
			logrus.Errorf("[PowerController] moving-control-handler Data type err: %s", "not PowerControlEvent struct")
		}
	})
}
