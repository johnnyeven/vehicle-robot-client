package main

import (
	"github.com/johnnyeven/libtools/servicex"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	"github.com/johnnyeven/vehicle-robot-client/modules/robot"
)

func main() {
	servicex.Execute()

	global.Config.ConfigAgent.BindConf(&global.Config.RobotConfiguration)
	global.Config.ConfigAgent.BindBus(global.Config.MessageBus)
	go global.Config.ConfigAgent.Start()

	broadcast := modules.NewBroadcastController()
	go broadcast.Start()

	r := robot.NewRobot(global.Config.RobotClient, global.Config.MessageBus, &global.Config.RobotConfiguration)
	r.Start()
}
