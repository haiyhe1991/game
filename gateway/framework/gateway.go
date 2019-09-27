package framework

import (
	"errors"

	"github.com/yamakiller/game/gateway/elements/clients"

	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/service/net"

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
	scriptLua *component.GatewayScirpt
	forward   *component.GatewayForward
	netListen *component.GatewayListener
}

//InitService init gateway system
func (gw *GatewayFrame) InitService() error {

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
	constant.GatewayLogoutName = util.GetEnvString(gatewayEnv, "login-logout-name", "sign/in/out")
	constant.GatewayLogoutPactum = util.GetEnvString(gatewayEnv, "logout-protocol-name", "UnLoginRequest")
	constant.GatewayConnectForwardErrMax = util.GetEnvInt(gatewayEnv, "forward-connect-fail-retry", 16)
	constant.GatewayConnectForwardInterval = util.GetEnvInt(gatewayEnv, "forward-reconnect-interval", 200)
	constant.GatewayConnectForwardAutoTick = util.GetEnvInt(gatewayEnv, "forward-check-connect-interval", 1000)
	constant.GatewayConnectLoaderReplicas = util.GetEnvInt(gatewayEnv, "forward-loader-replicas", 20)

	gw.scriptLua = &component.GatewayScirpt{}
	gw.scriptLua.Init()

	gw.forward = func() *component.GatewayForward {
		return service.Make(constant.ConstForwardServiceName, func() service.IService {
			h := &component.GatewayForward{}
			h.Init()
			return h
		}).(*component.GatewayForward)
	}()

	gw.netListen = func() *component.GatewayListener {
		return service.Make(constant.ConstNetworkServiceName, func() service.IService {

			h := &component.GatewayListener{NetListenService: implement.NetListenService{NetListen: &net.TCPListen{},
				NetDeleate: &component.GNetListenDeleate{},
				NetClients: clients.NewGClientManager(),
				ClientKeep: uint64(constant.GatewayConnectKleep)}}

			h.MaxClient = constant.GatewayMaxConnect
			h.CCMax = constant.GatewayCCMax
			h.Addr = constant.GatewayAddr
			h.NetClients.(*clients.GClientManager).Association(constant.GatewayID)
			h.Init()
			return h

		}).(*component.GatewayListener)
	}()

	return nil
}

//CloseService close gateway system
func (gw *GatewayFrame) CloseService() {

	if gw.forward == nil {
		gw.forward.Shutdown()
		gw.forward = nil
	}
	if gw.netListen != nil {
		gw.netListen.Shutdown()
		gw.netListen = nil
	}

	if gw.scriptLua != nil {
		gw.scriptLua.Shutdown()
		gw.scriptLua = nil
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
