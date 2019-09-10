package main

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules/controllers"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/opencv"
)

func CreateRobotFromConfig(config global.RobotConfiguration, messageBus *bus.MessageBus, robotClient *client.RobotClient) *gobot.Master {
	devices := make([]gobot.Device, 0)
	connections := make([]gobot.Connection, 0)
	moduleWorkers := make([]func(), 0)

	if config.ActivateFirmata.True() {
		firmataAdaptor := firmata.NewAdaptor(config.ArduinoDeviceID)
		connections = append(connections, firmataAdaptor)

		if config.ActivateCameraController.True() {
			servoHorizon := gpio.NewServoDriver(firmataAdaptor, config.ServoHorizonPin)
			servoVertical := gpio.NewServoDriver(firmataAdaptor, config.ServoVerticalPin)

			devices = append(devices, servoHorizon, servoVertical)
			moduleWorkers = append(moduleWorkers, func() {
				go controllers.CameraHolderController(servoHorizon, servoVertical)
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

	if config.ActivateCameraController.True() {
		window := opencv.NewWindowDriver()
		camera := opencv.NewCameraDriver(0)

		devices = append(devices, window, camera)
		moduleWorkers = append(moduleWorkers, func() {
			go controllers.ObjectDetectiveController(window, camera, robotClient)
		})
	}

	master := gobot.NewMaster()

	if config.ActivateApiSupport.True() {
		apiServer := api.NewAPI(master)
		apiServer.Port = config.APIServerPort
		apiServer.Start()
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

	master.AddRobot(robot)

	return master
}
