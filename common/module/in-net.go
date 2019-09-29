package module

import (
	"reflect"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/service/implement"
)

//InNetMethodEvent 内网数据包事件
type InNetMethodEvent struct {
	H uint64
	implement.NetMethodEvent
}

//InNetMethodClientEvent 内网客户端数据包事件
type InNetMethodClientEvent struct {
	Context   actor.Context
	ProtoType reflect.Type
	InNetMethodEvent
}
