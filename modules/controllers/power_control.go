package controllers

import (
	"github.com/johnnyeven/libtools/bus"
	bus2 "github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"time"
)

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

func (c *PowerController) Start() {

	c.message.RegisterHandler("moving-control-handler", "control.moving", func(e *bus2.Event) {
		logrus.Infof("%+v", e)
	})

	speed := byte(0)
	fadeAmount := byte(15)
	revert := false

	gobot.Every(100*time.Millisecond, func() {
		if !revert {
			c.Forward(speed)
		} else {
			c.Backward(speed)
		}
		speed = speed + fadeAmount
		if speed <= 0 || speed >= 255 {
			fadeAmount = -fadeAmount
		}
		if speed <= 0 {
			revert = !revert
		}
	})
}
