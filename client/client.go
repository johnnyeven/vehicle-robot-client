package client

import (
	"fmt"
	"github.com/henrylee2cn/teleport"
	"github.com/henrylee2cn/teleport/proto/pbproto"
	"github.com/johnnyeven/libtools/conf"
	"github.com/sirupsen/logrus"
)

type RobotClient struct {
	cli        tp.Peer
	sess       tp.Session
	RemoteAddr string `conf:"env"`
	NodeKey    string `conf:"env"`
}

func (c *RobotClient) Init() {
	if c.NodeKey == "" {
		logrus.Panicf("RobotClient NodeKey must not be empty")
	}
	tp.SetLoggerLevel2(tp.ERROR)()
	tp.SetDefaultProtoFunc(pbproto.NewPbProtoFunc())
	c.cli = tp.NewPeer(tp.PeerConfig{})
}

func (c *RobotClient) Start() {
	var stat *tp.Status

	c.sess, stat = c.cli.Dial(c.RemoteAddr)
	if !stat.OK() {
		panic(fmt.Sprintf("connection err, status: %v", stat))
	}

	token, err := c.AuthorizationAuth(c.NodeKey)
	if err != nil {
		logrus.Panicf("RobotClient.AuthorizationAuth err: %v", err)
	}
	fmt.Println(string(token))
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

func (c *RobotClient) RegisterCallRouter(route interface{}, plugins ...tp.Plugin) []string {
	return c.cli.RouteCall(route, plugins...)
}

func (c *RobotClient) RegisterPushRouter(route interface{}, plugins ...tp.Plugin) []string {
	return c.cli.RoutePush(route, plugins...)
}
