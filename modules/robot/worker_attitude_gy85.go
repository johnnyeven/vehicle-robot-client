package robot

import (
	"fmt"
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/robot-library/drivers"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/shantanubhadoria/go-kalmanfilter/kalmanfilter"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/firmata"
	"math"
	"time"
)

const (
	attitudeGY85WorkerID   = "attitude-gy85-worker"
	attitudeGravityRectify = math.MaxInt16 / 2   // 加速度计分辨率为 2 m/2s
	attitudeGyroRectify    = math.MaxInt16 / 250 // 陀螺仪分辨率为 250 degrees/sec
)

type AttitudeGY85Worker struct {
	accSensor     *drivers.ADXL345Driver
	gyroSensor    *drivers.ITG3200Driver
	compassSensor *drivers.HMC5883Driver
	bus           *bus.MessageBus

	calibrationTimes  int
	calibrationOffset Attitude
	data              Attitude

	lastTime    time.Time
	kalmanRoll  *kalmanfilter.FilterData
	kalmanPitch *kalmanfilter.FilterData
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
		accSensor:         accSensor,
		gyroSensor:        gyroSensor,
		compassSensor:     compassSensor,
		bus:               bus,
		calibrationOffset: Attitude{},
		data:              Attitude{},
		kalmanRoll:        &kalmanfilter.FilterData{},
		kalmanPitch:       &kalmanfilter.FilterData{},
	}
}

func (a *AttitudeGY85Worker) WorkerID() string {
	return attitudeGY85WorkerID
}

func (a *AttitudeGY85Worker) Start() {
	a.calibration()
	gobot.Every(10*time.Millisecond, func() {
		a.lastTime = time.Now()
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
		a.data.Accelerometer.X = float64(a.accSensor.Accelerometer.X)
		a.data.Accelerometer.Y = float64(a.accSensor.Accelerometer.Y)
		a.data.Accelerometer.Z = float64(a.accSensor.Accelerometer.Z)

		a.data.Gyroscope.X = float64(a.gyroSensor.Gyroscope.X)
		a.data.Gyroscope.Y = float64(a.gyroSensor.Gyroscope.Y)
		a.data.Gyroscope.Z = float64(a.gyroSensor.Gyroscope.Z)

		a.data.Temperature = float64(a.gyroSensor.Temperature)

		a.data.Compass.X = float64(a.compassSensor.Compass.X)
		a.data.Compass.Y = float64(a.compassSensor.Compass.Y)
		a.data.Compass.Z = float64(a.compassSensor.Compass.Z)

		a.rectify()
		a.calcAngle()
		a.bus.Emit(AttitudeBroadcastTopic, a.data, "")
	})
}

// 数据校准
func (a *AttitudeGY85Worker) calibration() error {
	logrus.Info("[AttitudeGY85Worker] calibration...")
	defer func() {
		logrus.Infof("[AttitudeGY85Worker] calibration complete, offset: %+v", a.calibrationOffset)
	}()
	totalCalibration := Attitude{
		Accelerometer: ThreeDDataCalibration{},
		Gyroscope:     ThreeDDataCalibration{},
		Compass:       ThreeDDataCalibration{},
		Temperature:   0,
	}
	for i := 0; i < a.calibrationTimes; i++ {
		err := a.accSensor.GetData()
		if err != nil {
			logrus.Errorf("[AttitudeGY85Worker] accSensor.GetData() err: %v", err)
			return err
		}
		err = a.gyroSensor.GetData()
		if err != nil {
			logrus.Errorf("[AttitudeGY85Worker] gyroSensor.GetData() err: %v", err)
			return err
		}
		err = a.compassSensor.GetData()
		if err != nil {
			logrus.Errorf("[AttitudeGY85Worker] compassSensor.GetData() err: %v", err)
			return err
		}
		totalCalibration.Accelerometer.X += float64(a.accSensor.Accelerometer.X)
		totalCalibration.Accelerometer.Y += float64(a.accSensor.Accelerometer.Y)
		totalCalibration.Accelerometer.Z += float64(a.accSensor.Accelerometer.Z)

		totalCalibration.Gyroscope.X += float64(a.gyroSensor.Gyroscope.X)
		totalCalibration.Gyroscope.Y += float64(a.gyroSensor.Gyroscope.Y)
		totalCalibration.Gyroscope.Z += float64(a.gyroSensor.Gyroscope.Z)

		totalCalibration.Temperature += float64(a.gyroSensor.Temperature)

		totalCalibration.Compass.X += float64(a.compassSensor.Compass.X)
		totalCalibration.Compass.Y += float64(a.compassSensor.Compass.Y)
		totalCalibration.Compass.Z += float64(a.compassSensor.Compass.Z)
	}

	a.calibrationOffset.Accelerometer.X = totalCalibration.Accelerometer.X / float64(a.calibrationTimes)
	a.calibrationOffset.Accelerometer.Y = totalCalibration.Accelerometer.Y / float64(a.calibrationTimes)
	a.calibrationOffset.Accelerometer.Z = totalCalibration.Accelerometer.Z/float64(a.calibrationTimes) + attitudeGravityRectify // 需要抵消初始垂直向下的地球重力加速度，因为传感器感知范围为2g，故最大值除以2为1g

	a.calibrationOffset.Gyroscope.X = totalCalibration.Gyroscope.X / float64(a.calibrationTimes)
	a.calibrationOffset.Gyroscope.Y = totalCalibration.Gyroscope.Y / float64(a.calibrationTimes)
	a.calibrationOffset.Gyroscope.Z = totalCalibration.Gyroscope.Z / float64(a.calibrationTimes)

	return nil
}

