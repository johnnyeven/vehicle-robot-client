package robot

import (
	"fmt"
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/libtools/config_agent"
	"github.com/johnnyeven/vehicle-robot-client/client"
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/johnnyeven/vehicle-robot-client/modules"
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
	devices        map[string]gobot.Device
	connections    map[string]gobot.Connection

	bus *bus.MessageBus
	cli *client.RobotClient
}

func NewRobot(cli *client.RobotClient, bus *bus.MessageBus, config *global.RobotConfiguration) *Robot {
	if bus == nil {
		logrus.Panic("MessageBus can not be nil")
	}
	r := &Robot{
		configurations: config,
		workers:        make(map[string]Worker),
		devices:        make(map[string]gobot.Device),
		connections:    make(map[string]gobot.Connection),
		bus:            bus,
		cli:            cli,
	}

	return r
}

func (r *Robot) Start() {
	r.bus.RegisterHandler("remote-address-handler", modules.RemoteAddressTopic, r.handleAddressEvent)
	r.bus.RegisterHandler("configuration-diff-handler", config_agent.DiffConfigTopic, r.handleDiffConfigEvent)

	r.gracefulRun()
}

func (r *Robot) startRobot() {
	r.master.Start()
}

func (r *Robot) Stop() error {
	return r.master.Stop()
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

func (r *Robot) AddDevice(d ...gobot.Device) {
	for _, dev := range d {
		r.devices[dev.Name()] = dev
	}
}

func (r *Robot) GetDevice(name string) gobot.Device {
	if d, ok := r.devices[name]; ok {
		return d
	}

	return nil
}

func (r *Robot) AddConnection(c ...gobot.Connection) {
	for _, conn := range c {
		r.connections[conn.Name()] = conn
	}
}

func (r *Robot) GetConnection(name string) gobot.Connection {
	if c, ok := r.connections[name]; ok {
		return c
	}

	return nil
}

func (r *Robot) handleAddressEvent(e *bus2.Event) {
	if addr, ok := e.Data.(*net.UDPAddr); ok {
		addr := net.TCPAddr{
			IP:   addr.IP,
			Port: addr.Port,
		}
		r.cli.RemoteAddr = addr.String()
		r.cli.Start()

		robots := CreateRobotFromConfig(r, &global.Config.RobotConfiguration, global.Config.MessageBus, global.Config.RobotClient)

		r.master = robots
		r.startRobot()
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
		r.Stop()
		break
	}
}
