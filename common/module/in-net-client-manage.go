package module

import (
	"sync"

	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/util"
)

//InNetClientAllocer Intranet client allocator base class
type InNetClientAllocer struct {
}

//Delete Release resources
func (inca *InNetClientAllocer) Delete(p implement.INetClient) {
	p.Shutdown()
}

//
type InNetClientManage struct {
	implement.NetClientManager
	sz    int
	iMaps map[int32]implement.INetClient
	sMaps map[int32]implement.INetClient
	sync  sync.Mutex
}

func (icm *InNetClientManage) Size() int {
	return icm.sz
}

func (icm *InNetClientManage) Register(sock int32, handle int32) {
	icm.sync.Lock()
	defer icm.sync.Unlock()

	c, ok := icm.sMaps[sock]
	if !ok {
		return
	}

	icm.iMaps[handle] = c
}

func (icm *InNetClientManage) Occupy(c implement.INetClient) (uint64, error) {
	icm.sync.Lock()
	defer icm.sync.Unlock()

	icm.sMaps[c.GetSocket()] = c
	c.SetRef(2)
	icm.sz++

	return 0, nil
}

func (icm *InNetClientManage) Grap(h uint64) implement.INetClient {
	icm.sync.Lock()
	defer icm.sync.Unlock()

	c, ok := icm.iMaps[int32(h)]
	if !ok {
		return nil
	}

	c.IncRef()
	return c
}

func (icm *InNetClientManage) GrapSocket(sock int32) implement.INetClient {
	icm.sync.Lock()
	defer icm.sync.Unlock()

	c, ok := icm.sMaps[sock]
	if !ok {
		return nil
	}
	c.IncRef()
	return c
}

func (icm *InNetClientManage) GetHandles() []uint64 {
	icm.sync.Lock()
	defer icm.sync.Unlock()
	if icm.sz == 0 {
		return nil
	}

	ick := 0
	result := make([]uint64, icm.sz)
	for k := range icm.iMaps {
		result[ick] = uint64(k)
		ick++
	}
	return result
}

func (icm *InNetClientManage) Erase(h *util.NetHandle) {
	icm.sync.Lock()
	defer icm.sync.Unlock()

	c, ok := icm.iMaps[int32(h.GetValue())]
	if !ok {
		return
	}

	delete(icm.iMaps, int32(h.GetValue()))
	icm.sz--
	if c.DecRef() <= 0 {
		icm.Allocer().Delete(c)
	}
}

//Release xxx
func (icm *InNetClientManage) Release(net implement.INetClient) {
	icm.sync.Lock()
	defer icm.sync.Unlock()

	if net.DecRef() <= 0 {
		icm.Allocer().Delete(net)
	}
}
