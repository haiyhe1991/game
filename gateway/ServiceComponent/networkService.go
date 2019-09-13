package ServiceComponent

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/gateway/elements"
	"github.com/yamakiller/game/gateway/elements/agreement"
	"github.com/yamakiller/game/gateway/elements/clients"
	"github.com/yamakiller/game/gateway/elements/route"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/engine/util"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
)

// NewTCPNetworkService Create a tcp network service
func NewTCPNetworkService() *NetworkService {
	return service.Make(elements.ConstNetworkServiceName, func() service.IService {

		handle := &NetworkService{TCPService: service.TCPService{
			Addr:  elements.GatewayAddr,
			CCMax: elements.GatewayCCMax},
			pm: &clients.PlayManager{},
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
	pm *clients.PlayManager
}

//Init Initialize the network service object
func (ns *NetworkService) Init() {
	ns.pm.Initial(elements.GatewayID)
	ns.TCPService.Init()
}

//Shutdown Termination of network services
func (ns *NetworkService) Shutdown() {
	logger.Info(ns.ID(), "Network Listen [TCP/IP] Service Closing connection")
	hs := ns.pm.GetHandls()
	for ns.pm.Size() > 0 {
		chk := 0
		for i := 0; i < len(hs); i++ {
			network.OperClose(hs[i].SocketID())
		}

		for {
			time.Sleep(time.Duration(500) * time.Microsecond)
			if ns.pm.Size() <= 0 {
				break
			}

			logger.Info(ns.ID(), "Network Listen [TCP/IP] Service The remaining %d connections need to be closed", ns.pm.Size())
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
	accepter := message.(network.NetAccept)
	if ns.pm.Size()+1 > elements.GatewayMaxConnect {
		network.OperClose(accepter.Handle)
		logger.Warning(self.Self().GetID(), "accept player fulled")
		return
	}

	ply, _, err := ns.pm.Occupy(accepter.Handle, accepter.Addr, accepter.Port)
	if err != nil {
		//close-socket
		network.OperClose(accepter.Handle)
		logger.Trace(self.Self().GetID(), "accept player closed: %v, %d-%s:%d", err,
			accepter.Handle,
			accepter.Addr.String(),
			accepter.Port)
		return
	}

	ns.pm.Release(ply)
	logger.Trace(self.Self().GetID(), "accept player %d-%s:%d\n", accepter.Handle, accepter.Addr.String(), accepter.Port)
}

func (ns *NetworkService) onRecv(self actor.Context, message interface{}) {
	data := message.(network.NetChunk)
	hid := ns.pm.ToHandleID(data.Handle)
	if hid == 0 {
		logger.Trace(self.Self().GetID(), "recv error closed unfind map-id socket %d", data.Handle)
		network.OperClose(data.Handle)
		return
	}

	recvHandle := util.NetHandle{}
	recvHandle.Generate(elements.GatewayID, 0, int32(hid), data.Handle)

	ply := ns.pm.Grap(&recvHandle)
	if ply == nil {
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
		space = elements.ConstPlayerBufferLimit - ply.DataLen()
		wby = len(data.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err := ply.DataWrite(data.Data[pos : pos+space])
			if err != nil {
				logger.Trace(self.Self().GetID(), "recv error %s socket %d", err.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto releaseplay
			}

			pos += space
			wby += space
		}

		for {
			pkName, pkData, pkErro = ply.DataAnalysis()
			if pkErro != nil {
				logger.Error(self.Self().GetID(), "recv error %s socket %d closing play", pkErro.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto releaseplay
			}

			if pkData != nil {
				pkErro = ns.onRoute(ply, pkName, pkData)
				if pkErro != nil {
					logger.Error(self.Self().GetID(), "route error %s socket %d closing play", pkErro.Error(), data.Handle)
					network.OperClose(data.Handle)
					goto releaseplay
				}

				continue
			}

			if wby == 0 {
				goto releaseplay
			}
		}
	}
releaseplay:
	ns.pm.Release(ply)
}

func (ns *NetworkService) onClose(self actor.Context, message interface{}) {
	closer := message.(network.NetClose)
	hid := ns.pm.ToHandleID(closer.Handle)
	if hid == 0 {
		logger.Trace(self.Self().GetID(), "close unfind map-id socket %d", closer.Handle)
		return
	}

	closeHandle := util.NetHandle{}
	closeHandle.Generate(elements.GatewayID, 0, int32(hid), closer.Handle)

	ply := ns.pm.Grap(&closeHandle)
	if ply == nil {
		logger.Trace(self.Self().GetID(), "close unfind player %d-%d-%d-%d",
			closeHandle.GatewayID(),
			closeHandle.WorldID(),
			closeHandle.HandleID(),
			closeHandle.SocketID())
		goto unline
	}

	closeHandle = ply.Handle
	ns.pm.Erase(&closeHandle)
	ns.pm.Release(ply)
unline:
	ns.pushOffline(&closeHandle)
}

func (ns *NetworkService) onRoute(ply *clients.Player, name string, data []byte) error {
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

	if adr.Auth && ply.Auth == 0 {
		logger.Error(ns.ID(), "route error Protocol needs to be verified, this connection is not verified and not verified")
		return route.ErrRoutePlayerUnverified
	}

	//! 需要将数据发送给连接服务
	if elements.ConnectServicePID.ID == 0 {
		logger.Error(ns.ID(), "route error Service has not started yet")
		return route.ErrRouteServiceNotStarted
	}

	//
	actor.DefaultSchedulerContext.Send(&elements.ConnectServicePID,
		&agreement.ForwardMessage{Handle: ply.Handle.GetValue(),
			AgreementName: name,
			ServerName:    adr.ServiceName,
			Data:          data})

	return nil
}

func (ns *NetworkService) pushOffline(h *util.NetHandle) {

}
