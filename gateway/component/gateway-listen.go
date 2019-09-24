package component

import (
	"fmt"
	"reflect"

	"github.com/yamakiller/game/gateway/elements/forward"
	"github.com/yamakiller/game/pactum"

	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/timer"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/gateway/constant"
	"github.com/yamakiller/game/gateway/elements"
	"github.com/yamakiller/game/gateway/elements/clients"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/service/net"
	"github.com/yamakiller/magicNet/util"
)

//NewGatewayListener Create a network listening service
func NewGatewayListener() *GatewayListener {
	return &GatewayListener{NetListenService: implement.NetListenService{
		NetListen:             &net.TCPListen{},
		NetDeleate:            &GNetListenDeleate{},
		NetClients:            clients.NewGClientManager(),
		ClientKeep:            uint64(constant.GatewayConnectKleep),
		ClientRecvBufferLimit: constant.ConstClientBufferLimit}}
}

//GNetListenDeleate network listening service, delegate logic
type GNetListenDeleate struct {
}

//Handshake Handshake processing
func (gnld *GNetListenDeleate) Handshake(c implement.INetClient) error {
	c.BuildKeyPair()
	publicKey := c.GetKeyPublic()
	shake, _ := proto.Marshal(&pactum.HandshakeResponse{Key: publicKey})
	shake = agreement.AgentParser(agreement.ConstExParser).Assemble(nil,
		agreement.ConstPactumVersion,
		0,
		"proto.HandshakeResponse",
		shake,
		int32(len(shake)))

	if err := network.OperWrite(c.GetSocket(), shake, len(shake)); err != nil {
		return err
	}
	return nil
}

//Analysis Client packet analysis from the Internet
func (gnld *GNetListenDeleate) Analysis(context actor.Context,
	nets *implement.NetListenService,
	c implement.INetClient) error {

	name, _, wrap, err := agreement.AgentParser(
		agreement.ConstExParser).Analysis(c.GetKeyPair(), c.GetRecvBuffer())
	if err != nil {
		return err
	}

	if wrap == nil {
		return implement.ErrAnalysisProceed
	}

	var unit *forward.Unit
	var fpid *actor.PID
	msgType := proto.MessageType(name)
	if msgType != nil {
		if f := nets.GetMethod(msgType); f != nil {
			f(context, wrap)
			goto end
		}
	}

	unit = elements.ForwardAddresses.Sreach(msgType)
	if unit == nil {
		return fmt.Errorf("Abnormal protocol, no protocol information defined")
	}

	if unit.Auth && c.GetAuth() == 0 {
		return forward.ErrForwardClientUnverified
	}

	fpid = elements.SSets.Sreach(constant.ConstForwardServiceName)
	if fpid == nil {
		return forward.ErrForwardServiceNotStarted
	}

	actor.DefaultSchedulerContext.Send(fpid,
		&agreement.ForwardServerEvent{Handle: c.GetID().GetValue(),
			PactumName: name,
			ServoName:  unit.ServoName,
			Data:       wrap})
end:
	//return name, data, err
	return implement.ErrAnalysisSuccess
}

//UnOnlineNotification Offline notification
func (gnld *GNetListenDeleate) UnOnlineNotification(h util.NetHandle) error {
	msgType := proto.MessageType(constant.GatewayLogoutPactum)
	if msgType == nil {
		return fmt.Errorf("An error occurred while processing the offline "+
			"notification. The %s protocol is not defined",
			constant.GatewayLogoutPactum)
	}

	wrap, err := proto.Marshal(
		reflect.Indirect(
			reflect.New(msgType.Elem())).Addr().Interface().(proto.Message))

	if err != nil {
		return err
	}

	fpid := elements.SSets.Sreach(constant.ConstForwardServiceName)
	if fpid == nil {
		return forward.ErrForwardServiceNotStarted
	}

	actor.DefaultSchedulerContext.Send(fpid,
		&agreement.ForwardServerEvent{Handle: h.GetValue(),
			PactumName: constant.GatewayLogoutPactum,
			ServoName:  constant.GatewayLogoutName,
			Data:       wrap})
	return nil
}

//GatewayListener Gateway Internet monitoring service
type GatewayListener struct {
	implement.NetListenService
}

//Init Gateway Internet listening service initialization
func (gnet *GatewayListener) Init() {
	gnet.NetClients.Init()
	gnet.NetListenService.Init()
	gnet.RegisterMethod(&agreement.ForwardClientEvent{}, gnet.onForwardClient)
}

//Started xxx
func (gnet *GatewayListener) Started(context actor.Context,
	message interface{}) {
	elements.SSets.Push(gnet.Key(), context.Self())
	gnet.NetListenService.Started(context, message)
}

//onForwardClient
func (gnet *GatewayListener) onForwardClient(context actor.Context,
	message interface{}) {
	msg := message.(*agreement.ForwardClientEvent)
	h := util.NetHandle{}
	h.SetValue(msg.Handle)
	c := gnet.NetClients.Grap(&h)
	if c == nil {
		gnet.LogError("Failed to send data, %+v client does not exist", msg.Handle)
		return
	}

	defer gnet.NetClients.Release(c)
	sd := agreement.AgentParser(agreement.ConstExParser).Assemble(
		c.GetKeyPair(),
		agreement.ConstPactumVersion,
		0,
		msg.PactumName,
		msg.Data,
		int32(len(msg.Data)))

	if sd == nil {
		gnet.LogError("Protocol package failed to be assembled: pactum name %s",
			msg.PactumName)
		return
	}

	if err := network.OperWrite(h.GetSocket(), sd, len(sd)); err != nil {
		gnet.LogError("Failed to send data to client socket: %+v pactum name %s",
			err, msg.PactumName)
		return
	}

	c.GetStat().UpdateRead(timer.Now(), uint64(len(sd)))
	gnet.LogDebug("Already sent data to the client: pactum name %s", msg.PactumName)
}

//Shutdown Turn off the gateway network listening service
func (gnet *GatewayListener) Shutdown() {
	name := gnet.Key()
	gnet.NetListenService.Shutdown()
	elements.SSets.Erase(name)
}
