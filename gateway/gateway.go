package gateway

import (
	"flag"
	"fmt"

	"github.com/yamakiller/magicNet/core"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
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

	logger.Info(0, "Gateway Start Connect Service")
	gw.cps.init(gw.id)
	//注册协议及路由
	//gw.rtb.register(xxxx, "login/service")
	//1.连接其它逻辑服务器
	//1-1.后面处理
	//2.启动网络服务
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
	//1.关闭所有的连接

	gw.nsrv.Shutdown()
	gw.dsrv.CloseService()
}

// VarValue : Command 变量绑定
func (gw *Gateway) VarValue() {
	gw.dcmd.VarValue()
	flag.StringVar(&gw.addr, "p", "0.0.0.0:7850", "gateway addr")
}

// LineOption :
func (gw *Gateway) LineOption() {
	gw.dcmd.LineOption()
}

/*func (gw *Gateway) appendSocket(c *client) {

}

func (gw *Gateway) removeSocket(sock int32) *client {
	var tmpID = ((uint64(gw.id) << 32) & uint64(gw.id))
	gw.csl.Lock()
	c := gw.csocks[tmpID]
	if c != nil {
		if c.playID > 0 && gw.cplays[c.playID] != nil {
			delete(gw.cplays, c.playID)
		}
		delete(gw.csocks, tmpID)
	}
	gw.csl.Unlock()
	return c
}*/

func (gw *Gateway) onAccept(self actor.Context, message interface{}) {
	accepter := message.(network.NetAccept)

	ply := gw.cps.alloc(accepter.Addr, accepter.Port)
	_, err := gw.cps.register(accepter.Handle, ply)
	if err != nil {
		//关闭连接
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
	//closer := message.(network.NetClose)

	/*c := gw.removeSocket(closer.Handle)
	if c != nil {
		logger.Trace(self.Self().GetID(),
			"close client socket:%d playID:%d %s:%d\n",
			closer.Handle,
			c.playID,
			c.addr.String(),
			c.port)

		if c.playID > 0 {
			//广播离线
		}

		clientPool.Put(c)
	}*/
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
