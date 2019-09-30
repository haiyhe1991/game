package module

import (
	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/magicNet/service/implement"
)

//GetLogicMethod 返回已经注册的方法
func GetLogicMethod(name string, md *implement.NetMethodDispatch) implement.NetMethodFun {
	evtType := proto.MessageType(name)
	if evtType == nil {
		return nil
	}

	return md.GetType(evtType)
}
