package constant

const (
	//ConstClientMax play number of max
	ConstClientMax = 65535
	//ConstClientIDMask play ID of mask
	ConstClientIDMask = 0xFFFFFF
	//ConstConnectGroupMax Maximum number of connectable services
	ConstConnectGroupMax = 128
	//ConstConnectChanMax Maximum chan data buffer limit for connection services
	ConstConnectChanMax = 256
	//ConstNetworkServiceName Network service name
	ConstNetworkServiceName = "Service/Gateway/Listen"
	//ConstForwardServiceName The name of the connection service
	ConstForwardServiceName = "Service/Gateway/Forward"
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
	//GatewayConnectKleep Gateway connection heartbeat event setting, in milliseconds
	GatewayConnectKleep int
	//GatewayLuaScriptPath Gateway script search path
	GatewayLuaScriptPath string
	//GatewayLuaScriptFile Gateway default script name
	GatewayLuaScriptFile string
	//GatewayLogoutName Login to the name of the logout server
	GatewayLogoutName string
	//GatewayLogoutPactum Login and sign-out agreement
	GatewayLogoutPactum string

	//GatewayConnectForwardErrMax Gateway connector automatic reconnection attempts
	GatewayConnectForwardErrMax int
	//GatewayConnectForwardInterval Automatic reconnection attempt interval, in milliseconds
	GatewayConnectForwardInterval int
	//GatewayConnectForwardAutoTick Automatic reconnection interval
	GatewayConnectForwardAutoTick int
	//GatewayConnectLoaderReplicas Connect the load balancer, copy factor, virtual node factor
	GatewayConnectLoaderReplicas int
)
