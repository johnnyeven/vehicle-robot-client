package robot

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/service-vehicle-robot/constants/types"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
)

func init() {
	factory.RegisterInitializer(types.ROBOT_MODE__MANUAL, createRobotManual)
}

func createRobotManual(robot *Robot, config *global.RobotConfiguration, messageBus *bus.MessageBus, robotClient *client.RobotClient) *gobot.Robot {
	logrus.Info("initial manual robot...")
	if config.ActivateFirmata.True() {
		if config.ActivateCameraHolderController.True() {
			cameraHolderWorker := NewCameraHolderWorker(robot, messageBus, config)
			robot.AddWorker(cameraHolderWorker)
		}

		if config.ActivatePowerController.True() {
			powerControlWorker := NewPowerWorker(robot, messageBus, config)
			robot.AddWorker(powerControlWorker)
		}
	}

	if config.ActivateCameraController.True() {
		cameraWorker := NewCameraManualWorker(robot, messageBus, robotClient, config)
		robot.AddWorker(cameraWorker)
	}

	r := gobot.NewRobot("VehicleRobotManual")
	return r
}
