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
		Host:               "service-configurations.profzone.service.profzone.net",
		PullConfigInterval: 60,
		StackID:            123,
	},

	RobotConfiguration: RobotConfiguration{},

	RobotClient: &client.RobotClient{
		RemoteAddr: "localhost:9090",
		NodeKey:    "123",
	},

	MessageBus: &bus.MessageBus{
		RegisterTopics: []string{
			"control.moving",
		},
	},
}
