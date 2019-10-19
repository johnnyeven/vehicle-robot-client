package robot

import (
	"fmt"
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/shantanubhadoria/go-kalmanfilter/kalmanfilter"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/firmata"
	"math"
	"time"
)

const (
	attitudeMPU6050WorkerID = "attitude-mpu6050-worker"
	AttitudeBroadcastTopic  = "attitude.broadcast"
	attitudeGravityRectify  = math.MaxInt16 / 2
)

type AttitudeMPU6050Worker struct {
	sensor *i2c.MPU6050Driver
	bus    *bus.MessageBus

	calibrationOffset Attitude
	calibrationTimes  int
	AttitudeWorker
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
		sensor:            sensor,
		bus:               bus,
		calibrationTimes:  1000,
		calibrationOffset: Attitude{},
		AttitudeWorker: AttitudeWorker{
			data:        Attitude{},
			kalmanRoll:  &kalmanfilter.FilterData{},
			kalmanPitch: &kalmanfilter.FilterData{},
		},
	}
}

func (a *AttitudeMPU6050Worker) WorkerID() string {
	return attitudeMPU6050WorkerID
}

func (a *AttitudeMPU6050Worker) Start() {
	a.calibration()
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

		a.rectify()
		a.calcAngle()
		a.bus.Emit(AttitudeBroadcastTopic, a.data, "")
	})
}

func (a *AttitudeMPU6050Worker) Restart() error {
	return nil
}

func (a *AttitudeMPU6050Worker) Stop() error {
	return nil
}

// 数据校准
func (a *AttitudeMPU6050Worker) calibration() error {
	logrus.Info("[AttitudeMPU6050Worker] calibration...")
	defer func() {
		logrus.Infof("[AttitudeMPU6050Worker] calibration complete, offset: %+v", a.calibrationOffset)
	}()
	totalCalibration := Attitude{
		Accelerometer: ThreeDDataCalibration{},
		Gyroscope:     ThreeDDataCalibration{},
		Compass:       ThreeDDataCalibration{},
		Temperature:   0,
	}
	for i := 0; i < a.calibrationTimes; i++ {
		err := a.sensor.GetData()
		if err != nil {
			logrus.Errorf("[AttitudeMPU6050Worker] sensor.GetData() err: %v", err)
			return err
		}
		totalCalibration.Accelerometer.X += float64(a.sensor.Accelerometer.X)
		totalCalibration.Accelerometer.Y += float64(a.sensor.Accelerometer.Y)
		totalCalibration.Accelerometer.Z += float64(a.sensor.Accelerometer.Z)

		totalCalibration.Gyroscope.X += float64(a.sensor.Gyroscope.X)
		totalCalibration.Gyroscope.Y += float64(a.sensor.Gyroscope.Y)
		totalCalibration.Gyroscope.Z += float64(a.sensor.Gyroscope.Z)
	}

	a.calibrationOffset.Accelerometer.X = totalCalibration.Accelerometer.X / float64(a.calibrationTimes)
	a.calibrationOffset.Accelerometer.Y = totalCalibration.Accelerometer.Y / float64(a.calibrationTimes)
	a.calibrationOffset.Accelerometer.Z = totalCalibration.Accelerometer.Z/float64(a.calibrationTimes) + attitudeGravityRectify // 需要抵消初始垂直向下的地球重力加速度，因为传感器感知范围为2g，故最大值除以2为1g

	a.calibrationOffset.Gyroscope.X = totalCalibration.Gyroscope.X / float64(a.calibrationTimes)
	a.calibrationOffset.Gyroscope.Y = totalCalibration.Gyroscope.Y / float64(a.calibrationTimes)
	a.calibrationOffset.Gyroscope.Z = totalCalibration.Gyroscope.Z / float64(a.calibrationTimes)

	return nil
}
