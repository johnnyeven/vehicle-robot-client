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

	// 主业务逻辑控制器
	//mainWorker := NewExploreMainWorker(robot, config, messageBus, robotClient)
	//robot.AddWorker(mainWorker)

	// 摄像头云台
	//cameraHolderWorker := NewCameraHolderWorker(robot, messageBus, config)
	//robot.AddWorker(cameraHolderWorker)

	// 摄像头
	//cameraWorker := NewCameraExploreWorker(robot, messageBus, robotClient, config)
	//robot.AddWorker(cameraWorker)

	// 马达控制器
	//powerControlWorker := NewPowerWorker(robot, messageBus, config)
	//robot.AddWorker(powerControlWorker)

	// 超声波测距
	//distanceWorker := NewDistanceHCSR04Worker(robot, messageBus, config)
	//robot.AddWorker(distanceWorker)

	// 姿态控制器
	attitudeWorker := NewAttitudeMPU6050Worker(robot, messageBus, config)
	robot.AddWorker(attitudeWorker)
	// attitudeWorker := NewAttitudeGY85Worker(robot, messageBus, config)
	//robot.AddWorker(attitudeWorker)

	r := gobot.NewRobot("VehicleRobotExplore")
	return r
}
