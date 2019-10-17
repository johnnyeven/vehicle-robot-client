package robot

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/robot-library/drivers"
	"github.com/johnnyeven/vehicle-robot-client/global"
	bus2 "github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"time"
)

const (
	distanceHCSR04WorkerID    = "distance-hcsr04-worker"
	DistanceServoTopic        = "distance.servo"
	DistanceServoEventHandler = "distance-servo-handler"
	DistanceServoResultTopic  = "distance.servo.result"
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

		distance, err := d.measure(e.Data.(uint8))
		if err != nil {
			return
		}
		d.bus.Emit(DistanceServoResultTopic, distance, "")
	})

	var offset uint8 = 5
	for {
		if d.currentHorizonAngle < 0 || d.currentHorizonAngle > 180 {
			offset = -offset
		}
		distance, err := d.measure(d.currentHorizonAngle + offset)
		if err != nil {
			return
		}
		logrus.Infof("angle: %d, distance: %.2f cm", d.currentHorizonAngle, distance)

		time.Sleep(10 * time.Millisecond)
	}
}

func (d *DistanceHCSR04Worker) measure(angle uint8) (float64, error) {
	d.currentHorizonAngle = servoAngle(angle)
	err := d.servoHorizon.Move(d.currentHorizonAngle)
	if err != nil {
		logrus.Errorf("[DistanceHCSR04Worker] %s servoHorizon.Move err: %v, angle: %d", DistanceServoEventHandler, err, d.currentHorizonAngle)
		return 0, err
	}
	distance, err := d.sensor.Measure()
	if err != nil {
		logrus.Errorf("[DistanceHCSR04Worker] %s sensor.Measure err: %v, angle: %d", DistanceServoEventHandler, err, d.currentHorizonAngle)
		return 0, err
	}
	return distance, nil
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
