package gateway

import (
	"errors"
	"fmt"
	"time"

	"github.com/yamakiller/magicNet/core"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/engine/util"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
)

// Gateway : 网关服务对象
type Gateway struct {
	core.DefaultStart
	core.DefaultEnv
	core.DefaultLoop
	//
	dsrv *core.DefaultService
	dcmd core.DefaultCMDLineOption
	nsrv *service.TCPService
	//
	id   int32
	addr string
	max  int
	//
	rtb routeTable
	//
	cps gatePlays
}

// InitService : 初始化网关服务
func (gw *Gateway) InitService() error {
	logger.Info(0, "Gateway Service Start")
	gw.dsrv = &core.DefaultService{}
	if err := gw.dsrv.InitService(); err != nil {
		return err
	}

	gatewayEnv := util.GetEnvMap(util.GetEnvRoot(), "gateway")
	if gatewayEnv == nil {
		return errors.New("gateway configuration information does not exist ")
	}

	gw.id = int32(util.GetEnvInt(gatewayEnv, "id", 1))
	gw.addr = util.GetEnvString(gatewayEnv, "addr", "0.0.0.0:7850")
	gw.max = util.GetEnvInt(gatewayEnv, "max", 1024)

	logger.Info(0, "Gateway Start Connect Service ID:%d", gw.id)
	gw.cps.init(gw.id)

	gw.nsrv = service.Make("Gateway/network/tcp", func() service.IService {
		srv := &service.TCPService{Addr: gw.addr,
			CCMax:    32,
			OnAccept: gw.onAccept,
			OnRecv:   gw.onRecv,
			OnClose:  gw.onClose,
		}

		srv.Init()
		//注册协议
		//srv.RegisterMethod
		return srv
	}).(*service.TCPService)

	return nil
}

// CloseService : 关闭网关服务
func (gw *Gateway) CloseService() {
	curID := gw.nsrv.ID()
	logger.Info(curID, "Start shutting down gateway services")
	hs := gw.cps.handles()
	logger.Info(curID, "Start closing the connection number of %d", len(hs))

	for gw.cps.size() > 0 {
		chk := 0
		for i := 0; i < len(hs); i++ {
			network.OperClose(hs[i].SocketID())
		}

		for {
			time.Sleep(time.Duration(500) * time.Microsecond)
			if gw.cps.size() <= 0 {
				break
			}

			logger.Info(curID, "Remaining connections to be closed number of %d", gw.cps.size())
			chk++
			if chk > 6 {
				break
			}
		}
	}

	logger.Info(curID, "Closing the connection has been completed")

	gw.nsrv.Shutdown()
	gw.dsrv.CloseService()

	logger.Info(curID, "Gateway Services closed")
}

// VarValue : Command 变量绑定
func (gw *Gateway) VarValue() {
	gw.dcmd.VarValue()
}

// LineOption :
func (gw *Gateway) LineOption() {
	gw.dcmd.LineOption()
}

func (gw *Gateway) onAccept(self actor.Context, message interface{}) {
	accepter := message.(network.NetAccept)
	if gw.cps.size()+1 > gw.max {
		network.OperClose(accepter.Handle)
		logger.Warning(self.Self().GetID(), "accept player fulled")
		return
	}

	ply := gw.cps.alloc(accepter.Addr, accepter.Port)
	_, err := gw.cps.register(accepter.Handle, ply)
	if err != nil {
		//close-socket
		network.OperClose(accepter.Handle)
		gw.cps.free(ply)
		logger.Trace(self.Self().GetID(), "accept player closed: %v, %d-%s:%d", err,
			accepter.Handle,
			accepter.Addr.String(),
			accepter.Port)
		return
	}
	gw.cps.free(ply)
	logger.Trace(self.Self().GetID(), "accept player %d-%s:%d\n", accepter.Handle, accepter.Addr.String(), accepter.Port)
}

func (gw *Gateway) onRecv(self actor.Context, message interface{}) {
	fmt.Println("onRecv...................")
}

func (gw *Gateway) onClose(self actor.Context, message interface{}) {
	closer := message.(network.NetClose)
	hid := gw.cps.tomap(closer.Handle)
	if hid == 0 {
		logger.Trace(self.Self().GetID(), "close unfind map-id socket %d", closer.Handle)
		return
	}

	closeHandle := util.NetHandle{}
	closeHandle.Generate(gw.id, 0, int32(hid), closer.Handle)

	ply := gw.cps.grap(&closeHandle)
	if ply == nil {
		logger.Trace(self.Self().GetID(), "close unfind player %d-%d-%d-%d",
			closeHandle.GatewayID(),
			closeHandle.WorldID(),
			closeHandle.HandleID(),
			closeHandle.SocketID())
		goto unline
	}

	closeHandle = ply.handle
	gw.cps.remove(&closeHandle)
	gw.cps.free(ply)

unline:
	//通知所有服务器，这个对象已下线
	//closeHandle
}

// 路由到登陆服务器
func (gw *Gateway) onRouteLogin(self actor.Context, message interface{}) {

}

func (gw *Gateway) onRouteWorld(self actor.Context, message interface{}) {

}

// 路由数据转发到客户端
func (gw *Gateway) onRouteClient(self actor.Context, message interface{}) {

}

func (gw *Gateway) onRouteLoginClient(self actor.Context, message interface{}) {

}
