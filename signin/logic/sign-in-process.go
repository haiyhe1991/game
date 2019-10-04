package logic

import (
	"github.com/yamakiller/game/common/module"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/service/implement"
)

//SignInProc Sign in Processing logic process
type SignInProc struct {
}

//OnProccess Sign in Processing logic method
func (proc *SignInProc) OnProccess(event implement.INetMethodEvent) {
	request := event.(*module.InNetMethodClientEvent)
	c := module.GetWarehouse().GrapSocket(request.Socket)
	if c == nil {
		logger.Error(request.Context.Self().GetID(), "Exception to find the target connection")
		return
	}
	defer module.GetWarehouse().Release(c)
	proc.defaultProccess(c, request)
}

func (proc *SignInProc) defaultProccess(c interface{}, request *module.InNetMethodClientEvent) {

}
