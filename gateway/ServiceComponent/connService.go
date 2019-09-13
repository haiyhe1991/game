package ServiceComponent

import (
	"time"

	"github.com/yamakiller/game/gateway/elements"
	"github.com/yamakiller/game/gateway/elements/agreement"
	"github.com/yamakiller/game/gateway/elements/servers"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/engine/util"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
)

var (
	connectManger *servers.ConnectionManager //Connection service manager
)

func init() {
	connectManger = servers.NewManager()
}

//NewConnService
func NewConnService() *ConnectService {
	return service.Make(elements.ConstConnectServiceName, func() service.IService {
		handle := &ConnectService{}

		handle.Init()
		return handle
	}).(*ConnectService)
}

//ConnectService Provide connection service
type ConnectService struct {
	service.Service
}

//Init Initialize connection service
func (cse *ConnectService) Init() {
	cse.Service.Init()
	cse.RegisterMethod(&actor.Started{}, cse.Started)
	cse.RegisterMethod(&actor.Stopped{}, cse.Stoped)
	cse.RegisterMethod(&network.NetChunk{}, cse.onRecv)
	cse.RegisterMethod(&network.NetClose{}, cse.onClose)
	cse.RegisterMethod(&agreement.ForwardMessage{}, cse.onForward)
	cse.RegisterMethod(&agreement.CheckConnectMessage{}, cse.onCheckConnect)
}

//Started Start connecting service
func (cse *ConnectService) Started(context actor.Context, message interface{}) {
	logger.Info(context.Self().GetID(), "Network[TCP/IP] Connect Service turning on")
	cse.Service.Started(context, message)
}

//Stoped Start connecting service
func (cse *ConnectService) Stoped(context actor.Context, message interface{}) {
	//关闭所有连接
}

// Shutdown TCP network service termination
func (cse *ConnectService) Shutdown() {

	cse.Service.Shutdown()
}

//onForward Push data to target service
func (cse *ConnectService) onForward(context actor.Context, message interface{}) {
	msg := message.(*agreement.ForwardMessage)
	grp := connectManger.GetGroup(msg.ServerName)
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
		forwardData []byte
		sock        int32
		err         error
	)

	//? 组装数据包协议

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
		if ick > elements.ConstConnectForwardErrMax {
			logger.Error(context.Self().GetID(), "Forward message error Not connected to the target service [%s]", msg.AgreementName)
			break
		}
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}

//onCheckConnect Detect all connection status and automatically connect to the service
func (cse *ConnectService) onCheckConnect(context actor.Context, message interface{}) {
	connectManger.CheckConnect(context)
}

func (cse *ConnectService) onRecv(context actor.Context, message interface{}) {
	///
}

func (cse *ConnectService) onClose(context actor.Context, message interface{}) {
	closer := message.(network.NetClose)
	conn := connectManger.GetHandle(closer.Handle)
	if conn == nil {
		logger.Error(context.Self().GetID(), "Close service connection error, no related connection found")
		return
	}

	conn.Sock = 0
	conn.ClearData()
}
