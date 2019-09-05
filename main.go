package main

import (
	"github.com/johnnyeven/libtools/courier/client"
	"github.com/johnnyeven/vehicle-robot-client/client_vehicle_robot"
	"github.com/johnnyeven/vehicle-robot-client/modules/controllers"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/opencv"
)

func main() {
	window := opencv.NewWindowDriver("Detective")
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
