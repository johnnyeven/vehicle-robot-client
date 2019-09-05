package main

import (
	"fmt"
	"github.com/johnnyeven/libtools/courier/client"
	"github.com/johnnyeven/libtools/servicex"
	"github.com/johnnyeven/vehicle-robot-client/client_vehicle_robot"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules/controllers"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
)

func main() {
	servicex.Execute()

	global.Config.ConfigAgent.BindConf(&global.Config.RobotConfiguration)
	global.Config.ConfigAgent.Start()

	fmt.Println(global.Config.RobotConfiguration.CameraMode)

	window := opencv.NewWindowDriver()
	camera := opencv.NewCameraDriver(0)

	cli := &client_vehicle_robot.ClientVehicleRobot{
		Client: client.Client{
			Host: "www.profzone.net",
			Port: 50999,
			Mode: "grpc",
		},
	}
	cli.MarshalDefaults(cli)

	robot := gobot.NewRobot("cameraBot",
		[]gobot.Device{window, camera},
		func() {
			controllers.ObjectDetectiveController(window, camera, cli)
		},
	)

	robot.Start()
}
