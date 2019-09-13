package elements

import (
	"github.com/yamakiller/game/gateway/elements/route"
	"github.com/yamakiller/magicNet/engine/actor"
)

const (
	//ConstPlayerMax play number of max
	ConstPlayerMax = 65535
	//ConstPlayerIDMask play ID of mask
	ConstPlayerIDMask = 0xFF
	//ConstPlayerBufferLimit Read buffer Max Cap
	ConstPlayerBufferLimit = 4096
	//ConstConnectGroupMax Maximum number of connectable services
	ConstConnectGroupMax = 128
	//ConstConnectChanMax Maximum chan data buffer limit for connection services
	ConstConnectChanMax = 256
	//
	ConstConnectForwardErrMax = 16
	//
	ConstNetworkServiceName = "Service/Gateway/Network"
	//
	ConstConnectServiceName = "Service/Gateway/Connection"
)

var (
	//RouteAddress route address informat
	RouteAddress *route.Table
	//GatewayID  gateway service id code
	GatewayID int32
	//GatewayMaxConnect Maximum number of connections for the gateway service
	GatewayMaxConnect int
	//GatewayAddr Gateway service address information
	GatewayAddr string
	//GatewayCCMax Gateway connects to the client pipe maximum buffer, the default is 32
	GatewayCCMax int
	//GatewayLuaScriptPath
	GatewayLuaScriptPath string
	//GatewayLuaScriptFile
	GatewayLuaScriptFile string
	//
	ConnectServicePID actor.PID
	//
	NetworkServicePID actor.PID
)

func init() {
	RouteAddress = route.NewTable()
}
