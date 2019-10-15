package robot

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/api"
)

func CreateRobotFromConfig(robot *Robot, config *global.RobotConfiguration, messageBus *bus.MessageBus, robotClient *client.RobotClient) *gobot.Master {
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
		cameraWorker := NewCameraWorker(robot, messageBus, robotClient, config)
		robot.AddWorker(cameraWorker)
	}

	master := gobot.NewMaster()

	if config.ActivateApiSupport.True() {
		apiServer := api.NewAPI(master)
		apiServer.Port = config.APIServerPort
		apiServer.Start()
	}

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

	master.AddRobot(r)

	return master
}
