package elements

import (
	"github.com/yamakiller/game/common/manager"
	"github.com/yamakiller/game/gateway/elements/forward"
	"github.com/yamakiller/game/gateway/elements/servers"
)

var (
	//ForwardAddresses Routing address table
	ForwardAddresses *forward.Table
	//TSets Connection configuration status information of the target server
	TSets *servers.TargetSet
	//SSets Service set
	SSets *manager.SSets
)

func init() {
	ForwardAddresses = forward.NewTable()
	TSets = servers.NewTargetSet()
	SSets = manager.NewSSets()
	//Clients = clients.NewGClientManager()
	//Conns = servers.NewManager()

	//preset.SetSingleLimit(constant.ConstPlayerBufferLimit >> 1)
}
