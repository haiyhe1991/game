package gateway

import (
	"github.com/yamakiller/magicNet/script/stack"

	"github.com/yamakiller/magicNet/core"
)

// Gateway : 网关服务对象
type Gateway struct {
	core.DefaultStart
	core.DefaultEnv
	core.DefaultLoop
	//
	dsrv *core.DefaultService
	dcmd core.DefaultCMDLineOption
	//
	id   int32
	addr string
	max  int
	// lua state
	spt *stack.LuaStack
	//  routing table
	//rtb *elements.RouteTable
}

/*func registerRouteProto(L *mlua.State) int {
	gwPtr := L.ToLightGoStruct(L.UpvalueIndex(1))
	if gwPtr == nil {
		logger.Fatal(0, "Gateway Object Lose")
		return 0
	}

	gw := (*Gateway)(gwPtr)
	argsNum := L.GetTop()
	if argsNum < 2 {
		return L.Error("register route proto need need 2-3 parameters")
	}

	name := L.ToCheckString(1)
	route := L.ToCheckString(2)
	auth := true
	if argsNum > 2 {
		auth = L.ToBoolean(3)
	}

	msgType := proto.MessageType(name)
	if msgType == nil {
		logger.Error(0, "Gateway Registration %s routing protocol error ", name)
		return 0
	}

	gw.rtb.register(msgType, route, auth)
	return 0
}

func (gw *Gateway) doInit() error {

	gw.dsrv = &core.DefaultService{}
	gw.rtb = &routeTable{make(map[interface{}]protoRegister, 32)}
	if err := gw.dsrv.InitService(); err != nil {
		return err
	}

	gatewayEnv := util.GetEnvMap(util.GetEnvRoot(), "gateway")
	if gatewayEnv == nil {
		return errors.New("Gateway configuration information does not exist ")
	}

	gw.id = int32(util.GetEnvInt(gatewayEnv, "id", 1))
	gw.addr = util.GetEnvString(gatewayEnv, "addr", "0.0.0.0:7850")
	gw.max = util.GetEnvInt(gatewayEnv, "max", 1024)
	gw.cps.init(gw.id)

	return nil
}

func (gw *Gateway) doScript() error {
	gw.spt = stack.NewLuaStack()
	gw.spt.GetLuaState().OpenLibs()
	gw.spt.AddSreachPath("./script")
	//register the registerRouteProto function and set gw
	gw.spt.GetLuaState().PushGoClosure(registerRouteProto, uintptr(unsafe.Pointer(gw)))
	gw.spt.GetLuaState().SetGlobal("registerRouteProto")

	if _, err := gw.spt.ExecuteScriptFile("./script/gateway.lua"); err != nil {
		return err
	}

	logger.Info(0, "Gateway Start Script[lua stack]")
	return nil
}

func (gw *Gateway) doNetwork() {
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
}

// InitService : 初始化网关服务
func (gw *Gateway) InitService() error {
	logger.Info(0, "Gateway Service Start")
	if err := gw.doInit(); err != nil {
		return err
	}

	if err := gw.doScript(); err != nil {
		return err
	}

	gw.doNetwork()

	logger.Info(0, "Gateway Start Connect Service ID:%d", gw.id)

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
	data := message.(network.NetChunk)
	hid := gw.cps.tomap(data.Handle)
	if hid == 0 {
		logger.Trace(self.Self().GetID(), "recv error closed unfind map-id socket %d", data.Handle)
		network.OperClose(data.Handle)
		return
	}

	recvHandle := util.NetHandle{}
	recvHandle.Generate(gw.id, 0, int32(hid), data.Handle)

	ply := gw.cps.grap(&recvHandle)
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
		space = constPlayerBufferLimit - ply.data.Len()
		wby = len(data.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err := ply.data.Write(data.Data[pos : pos+space])
			if err != nil {
				logger.Trace(self.Self().GetID(), "recv error %s socket %d", err.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto freeplay
			}

			pos += space
			wby += space
		}

		for {
			pkName, pkData, pkErro = extAgreeAnalysis(ply.data)
			if pkErro != nil {
				logger.Error(self.Self().GetID(), "recv error %s socket %d closing play", pkErro.Error(), data.Handle)
				network.OperClose(data.Handle)
				goto freeplay
			}

			if pkData != nil {
				pkErro = gw.onRoute(ply, pkName, pkData)
				if pkErro != nil {
					logger.Error(self.Self().GetID(), "route error %s socket %d closing play", pkErro.Error(), data.Handle)
					network.OperClose(data.Handle)
					goto freeplay
				}

				continue
			}

			if wby == 0 {
				goto freeplay
			}
		}
	}
freeplay:
	gw.cps.free(ply)
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
	gw.onPushOFFLine(&closeHandle)
}

//
func (gw *Gateway) onRoute(ply *player, name string, data []byte) error {
	msgType := proto.MessageType(name)
	if msgType == nil {
		logger.Error(gw.nsrv.ID(), "route error %s", errRouteAgreeUnDefined.Error())
		return errRouteAgreeUnDefined
	}

	p := gw.rtb.get(msgType)
	if p == nil {
		logger.Error(gw.nsrv.ID(), "route error %s", errRouteAgreeUnRegister.Error())
		return errRouteAgreeUnRegister
	}

	if p.auth && ply.auth == 0 {
		logger.Error(gw.nsrv.ID(), "route error Protocol needs to be verified, this connection is not verified and not verified")
		return errRoutePlayerUnverified
	}

	return nil
}

// 路由到登陆服务器
func (gw *Gateway) onRouteLogin(message interface{}) {

}

func (gw *Gateway) onRouteWorld(message interface{}) {

}

// 路由数据转发到客户端
func (gw *Gateway) onRouteClient(message interface{}) {

}

func (gw *Gateway) onRouteLoginClient(message interface{}) {

}

//
func (gw *Gateway) onPushOFFLine(handle *util.NetHandle) {

}*/
