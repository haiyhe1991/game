package component

import (
	"bytes"
	"time"

	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/timer"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/common/module"
	"github.com/yamakiller/game/gateway/constant"
	"github.com/yamakiller/game/gateway/elements"
	"github.com/yamakiller/game/gateway/elements/forward"
	"github.com/yamakiller/game/gateway/elements/servers"
	"github.com/yamakiller/game/pactum"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/service/net"
)

//GatewayHandleConnect  Gateway link handle
type GatewayHandleConnect struct {
	net.TCPConnection
	RecvBufferMax int

	auth     uint64
	rbuffer  *bytes.Buffer
	datastat implement.NetStat
}

//GetRecvBufferLimit Return this connection read buffer limit
func (ghc *GatewayHandleConnect) GetRecvBufferLimit() int {
	return ghc.RecvBufferMax
}

//GetRecvBuffer Return to read buffer
func (ghc *GatewayHandleConnect) GetRecvBuffer() *bytes.Buffer {
	return ghc.rbuffer
}

//GetDataStat Return data status information
func (ghc *GatewayHandleConnect) GetDataStat() net.INetConnectionDataStat {
	return &ghc.datastat
}

//GetAuth Return the authentication time
func (ghc *GatewayHandleConnect) GetAuth() uint64 {
	return ghc.auth
}

//SetAuth Set the authentication time
func (ghc *GatewayHandleConnect) SetAuth(auth uint64) {
	ghc.auth = auth
}

//Close Close thie connection
func (ghc *GatewayHandleConnect) Close() {
	ghc.TCPConnection.Close()
}

//GatewayConnectDeleate xxx
type GatewayConnectDeleate struct {
}

//Connected Connected proccess
func (gcd *GatewayConnectDeleate) Connected(context actor.Context,
	nets *implement.NetConnectService) error {
	return nil
}

//Analysis Packet decomposition
func (gcd *GatewayConnectDeleate) Analysis(context actor.Context,
	nets *implement.NetConnectService) error {
	name, h, wrap, err := agreement.AgentParser(agreement.ConstInParser).Analysis(nil,
		nets.Handle.GetRecvBuffer())
	if err != nil {
		return err
	}

	if wrap == nil {
		return implement.ErrAnalysisProceed
	}

	var fpid *actor.PID
	msgType := proto.MessageType(name)
	if msgType != nil {
		if f := nets.NetMethod.GetType(msgType); f != nil {
			f(&module.InNetMethodEvent{
				H:              h,
				NetMethodEvent: implement.NetMethodEvent{Name: name, Socket: nets.Handle.Socket(), Wrap: wrap},
			})
			goto end
		}
	}

	fpid = elements.SSets.Sreach(constant.ConstForwardServiceName)
	if fpid == nil {
		return forward.ErrForwardServiceNotStarted
	}

	actor.DefaultSchedulerContext.Send(fpid,
		&agreement.ForwardClientEvent{Handle: h,
			PactumName: name,
			Data:       wrap})

end:
	//return name, data, err
	return implement.ErrAnalysisSuccess
}

//GatewayConnect Gateway connector
type GatewayConnect struct {
	implement.NetConnectService
	GatewayID        int32
	AutoErrRetry     int //Automatic connection failure retries
	AutoErrRetryTime int //Automatic connection failure retry interval unit milliseconds
}

//Init Initialize connector
func (gconn *GatewayConnect) Init() {
	gconn.NetConnectService.Init()
	gconn.RegisterMethod(&network.NetClose{}, gconn.OnClose)
	gconn.NetMethod.Register(&pactum.HandshakeResponse{}, gconn.onNetHandshake)
	gconn.NetMethod.Register(&pactum.GatewayRegisterResponse{}, gconn.onNetRegisterResponse)
	gconn.RegisterMethod(&agreement.ForwardServerEvent{}, gconn.onForwardServer)
}

//OnClose Handling connection close events
func (gconn *GatewayConnect) OnClose(context actor.Context, message interface{}) {
	group := elements.TLSets.Get(gconn.Target.GetName())
	if group != nil {
		//Delete equalizer
		group.RemoveTarget(gconn.Name())
	} else {
		gconn.LogError("%s class equalizer does not exist", gconn.Target.GetName())
	}
	gconn.NetConnectService.OnClose(context, message)
}

