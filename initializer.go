package main

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules/controllers"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/opencv"
)

func CreateRobotFromConfig(config global.RobotConfiguration, messageBus *bus.MessageBus) *gobot.Robot {
	devices := make([]gobot.Device, 0)
	connections := make([]gobot.Connection, 0)
	moduleWorkers := make([]func(), 0)

	if config.ActivateFirmata.True() {
		firmataAdaptor := firmata.NewAdaptor(config.ArduinoDeviceID)
		connections = append(connections, firmataAdaptor)

		if config.ActivateCameraController.True() {
			servoHorizon := gpio.NewServoDriver(firmataAdaptor, config.ServoHorizonPin)
			servoVertical := gpio.NewServoDriver(firmataAdaptor, config.ServoVerticalPin)
			window := opencv.NewWindowDriver()
			camera := opencv.NewCameraDriver(0)

			devices = append(devices, window, camera, servoHorizon, servoVertical)
			moduleWorkers = append(moduleWorkers, func() {
				go controllers.CameraHolderController(servoHorizon, servoVertical)
				go controllers.ObjectDetectiveController(window, camera, global.Config.RobotClient)
			})
		}

		if config.ActivatePowerController.True() {
			motorLeft := gpio.NewMotorDriver(firmataAdaptor, config.LeftMotorSpeedPin)
			motorLeft.DirectionPin = config.LeftMotorDirectionPin
			motorRight := gpio.NewMotorDriver(firmataAdaptor, config.RightMotorSpeedPin)
			motorRight.DirectionPin = config.RightMotorDirectionPin
			powerController := controllers.NewPowerController(motorLeft, motorRight, messageBus)

			devices = append(devices, motorLeft, motorRight)
			moduleWorkers = append(moduleWorkers, func() {
				powerController.Start()
			})
		}
	}

	robot := gobot.NewRobot("VehicleRobot",
		connections,
		devices,
		func() {
			for _, worker := range moduleWorkers {
				go worker()
			}
		},
	)

	return robot
}
