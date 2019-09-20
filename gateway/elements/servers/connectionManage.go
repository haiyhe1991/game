package servers

import (
	"fmt"
	"time"

	"github.com/yamakiller/game/gateway/constant"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/util"
)

//NewManager Create a service connection manager
func NewManager() *ConnectionManager {
	return &ConnectionManager{serverGroup: make(map[string]ConnectionGroup, 32)}
}

// ConnectionManager Provide management for connection services
type ConnectionManager struct {
	serverGroup map[string]ConnectionGroup
}

//Register Register a service group
func (cmsr *ConnectionManager) Register(srvName string) {
	cmsr.serverGroup[srvName] = ConnectionGroup{g: util.NewConsistent(20)}
}

//GetGroup Get a service group
func (cmsr *ConnectionManager) GetGroup(srvName string) *ConnectionGroup {
	cgp, ok := cmsr.serverGroup[srvName]
	if !ok {
		return nil
	}
	return &cgp
}

//GetGroupNames returns the names of all servers
func (cmsr *ConnectionManager) GetGroupNames() []string {
	i := 0
	result := make([]string, len(cmsr.serverGroup))
	for name, _ := range cmsr.serverGroup {
		result[i] = name
		i++
	}
	return result
}

//GetHandle Get the connection object
func (cmsr *ConnectionManager) GetHandle(sock int32) *Connection {
	var conn *Connection
	for _, v := range cmsr.serverGroup {
		conn = v.FindSocket(sock)
		if conn != nil {
			return conn
		}
	}

	return nil
}

//CheckConnect checking connection state and auto connection service
func (cmsr *ConnectionManager) CheckConnect(context actor.Context) {
	f := func(v interface{}) {
		c := v.(*Connection)
		if c.GetID() == 0 || c.GetSocket() > 0 {
			return
		}
		err := AutoConnect(context, c)
		if err != nil {
			logger.Error(context.Self().GetID(), "Connection service failed %s %d-%s", err, c.GetID(), c.GetAddr())
		}
	}

	for _, v := range cmsr.serverGroup {
		v.g.Range(f)
	}
}

//AutoConnect Auto connection service
func AutoConnect(context actor.Context, c *Connection) error {

	h, err := network.OperTCPConnect(context.Self(), c.GetAddr(), constant.ConstConnectChanMax)
	if err != nil {
		return err
	}

	c.SetSocket(h)
	network.OperOpen(h)

	ick := 0
	for {
		if c.auth > 0 {
			break
		}

		ick++
		if ick > 100 {
			network.OperClose(h)
			return fmt.Errorf("Automatic reconnection timeout does not wait for handshake data:%d", h)
		}

		time.Sleep(time.Millisecond * time.Duration(100))
	}

	return nil
}
