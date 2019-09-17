package component

import (
	"time"
	"github.com/yamakiller/magicNet/timer"
	"github.com/yamakiller/magicNet/util"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/gateway/constant"
	"github.com/yamakiller/game/gateway/elements"
	"github.com/yamakiller/game/gateway/elements/clients"
	"github.com/yamakiller/game/gateway/elements/route"
	pkg "github.com/yamakiller/game/proto"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
)

// NewTCPNetworkService Create a tcp network service
func NewTCPNetworkService() *NetworkService {
	return service.Make(constant.ConstNetworkServiceName, func() service.IService {

		handle := &NetworkService{TCPService: service.TCPService{
			Addr:  constant.GatewayAddr,
			CCMax: constant.GatewayCCMax},
		}

		handle.OnAccept = handle.onAccept
		handle.OnRecv = handle.onRecv
		handle.OnClose = handle.onClose

		handle.Init()
		return handle
	}).(*NetworkService)
}

// NetworkService Components that provide network services
type NetworkService struct {
	service.TCPService
}

//Init Initialize the network service object
func (ns *NetworkService) Init() {
	elements.Clients.Initial(constant.GatewayID)
	ns.TCPService.Init()
	ns.RegisterMethod(agreement.CertificationConfirmation{}, ns.onConfirm)
}

//Shutdown Termination of network services
func (ns *NetworkService) Shutdown() {
	logger.Info(ns.ID(), "Network Listen [TCP/IP] Service Closing connection")
	hs := elements.Clients.GetHandls()
	for elements.Clients.Size() > 0 {
		chk := 0
		for i := 0; i < len(hs); i++ {
			network.OperClose(hs[i].SocketID())
		}

		for {
			time.Sleep(time.Duration(500) * time.Microsecond)
			if elements.Clients.Size() <= 0 {
				break
			}

			logger.Info(ns.ID(), "Network Listen [TCP/IP] Service The remaining %d connections need to be closed", elements.Clients.Size())
			chk++
			if chk > 6 {
				break
			}
		}
	}

	logger.Info(ns.ID(), "Network Listen [TCP/IP] Service All connections are closed")

	ns.TCPService.Shutdown()
}

func (ns *NetworkService) onAccept(self actor.Context, message interface{}) {
	accepter := message.(*network.NetAccept)
	if elements.Clients.Size()+1 > constant.GatewayMaxConnect {
		network.OperClose(accepter.Handle)
		logger.Warning(self.Self().GetID(), "accept client fulled")
		return
	}

	client, _, err := elements.Clients.Occupy(accepter.Handle, accepter.Addr, accepter.Port)
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
	network.OperSetKeep(accepter.Handle, uint64(constant.GatewayConnectKleep))
	//------------------------------------------------
	//First handshake and agree on the key
	shake, _ := proto.Marshal(&pkg.HandshakeResponse{Key: ""})
	shake = agreement.ExtAssemble(1, "proto.HandshakeResponse", shake, int32(len(shake)))
	network.OperWrite(accepter.Handle, shake, len(shake))

	//-------------------------------------------------
	client.Stat.UpdateOnline(timer.Now())
	elements.Clients.Release(client)
	logger.Debug(self.Self().GetID(), "accept client %d-%s:%d", accepter.Handle, accepter.Addr.String(), accepter.Port)
}

