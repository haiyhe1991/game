package clients

import (
	"bytes"
	"errors"
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
	sync.Mutex
}

// Initial Initialize Player management module
func (cms *ClientManager) Initial(id int32) {
	cms.d = id
	cms.s = make([]*Client, constant.ConstPlayerMax)
	cms.mps = make(map[int32]uint16, 64)
	cms.seqID = 1
}

// Occupy Assign and jion plays a player object
func (cms *ClientManager) Occupy(sock int32, addr net.IP, port int) (*Client, util.NetHandle, error) {
	client := clientPool.Get().(*Client)
	client.addr = addr
	client.port = port

	var i uint16
	cms.Lock()

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
			cms.Unlock()
			return client, handle, nil
		}
	}
	cms.Unlock()
	clientPool.Put(client)

	return nil, util.NetHandle{}, errPlayerFull
}

// Grap Pick up a player object and take it
func (cms *ClientManager) Grap(h *util.NetHandle) *Client {
	cms.Lock()
	defer cms.Unlock()
	hash := uint32(h.HandleID()) & uint32(constant.ConstPlayerMax-1)
	if cms.s[hash] != nil && cms.s[hash].Handle.HandleID() == h.HandleID() {
		pe := cms.s[hash]
		pe.ref++
		return pe
	}
	return nil
}

// Erase removes the Player from PlayManager
func (cms *ClientManager) Erase(h *util.NetHandle) {
	cms.Lock()
	defer cms.Unlock()
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
}

// Release  Release control
func (cms *ClientManager) Release(client *Client) {
	cms.Lock()
	defer cms.Unlock()
	client.ref--
	if client.ref <= 0 {
		clientPool.Put(client)
	}
}

// ToHandleID Socket conversion to the corresponding handle id
func (cms *ClientManager) ToHandleID(sock int32) (uint16, error) {
	cms.Lock()
	defer cms.Unlock()
	if v, ok := cms.mps[sock]; ok {
		return v, nil
	}
	return 0, errors.New("unknown id")
}

// Size the PlayManager of number
func (cms *ClientManager) Size() int {
	return cms.sz
}

// GetHandls Get all the player object handles
func (cms *ClientManager) GetHandls() []util.NetHandle {
	cms.Lock()
	defer cms.Unlock()
	i := 0
	hs := make([]util.NetHandle, cms.sz)
	for _, client := range cms.s {
		if client == nil {
			continue
		}

		hs[i] = client.Handle
		i++
	}

	return hs
}
