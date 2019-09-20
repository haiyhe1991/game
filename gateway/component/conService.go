package component

import (
	"strings"
	"sync"
	"time"

	"github.com/yamakiller/magicNet/timer"

	"github.com/gogo/protobuf/proto"
	pkg "github.com/yamakiller/game/proto"

	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/gateway/constant"
	"github.com/yamakiller/game/gateway/elements"
	"github.com/yamakiller/game/gateway/elements/servers"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/util"
)

//NewTCPConService Create a connection service
func NewTCPConService() *ConService {
	return service.Make(constant.ConstConnectServiceName, func() service.IService {
		handle := &ConService{}

		handle.Init()
		return handle
	}).(*ConService)
}

//FromServiceMessage Message from the server
type FromServiceMessage struct {
	Handle uint64
	Sock   int32
	Name   string
	Data   interface{}
}

//ConService Provide connection service
type ConService struct {
	service.Service
	stopTicker     chan bool
	workTicker     sync.WaitGroup
	autoConnecting bool
	isShutdown     bool
}

//Init Initialize connection service
func (cse *ConService) Init() {
	cse.isShutdown = true
	cse.autoConnecting = false
	cse.Service.Init()
	cse.RegisterMethod(&actor.Started{}, cse.Started)
	cse.RegisterMethod(&actor.Stopped{}, cse.Stoped)
	cse.RegisterMethod(&network.NetChunk{}, cse.onRecv)
	cse.RegisterMethod(&network.NetClose{}, cse.onClose)
	cse.RegisterMethod(&agreement.ForwardMessage{}, cse.onForwardService)
	cse.RegisterMethod(&agreement.CheckConnectMessage{}, cse.onCheckConnect)
	cse.RegisterMethod(&pkg.HandshakeResponse{}, cse.onHandshake)
	cse.RegisterMethod(&pkg.GatewayRegisterResponse{}, cse.onRegisterResponse)
	cse.RegisterMethod(&pkg.LoginResponse{}, cse.onLoginResponse)

}

//Started Start connecting service
func (cse *ConService) Started(context actor.Context, message interface{}) {
	logger.Info(context.Self().GetID(), "Network Connect [TCP/IP] Service Startup")
	constant.ConnectServicePID = *context.Self()
	cse.workTicker.Add(1)
	cse.stopTicker = make(chan bool)
	go func(t *time.Ticker) {
		defer t.Stop()
		defer cse.workTicker.Done()
		for {
			select {
			case <-t.C:
				if !cse.autoConnecting {
					actor.DefaultSchedulerContext.Send(&constant.ConnectServicePID,
						&agreement.CheckConnectMessage{})
				}
			case stop := <-cse.stopTicker:
				if stop {
					return
				}
			}
		}
	}(time.NewTicker(time.Millisecond * time.Duration(constant.ConstConnectAutoTick)))

	cse.Service.Started(context, message)

	logger.Info(context.Self().GetID(), "Network Connect [TCP/IP] Service Startup completed")
}

//Stoped Start connecting service
func (cse *ConService) Stoped(context actor.Context, message interface{}) {
	cse.isShutdown = true
	cse.stopTicker <- true
	cse.workTicker.Wait()
	close(cse.stopTicker)
}

// Shutdown TCP network connection service termination
func (cse *ConService) Shutdown() {
	cse.Service.Shutdown()
}

