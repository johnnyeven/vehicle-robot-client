package robot

const (
	ServoMaxAngle    uint8 = 180
	ServoCentreAngle uint8 = 90
)

type Distance struct {
	Angle    uint8
	Distance float64
}

type Attitude struct {
	// 加速度
	Accelerometer ThreeDDataCalibration
	// 陀螺仪
	Gyroscope ThreeDDataCalibration
	// 磁力计
	Compass ThreeDDataCalibration
	// 欧拉角
	EulerAngle ThreeDDataCalibration
	// 温度
	Temperature float64
}

type ThreeDDataCalibration struct {
	X float64
	Y float64
	Z float64
}

func servoAngleChange(current uint8, offset float64) uint8 {
	current = uint8(float64(current) + offset)
	if current < 0 {
		current = 0
	} else if current > ServoMaxAngle {
		current = ServoMaxAngle
	}
	return current
}

func servoAngle(target uint8) uint8 {
	if target < 0 {
		target = 0
	} else if target > ServoMaxAngle {
		target = ServoMaxAngle
	}
	return target
}
