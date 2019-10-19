package robot

import (
	"github.com/shantanubhadoria/go-kalmanfilter/kalmanfilter"
	"math"
	"time"
)

const (
	attitudeGravityRectify = 256    // 根据芯片手册查询加速度计比例因子
	attitudeGyroRectify    = 14.375 // 根据芯片手册查询陀螺仪比例因子
)

type AttitudeWorker struct {
	calibrationOffset Attitude
	data              Attitude

	lastTime    time.Time
	kalmanRoll  *kalmanfilter.FilterData
	kalmanPitch *kalmanfilter.FilterData
}

// 重力加速度转换为N * 1g，角度转换为N * 1degress/sec
// 计算欧拉角
func (a *AttitudeWorker) rectify() {
	a.data.Accelerometer.X = (a.data.Accelerometer.X - a.calibrationOffset.Accelerometer.X) / attitudeGravityRectify
	a.data.Accelerometer.Y = (a.data.Accelerometer.Y - a.calibrationOffset.Accelerometer.Y) / attitudeGravityRectify
	a.data.Accelerometer.Z = (a.data.Accelerometer.Z - a.calibrationOffset.Accelerometer.Z) / attitudeGravityRectify

	a.data.Temperature = a.data.Temperature/340.0 + 36.53

	a.data.Gyroscope.X = (a.data.Gyroscope.X - a.calibrationOffset.Gyroscope.X) / attitudeGyroRectify
	a.data.Gyroscope.Y = (a.data.Gyroscope.Y - a.calibrationOffset.Gyroscope.Y) / attitudeGyroRectify
	a.data.Gyroscope.Z = (a.data.Gyroscope.Z - a.calibrationOffset.Gyroscope.Z) / attitudeGyroRectify
}

func (a *AttitudeWorker) calcAngle() {
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

func (a *AttitudeWorker) getAngle(x, y float64, normal float64) float64 {
	normalXY := math.Sqrt(x*x + y*y)
	return math.Acos(normalXY/normal) * 180 / math.Pi
}