// 重力加速度转换为N * 1g，角度转换为N * 1degress/sec
// 计算欧拉角
func (a *AttitudeGY85Worker) rectify() {
	a.data.Accelerometer.X = (a.data.Accelerometer.X - a.calibrationOffset.Accelerometer.X) / attitudeGravityRectify
	a.data.Accelerometer.Y = (a.data.Accelerometer.Y - a.calibrationOffset.Accelerometer.Y) / attitudeGravityRectify
	a.data.Accelerometer.Z = (a.data.Accelerometer.Z - a.calibrationOffset.Accelerometer.Z) / attitudeGravityRectify

	a.data.Temperature = a.data.Temperature/340.0 + 36.53

	a.data.Gyroscope.X = (a.data.Gyroscope.X - a.calibrationOffset.Gyroscope.X) / attitudeGyroRectify
	a.data.Gyroscope.Y = (a.data.Gyroscope.Y - a.calibrationOffset.Gyroscope.Y) / attitudeGyroRectify
	a.data.Gyroscope.Z = (a.data.Gyroscope.Z - a.calibrationOffset.Gyroscope.Z) / attitudeGyroRectify
}

func (a *AttitudeGY85Worker) calcAngle() {
	// 加速度向量模长
	accNormal := math.Sqrt(a.data.Accelerometer.X*a.data.Accelerometer.X + a.data.Accelerometer.Y*a.data.Accelerometer.Y + a.data.Accelerometer.Z*a.data.Accelerometer.Z)

	// 计算滚转角X
	rollAngle := a.getAngle(a.data.Accelerometer.X, a.data.Accelerometer.Z, accNormal)
	if a.data.Accelerometer.Y > 0 {
		rollAngle = -rollAngle
	}

	// 计算俯仰角Y
	pitchAngle := a.getAngle(a.data.Accelerometer.Y, a.data.Accelerometer.Z, accNormal)
	if a.data.Accelerometer.X < 0 {
		pitchAngle = -pitchAngle
	}

	currentTime := time.Now()
	duration := currentTime.Sub(a.lastTime)
	rollAngle = a.kalmanRoll.Update(rollAngle, a.data.Gyroscope.Y, float64(duration/time.Second))
	pitchAngle = a.kalmanPitch.Update(pitchAngle, a.data.Gyroscope.Z, float64(duration/time.Second))

	a.data.EulerAngle.X = rollAngle
	a.data.EulerAngle.Y = pitchAngle
}

func (a *AttitudeGY85Worker) getAngle(x, y float64, normal float64) float64 {
	normalXY := math.Sqrt(x*x + y*y)
	return math.Acos(normalXY/normal) * 180 / math.Pi
}

func (a *AttitudeGY85Worker) Restart() error {
	return nil
}

func (a *AttitudeGY85Worker) Stop() error {
	return nil
}
