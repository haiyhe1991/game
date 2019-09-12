package gateway

import (
	"sync"
	"time"

	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/engine/util"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
)

type pushEvent struct {
	route  string
	handle uint64
	agree  string
	data   []byte
}

type autoConnEvent struct {
}

type pusher struct {
	service.Service
	srvs       *servers
	isShutdown bool
	autoWait   sync.WaitGroup
}

func (psr *pusher) Init() {
	psr.Service.Init()
	psr.RegisterMethod(&actor.Started{}, psr.Started)
	psr.RegisterMethod(&actor.Stopped{}, psr.Stoped)
	psr.RegisterMethod(&pushEvent{}, psr.onPusher)
	psr.RegisterMethod(&autoConnEvent{}, psr.onAutoConnect)
	psr.RegisterMethod(&network.NetChunk{}, psr.onRecv)
	psr.RegisterMethod(&network.NetClose{}, psr.onClose)
}

// Started network push service is enabled
func (psr *pusher) Started(context actor.Context, message interface{}) {
	logger.Info(context.Self().GetID(), "Pusher Service Start")

	psr.Service.Started(context, message)
	logger.Info(context.Self().GetID(), "Pusher Service Success")
}

// Stoped network push service stops
func (psr *pusher) Stoped(context actor.Context, message interface{}) {
	logger.Info(context.Self().GetID(), "Pusher Service Stoping")

	logger.Info(context.Self().GetID(), "Pusher Service Stoped")
}

// Shutdown network push service termination
func (psr *pusher) Shutdown() {

	psr.Service.Shutdown()
}

func (psr *pusher) InitiateAutoConnect() {
	if !psr.isShutdown {
		psr.autoWait.Wait()
		psr.autoWait.Add(1)
	}
}

// onPusher
func (psr *pusher) onPusher(context actor.Context, message interface{}) {
	evt := message.(*pushEvent)
	sv := psr.srvs.get(evt.route)
	if sv == nil {
		logger.Error(context.Self().GetID(), "Push data error No corresponding service connection found")
		return
	}

	handle := util.NetHandle{}
	handle.SetValue(evt.handle)

	worldID := handle.WorldID()

	c := sv.get(worldID)
	if c.id == 0 {
		logger.Error(context.Self().GetID(), "Push data error Target service does not exist [%s]", evt.agree)
		return
	}

	var (
		pushData []byte
		sock     int32
		err      error
	)

	ick := 0
	for {
		c.sync.Lock()
		if c.sock == 0 {
			//Auto Connect
			err := psr.autoConnect(context, c)
			c.sync.Unlock()
			if err != nil {
				goto loop_slp
			}
		}
		sock = c.sock
		c.sync.Unlock()

		err = network.OperWrite(sock, pushData, len(pushData))
		if err != nil {
			logger.Error(context.Self().GetID(), "Push data error write fail %s %d-%s[%s]", err.Error(), c.id, evt.route, evt.agree)
		}
		break
	loop_slp:
		ick++
		if ick > constConnectPushErrMax {
			logger.Error(context.Self().GetID(), "Push data error Not connected to the target service [%s]", evt.agree)
			break
		}
		time.Sleep(time.Millisecond * time.Duration(100))
	}
}

// Automatic connection automatic reconnection
func (psr *pusher) onAutoConnect(context actor.Context, message interface{}) {
	defer psr.autoWait.Done()
	if psr.isShutdown {
		return
	}

	if psr.srvs == nil {
		logger.Error(context.Self().GetID(), "No target service information is associated, no automatic connection is required")
		return
	}

	var err error
	for k, v := range psr.srvs.ms {
		for i := 0; i < constConnectMax; i++ {
			v.cs[i].sync.Lock()
			if v.cs[i].id == 0 || v.cs[i].sock > 0 {
				v.cs[i].sync.Unlock()
				continue
			}
			err = psr.autoConnect(context, &v.cs[i])
			v.cs[i].sync.Unlock()

			if err == nil {
				logger.Error(context.Self().GetID(), "Connection service failed %s %s-%d-%s", err, k, v.cs[i].id, v.cs[i].addr)
			}
		}

		psr.srvs.ms[k] = v
	}
}

func (psr *pusher) onRecv(context actor.Context, message interface{}) {

}

func (psr *pusher) onClose(context actor.Context, message interface{}) {

	closer := message.(network.NetClose)
	c := psr.srvs.getConnector(closer.Handle)
	if c == nil {
		logger.Error(context.Self().GetID(), "Close service connection error, no related connection found")
		return
	}

	c.sock = 0
	n := c.data.Len()
	if n > 0 {
		c.data.Next(n)
	}
}

func (psr *pusher) autoConnect(context actor.Context, c *connector) error {

	h, err := network.OperTCPConnect(context.Self(), c.addr, constConnectChanMax)
	if err != nil {
		return err
	}

	c.sock = h
	network.OperOpen(h)

	return nil
}
