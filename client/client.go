package client

import (
	"fmt"
	"github.com/henrylee2cn/teleport"
	"github.com/johnnyeven/libtools/conf"
)

type RobotClient struct {
	cli        tp.Peer
	sess       tp.Session
	RemoteAddr string `conf:"env"`
}

func (c *RobotClient) Init() {
	var stat *tp.Status
	c.cli = tp.NewPeer(tp.PeerConfig{})

	c.sess, stat = c.cli.Dial(c.RemoteAddr)
	if !stat.OK() {
		panic(fmt.Sprintf("connection err, status: %v", stat))
	}
}

func (c RobotClient) MarshalDefaults(v interface{}) {
	if h, ok := v.(*RobotClient); ok {
		if h.RemoteAddr == "" {
			h.RemoteAddr = "127.0.0.1:9090"
		}
	}
}

func (c *RobotClient) DockerDefaults() conf.DockerDefaults {
	return conf.DockerDefaults{
		"RemoteAddr": "127.0.0.1:9090",
	}
}
