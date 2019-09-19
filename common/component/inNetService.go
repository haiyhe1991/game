package component

import (
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/common/elements/inclients"
	"github.com/yamakiller/game/common/elements/visitors"
	"github.com/yamakiller/game/gateway/constant"
	pkg "github.com/yamakiller/game/proto"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/timer"
)

/*
func NewTCPNetService() *OutNetService {
	return service.Make(constant.ConstNetworkServiceName, func() service.IService {

		handle := &OutNetService{
			NetService: *cmcm.NewTCPService(constant.GatewayAddr, constant.GatewayCCMax, elements.Clients),
		}

		handle.OnAccept = handle.onAccept
		handle.OnRecv = handle.onRecv
		handle.OnClose = handle.onClose

		handle.Init()
		return handle
	}).(*OutNetService)
}
*/

//NewTCPINNetService Internal network service
func NewTCPINNetService(addr string, ccmax int, cs visitors.IVisitorManager) *InNetService {
	result := &InNetService{NetService: *NewTCPService(addr, ccmax, cs)}
	result.OnAccept = result.onAccept
	result.OnRecv = result.onRecv
	result.OnClose = result.onClose
	return result
}

//InNetService Internal network service
type InNetService struct {
	NetService
	netMethod map[interface{}]func(self actor.Context, message interface{})
}

//Init Initialize the network service object
func (ns *InNetService) Init() {
	ns.CS.(*inclients.InClientManager).Spawned()
	ns.netMethod = make(map[interface{}]func(self actor.Context, message interface{}))
	ns.TCPService.Init()
}

//RegisterNetMethod Registration network method
func (ns *InNetService) RegisterNetMethod(agree interface{}, f func(self actor.Context, message interface{})) {
	ns.netMethod[reflect.TypeOf(agree)] = f
}

func (ns *InNetService) getNetMethod(agree interface{}) func(self actor.Context, message interface{}) {
	if f, ok := ns.netMethod[agree]; ok {
		return f
	}
	return nil
}

func (ns *InNetService) onAccept(self actor.Context, message interface{}) {
	accepter := message.(*network.NetAccept)
	if ns.CS.Size()+1 > constant.GatewayMaxConnect {
		network.OperClose(accepter.Handle)
		logger.Warning(self.Self().GetID(), "accept client fulled")
		return
	}

	client, _, err := ns.CS.Occupy(accepter.Handle, accepter.Addr, accepter.Port)
	if err != nil {
		//close-socket
		network.OperClose(accepter.Handle)
		logger.Trace(self.Self().GetID(), "accept client closed: %v, %d-%s:%d", err,
			accepter.Handle,
			accepter.Addr.String(),
			accepter.Port)
		return
	}

	network.OperOpen(accepter.Handle)
	network.OperSetKeep(accepter.Handle, uint64(200000))

	//------------------------------------------------
	//First handshake and agree on the key
	shake, _ := proto.Marshal(&pkg.HandshakeResponse{Key: ""})
	shake = agreement.AgentParser(agreement.ConstExParser).Assemble(1, 0, "proto.HandshakeResponse", shake, int32(len(shake)))
	network.OperWrite(accepter.Handle, shake, len(shake))

	//-------------------------------------------------
	client.GetStat().UpdateOnline(timer.Now())
	ns.CS.Release(client)
	logger.Debug(self.Self().GetID(), "accept client %d-%s:%d", accepter.Handle, accepter.Addr.String(), accepter.Port)
}

func (ns *InNetService) onRecv(self actor.Context, message interface{}) {
}

func (ns *InNetService) onClose(self actor.Context, message interface{}) {

}
