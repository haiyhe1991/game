package component

import (
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
	"github.com/yamakiller/magicNet/util"
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

//FromGatewayMessage Message from the gateway
type FromGatewayMessage struct {
	Handle uint64
	Sock   int32
	Name   string
	Data   interface{}
}

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
}

//Init Initialize the network service object
func (ns *InNetService) Init() {
	ns.CS.(*inclients.InClientManager).Spawned()
	ns.TCPService.Init()
	ns.RegisterMethod(&pkg.GatewayRegisterRequest{}, ns.onGatewayRegister)
}

func (ns *InNetService) onAccept(self actor.Context, message interface{}) {
	accepter := message.(*network.NetAccept)
	if ns.CS.Size()+1 > constant.GatewayMaxConnect {
		network.OperClose(accepter.Handle)
		logger.Warning(self.Self().GetID(), "accept client fulled")
		return
	}

	client, k, err := ns.CS.Occupy(accepter.Handle, accepter.Addr, accepter.Port)
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
	h := util.NetHandle{}
	h.Generate(0, 0, k, accepter.Handle)

	shake, _ := proto.Marshal(&pkg.HandshakeResponse{Key: ""})
	shake = agreement.AgentParser(agreement.ConstInParser).Assemble(agreement.ConstArgeeVersion, h.GetValue(), "proto.HandshakeResponse", shake, int32(len(shake)))
	network.OperWrite(accepter.Handle, shake, len(shake))

	//-------------------------------------------------
	client.GetStat().UpdateOnline(timer.Now())
	ns.CS.Release(client)
	logger.Debug(self.Self().GetID(), "accept client %d-%s:%d", accepter.Handle, accepter.Addr.String(), accepter.Port)
}

func (ns *InNetService) onRecv(self actor.Context, message interface{}) {
	data := message.(*network.NetChunk)
	k, err := ns.CS.ToKey(data.Handle)
	if err != nil {
		logger.Error(self.Self().GetID(), "Receive data error did not find service connection:%+v", err)
		return
	}

	client := ns.CS.Grap(int32(k))
	if client == nil {
		logger.Error(self.Self().GetID(), "Receive data error did not find service connection:%+v", k)
		return
	}

	var (
		space  int
		writed int
		wby    int
		pos    int

		pkHandle uint64
		pkName   string
		pkData   []byte
		pkErro   error
	)

	for {
		space = constant.ConstPlayerBufferLimit - client.GetData().Len()
		wby = len(data.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err := client.GetData().Write(data.Data[pos : pos+space])
			if err != nil {
				logger.Trace(self.Self().GetID(), "recv error %s socket %d", err.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto releaseclient
			}

			pos += space
			writed += space

			client.GetStat().UpdateRecv(timer.Now(), uint64(space))
		}

		c := client.(*inclients.InClient)

		for {
			pkName, pkHandle, pkData, pkErro = c.Analysis()
			if pkErro != nil {
				logger.Error(self.Self().GetID(), "recv error %s socket %d closing gateway", pkErro.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto releaseclient
			}

			if pkData != nil {
				msgType := proto.MessageType(pkName)
				if msgType == nil {
					logger.Error(self.Self().GetID(), "recv error unknown protocol %s socket %d from gateway", pkName, data.Handle)
					continue
				}

				if f := ns.GetMethod(msgType); f != nil {
					f(self, &FromGatewayMessage{Handle: pkHandle, Name: pkName, Sock: data.Handle, Data: pkData})
				}
				continue
			}

			if writed >= len(data.Data) {
				goto releaseclient
			} else {
				break
			}
		}
	}
releaseclient:
	ns.CS.Release(client)
	logger.Debug(self.Self().GetID(), "Exit onRecv")
}

func (ns *InNetService) onGatewayRegister(self actor.Context, message interface{}) {

}

func (ns *InNetService) onClose(self actor.Context, message interface{}) {

}
