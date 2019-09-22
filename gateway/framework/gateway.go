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
	netListen *component.GatewayListener
	//conService *component.ConService
	//netService *component.OutNetService
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

	gw.scriptLua = &component.GatewayScirpt{}
	gw.scriptLua.Init()

	//gw.conService = component.NewTCPConService()
	gw.netListen = func() *component.GatewayListener {
		return service.Make(constant.ConstNetworkServiceName, func() service.IService {

			h := &component.GatewayListener{NetListenService: implement.NetListenService{NetListen: &net.TCPListen{},
				NetDeleate:            &component.GNetListenDeleate{},
				NetClients:            clients.NewGClientManager(),
				ClientKeep:            uint64(constant.GatewayConnectKleep),
				ClientRecvBufferLimit: constant.ConstClientBufferLimit}}

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

	if gw.netListen != nil {
		gw.netListen.Shutdown()
		gw.netListen = nil
	}

	/*if gw.conService != nil {
		gw.conService.Shutdown()
		gw.conService = nil
	}*/

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

//NewTCPConService Create a connection service
/*func NewTCPConService() *ConService {
	return service.Make(constant.ConstConnectServiceName, func() service.IService {
		handle := &ConService{}

		handle.Init()
		return handle
	}).(*ConService)
}*/
