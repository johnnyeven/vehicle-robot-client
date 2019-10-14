package robot

import (
	"fmt"
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/libtools/config_agent"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules"
	"github.com/johnnyeven/vehicle-robot-client/routes"
	bus2 "github.com/mustafaturan/bus"
	"github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Robot struct {
	configurations *global.RobotConfiguration
	master         *gobot.Master
	workers        map[string]Worker

	bus    *bus.MessageBus
	server *client.RobotClient
}

func NewRobot(server *client.RobotClient, bus *bus.MessageBus, config *global.RobotConfiguration) *Robot {
	if bus == nil {
		logrus.Panic("MessageBus can not be nil")
	}
	r := &Robot{
		configurations: config,
		workers:        make(map[string]Worker),
		bus:            bus,
		server:         server,
	}

	return r
}

func (r *Robot) Start() {
	r.bus.RegisterHandler("remote-address-handler", modules.RemoteAddressTopic, r.handleAddressEvent)
	r.bus.RegisterHandler("configuration-diff-handler", config_agent.DiffConfigTopic, r.handleDiffConfigEvent)

	r.gracefulRun()
}

func (r *Robot) AddWorker(w Worker) error {
	if _, ok := r.workers[w.WorkerID()]; ok {
		return fmt.Errorf("worker id duplicated: %s", w.WorkerID())
	}

	r.workers[w.WorkerID()] = w
	return nil
}

func (r *Robot) RemoveWorker(workerID string) error {
	w, ok := r.workers[workerID]
	if !ok {
		return fmt.Errorf("worker id not found: %s", workerID)
	}

	if err := w.Stop(); err != nil {
		return err
	}
	delete(r.workers, workerID)
	return nil
}

func (r *Robot) RestartWorker(workerID string) error {
	w, ok := r.workers[workerID]
	if !ok {
		return fmt.Errorf("worker id not found: %s", workerID)
	}

	return w.Restart()
}

func (r *Robot) handleAddressEvent(e *bus2.Event) {
	if addr, ok := e.Data.(*net.UDPAddr); ok {
		addr := net.TCPAddr{
			IP:   addr.IP,
			Port: addr.Port,
		}
		r.server.RemoteAddr = addr.String()
		r.server.Start()
		routes.InitRouters()

		robots := CreateRobotFromConfig(global.Config.RobotConfiguration, global.Config.MessageBus, global.Config.RobotClient)
		robots.Start()
	}
}

func (r *Robot) handleDiffConfigEvent(e *bus2.Event) {
	if config, ok := e.Data.(config_agent.DiffConfig); ok {
		fmt.Println(config.Key)
	}
}

func (r *Robot) gracefulRun() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

	select {
	case <-ch:
		signal.Stop(ch)
		break
	}
}
