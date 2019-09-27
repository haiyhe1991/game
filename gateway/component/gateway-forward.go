package component

import (
	"bytes"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/gateway/constant"
	"github.com/yamakiller/game/gateway/elements"
	"github.com/yamakiller/game/gateway/elements/servers"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/service/implement"
)

type checkConnectEvent struct {
}

//GatewayForward Gateway forwarding service
type GatewayForward struct {
	service.Service
	conns      []*GatewayConnect
	authTicker *time.Ticker
	authSign   chan bool
	authWait   sync.WaitGroup
	isChecking bool
	isShutdown bool
}

//Init Initialize the repeater
func (gf *GatewayForward) Init() {
	gf.Service.Init()
	gf.RegisterMethod(&actor.Started{}, gf.Started)
	gf.RegisterMethod(&actor.Stopped{}, gf.Stoped)
	gf.RegisterMethod(&checkConnectEvent{}, gf.onCheckConnect)
	gf.RegisterMethod(&agreement.ForwardClientEvent{}, gf.onForwardClient)
	gf.RegisterMethod(&agreement.ForwardServerEvent{}, gf.onForwardServer)
}

//Started Start forwarding service
func (gf *GatewayForward) Started(context actor.Context, message interface{}) {
	gf.Service.Assignment(context)
	gf.LogInfo("Service Startup %s", gf.Name())
	tset := elements.TSets.GetValues()
	if tset == nil {
		goto end
	}

	gf.LogInfo("%d connection target numbers", len(tset))
	gf.LogInfo("Startup trying to create a machine")
	for _, v := range tset {
		t := v.(*servers.TargetConnection)
		name := t.Name + "#" + strconv.Itoa(int(t.ID))
		gf.LogInfo("Start generating %s connectors address:%s", name, t.Addr)
		con := service.Make(name, func() service.IService {
			h := &GatewayConnect{GatewayID: constant.GatewayID,
				AutoErrRetry:     constant.GatewayConnectForwardErrMax,
				AutoErrRetryTime: constant.GatewayConnectForwardInterval,
				NetConnectService: implement.NetConnectService{
					Handle: &GatewayHandleConnect{RecvBufferMax: common.ConstClientBufferLimit,
						rbuffer: bytes.NewBuffer([]byte{})},
					Deleate: &GatewayConnectDeleate{},
					Target:  t}}
			h.Handle.GetRecvBuffer().Grow(h.Handle.GetRecvBufferLimit())
			h.Init()
			return h
		})

		gf.conns = append(gf.conns, con.(*GatewayConnect))
		gf.LogInfo("Start generating %s connectors address:%s complete", name, t.Addr)
	}

	//auto connect==========================================
	gf.isChecking = true
	actor.DefaultSchedulerContext.Send(gf.GetPID(),
		&checkConnectEvent{})
	//======================================================

	gf.authWait.Add(1)
	gf.authTicker = time.NewTicker(time.Duration(constant.GatewayConnectForwardAutoTick) * time.Millisecond)
	gf.authSign = make(chan bool, 1)
	go func(t *time.Ticker) {
		defer gf.authWait.Done()
		for {
			select {
			case <-t.C:
				if !gf.isChecking {
					gf.isChecking = true
					actor.DefaultSchedulerContext.Send(gf.GetPID(),
						&checkConnectEvent{})
				}
			case stop := <-gf.authSign:
				if stop {
					return
				}
			}
		}
	}(gf.authTicker)

	gf.LogInfo("%s Service Startup completed", gf.Name())
end:
	gf.Service.Started(context, message)
}

//Stoped Stop forwarding service
func (gf *GatewayForward) Stoped(context actor.Context, message interface{}) {
	n := len(gf.conns)
	gf.LogInfo("Service Stopping [connecting:%d]", n)
	for _, v := range gf.conns {
		gf.LogInfo("Connection Stopping %d name:%s address:%s", n, v.Target.GetName(), v.Target.GetAddr())
		v.Shutdown()
		gf.LogInfo("Connection Stoped name:%s address:%s", v.Target.GetName(), v.Target.GetAddr())
		n--
	}
	gf.conns = gf.conns[:0]
	gf.LogInfo("Service Stoped")
}

//Shutdown Termination of service
func (gf *GatewayForward) Shutdown() {
	gf.isShutdown = false
	gf.authSign <- true
	gf.authWait.Wait()
	close(gf.authSign)
	gf.authTicker.Stop()
	gf.Service.Shutdown()
}

func (gf *GatewayForward) restCheckStatus() {
	gf.isChecking = false
}

func (gf *GatewayForward) onCheckConnect(context actor.Context, message interface{}) {
	defer gf.restCheckStatus()
	for _, v := range gf.conns {
		//End and exit
		if gf.isShutdown {
			return
		}

		switch v.Target.GetEtat() {
		case implement.Connected:
		case implement.Connecting:
			fallthrough
		case implement.Verify:
			if outTm := v.Target.IsTimeout(); outTm > 0 {
				v.Handle.Close()
			}
		case implement.UnConnected:
			if v.GetPID() != nil {
				v.Target.SetEtat(implement.Connecting)
				actor.DefaultSchedulerContext.Send(v.GetPID(),
					&implement.NetConnectEvent{})
			}
		default:
			gf.LogDebug("Exception non-existent logic")
		}
	}
}

func (gf *GatewayForward) onForwardClient(context actor.Context, message interface{}) {
	msg := message.(*agreement.ForwardClientEvent)
	msgType := proto.MessageType(msg.PactumName)
	if msgType == nil {
		gf.LogError("The %s protocol is not defined, and the data is discarded depending on the abnormal operation.",
			msg.PactumName)
		return
	}

	adr := elements.ForwardAddresses.Sreach(msgType)
	if adr == nil || !(adr.ServoName == "client") {
		gf.LogError("Protocol not registered route or routing address error: pactum name:%s",
			msg.PactumName)
		return
	}

	pid := elements.SSets.Sreach(constant.ConstNetworkServiceName)
	if pid == nil {
		gf.LogError("Network Service Department exists")
		return
	}

	actor.DefaultSchedulerContext.Send(pid, msg)
	gf.LogDebug("The send request has been handed over to the network service module: pactum name:%s",
		msg.PactumName)
}

func (gf *GatewayForward) onForwardServer(context actor.Context, message interface{}) {
	msg := message.(*agreement.ForwardServerEvent)
	loader := elements.TLSets.Get(msg.ServoName)
	if loader == nil {
		gf.LogError("The %s target server was not found"+
			" and the data was discarded.", msg.ServoName)
		return
	}

	ick := 0
	var to *servers.TargeObject
	for {
		to = loader.GetTarget(strconv.Itoa(rand.Intn(10000)))
		if to != nil {
			break
		}

		ick++
		if ick >= 6 {
			gf.LogError("[%s]No attempts were made to find"+
				" available nodes and data was dropped", msg.ServoName)
			return
		}
	}

	actor.DefaultSchedulerContext.Send(to.Target, msg)

	gf.LogDebug("Data send request has been pushed successfully")
}
