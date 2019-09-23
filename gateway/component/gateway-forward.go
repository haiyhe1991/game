package component

import (
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/service"
)

//GatewayForward Gateway forwarding service
type GatewayForward struct {
	service.Service
}

func (gf *GatewayForward) Init() {
	gf.Service.Init()
	gf.RegisterMethod(&actor.Started{}, gf.Started)
	gf.RegisterMethod(&actor.Stopped{}, gf.Stoped)
	gf.RegisterMethod(&agreement.ForwardClientEvent{}, gf.onForwardClient)
	gf.RegisterMethod(&agreement.ForwardServerEvent{}, gf.onForwardServer)
}

func (gf *GatewayForward) Started(context actor.Context, message interface{}) {
	gf.Service.Assignment(context)
	gf.LogInfo("Service Startup %s", gf.Name())
	//check link
	gf.Service.Started(context, message)
}

func (gf *GatewayForward) onForwardClient(context actor.Context, message interface{}) {

}

func (gf *GatewayForward) onForwardServer(context actor.Context, message interface{}) {

}
