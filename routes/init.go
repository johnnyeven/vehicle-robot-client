package routes

import (
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/routes/power"
)

func InitRouters() {
	cli := global.Config.RobotClient
	cli.RegisterPushRouter(power.NewPowerRouter(global.Config.MessageBus))
}
