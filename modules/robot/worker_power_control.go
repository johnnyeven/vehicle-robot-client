package robot

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/service-vehicle-robot/constants/types"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	bus2 "github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

const (
	powerWorkerID                    = "power-worker"
	PowerControlTopic                = "power.moving"
	powerControlEventHandler         = "camera-moving-handler"
	MaxPower                 float64 = 255
)

type PowerWorker struct {
	motorLeft  *gpio.MotorDriver
	motorRight *gpio.MotorDriver
	bus        *bus.MessageBus
}

func NewPowerWorker(robot *Robot, bus *bus.MessageBus, config *global.RobotConfiguration) *PowerWorker {
	var firmataAdaptor *firmata.Adaptor
	var ok bool
	conn := robot.GetConnection(config.FirmataConnectionName)
	if conn == nil {
		firmataAdaptor = firmata.NewAdaptor(config.ArduinoDeviceID)
		firmataAdaptor.SetName(config.FirmataConnectionName)
		robot.AddConnection(firmataAdaptor)
	} else {
		if firmataAdaptor, ok = conn.(*firmata.Adaptor); !ok {
			logrus.Panicf("[PowerWorker] 连接器已存在，但并不是 *firmata.Adaptor 类型")
		}
	}

	motorLeft := gpio.NewMotorDriver(firmataAdaptor, config.LeftMotorSpeedPin)
	motorLeft.SetName(config.LeftMotorName)
	motorLeft.DirectionPin = config.LeftMotorDirectionPin
	motorRight := gpio.NewMotorDriver(firmataAdaptor, config.RightMotorSpeedPin)
	motorRight.SetName(config.RightMotorName)
	motorRight.DirectionPin = config.RightMotorDirectionPin

	robot.AddDevice(motorLeft, motorRight)
	return &PowerWorker{
		motorLeft:  motorLeft,
		motorRight: motorRight,
		bus:        bus,
	}
}

func (c *PowerWorker) forward(speed uint8) error {
	err := c.motorLeft.Forward(speed)
	if err != nil {
		return err
	}
	err = c.motorRight.Forward(speed)
	return err
}

func (c *PowerWorker) backward(speed uint8) error {
	err := c.motorLeft.Backward(speed)
	if err != nil {
		return err
	}
	err = c.motorRight.Backward(speed)
	return err
}

func (c *PowerWorker) turnLeft(speed uint8) error {
	err := c.motorLeft.Off()
	if err != nil {
		return err
	}
	err = c.motorRight.Forward(speed)
	return err
}

func (c *PowerWorker) turnRight(speed uint8) error {
	err := c.motorRight.Off()
	if err != nil {
		return err
	}
	err = c.motorLeft.Forward(speed)
	return err
}

func (c *PowerWorker) Stop() error {
	err := c.motorRight.Off()
	if err != nil {
		return err
	}
	err = c.motorLeft.Off()
	return err
}

func (c *PowerWorker) Start() {
	c.bus.RegisterTopic(PowerControlTopic)
	c.bus.RegisterHandler(powerControlEventHandler, PowerControlTopic, func(e *bus2.Event) {
		var err error
		if evt, ok := e.Data.(*client.PowerMovingRequest); ok {
			switch evt.Direction {
			case types.MOVING_DIRECTION__FORWARD:
				err = c.forward(uint8(evt.Speed * MaxPower))
			case types.MOVING_DIRECTION__BACKWARD:
				err = c.backward(uint8(evt.Speed * MaxPower))
			case types.MOVING_DIRECTION__TURN_LEFT:
				err = c.turnLeft(uint8(evt.Speed * MaxPower))
			case types.MOVING_DIRECTION__TURN_RIGHT:
				err = c.turnRight(uint8(evt.Speed * MaxPower))
			case types.MOVING_DIRECTION__STOP:
				err = c.Stop()
			}

			if err != nil {
				logrus.Errorf("[PowerWorker] camera-moving-handler moving err: %v, event: %+v", err, evt)
			}
		} else {
			logrus.Errorf("[PowerWorker] camera-moving-handler Data type err: %s", "not PowerMovingRequest struct")
		}
	})
}

func (c *PowerWorker) WorkerID() string {
	return powerWorkerID
}

func (c *PowerWorker) Restart() error {
	return nil
}
