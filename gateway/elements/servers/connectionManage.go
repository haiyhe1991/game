package servers

import (
	"github.com/yamakiller/game/gateway/constant"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/network"
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
	cmsr.serverGroup[srvName] = ConnectionGroup{}
}

//GetGroup Get a service group
func (cmsr *ConnectionManager) GetGroup(srvName string) *ConnectionGroup {
	cgp, ok := cmsr.serverGroup[srvName]
	if !ok {
		return nil
	}
	return &cgp
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
	var err error
	for k, v := range cmsr.serverGroup {
		for i := 0; i < len(v.group); i++ {
			if v.group[i].ID == 0 || v.group[i].Sock > 0 {
				continue
			}
			err = AutoConnect(context, &v.group[i])
			if err != nil {
				logger.Error(context.Self().GetID(), "Connection service failed %s %s-%d-%s", err, k, v.group[i].ID, v.group[i].Addr)
			}
		}

		cmsr.serverGroup[k] = v
	}
}

//AutoConnect Auto connection service
func AutoConnect(context actor.Context, c *Connection) error {

	h, err := network.OperTCPConnect(context.Self(), c.Addr, constant.ConstConnectChanMax)
	if err != nil {
		return err
	}

	c.Sock = h
	network.OperOpen(h)

	return nil
}
