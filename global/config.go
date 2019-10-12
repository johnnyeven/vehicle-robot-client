package global

import (
	"github.com/johnnyeven/libtools/bus"
	"github.com/johnnyeven/libtools/config_agent"
	"github.com/johnnyeven/libtools/servicex"
	"github.com/johnnyeven/vehicle-robot-client/client"
)

func init() {
	servicex.SetServiceName("vehicle-robot-client")
	servicex.ConfP(&Config)
}

var Config = struct {
	ConfigAgent        *config_agent.Agent
	RobotConfiguration RobotConfiguration
	RobotClient        *client.RobotClient
	MessageBus         *bus.MessageBus
}{
	ConfigAgent: &config_agent.Agent{
		Host:               "localhost",
		Port:               8002,
		PullConfigInterval: 10,
		StackID:            123,
	},

	RobotConfiguration: RobotConfiguration{},

	RobotClient: &client.RobotClient{
		NodeKey: "123",
	},

	MessageBus: &bus.MessageBus{},
}
