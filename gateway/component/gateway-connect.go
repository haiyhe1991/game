package component

import (
	"bytes"

	"github.com/yamakiller/magicNet/timer"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/gateway/constant"
	"github.com/yamakiller/game/gateway/elements"
	"github.com/yamakiller/game/gateway/elements/forward"
	"github.com/yamakiller/game/pactum"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/service/net"
)

//ForwardPacket Packets from internal network communication
type ForwardPacket struct {
	H    uint64
	S    int32
	Name string
	Wrap []byte
}

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
func (gcd *GatewayConnectDeleate) Connected(context actor.Context, nets *implement.NetConnectService) error {
	return nil
}

//Analysis Packet decomposition
func (gcd *GatewayConnectDeleate) Analysis(context actor.Context, nets *implement.NetConnectService) error {
	name, h, wrap, err := agreement.AgentParser(agreement.ConstInParser).Analysis(nets.Handle.GetRecvBuffer())
	if err != nil {
		return err
	}

	if wrap == nil {
		return implement.ErrAnalysisProceed
	}

	var fpid *actor.PID
	msgType := proto.MessageType(name)
	if msgType != nil {
		if f := nets.GetMethod(msgType); f != nil {
			f(context, &ForwardPacket{H: h, Name: name, Wrap: wrap})
			goto end
		}
	}

	fpid = elements.SSets.Sreach(constant.ConstForwardServiceName)
	if fpid == nil {
		return forward.ErrForwardServiceNotStarted
	}

	actor.DefaultSchedulerContext.Send(fpid,
		&agreement.ForwardClientEvent{Handle: h,
			PactunName: name,
			Data:       wrap})

end:
	//return name, data, err
	return implement.ErrAnalysisSuccess
}

//GatewayConnect Gateway connector
type GatewayConnect struct {
	implement.NetConnectService
	GatewayID int32
}

//Init Initialize connector
func (gconn *GatewayConnect) Init() {
	gconn.NetConnectService.Init()
	gconn.RegisterMethod(&pactum.HandshakeResponse{}, gconn.onNetHandshake)
	gconn.RegisterMethod(&pactum.GatewayRegisterResponse{}, gconn.onNetRegisterResponse)
}

func (gconn *GatewayConnect) onNetHandshake(context actor.Context, message interface{}) {
	//Internal communication does not consider encrypted communication
	request := pactum.GatewayRegisterRequest{}
	request.Id = gconn.GatewayID
	var requestData []byte
	var err error

	if gconn.Target.GetEtat() != implement.Connecting {
		gconn.LogError("onNetHandshake: handshake fail: current status %+v,%+v", gconn.Target.GetEtat(), implement.Connecting)
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

func (gconn *GatewayConnect) onNetRegisterResponse(context actor.Context, message interface{}) {
	response := message.(*ForwardPacket)
	gconn.LogDebug("onNetRegisterResponse: remote handle:%+v %s", response.H, response.Name)
	responseMsg := pactum.GatewayRegisterResponse{}
	err := proto.Unmarshal(response.Wrap, &responseMsg)
	now := timer.Now()
	if err != nil {
		gconn.LogError("onNetRegisterResponse: unmarshal fail:%+v", err)
		goto unend
	}

	if gconn.Target.GetEtat() != implement.Verify {
		gconn.LogError("onNetRegisterResponse: register fail: current status %+v,%+v", gconn.Target.GetEtat(), implement.Verify)
		return
	}

	gconn.Handle.SetAuth(now)
	gconn.Target.SetEtat(implement.Connected)

	gconn.LogInfo("onNetRegisterResponse: connected address:%s time:%+v success ", gconn.Target.GetAddr(), now)
	return
unend:
	gconn.Target.SetEtat(implement.UnConnected)
}
