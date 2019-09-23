package elements

import (
	"github.com/yamakiller/game/common/manager"
	"github.com/yamakiller/game/gateway/elements/forward"
	"github.com/yamakiller/game/gateway/elements/servers"
)

var (
	//ForwardAddresses Routing address table
	ForwardAddresses *forward.Table
	//TargetRecord Connection configuration status information of the target server
	TargetRecord *servers.TargetGroup
	//Conns *servers.ConnectionManager
	//SSets Service set
	SSets *manager.SSets
)

func init() {
	ForwardAddresses = forward.NewTable()
	TargetRecord = servers.NewTargetGroup()
	TargetRecord.Init()
	SSets = manager.NewSSets()
	//Clients = clients.NewGClientManager()
	//Conns = servers.NewManager()

	//preset.SetSingleLimit(constant.ConstPlayerBufferLimit >> 1)
}
