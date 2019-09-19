package component

import (
	"time"

	"github.com/yamakiller/game/common/elements/visitors"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/network"
	net "github.com/yamakiller/magicNet/service/net"
)

//NewTCPService Create a new TCP network service object
func NewTCPService(addr string, ccmax int, cs visitors.IVisitorManager) *NetService {
	return &NetService{CS: cs, TCPService: net.TCPService{Addr: addr, CCMax: ccmax}}
}

//NetService Internet service
type NetService struct {
	CS visitors.IVisitorManager
	net.TCPService
}

//Shutdown Termination of network services
func (ns *NetService) Shutdown() {
	logger.Info(ns.ID(), "Network Listen [TCP/IP] Service Closing connection")
	hs := ns.CS.GetKeys()
	for ns.CS.Size() > 0 {
		chk := 0
		for i := 0; i < len(hs); i++ {
			c := ns.CS.Grap(hs[i])
			if c == nil {
				continue
			}
			sck := c.GetSocket()
			ns.CS.Release(c)
			network.OperClose(sck)
		}

		for {
			time.Sleep(time.Duration(500) * time.Microsecond)
			if ns.CS.Size() <= 0 {
				break
			}

			logger.Info(ns.ID(), "Network Listen [TCP/IP] Service The remaining %d connections need to be closed", ns.CS.Size())
			chk++
			if chk > 6 {
				break
			}
		}
	}

	logger.Info(ns.ID(), "Network Listen [TCP/IP] Service All connections are closed")

	ns.TCPService.Shutdown()
}
