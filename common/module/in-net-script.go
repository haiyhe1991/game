package module

import (
	"reflect"
	"unsafe"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/magicNet/script/stack"
	"github.com/yamakiller/magicNet/util"
	"github.com/yamakiller/mgolua/mlua"
)

//InNetScript No listening service provides registration operation service
type InNetScript struct {
	handle *stack.LuaStack
	parentLink   reflect.Type
}

//Execution Execution script
func (ins *InNetScript) Execution(fileName string,
	parent unsafe.Pointer,
	parentLink reflect.Type) {

	ins.handle = stack.NewLuaStack()
	ins.parentLink = parentLink
	defer ins.handle.GetLuaState().Close()

	ins.handle.GetLuaState().OpenLibs()

	ins.handle.GetLuaState().PushGoClosure(ins.luaRegisterProtobuf, (uintptr)(parent))
	ins.handle.GetLuaState().SetGlobal("register_proto_method")

	if _, err := ins.handle.ExecuteScriptFile(fileName); err != nil {
		panic(err)
	}
}

func (ins *InNetScript) luaRegisterProtobuf(L *mlua.State) int {
	p  := L.ToLightGoStruct(L.UpvalueIndex(1))

	if p == nil {
		return L.Error("Upvalue index is empty")
	}

	c  := reflect.NewAt(ins.parentLink, p)
	util.Assert(!c.IsNil(), "Lua Client Service is Null")

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

	methodRegister := c.MethodByName("RegisterMethod")
	util.Assert(!methodRegister.IsNil(), "Registration method function does not exist")
	regObject := c.FieldByName(methodObjectName)
	util.Assert(!regObject.IsNil(), methodObjectName + " Member not found")
	methodFunc := regObject.MethodByName(methodName)
	util.Assert(!methodFunc.IsNil(), methodObjectName + " " +  methodName + " method not found")

	protokey := reflect.Indirect(reflect.New(protoType.Elem())).Addr().Interface().(proto.Message)

	params := make([]reflect.Value, 2)
	params[0] = reflect.ValueOf(protokey)
	params[1] = methodFunc
	methodRegister.Call(params)

	return 0
}
