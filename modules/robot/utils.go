package robot

const (
	ServoMaxAngle    uint8 = 180
	ServoCentreAngle uint8 = 90
)

func servoAngleChange(current uint8, offset float64) uint8 {
	current = uint8(float64(current) + offset)
	if current < 0 {
		current = 0
	} else if current > ServoMaxAngle {
		current = ServoMaxAngle
	}
	return current
}
