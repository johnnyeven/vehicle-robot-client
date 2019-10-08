package main

import (
	"github.com/johnnyeven/libtools/servicex"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	"github.com/johnnyeven/vehicle-robot-client/routes"
	"github.com/mustafaturan/bus"
	"net"
)

func main() {
	servicex.Execute()

	global.Config.ConfigAgent.BindConf(&global.Config.RobotConfiguration)
	global.Config.ConfigAgent.Start()

	global.Config.MessageBus.RegisterHandler("remote-address-handler", modules.RemoteAddressTopic, handleAddressEvent)

	broadcast := modules.NewBroadcastController()
	defer broadcast.Close()
	go broadcast.Start()

	select {}
}

func handleAddressEvent(e *bus.Event) {
	if ip, ok := e.Data.(net.IP); ok {
		addr := net.TCPAddr{
			IP:   ip,
			Port: 9090,
		}
		global.Config.RobotClient.RemoteAddr = addr.String()
		global.Config.RobotClient.Start()
		routes.InitRouters()

		robots := CreateRobotFromConfig(global.Config.RobotConfiguration, global.Config.MessageBus, global.Config.RobotClient)
		robots.Start()
	}
}
