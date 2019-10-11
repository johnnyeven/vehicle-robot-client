package global

import (
	"github.com/johnnyeven/libtools/courier/enumeration"
	"github.com/johnnyeven/service-vehicle-robot/constants/types"
)

type RobotConfiguration struct {
	// 总控开关（不启用无法启用任何其他模块，仅用于调试）
	ActivateFirmata enumeration.Bool `json:"activateFirmata"`
	// 启用API
	ActivateApiSupport enumeration.Bool `json:"activateApiSupport"`
	// 是否启用摄像头模块（不启用无法开启视频同步及物体识别）
	ActivateCameraController enumeration.Bool `json:"activateCameraController"`
	// 是否启用摄像头云台控制模块
	ActivateCameraHolderController enumeration.Bool `json:"activateCameraHolderController"`
	// 是否启用动力系统模块（不启用无法行走）
	ActivatePowerController enumeration.Bool `json:"activatePowerController"`

	// API服务端口号
	APIServerPort string `json:"apiServerPort"`

	// 摄像头模式（仅视频同步或者开启物体识别）
	CameraMode types.CameraMode `json:"cameraMode"`

	// Arduino设备文件名
	ArduinoDeviceID string `json:"arduinoDeviceID"`
	// 摄像头云台水平舵机信号针脚编号
	ServoHorizonPin string `json:"servoHorizonPin" default:"8"`
	// 摄像头云台垂直舵机信号针脚编号
	ServoVerticalPin string `json:"servoVerticalPin" default:"9"`
	// 动力系统左电机转速控制针脚编号
	LeftMotorSpeedPin string `json:"leftMotorSpeedPin" default:"11"`
	// 动力系统左电机正反转控制针脚编号
	LeftMotorDirectionPin string `json:"leftMotorDirectionPin" default:"13"`
	// 动力系统右电机转速控制针脚编号
	RightMotorSpeedPin string `json:"rightMotorSpeedPin" default:"10"`
	// 动力系统右电机正反转控制针脚编号
	RightMotorDirectionPin string `json:"rightMotorDirectionPin" default:"12"`
}
