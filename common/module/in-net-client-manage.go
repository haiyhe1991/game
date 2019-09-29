package module

import (
	"sync"

	"github.com/yamakiller/magicNet/service/implement"
)

//InNetClientAllocer Intranet client allocator base class
type InNetClientAllocer struct {
}

//Delete Release resources
func (inca *InNetClientAllocer) Delete(p implement.INetClient) {
	p.Shutdown()
}

//InNetClientManage Internal network client manager
type InNetClientManage struct {
	implement.NetClientManager
	sz    int
	iMaps map[int32]implement.INetClient
	sMaps map[int32]implement.INetClient
	sync  sync.Mutex
}

//Init Internal network client manager initialization
func (icm *InNetClientManage) Init() {
	icm.iMaps = make(map[int32]implement.INetClient)
	icm.sMaps = make(map[int32]implement.INetClient)
}

//Size Returns client is number
func (icm *InNetClientManage) Size() int {
	return icm.sz
}

//Register Register the connector ID and put back the conflicting old object ID
func (icm *InNetClientManage) Register(sock int32, handle int32) implement.INetClient {
	var result implement.INetClient
	icm.sync.Lock()
	defer icm.sync.Unlock()

	c, ok := icm.sMaps[sock]
	if !ok {
		return result
	}

	c, ok = icm.iMaps[handle]
	if ok && c.GetSocket() != sock {
		result = c
	}

	icm.iMaps[handle] = c
	return result
}

//Occupy Register or occupy a client resource
func (icm *InNetClientManage) Occupy(c implement.INetClient) (uint64, error) {
	icm.sync.Lock()
	defer icm.sync.Unlock()

	icm.sMaps[c.GetSocket()] = c
	c.SetRef(2)
	icm.sz++

	return 0, nil
}

//Grap Return a client object based on the ID, and add a reference counter
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

//GrapSocket Return a client object according to SOCKET, and add a reference counter
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

//GetHandles Return to the client Socket group
func (icm *InNetClientManage) GetHandles() []uint64 {
	icm.sync.Lock()
	defer icm.sync.Unlock()
	if icm.sz == 0 {
		return nil
	}

	ick := 0
	result := make([]uint64, icm.sz)
	for k := range icm.sMaps {
		result[ick] = uint64(k)
		ick++
	}
	return result
}

//Erase xxx
func (icm *InNetClientManage) Erase(h uint64) {
	icm.sync.Lock()

	hid := int32(h & 0xFFFFFFFF)
	sid := int32((h >> 32) & 0xFFFFFFFF)

	c, ok := icm.iMaps[hid]
	if ok {
		delete(icm.iMaps, hid)
	}

	c, ok = icm.sMaps[sid]
	if ok {
		delete(icm.sMaps, sid)
	}

	icm.sz--
	if c.DecRef() <= 0 {
		icm.sync.Unlock()
		icm.Allocer().Delete(c)
	} else {
		icm.sync.Unlock()
	}
}

//Release xxx
func (icm *InNetClientManage) Release(net implement.INetClient) {
	icm.sync.Lock()

	if net.DecRef() <= 0 {
		icm.sync.Unlock()
		icm.Allocer().Delete(net)
		return
	}
	icm.sync.Unlock()
}
