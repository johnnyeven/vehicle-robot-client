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
	factory.RegisterInitializer(types.ROBOT_MODE__EXPLORE, createRobotExplore)
}

func createRobotExplore(robot *Robot, config *global.RobotConfiguration, messageBus *bus.MessageBus, robotClient *client.RobotClient) *gobot.Robot {
	logrus.Info("initial explore robot...")
	mainWorker := NewExploreMainWorker(robot, config, messageBus, robotClient)
	robot.AddWorker(mainWorker)

	cameraHolderWorker := NewCameraHolderWorker(robot, messageBus, config)
	robot.AddWorker(cameraHolderWorker)

	powerControlWorker := NewPowerWorker(robot, messageBus, config)
	robot.AddWorker(powerControlWorker)

	cameraWorker := NewCameraExploreWorker(robot, messageBus, robotClient, config)
	robot.AddWorker(cameraWorker)

	distanceWorker := NewDistanceHCSR04Worker(robot, messageBus, config)
	robot.AddWorker(distanceWorker)

	attitudeWorker := NewAttitudeGY85Worker(robot, messageBus, config)
	robot.AddWorker(attitudeWorker)

	r := gobot.NewRobot("VehicleRobotExplore")
	return r
}
