package robot

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/robot-library/drivers"
	"github.com/johnnyeven/vehicle-robot-client/global"
	bus2 "github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

const (
	distanceHCSR04WorkerID    = "distance-hcsr04-worker"
	DistanceServoTopic        = "distance.servo"
	DistanceServoEventHandler = "distance-servo-handler"
)

type DistanceHCSR04Worker struct {
	sensor       *drivers.HCSR04Driver
	servoHorizon *gpio.ServoDriver
	bus          *bus.MessageBus

	currentHorizonAngle uint8
	manualControl       bool
}

func (d *DistanceHCSR04Worker) WorkerID() string {
	return distanceHCSR04WorkerID
}

func (d *DistanceHCSR04Worker) Start() {
	logrus.Infof("[DistanceHCSR04Worker] Init servo to center angle: %d", ServoCentreAngle)
	err := d.servoHorizon.Move(ServoCentreAngle)
	if err != nil {
		logrus.Errorf("[DistanceHCSR04Worker] horizon servo move failed with err: %v", err)
		return
	}
	d.bus.RegisterTopic(DistanceServoTopic)
	d.bus.RegisterHandler(DistanceServoEventHandler, DistanceServoTopic, func(e *bus2.Event) {
		d.manualControl = true
		defer func() {
			d.manualControl = false
		}()
	})
}

func (d *DistanceHCSR04Worker) Restart() error {
	return nil
}

func (d *DistanceHCSR04Worker) Stop() error {
	return nil
}

func NewDistanceHCSR04Worker(robot *Robot, bus *bus.MessageBus, config *global.RobotConfiguration) *DistanceHCSR04Worker {
	var firmataAdaptor *firmata.Adaptor
	var ok bool
	conn := robot.GetConnection(config.FirmataConnectionName)
	if conn == nil {
		firmataAdaptor = firmata.NewAdaptor(config.ArduinoDeviceID)
		firmataAdaptor.SetName(config.FirmataConnectionName)
		robot.AddConnection(firmataAdaptor)
	} else {
		if firmataAdaptor, ok = conn.(*firmata.Adaptor); !ok {
			logrus.Panicf("[CameraHolderWorker] 连接器已存在，但并不是 *firmata.Adaptor 类型")
		}
	}

	sensor := drivers.NewHCSR04Driver(firmataAdaptor, config.DistanceTrigPin, config.DistanceEchoPin)
	sensor.SetName(config.DistanceName)
	servo := gpio.NewServoDriver(firmataAdaptor, config.DistanceServoHorizonPin)
	servo.SetName(config.DistanceServoHorizonName)

	robot.AddDevice(sensor, servo)

	return &DistanceHCSR04Worker{
		sensor:              sensor,
		servoHorizon:        servo,
		bus:                 bus,
		currentHorizonAngle: ServoCentreAngle,
		manualControl:       false,
	}
}
