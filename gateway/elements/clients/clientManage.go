package clients

import (
	"bytes"
	"net"
	"sync"

	"github.com/yamakiller/game/gateway/constant"
	"github.com/yamakiller/magicNet/util"
)

var clientPool = sync.Pool{
	New: func() interface{} {
		b := new(Client)
		if b.data == nil {
			b.data = bytes.NewBuffer([]byte{})
			b.data.Grow(constant.ConstPlayerBufferLimit)
		} else {
			b.data.Reset()
		}
		b.ref = 0
		return b
	},
}

//ClientManager play Manager
type ClientManager struct {
	d     int32
	s     []*Client
	mps   map[int32]uint16
	sz    int
	seqID uint16
	sync  sync.Mutex
}

// Initial Initialize Player management module
func (cms *ClientManager) Initial(id int32) {

	cms.d = id
	cms.s = make([]*Client, constant.ConstPlayerMax)
	cms.mps = make(map[int32]uint16, 64)
}

// Occupy Assign and jion plays a player object
func (cms *ClientManager) Occupy(sock int32, addr net.IP, port int) (*Client, util.NetHandle, error) {
	client := clientPool.Get().(*Client)
	client.addr = addr
	client.port = port

	var i uint16
	cms.sync.Lock()
	for i = 0; i < constant.ConstPlayerMax; i++ {
		key := ((i + cms.seqID) & constant.ConstPlayerIDMask)
		hash := key & (constant.ConstPlayerMax - 1)
		if cms.s[hash] == nil {
			handle := util.NetHandle{}
			handle.Generate(cms.d, 0, int32(key), sock)
			client.Handle = handle
			cms.seqID = key + 1
			cms.s[hash] = client
			cms.s[hash].ref = 2
			cms.mps[sock] = key
			cms.sz++
			cms.sync.Unlock()
			return client, handle, nil
		}
	}

	cms.sync.Unlock()
	clientPool.Put(client)

	return nil, util.NetHandle{}, errPlayerFull
}

// Grap Pick up a player object and take it
func (cms *ClientManager) Grap(h *util.NetHandle) *Client {
	cms.sync.Lock()
	hash := uint32(h.HandleID()) & uint32(constant.ConstPlayerMax-1)
	if cms.s[hash] != nil && cms.s[hash].Handle.HandleID() == h.HandleID() {
		pe := cms.s[hash]
		pe.ref++
		cms.sync.Unlock()
		return pe
	}
	cms.sync.Unlock()
	return nil
}

// Erase removes the Player from PlayManager
func (cms *ClientManager) Erase(h *util.NetHandle) {
	cms.sync.Lock()
	hash := uint32(h.HandleID()) & uint32(constant.ConstPlayerMax-1)
	if cms.s[hash] != nil && cms.s[hash].Handle.HandleID() == h.HandleID() {
		pe := cms.s[hash]
		cms.s[hash] = nil
		if _, ok := cms.mps[pe.Handle.SocketID()]; ok {
			delete(cms.mps, pe.Handle.SocketID())
		}
		cms.sz--
		pe.ref--
		if pe.ref <= 0 {
			clientPool.Put(pe)
		}

	}
	cms.sync.Unlock()
}

// Release  Release control
func (cms *ClientManager) Release(client *Client) {
	cms.sync.Lock()
	client.ref--
	if client.ref <= 0 {
		clientPool.Put(client)
	}
	cms.sync.Unlock()
}

// ToHandleID Socket conversion to the corresponding handle id
func (cms *ClientManager) ToHandleID(sock int32) uint16 {
	cms.sync.Lock()
	if v, ok := cms.mps[sock]; ok {
		cms.sync.Unlock()
		return v
	}
	cms.sync.Unlock()
	return 0
}

// Size the PlayManager of number
func (cms *ClientManager) Size() int {
	return cms.sz
}

// GetHandls Get all the player object handles
func (cms *ClientManager) GetHandls() []util.NetHandle {

	cms.sync.Lock()
	i := 0
	hs := make([]util.NetHandle, cms.sz)
	for _, client := range cms.s {
		if client == nil {
			continue
		}

		hs[i] = client.Handle
		i++
	}
	cms.sync.Unlock()
	return hs
}
