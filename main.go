package main

import (
	"github.com/johnnyeven/libtools/servicex"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	"github.com/johnnyeven/vehicle-robot-client/routes"
	"github.com/mustafaturan/bus"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	servicex.Execute()

	global.Config.ConfigAgent.BindConf(&global.Config.RobotConfiguration)
	global.Config.ConfigAgent.Start()

	global.Config.MessageBus.RegisterHandler("remote-address-handler", modules.RemoteAddressTopic, handleAddressEvent)

	broadcast := modules.NewBroadcastController()
	defer broadcast.Close()
	go broadcast.Start()

	gracefulClose()
}

func handleAddressEvent(e *bus.Event) {
	if addr, ok := e.Data.(*net.UDPAddr); ok {
		addr := net.TCPAddr{
			IP:   addr.IP,
			Port: addr.Port,
		}
		global.Config.RobotClient.RemoteAddr = addr.String()
		global.Config.RobotClient.Start()
		routes.InitRouters()

		robots := CreateRobotFromConfig(global.Config.RobotConfiguration, global.Config.MessageBus, global.Config.RobotClient)
		robots.Start()
	}
}

func gracefulClose() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

	select {
	case <-ch:
		signal.Stop(ch)
		break
	}
}
