package framework

import (
	"errors"

	"github.com/yamakiller/game/gateway/component"
	"github.com/yamakiller/game/gateway/constant"

	"github.com/yamakiller/magicNet/core"
	"github.com/yamakiller/magicNet/util"
)

// GatewayFrame Gateway main frame
type GatewayFrame struct {
	core.DefaultStart
	core.DefaultEnv
	core.DefaultLoop
	//
	core.DefaultService
	dcmd core.DefaultCMDLineOption
	//
	id   int32
	addr string
	max  int

	//
	luaService *component.ScriptService
	conService *component.ConService
	netService *component.NetService
}

//InitService init gateway system
func (gw *GatewayFrame) InitService() error {
	//gw.dsrv = &core.DefaultService{}
	if err := gw.DefaultService.InitService(); err != nil {
		return err
	}

	gatewayEnv := util.GetEnvMap(util.GetEnvRoot(), "gateway")
	if gatewayEnv == nil {
		return errors.New("Gateway configuration information does not exist ")
	}

	constant.GatewayID = int32(util.GetEnvInt(gatewayEnv, "id", 1))
	constant.GatewayAddr = util.GetEnvString(gatewayEnv, "addr", "0.0.0.0:7850")
	constant.GatewayMaxConnect = util.GetEnvInt(gatewayEnv, "max", 1024)
	constant.GatewayCCMax = util.GetEnvInt(gatewayEnv, "chan-max", 32)
	constant.GatewayConnectKleep = util.GetEnvInt(gatewayEnv, "connection-kleep", 1000*30)
	constant.GatewayLuaScriptPath = util.GetEnvString(gatewayEnv, "lua-script-path", "./script")
	constant.GatewayLuaScriptFile = util.GetEnvString(gatewayEnv, "lua-script-file", "./script/gateway.lua")

	gw.luaService = &component.ScriptService{}
	gw.luaService.Init()

	gw.conService = component.NewTCPConService()
	gw.netService = component.NewTCPNetService()

	return nil
}

//CloseService close gateway system
func (gw *GatewayFrame) CloseService() {

	if gw.netService != nil {
		gw.netService.Shutdown()
		gw.netService = nil
	}

	if gw.conService != nil {
		gw.conService.Shutdown()
		gw.conService = nil
	}

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
