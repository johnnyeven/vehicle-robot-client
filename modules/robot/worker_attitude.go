package robot

import (
	"github.com/johnnyeven/robot-library/drivers"
	"github.com/shantanubhadoria/go-kalmanfilter/kalmanfilter"
	"math"
	"time"
)

type AttitudeWorker struct {
	Data Attitude

	lastTime    time.Time
	kalmanRoll  *kalmanfilter.FilterData
	kalmanPitch *kalmanfilter.FilterData
}

// 重力加速度转换为N * 1g，角度转换为N * 1degress/sec
// 计算欧拉角
func (a *AttitudeWorker) rectify() {
	a.Data.Accelerometer.X = a.Data.Accelerometer.X / drivers.Adxl345LSB
	a.Data.Accelerometer.Y = a.Data.Accelerometer.Y / drivers.Adxl345LSB
	a.Data.Accelerometer.Z = a.Data.Accelerometer.Z / drivers.Adxl345LSB

	a.Data.Temperature = a.Data.Temperature/340.0 + 36.53

	a.Data.Gyroscope.X = a.Data.Gyroscope.X / drivers.ITG3200LSB
	a.Data.Gyroscope.Y = a.Data.Gyroscope.Y / drivers.ITG3200LSB
	a.Data.Gyroscope.Z = a.Data.Gyroscope.Z / drivers.ITG3200LSB
}

func (a *AttitudeWorker) calcAngle() {
	// 加速度向量模长
	accNormal := math.Sqrt(a.Data.Accelerometer.X*a.Data.Accelerometer.X + a.Data.Accelerometer.Y*a.Data.Accelerometer.Y + a.Data.Accelerometer.Z*a.Data.Accelerometer.Z)

	// 计算滚转角X
	rollAngle := a.getAngle(a.Data.Accelerometer.X, a.Data.Accelerometer.Z, accNormal)
	if a.Data.Accelerometer.Y > 0 {
		rollAngle = -rollAngle
	}

	// 计算俯仰角Y
	pitchAngle := a.getAngle(a.Data.Accelerometer.Y, a.Data.Accelerometer.Z, accNormal)
	if a.Data.Accelerometer.X < 0 {
		pitchAngle = -pitchAngle
	}

	currentTime := time.Now()
	duration := currentTime.Sub(a.lastTime)
	rollAngle = a.kalmanRoll.Update(rollAngle, a.Data.Gyroscope.Y, float64(duration/time.Second))
	pitchAngle = a.kalmanPitch.Update(pitchAngle, a.Data.Gyroscope.Z, float64(duration/time.Second))

	a.Data.EulerAngle.X = rollAngle
	a.Data.EulerAngle.Y = pitchAngle
}

func (a *AttitudeWorker) getAngle(x, y float64, normal float64) float64 {
	normalXY := math.Sqrt(x*x + y*y)
	return math.Acos(normalXY/normal) * 180 / math.Pi
}
