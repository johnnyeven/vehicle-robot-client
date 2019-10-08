package modules

import (
	"github.com/johnnyeven/vehicle-robot-client/global"
	"github.com/sirupsen/logrus"
	"net"
)

const RemoteAddressTopic = "remote.address"

type BroadcastController struct {
	conn *net.UDPConn
}

func (c *BroadcastController) init() {
	laddr := &net.UDPAddr{
		IP:   nil,
		Port: 9091,
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		logrus.Panic(err)
	}

	c.conn = conn
	logrus.Infof("[BroadcastController] listen at %s", laddr.String())
}

func (c *BroadcastController) Close() error {
	return c.conn.Close()
}

func (c *BroadcastController) Start() {
	buffer := make([]byte, 1024)
	for {
		count, addr, err := c.conn.ReadFromUDP(buffer)
		if err != nil {
			logrus.Warningf("[BroadcastController] conn.ReadFromUDP err: %v", err)
		}

		logrus.Infof("received udp: length=%d, address=%s", count, addr.String())
		_, err = global.Config.MessageBus.Emit(RemoteAddressTopic, addr.IP, "")
		if err != nil {
			logrus.Warningf("[BroadcastController] MessageBus.Emit err: %v", err)
		}
	}
}

func NewBroadcastController() *BroadcastController {
	c := &BroadcastController{}
	c.init()
	return c
}