func (gconn *GatewayConnect) onForwardServer(context actor.Context,
	message interface{}) {

	ick := 0
	var err error
	request := message.(*agreement.ForwardServerEvent)
	for {
		if gconn.IsShutdown() {
			gconn.LogWarning("Service has been terminated, data is discarded:%s",
				request.PactumName)
			return
		}

		if gconn.Target.GetEtat() != implement.Connected {
			if gconn.Target.GetEtat() == implement.UnConnected {
				err = gconn.NetConnectService.AutoConnect(context)
				if err != nil {
					gconn.LogError("onForwardServer: AutoConnect fail:%+v", err)

				}
			} else {
				if outTime := gconn.Target.IsTimeout(); outTime > 0 {
					gconn.LogWarning("onForwardServer: AutoConnect TimeOut:%d millisecond",
						outTime)
					gconn.Target.RestTimeout()
					gconn.Handle.Close()
				}
			}
		} else {
			break
		}

		ick++
		if ick > gconn.AutoErrRetry {
			gconn.LogError("OnForwardServer AutoConnect fail, "+
				"Data is discarded %+v %s %s %d check num:%d",
				request.Handle,
				request.PactumName,
				request.ServoName,
				len(request.Data),
				ick)
			return
		}
		time.Sleep(time.Duration(gconn.AutoErrRetryTime) * time.Millisecond)
	}

	//Assembling the inner net package
	fdata := agreement.AgentParser(agreement.ConstInParser).Assemble(nil,
		agreement.ConstPactumVersion,
		request.Handle,
		request.PactumName,
		request.Data,
		int32(len(request.Data)))

	//Send
	if err = gconn.Handle.Write(fdata, len(fdata)); err != nil {
		gconn.LogError("OnForwardServer Send fail, Data is discarded %+v %s %s %d",
			request.Handle,
			request.PactumName,
			request.ServoName,
			len(request.Data))
		return
	}
	gconn.LogDebug("OnForwardServer Send %s success", request.PactumName)
}

func (gconn *GatewayConnect) onNetHandshake(event implement.INetMethodEvent) {
	//Internal communication does not consider encrypted communication
	request := pactum.GatewayRegisterRequest{}
	request.Id = gconn.GatewayID
	var requestData []byte
	var err error

	if gconn.Target.GetEtat() != implement.Connecting {
		gconn.LogError("onNetHandshake: handshake fail: current status %+v,%+v",
			gconn.Target.GetEtat(), implement.Connecting)
		return
	}

	requestData, err = proto.Marshal(&request)
	if err != nil {
		gconn.LogError("onNetHandshake: handshake fail:%+v", err)
		goto unend
	}

	err = gconn.Handle.Write(requestData, len(requestData))
	if err != nil {
		gconn.LogError("onNetHandshake: Register ID fail:%+v", err)
		goto unend
	}

	gconn.Target.SetEtat(implement.Verify)
	return
unend:
	gconn.Target.SetEtat(implement.UnConnected)
}

func (gconn *GatewayConnect) onNetRegisterResponse(event implement.INetMethodEvent) {
	response := event.(*module.InNetMethodEvent)
	gconn.LogDebug("onNetRegisterResponse: remote handle:%+v %s",
		response.H, response.Name)
	responseMsg := pactum.GatewayRegisterResponse{}
	err := proto.Unmarshal(response.Wrap, &responseMsg)
	now := timer.Now()
	if err != nil {
		gconn.LogError("onNetRegisterResponse: unmarshal fail:%+v", err)
		gconn.Handle.Close()
		goto unend
	}

	if gconn.Target.GetEtat() != implement.Verify {
		gconn.LogError("onNetRegisterResponse: register fail: current status %+v,%+v",
			gconn.Target.GetEtat(), implement.Verify)
		return
	}

	gconn.Handle.SetAuth(now)
	gconn.Target.SetEtat(implement.Connected)
	//Register the connection to the loader======================================
	if group := elements.TLSets.Get(gconn.Target.GetName()); group != nil {
		group.AddTarget(gconn.Name(),
			&servers.TargeObject{ID: gconn.Target.(*servers.TargetConnection).GetVirtualID(),
				Target: gconn.GetPID()})
	} else { //registration failed
		gconn.Handle.Close()
		gconn.LogError("Registration to the loader failed, %s such loader does not exist",
			gconn.Target.GetName())
		goto unend
	}
	//===========================================================================
	gconn.LogInfo("onNetRegisterResponse: connected address:%s time:%+v success ",
		gconn.Target.GetAddr(), now)
	return
unend:
	gconn.Target.SetEtat(implement.UnConnected)
}
