package robot

import "gobot.io/x/gobot/drivers/i2c"

const (
	ServoMaxAngle    uint8 = 180
	ServoCentreAngle uint8 = 90
)

type Attitude struct {
	Accelerometer i2c.ThreeDData
	Gyroscope     i2c.ThreeDData
	Temperature   int16
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