func (ns *NetworkService) onRecv(self actor.Context, message interface{}) {
	data := message.(*network.NetChunk)
	hid, herr := elements.Clients.ToHandleID(data.Handle)
	if herr != nil {
		logger.Trace(self.Self().GetID(), "recv error closed unfind map-id socket %d", data.Handle)
		network.OperClose(data.Handle)
		return
	}

	recvHandle := util.NetHandle{}
	recvHandle.Generate(constant.GatewayID, 0, int32(hid), data.Handle)

	client := elements.Clients.Grap(&recvHandle)
	if client == nil {
		logger.Trace(self.Self().GetID(), "recv unfind player %d-%d-%d-%d",
			recvHandle.GatewayID(),
			recvHandle.WorldID(),
			recvHandle.HandleID(),
			recvHandle.SocketID())
		return
	}

	var (
		space  int
		writed int
		wby    int
		pos    int

		pkName string
		pkData []byte
		pkErro error
	)

	for {
		space = constant.ConstPlayerBufferLimit - client.DataLen()
		wby = len(data.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err := client.DataWrite(data.Data[pos : pos+space])
			if err != nil {
				logger.Trace(self.Self().GetID(), "recv error %s socket %d", err.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto releaseclient
			}

			pos += space
			writed += space

			client.Stat.UpdateRecv(timer.Now(), uint64(space))
		}

		for {
			pkName, pkData, pkErro = client.DataAnalysis()
			if pkErro != nil {
				logger.Error(self.Self().GetID(), "recv error %s socket %d closing play", pkErro.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto releaseclient
			}

			if pkData != nil {
				pkErro = ns.onRoute(client, pkName, pkData)
				if pkErro != nil {
					logger.Error(self.Self().GetID(), "route error %s socket %d closing play", pkErro.Error(), data.Handle)
					network.OperClose(data.Handle)
					goto releaseclient
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
	elements.Clients.Release(client)
	logger.Debug(self.Self().GetID(), "Exit onRecv")
}

func (ns *NetworkService) onClose(self actor.Context, message interface{}) {
	closer := message.(*network.NetClose)
	logger.Debug(self.Self().GetID(), "关闭")
	hid, herr := elements.Clients.ToHandleID(closer.Handle)
	if herr != nil {
		logger.Trace(self.Self().GetID(), "close unfind map-id socket %d", closer.Handle)
		return
	}

	closeHandle := util.NetHandle{}
	closeHandle.Generate(constant.GatewayID, 0, int32(hid), closer.Handle)

	client := elements.Clients.Grap(&closeHandle)
	if client == nil {
		logger.Trace(self.Self().GetID(), "close unfind client %d-%d-%d-%d",
			closeHandle.GatewayID(),
			closeHandle.WorldID(),
			closeHandle.HandleID(),
			closeHandle.SocketID())
		goto unline
	}

	closeHandle = client.Handle
	elements.Clients.Erase(&closeHandle)
	elements.Clients.Release(client)

	logger.Debug(self.Self().GetID(), "closed client %d-%d-%d-%d", closeHandle.GatewayID(),
																closeHandle.WorldID(),
																closeHandle.HandleID(),
																closeHandle.SocketID())
unline:
	ns.pushOffline(&closeHandle)
}

func (ns *NetworkService) onConfirm(self actor.Context, message interface{}) {
	confirm := message.(*agreement.CertificationConfirmation)
	confirmHandle := util.NetHandle{}
	confirmHandle.SetValue(confirm.Handle)
	client := elements.Clients.Grap(&confirmHandle)
	if client == nil {
		logger.Trace(self.Self().GetID(), "Authentication confirmation failed, target connection not found")
		return
	}

	client.Auth = timer.Now()

	logger.Debug(self.Self().GetID(), "Connection authentication succeeded %+v", confirmHandle.HandleID)
}

func (ns *NetworkService) onRoute(client *clients.Client, name string, data []byte) error {
	msgType := proto.MessageType(name)
	if msgType == nil {
		logger.Error(ns.ID(), "route error %s", route.ErrRouteAgreeUnDefined.Error())
		return route.ErrRouteAgreeUnDefined
	}

	adr := elements.RouteAddress.Sreach(msgType)
	if adr == nil {
		logger.Error(ns.ID(), "route error %s", route.ErrRouteAgreeUnRegister.Error())
		return route.ErrRouteAgreeUnRegister
	}

	if adr.Auth && client.Auth == 0 {
		logger.Error(ns.ID(), "route error Protocol needs to be verified, this connection is not verified and not verified")
		return route.ErrRoutePlayerUnverified
	}

	if constant.ConnectServicePID.ID == 0 {
		logger.Error(ns.ID(), "route error Service has not started yet")
		return route.ErrRouteServiceNotStarted
	}

	//
	
	actor.DefaultSchedulerContext.Send(&constant.ConnectServicePID,
		&agreement.ForwardMessage{Handle: client.Handle.GetValue(),
			AgreementName: name,
			ServerName:    adr.ServiceName,
			Data:          data})

	logger.Debug(ns.ID(), "forward message agreement name:%s server name:%d data length:%d", name, adr.ServiceName, len(data))

	return nil
}

func (ns *NetworkService) pushOffline(h *util.NetHandle) {

}
