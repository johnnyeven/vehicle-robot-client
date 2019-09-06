package global

import (
	"github.com/johnnyeven/service-vehicle-robot/constants/types"
)

type RobotConfiguration struct {
	CameraMode       types.CameraMode `json:"cameraMode"`
	ArduinoDeviceID  string           `json:"arduinoDeviceID"`
	ServoHorizonPin  string           `json:"servoHorizonPin" default:"8"`
	ServoVerticalPin string           `json:"servoVerticalPin" default:"9"`
}

func (c RobotConfiguration) Init() {

}
