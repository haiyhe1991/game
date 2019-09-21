package elements

import (
	"github.com/yamakiller/game/gateway/elements/route"
	"github.com/yamakiller/game/gateway/elements/servers"
)

var (
	//RouteAddress Routing address table
	RouteAddress *route.Table
	//Clients Manage all externally connected players
	//Clients *clients.GClientManager
	//Conns Connection service manager
	Conns *servers.ConnectionManager
)

func init() {
	RouteAddress = route.NewTable()
	//Clients = clients.NewGClientManager()
	Conns = servers.NewManager()

	//preset.SetSingleLimit(constant.ConstPlayerBufferLimit >> 1)
}
