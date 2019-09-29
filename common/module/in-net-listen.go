package module

import (
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/pactum"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/service/net"
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

//SpawnInNetListen xxx
func SpawnInNetListen(clients implement.INetClientManager,
	deleate implement.INetListenDeleate,
	addr string,
	ccmax int,
	max int,
	keep uint64) InNetListen {
	return InNetListen{NetListenService: implement.NetListenService{
		NetListen:  &net.TCPListen{},
		NetClients: clients,
		NetDeleate: deleate,
		Addr:       addr,
		CCMax:      ccmax,
		MaxClient:  max,
		ClientKeep: keep}}
}

//InNetListen Intranet listening service base class
type InNetListen struct {
	implement.NetListenService
}

//Init InNetListen Internet listening service initialization
func (inet *InNetListen) Init() {
	inet.NetClients.Init()
	inet.NetListenService.Init()
	inet.RegisterMethod(&actor.Stopped{}, inet.Stoped)
	inet.RegisterMethod(&network.NetChunk{}, inet.OnRecv)
	inet.RegisterMethod(&network.NetClose{}, inet.OnClose)
}

//Stoped Turn off network monitoring service
func (inet *InNetListen) Stoped(context actor.Context, message interface{}) {
	inet.LogInfo("Service Stoping %s", inet.Addr)
	hls := inet.NetClients.GetHandles()
	if hls != nil && len(hls) > 0 {
		for inet.NetClients.Size() > 0 {
			ick := 0
			for i := 0; i < len(hls); i++ {
				c := inet.NetClients.GrapSocket(int32(hls[i]))
				if c == nil {
					continue
				}
				sck := c.GetSocket()
				inet.NetClients.Release(c)
				network.OperClose(sck)
			}

			for {
				time.Sleep(time.Duration(500) * time.Microsecond)
				if inet.NetClients.Size() <= 0 {
					break
				}

				inet.LogInfo("Service The remaining %d connections need to be closed", inet.NetClients.Size())
				ick++
				if ick > 6 {
					break
				}
			}
		}
	}
	inet.NetListen.Close()
	inet.LogInfo("Service Stoped")
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

//OnClose Close connection event
func (inet *InNetListen) OnClose(context actor.Context, message interface{}) {
	closer := message.(*network.NetClose)
	inet.LogDebug("close socket:%d", closer.Handle)
	c := inet.NetClients.GrapSocket(closer.Handle)
	if c == nil {
		inet.LogError("close unfind map-id socket %d", closer.Handle)
		return
	}

	defer inet.NetClients.Release(c)

	hClose := c.GetID()
	hClose |= (uint64(closer.Handle) << 32)

	inet.NetClients.Erase(hClose)

	if err := inet.NetDeleate.UnOnlineNotification(hClose); err != nil {
		inet.LogDebug("closed client Notification %+v", err)
	}

	inet.LogDebug("closed client: %+v", hClose)
}
