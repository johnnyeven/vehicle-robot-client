package robot

import (
	"fmt"
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
	"time"
)

const (
	attitudeMPU6050WorkerID = "attitude-mpu6050-worker"
	AttitudeTopic           = "attitude"
)

type AttitudeMPU6050Worker struct {
	sensor *i2c.MPU6050Driver
	bus    *bus.MessageBus
}

func NewAttitudeMPU6050Worker(robot *Robot, bus *bus.MessageBus, config *global.RobotConfiguration) *AttitudeMPU6050Worker {
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

	sensor := i2c.NewMPU6050Driver(firmataAdaptor)
	sensor.SetName(config.AttitudeName)
	robot.AddDevice(sensor)

	return &AttitudeMPU6050Worker{
		sensor: sensor,
		bus:    bus,
	}
}

func (a *AttitudeMPU6050Worker) WorkerID() string {
	return attitudeMPU6050WorkerID
}

func (a *AttitudeMPU6050Worker) Start() {
	gobot.Every(10*time.Millisecond, func() {
		err := a.sensor.GetData()
		if err != nil {
			logrus.Errorf("[AttitudeMPU6050Worker] sensor.GetData() err: %v", err)
			return
		}
		fmt.Printf("\rAcc: %v, Gyr: %v, Temp: %d", a.sensor.Accelerometer, a.sensor.Gyroscope, a.sensor.Temperature)
		data := Attitude{
			Accelerometer: a.sensor.Accelerometer,
			Gyroscope:     a.sensor.Gyroscope,
			Temperature:   a.sensor.Temperature,
		}
		a.bus.Emit(AttitudeTopic, data, "")
	})
}

func (a *AttitudeMPU6050Worker) Restart() error {
	return nil
}

func (a *AttitudeMPU6050Worker) Stop() error {
	return nil
}
