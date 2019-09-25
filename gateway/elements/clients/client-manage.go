package clients

import (
	"bytes"
	"sync"

	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/st/table"
	"github.com/yamakiller/magicNet/util"

	"github.com/yamakiller/game/gateway/constant"
)

var gClientPool = sync.Pool{
	New: func() interface{} {
		b := new(GClient)
		if b.GetRecvBuffer() == nil {
			b.SetRecvBuffer(bytes.NewBuffer([]byte{}))
			b.GetRecvBuffer().Grow(constant.ConstClientBufferLimit)
		} else {
			b.GetRecvBuffer().Reset()
		}
		b.SetAuth(0)
		b.SetRef(0)
		return b
	},
}

//GClientAllocer Client memory allocator
type GClientAllocer struct {
}

//New resource allocation
func (cga *GClientAllocer) New() implement.INetClient {
	return gClientPool.Get().(implement.INetClient)
}

//Delete Release resources
func (cga *GClientAllocer) Delete(p implement.INetClient) {
	gClientPool.Put(p)
}

//NewGClientManager xxx
func NewGClientManager() *GClientManager {
	return &GClientManager{NetClientManager: implement.NetClientManager{Malloc: &GClientAllocer{}},
		HashTable: table.HashTable{Mask: 0xFFFFFF, Max: constant.ConstClientMax, Comp: clientComparator}}
}

func clientComparator(a, b interface{}) int {
	c := a.(*GClient)
	if c.GetID().GetHandle() == int32(b.(uint32)) {
		return 0
	}
	return 1
}

//
//GClientManager Gateway client manager
type GClientManager struct {
	serverid int32
	table.HashTable
	implement.NetClientManager
	smp map[int32]int32
	sync.Mutex
}

//Init Initialize the gateway client manager
func (gcm *GClientManager) Init() {
	gcm.smp = make(map[int32]int32, 128)
	gcm.HashTable.Init()
}

//Association Associate current server ID
func (gcm *GClientManager) Association(id int32) {
	gcm.serverid = id
}

//Occupy xxxx
func (gcm *GClientManager) Occupy(c implement.INetClient) (uint64, error) {
	gcm.Lock()
	defer gcm.Unlock()
	key, err := gcm.Push(c)
	if err != nil {
		return 0, err
	}

	c.SetRef(2)
	h := c.GetID()
	nh := util.NetHandle{}
	nh.SetValue(h)
	nh.Generate(gcm.serverid, int32(key), nh.GetSocket())
	c.SetID(nh.GetValue())

	return nh.GetValue(), nil
}

//Grap xxx
func (gcm *GClientManager) Grap(h uint64) implement.INetClient {
	ea := util.NetHandle{}
	ea.SetValue(h)

	gcm.Lock()
	defer gcm.Unlock()
	return gcm.getClient(ea.GetHandle())
}

//GrapSocket xxx
func (gcm *GClientManager) GrapSocket(sock int32) implement.INetClient {
	gcm.Lock()
	defer gcm.Unlock()

	k, ok := gcm.smp[sock]
	if !ok {
		return nil
	}

	return gcm.getClient(k)
}

func (gcm *GClientManager) getClient(key int32) implement.INetClient {
	c := gcm.Get(uint32(key))
	if c == nil {
		return nil
	}

	c.(implement.INetClient).IncRef()
	return c.(implement.INetClient)
}

//Erase xxxx
func (gcm *GClientManager) Erase(h uint64) {
	gcm.Lock()
	defer gcm.Unlock()

	ea := util.NetHandle{}
	ea.SetValue(h)
	if ea.GetSocket() > 0 {
		if _, ok := gcm.smp[ea.GetSocket()]; ok {
			delete(gcm.smp, ea.GetSocket())
		}
	}

	c := gcm.Get(uint32(ea.GetHandle()))
	if c == nil {
		return
	}

	gcm.Remove(uint32(ea.GetHandle()))

	if c.(implement.INetClient).DecRef() <= 0 {
		gcm.Allocer().Delete(c.(implement.INetClient))
	}
}

//Release xxx
func (gcm *GClientManager) Release(net implement.INetClient) {
	gcm.Lock()
	defer gcm.Unlock()

	if net.DecRef() <= 0 {
		gcm.Allocer().Delete(net)
	}
}

//GetHandles xxx
func (gcm *GClientManager) GetHandles() []uint64 {
	gcm.Lock()
	defer gcm.Unlock()

	cs := gcm.GetValues()
	if cs == nil {
		return nil
	}

	i := 0
	result := make([]uint64, len(cs))
	for _, v := range cs {
		result[i] = v.(implement.INetClient).GetID()
		i++
	}

	return result
}
