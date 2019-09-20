package component

import (
	"github.com/yamakiller/magicNet/timer"
	"github.com/yamakiller/magicNet/util"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	cmcm "github.com/yamakiller/game/common/component"
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

// NewTCPNetService Create a tcp network service
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

var (
	localTargetServiceNames []string = nil
)

// OutNetService Components that provide network services
type OutNetService struct {
	cmcm.NetService
}

//Init Initialize the network service object
func (ns *OutNetService) Init() {
	ns.CS.(*clients.ClientManager).Spawned()
	ns.TCPService.Init()
	ns.RegisterMethod(agreement.CertificationConfirmation{}, ns.onConfirm)
}

func (ns *OutNetService) onAccept(self actor.Context, message interface{}) {
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
	//Bind gateway ID
	client.(*clients.Client).SetGateway(constant.GatewayID)

	network.OperOpen(accepter.Handle)
	network.OperSetKeep(accepter.Handle, uint64(constant.GatewayConnectKleep))
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

func (ns *OutNetService) onRecv(self actor.Context, message interface{}) {
	data := message.(*network.NetChunk)
	hKey, herr := ns.CS.ToKey(data.Handle)
	if herr != nil {
		logger.Trace(self.Self().GetID(), "recv error closed unfind map-id socket %d", data.Handle)
		network.OperClose(data.Handle)
		return
	}

	recvHandle := util.NetHandle{}
	recvHandle.Generate(constant.GatewayID, 0, int32(hKey), data.Handle)

	client := ns.CS.Grap(recvHandle.HandleID())
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

		c := client.(*clients.Client)

		for {
			pkName, pkData, pkErro = c.Analysis()
			if pkErro != nil {
				logger.Error(self.Self().GetID(), "recv error %s socket %d closing client", pkErro.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto releaseclient
			}

			if pkData != nil {
				pkErro = ns.onRoute(c, pkName, pkData)
				if pkErro != nil {
					logger.Error(self.Self().GetID(), "route error %s socket %d closing client", pkErro.Error(), data.Handle)
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
	ns.CS.Release(client)
	logger.Debug(self.Self().GetID(), "Exit onRecv")
}

func (ns *OutNetService) onClose(self actor.Context, message interface{}) {
	closer := message.(*network.NetClose)
	logger.Debug(self.Self().GetID(), "close socket:%d", closer.Handle)
	hid, herr := ns.CS.ToKey(closer.Handle)
	if herr != nil {
		logger.Trace(self.Self().GetID(), "close unfind map-id socket %d", closer.Handle)
		return
	}

	closeHandle := util.NetHandle{}
	closeHandle.Generate(constant.GatewayID, 0, int32(hid), closer.Handle)

	client := ns.CS.Grap(closeHandle.HandleID())
	if client == nil {
		logger.Trace(self.Self().GetID(), "close unfind client %d-%d-%d-%d",
			closeHandle.GatewayID(),
			closeHandle.WorldID(),
			closeHandle.HandleID(),
			closeHandle.SocketID())
		goto unline
	}

	ns.CS.Erase(client.GetKey())
	ns.CS.Release(client)

	logger.Debug(self.Self().GetID(), "closed client %d-%d-%d-%d", closeHandle.GatewayID(),
		closeHandle.WorldID(),
		closeHandle.HandleID(),
		closeHandle.SocketID())
unline:
	ns.blukOffline(&closeHandle)
}

func (ns *OutNetService) onConfirm(self actor.Context, message interface{}) {
	confirm := message.(*agreement.CertificationConfirmation)
	confirmHandle := util.NetHandle{}
	confirmHandle.SetValue(confirm.Handle)
	client := ns.CS.Grap(confirmHandle.HandleID())
	if client == nil {
		logger.Trace(self.Self().GetID(), "Authentication confirmation failed, target connection not found")
		return
	}

	client.SetAuth(timer.Now())

	logger.Debug(self.Self().GetID(), "Connection authentication succeeded %+v", confirmHandle.HandleID)
}

func (ns *OutNetService) onRoute(client *clients.Client, name string, data []byte) error {
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

	if adr.Auth && client.GetAuth() == 0 {
		logger.Error(ns.ID(), "route error Protocol needs to be verified, this connection is not verified and not verified")
		return route.ErrRoutePlayerUnverified
	}

	if constant.ConnectServicePID.ID == 0 {
		logger.Error(ns.ID(), "route error Service has not started yet")
		return route.ErrRouteServiceNotStarted
	}

	actor.DefaultSchedulerContext.Send(&constant.ConnectServicePID,
		&agreement.ForwardMessage{Handle: client.GetKeyValue(),
			AgreementName: name,
			ServerName:    adr.ServiceName,
			Data:          data})

	logger.Debug(ns.ID(), "forward message agreement name:%s server name:%d data length:%d", name, adr.ServiceName, len(data))

	return nil
}

func (ns *OutNetService) blukOffline(h *util.NetHandle) {
	if h.WorldID() == 0 {
		logger.Debug(ns.ID(), "No need to log in to send offline notifications:[socket-%d]", h.SocketID())
		return
	}

	pak := pkg.UnLoginRequest{}
	pak.Tts = timer.Now()
	data, err := proto.Marshal(&pak)

	if err != nil {
		logger.Error(ns.ID(), "forward offline message fail:%s", err)
		return
	}

	if localTargetServiceNames == nil {
		localTargetServiceNames = elements.Conns.GetGroupNames()
	}

	//Bulk offline notification
	for _, name := range localTargetServiceNames {
		actor.DefaultSchedulerContext.Send(&constant.ConnectServicePID,
			&agreement.ForwardMessage{Handle: h.GetValue(),
				AgreementName: "proto.UnLoginRequest",
				ServerName:    name,
				Data:          data})
	}
}
