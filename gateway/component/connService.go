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

//NewConnService Create a connection service
func NewConnService() *ConnectService {
	return service.Make(constant.ConstConnectServiceName, func() service.IService {
		handle := &ConnectService{}

		handle.Init()
		return handle
	}).(*ConnectService)
}

//ConnectService Provide connection service
type ConnectService struct {
	service.Service
	stopTicker     chan bool
	workTicker     sync.WaitGroup
	autoConnecting bool
	isShutdown     bool
}

//Init Initialize connection service
func (cse *ConnectService) Init() {
	cse.isShutdown = true
	cse.autoConnecting = false
	cse.Service.Init()
	cse.RegisterMethod(&actor.Started{}, cse.Started)
	cse.RegisterMethod(&actor.Stopped{}, cse.Stoped)
	cse.RegisterMethod(&network.NetChunk{}, cse.onRecv)
	cse.RegisterMethod(&network.NetClose{}, cse.onClose)
	cse.RegisterMethod(&agreement.ForwardMessage{}, cse.onForwardService)
	cse.RegisterMethod(&agreement.CheckConnectMessage{}, cse.onCheckConnect)
}

//Started Start connecting service
func (cse *ConnectService) Started(context actor.Context, message interface{}) {
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
func (cse *ConnectService) Stoped(context actor.Context, message interface{}) {
	cse.isShutdown = true
	cse.stopTicker <- true
	cse.workTicker.Wait()
	close(cse.stopTicker)
}

// Shutdown TCP network connection service termination
func (cse *ConnectService) Shutdown() {
	cse.Service.Shutdown()
}

//onForward Push data to target service
func (cse *ConnectService) onForwardService(context actor.Context, message interface{}) {
	msg := message.(*agreement.ForwardMessage)
	grp := elements.Conns.GetGroup(msg.ServerName)
	if grp == nil {
		logger.Error(context.Self().GetID(), "Forward message error No corresponding service connection found")
		return
	}

	handle := util.NetHandle{}
	handle.SetValue(msg.Handle)
	worldID := handle.WorldID()
	conn := grp.HashConnection(worldID)
	if conn.ID == 0 {
		logger.Error(context.Self().GetID(), "Forward message error Target service does not exist [%s]", msg.AgreementName)
		return
	}

	var (
		sock int32
		err  error
	)

	forwardData := agreement.InsideAssemble(1, handle.GetValue(), msg.AgreementName, msg.Data, int32(len(msg.Data)))
	if forwardData == nil {
		logger.Error(context.Self().GetID(), "Forward message error Failed to assemble internal data packets")
		return
	}

	ick := 0
	for {
		if conn.Sock == 0 {
			//Auto Connect
			err := servers.AutoConnect(context, conn)

			if err != nil {
				goto loop_slp
			}
		}
		sock = conn.Sock

		err = network.OperWrite(sock, forwardData, len(forwardData))
		if err != nil {
			logger.Error(context.Self().GetID(),
				"Forward message error write fail %s %d-%s[%s]",
				err.Error(), conn.ID, msg.ServerName, msg.AgreementName)
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
func (cse *ConnectService) onCheckConnect(context actor.Context, message interface{}) {
	cse.autoConnecting = true
	cse.workTicker.Add(1)
	defer cse.workTicker.Done()
	defer cse.restAutoConnection()

	if cse.isShutdown {
		return
	}
	elements.Conns.CheckConnect(context)
}

func (cse *ConnectService) onRecv(self actor.Context, message interface{}) {
	data := message.(network.NetChunk)
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
		space = constant.ConstPlayerBufferLimit - csrv.DataLen()
		wby = len(data.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err := csrv.DataWrite(data.Data[pos : pos+space])
			if err != nil {
				logger.Trace(self.Self().GetID(), "recv error %s socket %d", err.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto releaseconnect
			}

			pos += space
			wby += space
		}

		for {
			pkName, pkHandle, pkData, pkErro = csrv.DataAnalysis()
			if pkErro != nil {
				logger.Error(self.Self().GetID(), "recv error %s socket %d closing play", pkErro.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto releaseconnect
			}

			if pkData != nil {
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

func (cse *ConnectService) onClose(context actor.Context, message interface{}) {
	closer := message.(network.NetClose)
	conn := elements.Conns.GetHandle(closer.Handle)
	if conn == nil {
		logger.Error(context.Self().GetID(), "Close service connection error, no related connection found")
		return
	}

	conn.Sock = 0
	conn.ClearData()
}

func (cse *ConnectService) onForwardClient(context actor.Context, handle uint64, agreementName string, data []byte) {
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
		client := elements.Clients.Grap(&h)
		if client == nil {
			logger.Error(context.Self().GetID(), "Forward data error %s Target player does not exist", agreementName)
			return
		}

		if client.Auth == 0 {
			logger.Error(context.Self().GetID(), "Forward data error %s Target player need to login authentication", agreementName)
			elements.Clients.Release(client)
			return
		}

		client.Stat.UpdateWrite(timer.Now(), uint64(len(data)))
		elements.Clients.Release(client)
	}
	//=======================================================================================================================================
	if !(strings.ToLower(re.ServiceName) == "client") {
		logger.Error(context.Self().GetID(), "Forward data error %s protocol no forwarding to client permissions", agreementName)
		return
	}

	forwardData := agreement.ExtAssemble(1, agreementName, data, int32(len(data)))
	if forwardData == nil {
		logger.Error(context.Self().GetID(), "Forward data error %s protocol data packaging failed", agreementName)
		return
	}

	network.OperWrite(h.SocketID(), forwardData, len(forwardData))

	if !re.Confirm {
		return
	}
	var response pkg.LoginResponse
	err := proto.Unmarshal(data, &response)
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

func (cse *ConnectService) restAutoConnection() {
	cse.autoConnecting = false
}
