package source

import (
	"flag"
	"fmt"

	"github.com/yamakiller/magicNet/core"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
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
	addr string
}

// InitService : 初始化网关服务
func (gw *Gateway) InitService() error {
	logger.Info(0, "Gateway Service Start")
	gw.dsrv = &core.DefaultService{}
	if err := gw.dsrv.InitService(); err != nil {
		return err
	}

	logger.Info(0, "Gateway Start Connect Service")
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
func (gw *Gateway) onAccept(self actor.Context, message interface{}) {
	fmt.Println("accept.....")
}

func (gw *Gateway) onRecv(self actor.Context, message interface{}) {
	fmt.Println("recv.....")
}

func (gw *Gateway) onClose(self actor.Context, message interface{}) {
	fmt.Println("close.....")
}
