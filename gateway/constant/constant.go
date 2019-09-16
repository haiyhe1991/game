package constant

import (
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
	//ConstConnectForwardErrMax  Push failure maximum retries
	ConstConnectForwardErrMax = 16
	//ConstNetworkServiceName Network service name
	ConstNetworkServiceName = "Service/Gateway/Network"
	//ConstConnectServiceName The name of the connection service
	ConstConnectServiceName = "Service/Gateway/Connection"
	//ConstConnectAutoTick Automatic connection detection interval event
	ConstConnectAutoTick = 100 // Unit millisecond
)

var (
	//GatewayID  gateway service id code
	GatewayID int32
	//GatewayMaxConnect Maximum number of connections for the gateway service
	GatewayMaxConnect int
	//GatewayAddr Gateway service address information
	GatewayAddr string
	//GatewayCCMax Gateway connects to the client pipe maximum buffer, the default is 32
	GatewayCCMax int
	//GatewayLuaScriptPath Gateway script search path
	GatewayLuaScriptPath string
	//GatewayLuaScriptFile Gateway default script name
	GatewayLuaScriptFile string
	//ConnectServicePID Connected service ACTOR PID
	ConnectServicePID actor.PID
	//NetworkServicePID ACTOR PID for network services
	NetworkServicePID actor.PID
)
