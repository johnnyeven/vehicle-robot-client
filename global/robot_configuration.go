package global

import (
	"github.com/johnnyeven/service-vehicle-robot/constants/types"
)

type RobotConfiguration struct {
	CameraMode types.CameraMode `json:"cameraMode"`
}

func (c RobotConfiguration) Init() {

}
