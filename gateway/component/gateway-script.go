package component

import (
	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/gateway/constant"
	"github.com/yamakiller/game/gateway/elements"
	"github.com/yamakiller/game/gateway/elements/servers"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/script/stack"
	"github.com/yamakiller/mgolua/mlua"
)

// GatewayScirpt  lua script service
type GatewayScirpt struct {
	handle *stack.LuaStack
}

//Init Initialize the lua script service
func (gsc *GatewayScirpt) Init() {
	gsc.handle = stack.NewLuaStack()
	gsc.handle.GetLuaState().OpenLibs()
	gsc.handle.AddSreachPath(constant.GatewayLuaScriptPath)
	//register the registerRouteProto function and set gw
	gsc.handle.GetLuaState().Register("register_forward", luaRegisterForward)
	gsc.handle.GetLuaState().Register("register_target_connect", luaRegisterTargetService)

	if _, err := gsc.handle.ExecuteScriptFile(constant.GatewayLuaScriptFile); err != nil {
		panic(err)
	}
}

//Shutdown Close the script service
func (gsc *GatewayScirpt) Shutdown() {
	if gsc.handle == nil {
		return
	}

	gsc.handle.GetLuaState().Close()
}

func luaRegisterForward(L *mlua.State) int {

	argsNum := L.GetTop()
	if argsNum < 2 {
		return L.Error("register route need  2-3 parameters")
	}

	protocolName := L.ToCheckString(1)
	serverName := L.ToCheckString(2)
	auth := true
	if argsNum > 2 {
		auth = L.ToBoolean(3)
	}

	protocolType := proto.MessageType(protocolName)
	if protocolType == nil {
		logger.Error(0, "Gateway Registration %s forward agreement error ", protocolName)
		return 0
	}

	elements.ForwardAddresses.Register(protocolType, protocolName, serverName, auth)

	logger.Debug(0, "Gateway Registration Forward Address %s,%s,%+v", protocolName, serverName, auth)

	return 0
}

func luaRegisterTargetService(L *mlua.State) int {
	argsNum := L.GetTop()
	if argsNum < 3 {
		return L.Error("append service connection error need 5-6 parameters[ID,Name,Address,Timeout, outChanMax]")
	}
	targetID := int32(L.ToCheckInteger(1))
	targetName := L.ToCheckString(2)
	targetAddr := L.ToCheckString(3)
	targetTimeout := L.ToCheckInteger(4)
	targetOutChanMax := L.ToCheckInteger(5)

	targetDesc := ""
	if argsNum > 5 {
		targetDesc = L.ToCheckString(6)
	}

	err := elements.TSets.Push(&servers.TargetConnection{ID: targetID,
		Name:       targetName,
		Addr:       targetAddr,
		Desc:       targetDesc,
		TimeOut:    uint64(targetTimeout),
		OutChanMax: int(targetOutChanMax)})

	if err != nil {
		logger.Error(0, "Gateway Registration TargetConnection fail error: %+v", err)
	}

	logger.Debug(0, "Gateway Registration TargetConnection ID:%d Name:%s Addr:%s,Timeout:%d milli ,out-chan-max:%d",
		targetID,
		targetName,
		targetAddr,
		targetTimeout,
		targetOutChanMax)
	return 0
}
