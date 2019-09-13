package framework

import (
	"errors"

	"github.com/yamakiller/game/gateway/ServiceComponent"
	"github.com/yamakiller/game/gateway/elements"

	"github.com/yamakiller/magicNet/core"
	"github.com/yamakiller/magicNet/engine/util"
)

// GatewayFrame Gateway main frame
type GatewayFrame struct {
	core.DefaultStart
	core.DefaultEnv
	core.DefaultLoop
	//
	dsrv *core.DefaultService
	dcmd core.DefaultCMDLineOption
	//
	id   int32
	addr string
	max  int

	//
	luaService *ServiceComponent.ScriptService
	conService *ServiceComponent.ConnectService
	netService *ServiceComponent.NetworkService
}

//InitService init gateway system
func (gw *GatewayFrame) InitService() error {
	gw.dsrv = &core.DefaultService{}
	gatewayEnv := util.GetEnvMap(util.GetEnvRoot(), "gateway")
	if gatewayEnv == nil {
		return errors.New("Gateway configuration information does not exist ")
	}

	elements.GatewayID = int32(util.GetEnvInt(gatewayEnv, "id", 1))
	elements.GatewayAddr = util.GetEnvString(gatewayEnv, "addr", "0.0.0.0:7850")
	elements.GatewayMaxConnect = util.GetEnvInt(gatewayEnv, "max", 1024)
	elements.GatewayCCMax = util.GetEnvInt(gatewayEnv, "chan-max", 32)
	elements.GatewayLuaScriptPath = util.GetEnvString(gatewayEnv, "lua-script-path", "./script")
	elements.GatewayLuaScriptFile = util.GetEnvString(gatewayEnv, "lua-script-file", "./script/gateway.lua")

	gw.luaService = &ServiceComponent.ScriptService{}
	gw.luaService.Init()

	gw.conService = ServiceComponent.NewConnService()
	gw.netService = ServiceComponent.NewTCPNetworkService()

	return nil
}

//CloseService close gateway system
func (gw *GatewayFrame) CloseService() {
	//1.停止内部连接服务
	//2.停止外部监听服务
	//3.停止脚本服务
	if gw.luaService != nil {
		gw.luaService.Shutdown()
		gw.luaService = nil
	}
}

// VarValue : Command bind
func (gw *GatewayFrame) VarValue() {
	gw.dcmd.VarValue()
}

// LineOption :
func (gw *GatewayFrame) LineOption() {
	gw.dcmd.LineOption()
}
