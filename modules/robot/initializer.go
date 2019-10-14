package robot

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules/robot/workers"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gocv.io/x/gocv"
)

func CreateRobotFromConfig(config global.RobotConfiguration, messageBus *bus.MessageBus, robotClient *client.RobotClient) *gobot.Master {
	devices := make([]gobot.Device, 0)
	connections := make([]gobot.Connection, 0)
	moduleWorkers := make([]func(), 0)

	if config.ActivateFirmata.True() {
		firmataAdaptor := firmata.NewAdaptor(config.ArduinoDeviceID)
		connections = append(connections, firmataAdaptor)

		if config.ActivateCameraHolderController.True() {
			servoHorizon := gpio.NewServoDriver(firmataAdaptor, config.ServoHorizonPin)
			servoVertical := gpio.NewServoDriver(firmataAdaptor, config.ServoVerticalPin)

			devices = append(devices, servoHorizon, servoVertical)
			moduleWorkers = append(moduleWorkers, func() {
				workers.CameraHolderController(servoHorizon, servoVertical, messageBus)
			})
		}

		if config.ActivatePowerController.True() {
			motorLeft := gpio.NewMotorDriver(firmataAdaptor, config.LeftMotorSpeedPin)
			motorLeft.DirectionPin = config.LeftMotorDirectionPin
			motorRight := gpio.NewMotorDriver(firmataAdaptor, config.RightMotorSpeedPin)
			motorRight.DirectionPin = config.RightMotorDirectionPin
			powerController := workers.NewPowerController(motorLeft, motorRight, messageBus)

			devices = append(devices, motorLeft, motorRight)
			moduleWorkers = append(moduleWorkers, func() {
				powerController.Start()
			})
		}
	}

	if config.ActivateCameraController.True() {
		camera, err := gocv.VideoCaptureDevice(0)
		if err != nil {
			logrus.Panicf("gocv.VideoCaptureDevice err: %v", err)
		}
		camera.Set(gocv.VideoCaptureFrameWidth, 320)
		camera.Set(gocv.VideoCaptureFrameHeight, 240)
		camera.Set(gocv.VideoCaptureFPS, 1)

		moduleWorkers = append(moduleWorkers, func() {
			workers.ObjectDetectiveController(config, camera, robotClient)
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
