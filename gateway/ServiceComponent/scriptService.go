package ServiceComponent

import (
	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/gateway/elements"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/script/stack"
	"github.com/yamakiller/mgolua/mlua"
)

// ScriptService  lua script service
type ScriptService struct {
	handle *stack.LuaStack
}

//Init
func (sse *ScriptService) Init() {
	sse.handle = stack.NewLuaStack()
	sse.handle.GetLuaState().OpenLibs()
	sse.handle.AddSreachPath(elements.GatewayLuaScriptPath)
	//register the registerRouteProto function and set gw
	sse.handle.GetLuaState().Register("register_route", luaRegisterRoute)
	sse.handle.GetLuaState().Register("register_service_group", luaRegisterServiceGroup)
	sse.handle.GetLuaState().Register("append_connection", luaAppendServiceConnection)

	if _, err := sse.handle.ExecuteScriptFile(elements.GatewayLuaScriptFile); err != nil {
		panic(err)
	}
}

//Shutdown
func (sse *ScriptService) Shutdown() {
	if sse.handle == nil {
		return
	}

	sse.handle.GetLuaState().Close()
}

func luaRegisterRoute(L *mlua.State) int {

	argsNum := L.GetTop()
	if argsNum < 2 {
		return L.Error("register route need  2-3 parameters")
	}

	agreementName := L.ToCheckString(1)
	serviceName := L.ToCheckString(2)
	auth := true
	if argsNum > 2 {
		auth = L.ToBoolean(3)
	}

	agreementType := proto.MessageType(agreementName)
	if agreementType == nil {
		logger.Error(0, "Gateway Registration %s routing agreement error ", agreementName)
		return 0
	}

	elements.RouteAddress.Register(agreementType, agreementName, serviceName, auth)
	return 0
}

func luaRegisterServiceGroup(L *mlua.State) int {
	serverName := L.ToCheckString(1)
	connectManger.Register(serverName)
	return 0
}

func luaAppendServiceConnection(L *mlua.State) int {
	argsNum := L.GetTop()
	if argsNum < 3 {
		return L.Error("append service connection error need 3 parameters")
	}

	serverName := L.ToCheckString(1)
	serverID := L.ToCheckInteger(2)
	serverAddr := L.ToCheckString(3)
	grp := connectManger.GetGroup(serverName)
	if grp == nil {
		return L.Error("pend service connection error unfind server group")
	}

	grp.Register(int32(serverID), serverAddr)
	return 0
}
