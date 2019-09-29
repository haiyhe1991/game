package module

import (
	"reflect"
	"unsafe"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/magicNet/script/stack"
	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/util"
	"github.com/yamakiller/mgolua/mlua"
)

//InNetScript No listening service provides registration operation service
type InNetScript struct {
	handle *stack.LuaStack
	parent interface{}
	//parentLink reflect.Type
}

//Execution Execution script
func (ins *InNetScript) Execution(fileName string,
	parent interface{}) {

	ins.handle = stack.NewLuaStack()
	ins.parent = parent

	//ins.parentLink = reflect.TypeOf(parent)
	defer ins.handle.Shutdown()

	ins.handle.OpenLibs()
	ins.handle.Register(luaRegisterProtobuf, "register_proto_method",
		uintptr(unsafe.Pointer(ins)))

	if _, err := ins.handle.ExecuteScriptFile(fileName); err != nil {
		panic(err)
	}
}

func luaRegisterProtobuf(L *mlua.State) int {
	p := L.ToLightGoStruct(L.UpvalueIndex(1))
	if p == nil {
		return L.Error("Upvalue index is empty")
	}

	scriptH := (*InNetScript)(p)
	pointer := scriptH.parent.(*implement.NetMethodDispatch)

	argsNum := L.GetTop()
	if argsNum < 3 {
		return L.Error("param error: param 1 is protobuf name, param 2 is method object, param 3 is method name")
	}

	protoName := L.CheckString(1)
	methodObjectName := L.CheckString(2)
	methodName := L.CheckString(3)

	protoType := proto.MessageType(protoName)
	if protoType == nil {
		return L.Error("The %s protocol does not exist.", protoName)
	}

	regObject := FactoryInstance().Get(methodObjectName)
	util.Assert(!(regObject == nil), methodObjectName+" Member not found")
	methodFunc := reflect.ValueOf(regObject).MethodByName(methodName)
	util.Assert(methodFunc.IsValid(), methodObjectName+" "+methodName+" method not found")

	pointer.RegisterType(protoType, methodFunc.Interface().(func(event implement.INetMethodEvent)))

	return 0
}
