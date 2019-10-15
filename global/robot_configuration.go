package global

import (
	"github.com/johnnyeven/libtools/courier/enumeration"
	"github.com/johnnyeven/service-vehicle-robot/constants/types"
)

type RobotConfiguration struct {
	// 机器人模式
	RobotMode types.RobotMode `json:"robotMode"`
	// 总控开关（不启用无法启用任何其他模块，仅用于调试）
	ActivateFirmata enumeration.Bool `json:"activateFirmata"`
	// Firmata连接名称
	FirmataConnectionName string `json:"firmataConnectionName" default:"firmataConnection"`
	// 启用API
	ActivateApiSupport enumeration.Bool `json:"activateApiSupport"`
	// 是否启用摄像头模块（不启用无法开启视频同步及物体识别）
	ActivateCameraController enumeration.Bool `json:"activateCameraController"`
	// 是否启用图传
	ActivateCameraTransfer enumeration.Bool `json:"activateCameraTransfer"`
	// 摄像头图像宽度
	CameraCaptureWidth uint64 `json:"cameraCaptureWidth,string" default:"640"`
	// 摄像头图像高度
	CameraCaptureHeight uint64 `json:"cameraCaptureHeight,string" default:"480"`
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
	// 摄像头云台水平舵机名称
	ServoHorizonName string `json:"servoHorizonName" default:"horizonServo"`
	// 摄像头云台垂直舵机信号针脚编号
	ServoVerticalPin string `json:"servoVerticalPin" default:"9"`
	// 摄像头云台垂直舵机名称
	ServoVerticalName string `json:"servoVerticalName" default:"verticalServo"`
	// 动力系统左电机转速控制针脚编号
	LeftMotorSpeedPin string `json:"leftMotorSpeedPin" default:"11"`
	// 动力系统左电机正反转控制针脚编号
	LeftMotorDirectionPin string `json:"leftMotorDirectionPin" default:"13"`
	// 动力系统左电机名称
	LeftMotorName string `json:"leftMotorName" default:"leftMotor"`
	// 动力系统右电机转速控制针脚编号
	RightMotorSpeedPin string `json:"rightMotorSpeedPin" default:"10"`
	// 动力系统右电机正反转控制针脚编号
	RightMotorDirectionPin string `json:"rightMotorDirectionPin" default:"12"`
	// 动力系统右电机名称
	RightMotorName string `json:"rightMotorName" default:"rightMotor"`
}