//onForward Push data to target service
func (cse *ConService) onForwardService(context actor.Context, message interface{}) {

	msg := message.(*agreement.ForwardMessage)
	grp := elements.Conns.GetGroup(msg.ServerName)
	if grp == nil {
		logger.Error(context.Self().GetID(), "Forward message error No corresponding service connection found")
		return
	}

	handle := util.NetHandle{}
	handle.SetValue(msg.Handle)
	worldID := handle.WorldID()

	conn, connerr := grp.HashConnection(msg.ServerName, worldID)
	if connerr != nil {
		logger.Error(context.Self().GetID(), "Forward message error Target service does not exist [%s]", msg.AgreementName)
		return
	}

	var (
		sock int32
		err  error
	)

	forwardData := agreement.AgentParser(0).Assemble(1, handle.GetValue(), msg.AgreementName, msg.Data, int32(len(msg.Data)))
	if forwardData == nil {
		logger.Error(context.Self().GetID(), "Forward message error Failed to assemble internal data packets")
		return
	}

	ick := 0
	for {
		if conn.GetSocket() == 0 {
			//Auto Connect
			err := servers.AutoConnect(context, conn)

			if err != nil {
				logger.Error(context.Self().GetID(), "Automatic reconnection failed:%s,[%s]", err.Error(), msg.ServerName)
				goto loop_slp
			}
		}
		sock = conn.GetSocket()

		err = network.OperWrite(sock, forwardData, len(forwardData))
		if err != nil {
			logger.Error(context.Self().GetID(),
				"Forward message error write fail %s %d-%s[%s]",
				err.Error(), conn.GetID(), msg.ServerName, msg.AgreementName)
		}

		break
	loop_slp:
		ick++
		if ick > constant.ConstConnectForwardErrMax {
			logger.Error(context.Self().GetID(), "Forward message error Not connected to the target service [%s]", msg.AgreementName)
			break
		}
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}

//onCheckConnect Detect all connection status and automatically connect to the service
func (cse *ConService) onCheckConnect(context actor.Context, message interface{}) {
	cse.autoConnecting = true
	cse.workTicker.Add(1)
	defer cse.workTicker.Done()
	defer cse.restAutoConnection()

	if cse.isShutdown {
		return
	}
	elements.Conns.CheckConnect(context)
}

func (cse *ConService) onRecv(self actor.Context, message interface{}) {
	data := message.(*network.NetChunk)
	csrv := elements.Conns.GetHandle(data.Handle)
	if csrv == nil {
		logger.Error(self.Self().GetID(), "Receive data error did not find service connectio")
		return
	}

	var (
		space  int
		writed int
		wby    int
		pos    int

		pkName   string
		pkHandle uint64
		pkData   []byte
		pkErro   error
	)

	for {
		space = constant.ConstPlayerBufferLimit - csrv.GetData().Len()
		wby = len(data.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err := csrv.GetData().Write(data.Data[pos : pos+space])
			if err != nil {
				logger.Trace(self.Self().GetID(), "recv error %s socket %d", err.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto releaseconnect
			}

			pos += space
			wby += space
		}

		for {
			pkName, pkHandle, pkData, pkErro = csrv.Analysis()
			if pkErro != nil {
				logger.Error(self.Self().GetID(), "recv error %s socket %d closing play", pkErro.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto releaseconnect
			}

			if pkData != nil {
				if msgType := proto.MessageType(pkName); msgType != nil {
					if f := cse.GetMethod(msgType); f != nil {
						f(self, &FromServiceMessage{Handle: pkHandle, Sock: data.Handle, Name: pkName, Data: pkData})
						continue
					}
				}

				cse.onForwardClient(self, pkHandle, pkName, pkData)
				continue
			}

			if wby == 0 {
				goto releaseconnect
			}
		}
	}
releaseconnect:
}

func (cse *ConService) onClose(context actor.Context, message interface{}) {
	closer := message.(network.NetClose)
	conn := elements.Conns.GetHandle(closer.Handle)
	if conn == nil {
		logger.Error(context.Self().GetID(), "Close service connection error, no related connection found")
		return
	}

	conn.SetSocket(0)
	conn.SetAuth(0)
	conn.Clear()
}

func (cse *ConService) onHandshake(context actor.Context, message interface{}) {
	msg := message.(*FromServiceMessage)
	h := util.NetHandle{}
	h.SetValue(msg.Handle)

	conn := elements.Conns.GetHandle(msg.Sock)

	if conn == nil {
		logger.Error(context.Self().GetID(), "Handshake message did not find the target connection[socket:%d][%d][%s]", msg.Sock, conn.GetID(), conn.GetAddr())
		return
	}

	logger.Debug(context.Self().GetID(), "Handshake message process target connection[socket:%d][%d][%s]", msg.Sock, conn.GetID(), conn.GetAddr())

	pk := pkg.GatewayRegisterRequest{}
	pk.Id = constant.GatewayID

	regData, err := proto.Marshal(&pk)
	if err != nil {
		logger.Error(context.Self().GetID(), "Registration package build failed:%s,[socket:%d][%d][%s]", err.Error(), msg.Sock, conn.GetID(), conn.GetAddr())
		return
	}

	regData = agreement.AgentParser(agreement.ConstInParser).Assemble(agreement.ConstArgeeVersion, h.GetValue(),
		"proto.GatewayRegisterRequest",
		regData,
		int32(len(regData)))

	err = network.OperWrite(conn.GetSocket(), regData, len(regData))
	if err != nil {
		logger.Error(context.Self().GetID(), "Registration package send failed:%s,[socket:%d][%d][%s]", err.Error(), msg.Sock, conn.GetID(), conn.GetAddr())
		return
	}

	logger.Debug(context.Self().GetID(), "Registration package send success[socket:%d][%d][%s]", msg.Sock, conn.GetID(), conn.GetAddr())
}

func (cse *ConService) onRegisterResponse(context actor.Context, message interface{}) {
	msg := message.(*FromServiceMessage)
	h := util.NetHandle{}
	h.SetValue(msg.Handle)

	conn := elements.Conns.GetHandle(msg.Sock)

	if conn == nil {
		logger.Error(context.Self().GetID(), "Register Response Fail did not find the target connection [socket:%d]", msg.Sock)
		return
	}

	var response pkg.GatewayRegisterResponse
	err := proto.Unmarshal(msg.Data.([]byte), &response)

	if err != nil {
		logger.Error(context.Self().GetID(), "Register Response confirmation failed, resolution protocol failed %+v", err)
		return
	}

	if response.Code != 0 {
		logger.Error(context.Self().GetID(), "Register Response confirmation failed, %+v", err)
		return
	}

	now := timer.Now()
	conn.SetAuth(now)

	logger.Debug(context.Self().GetID(), "Register Response confirmation success %+v, %s", conn.GetID(), conn.GetAddr())
}

func (cse *ConService) onLoginResponse(context actor.Context, message interface{}) {
	msg := message.(*FromServiceMessage)

	h := util.NetHandle{}
	h.SetValue(msg.Handle)

	tmpData := msg.Data.([]byte)

	forwardData := agreement.AgentParser(agreement.ConstExParser).Assemble(agreement.ConstArgeeVersion, 0, msg.Name, tmpData, int32(len(tmpData)))
	if forwardData == nil {
		logger.Error(context.Self().GetID(), "Login Response Forward data error %s protocol data packaging failed", msg.Name)
		return
	}

	network.OperWrite(h.SocketID(), forwardData, len(forwardData))

	var response pkg.LoginResponse
	err := proto.Unmarshal(msg.Data.([]byte), &response)
	if err != nil {
		logger.Error(context.Self().GetID(), "Authentication confirmation failed, resolution protocol failed %+v", err)
		return
	}

	if response.Rep.State != 0 {
		logger.Trace(context.Self().GetID(), "Authentication failed, the connection will be closed:%d", response.GetHandle())
		network.OperClose(h.SocketID())
		return
	}

	actor.DefaultSchedulerContext.Send(&constant.NetworkServicePID,
		&agreement.CertificationConfirmation{Handle: h.GetValue()})
}

func (cse *ConService) onForwardClient(context actor.Context, handle uint64, agreementName string, data []byte) {
	h := util.NetHandle{}
	h.SetValue(handle)

	msgType := proto.MessageType(agreementName)
	if msgType == nil {
		logger.Error(context.Self().GetID(), "Forward data error %s protocol does not exist", agreementName)
		return
	}

	re := elements.RouteAddress.Sreach(msgType)
	if re == nil {
		logger.Error(context.Self().GetID(), "Forward data error %s protocol did not find the corresponding routing relationship", agreementName)
		return
	}
	//======================================================================================================================================
	if re.Auth {
		client := elements.Clients.Grap(h.HandleID())
		if client == nil {
			logger.Error(context.Self().GetID(), "Forward data error %s Target player does not exist", agreementName)
			return
		}

		if client.GetAuth() == 0 {
			logger.Error(context.Self().GetID(), "Forward data error %s Target player need to login authentication", agreementName)
			elements.Clients.Release(client)
			return
		}

		client.GetStat().UpdateWrite(timer.Now(), uint64(len(data)))
		elements.Clients.Release(client)
	}
	//=======================================================================================================================================
	if !(strings.ToLower(re.ServiceName) == "client") {
		logger.Error(context.Self().GetID(), "Forward data error %s protocol no forwarding to client permissions", agreementName)
		return
	}

	forwardData := agreement.AgentParser(agreement.ConstExParser).Assemble(1, 0, agreementName, data, int32(len(data)))
	if forwardData == nil {
		logger.Error(context.Self().GetID(), "Forward data error %s protocol data packaging failed", agreementName)
		return
	}

	network.OperWrite(h.SocketID(), forwardData, len(forwardData))
}

func (cse *ConService) restAutoConnection() {
	cse.autoConnecting = false
}
