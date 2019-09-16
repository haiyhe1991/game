package elements

import (
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/gateway/constant"
	"github.com/yamakiller/game/gateway/elements/clients"
	"github.com/yamakiller/game/gateway/elements/route"
	"github.com/yamakiller/game/gateway/elements/servers"
)

var (
	//RouteAddress Routing address table
	RouteAddress *route.Table
	//Clients Manage all externally connected players
	Clients *clients.ClientManager
	//Conns Connection service manager
	Conns *servers.ConnectionManager
)

func init() {
	RouteAddress = route.NewTable()
	Clients = &clients.ClientManager{}
	Conns = servers.NewManager()

	agreement.SetSingleLimit(constant.ConstPlayerBufferLimit >> 1)
}
