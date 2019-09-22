package elements

import (
	"github.com/yamakiller/game/gateway/elements/forward"
	"github.com/yamakiller/game/gateway/elements/servers"
	"github.com/yamakiller/magicNet/engine/actor"
)

var (
	//ForwardAddresses Routing address table
	ForwardAddresses *forward.Table
	//TargetRecord Connection configuration status information of the target server
	TargetRecord *servers.TargetGroup
	//Conns *servers.ConnectionManager
	//ForwardServer ACTOR PID for Service module for data interaction with other servers
	ForwardServer actor.PID
	//ListenServer ACTOR PID for network listen services
	ListenServer actor.PID
)

func init() {
	ForwardAddresses = forward.NewTable()
	TargetRecord = servers.NewTargetGroup()
	TargetRecord.Init()
	//Clients = clients.NewGClientManager()
	//Conns = servers.NewManager()

	//preset.SetSingleLimit(constant.ConstPlayerBufferLimit >> 1)
}
