package logic

import (
	"github.com/yamakiller/game/common/module"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/service/implement"
)

//SignOutProc Logout Processing logic process
type SignOutProc struct {
}

//OnProccess Logout Processing logic method
func (proc *SignOutProc) OnProccess(event implement.INetMethodEvent) {
	request := event.(*module.InNetMethodClientEvent)
	c := module.GetWarehouse().GrapSocket(request.Socket)
	if c == nil {
		logger.Error(request.Context.Self().GetID(), "Exception to find the target connection")
		return
	}
	defer module.GetWarehouse().Release(c)
	proc.defaultProccess(c, request)
}

func (proc *SignOutProc) defaultProccess(c interface{}, request *module.InNetMethodClientEvent) {

}
