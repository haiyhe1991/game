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
	//TLSets Service clusters that are not connected to the gateway provide load balancing
	TLSets *servers.TargetLoadSet
	//SSets Service set
	SSets *manager.SSets
)

func init() {
	ForwardAddresses = forward.NewTable()
	TSets = servers.NewTargetSet()
	TLSets = servers.NewLoadSet()
	SSets = manager.NewSSets()
}
