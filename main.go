package main

import (
	"github.com/johnnyeven/libtools/servicex"
	"github.com/johnnyeven/vehicle-robot-client/global"
)

func main() {
	servicex.Execute()

	global.Config.ConfigAgent.BindConf(&global.Config.RobotConfiguration)
	global.Config.ConfigAgent.Start()

	robots := CreateRobotFromConfig(global.Config.RobotConfiguration, global.Config.MessageBus, global.Config.RobotClient)
	robots.Start()
}
