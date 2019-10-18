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
	AttitudeBroadcastTopic  = "attitude.broadcast"
)

type AttitudeMPU6050Worker struct {
	sensor *i2c.MPU6050Driver
	bus    *bus.MessageBus

	data Attitude
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
		data:   Attitude{},
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
		a.data.Accelerometer.X = float64(a.sensor.Accelerometer.X)
		a.data.Accelerometer.Y = float64(a.sensor.Accelerometer.Y)
		a.data.Accelerometer.Z = float64(a.sensor.Accelerometer.Z)

		a.data.Gyroscope.X = float64(a.sensor.Gyroscope.X)
		a.data.Gyroscope.Y = float64(a.sensor.Gyroscope.Y)
		a.data.Gyroscope.Z = float64(a.sensor.Gyroscope.Z)

		a.data.Temperature = float64(a.sensor.Temperature)

		a.bus.Emit(AttitudeBroadcastTopic, a.data, "")
	})
}

func (a *AttitudeMPU6050Worker) Restart() error {
	return nil
}

func (a *AttitudeMPU6050Worker) Stop() error {
	return nil
}
