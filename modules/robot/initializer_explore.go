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

func createRobotExplore(robot *Robot, config *global.RobotConfiguration, messageBus *bus.MessageBus, robotClient *client.RobotClient) *gobot.Master {
	logrus.Info("initial explore robot...")
	cameraHolderWorker := NewCameraHolderWorker(robot, messageBus, config)
	robot.AddWorker(cameraHolderWorker)

	powerControlWorker := NewPowerWorker(robot, messageBus, config)
	robot.AddWorker(powerControlWorker)

	cameraWorker := NewCameraWorker(robot, messageBus, robotClient, config)
	robot.AddWorker(cameraWorker)

	r := gobot.NewRobot("VehicleRobot")
	for _, c := range robot.connections {
		r.AddConnection(c)
	}
	for _, d := range robot.devices {
		r.AddDevice(d)
	}
	r.Work = func() {
		for _, worker := range robot.workers {
			worker.Start()
		}
	}

	master := gobot.NewMaster()
	master.AddRobot(r)

	return master
}
