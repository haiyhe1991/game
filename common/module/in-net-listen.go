package module

import (
	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/pactum"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service/implement"
)

//InNetListenDeleate Intranet listening delegation method
type InNetListenDeleate struct {
}

//Handshake Intranet listening service delegation base interface
func (inld *InNetListenDeleate) Handshake(c implement.INetClient) error {
	shake, _ := proto.Marshal(&pactum.HandshakeResponse{Key: ""})
	shake = agreement.AgentParser(agreement.ConstInParser).Assemble(nil,
		agreement.ConstPactumVersion,
		uint64(c.GetSocket()),
		"proto.HandshakeResponse",
		shake,
		int32(len(shake)))

	if err := network.OperWrite(c.GetSocket(), shake, len(shake)); err != nil {
		return err
	}
	return nil
}

//Analysis Client packet analysis from the Internet
func (inld *InNetListenDeleate) Analysis(context actor.Context,
	nets *implement.NetListenService,
	c implement.INetClient) error {
	//TODO:
	return nil
}

//UnOnlineNotification xxx
func (inld *InNetListenDeleate) UnOnlineNotification(h uint64) error {
	return nil
}

//InNetListen Intranet listening service base class
type InNetListen struct {
	implement.NetListenService
}

//Init InNetListen Internet listening service initialization
func (inet *InNetListen) Init() {
	inet.NetClients.Init()
	inet.NetListenService.Init()
}

//Started Re-register the OnRecv method
func (inet *InNetListen) Started(context actor.Context,
	message interface{}) {
	inet.NetListenService.Started(context, message)
	inet.RegisterMethod(&network.NetChunk{}, inet.OnRecv)
}

//OnRecv Overloaded OnRecv method
func (inet *InNetListen) OnRecv(context actor.Context,
	message interface{}) {
	defer inet.LogDebug("onRecv: complete")

	wrap := message.(*network.NetChunk)
	c := inet.NetClients.GrapSocket(wrap.Handle)
	if c == nil {
		inet.LogError("OnRecv: No target [%d] client service was found", wrap.Handle)
		return
	}

	csrv, conv := c.(*InNetClient)
	if !conv {
		inet.LogError("OnRecv: Failed to convert to service object does not work properly")
		return
	}

	if csrv.GetPID() == nil {
		inet.LogError("OnRecv: The target service is not running and is not working properly")
		return
	}

	actor.DefaultSchedulerContext.Send(csrv.GetPID(), wrap)
}
