package robot

import (
	"fmt"
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/robot-library/drivers"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/firmata"
	"time"
)

const (
	attitudeGY85WorkerID = "attitude-gy85-worker"
)

type AttitudeGY85Worker struct {
	accSensor     *drivers.ADXL345Driver
	gyroSensor    *drivers.ITG3200Driver
	compassSensor *drivers.HMC5883Driver
	bus           *bus.MessageBus

	data Attitude
}

func NewAttitudeGY85Worker(robot *Robot, bus *bus.MessageBus, config *global.RobotConfiguration) *AttitudeGY85Worker {
	var firmataAdaptor *firmata.Adaptor
	var ok bool
	conn := robot.GetConnection(config.FirmataConnectionName)
	if conn == nil {
		firmataAdaptor = firmata.NewAdaptor(config.ArduinoDeviceID)
		firmataAdaptor.SetName(config.FirmataConnectionName)
		robot.AddConnection(firmataAdaptor)
	} else {
		if firmataAdaptor, ok = conn.(*firmata.Adaptor); !ok {
			logrus.Panicf("[NewAttitudeGY85Worker] 连接器已存在，但并不是 *firmata.Adaptor 类型")
		}
	}

	accSensor := drivers.NewADXL345Driver(firmataAdaptor)
	accSensor.SetName(config.AttitudeName + "_ACC")
	gyroSensor := drivers.NewITG3200Driver(firmataAdaptor)
	gyroSensor.SetName(config.AttitudeName + "_GYRO")
	compassSensor := drivers.NewHMC5883Driver(firmataAdaptor)
	compassSensor.SetName(config.AttitudeName + "_COMPASS")
	robot.AddDevice(accSensor, gyroSensor, compassSensor)

	return &AttitudeGY85Worker{
		accSensor:     accSensor,
		gyroSensor:    gyroSensor,
		compassSensor: compassSensor,
		bus:           bus,
		data:          Attitude{},
	}
}

func (a *AttitudeGY85Worker) WorkerID() string {
	return attitudeGY85WorkerID
}

func (a *AttitudeGY85Worker) Start() {
	gobot.Every(10*time.Millisecond, func() {
		err := a.accSensor.GetData()
		if err != nil {
			logrus.Errorf("[AttitudeGY85Worker] accSensor.GetData() err: %v", err)
			return
		}
		err = a.gyroSensor.GetData()
		if err != nil {
			logrus.Errorf("[AttitudeGY85Worker] gyroSensor.GetData() err: %v", err)
			return
		}
		err = a.compassSensor.GetData()
		if err != nil {
			logrus.Errorf("[AttitudeGY85Worker] compassSensor.GetData() err: %v", err)
			return
		}
		fmt.Printf("\rAcc: %v, Gyr: %v, Temp: %d, Compass: %v", a.accSensor.Accelerometer, a.gyroSensor.Gyroscope, a.gyroSensor.Temperature, a.compassSensor.Compass)
		a.data.Accelerometer = a.accSensor.Accelerometer
		a.data.Gyroscope = a.gyroSensor.Gyroscope
		a.data.Temperature = a.gyroSensor.Temperature
		a.data.Compass = a.compassSensor.Compass
		a.bus.Emit(AttitudeBroadcastTopic, a.data, "")
	})
}

func (a *AttitudeGY85Worker) Restart() error {
	return nil
}

func (a *AttitudeGY85Worker) Stop() error {
	return nil
}
