package module

import (
	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/pactum"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service/implement"
)

type InNetListenDeleate struct {
}

//Handshake Intranet listening service delegation base interface
func (inld *InNetListenDeleate) Handshake(c implement.INetClient) error {
	shake, _ := proto.Marshal(&pactum.HandshakeResponse{Key: ""})
	shake = agreement.AgentParser(agreement.ConstInParser).Assemble(nil,
		agreement.ConstPactumVersion,
		0, //修改
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

//
type InNetListen struct {
	implement.NetListenService
}

//Init InNetListen Internet listening service initialization
func (inet *InNetListen) Init() {
	inet.NetClients.Init()
	inet.NetListenService.Init()
}

//Started xxx
func (inet *InNetListen) Started(context actor.Context,
	message interface{}) {
	inet.NetListenService.Started(context, message)
}
