package clients

import (
	"bytes"
	"net"
	"sync"

	"github.com/yamakiller/game/gateway/elements"
	"github.com/yamakiller/magicNet/engine/util"
)

var playerPool = sync.Pool{
	New: func() interface{} {
		b := new(Player)
		if b.data == nil {
			b.data = bytes.NewBuffer([]byte{})
			b.data.Grow(elements.ConstPlayerBufferLimit)
		} else {
			b.data.Reset()
		}
		b.ref = 0
		return b
	},
}

//PlayManager play Manager
type PlayManager struct {
	d     int32
	s     []*Player
	mps   map[int32]uint16
	sz    int
	seqID uint16
	sync  sync.Mutex
}

// Initial Initialize Player management module
func (pms *PlayManager) Initial(id int32) {

	pms.d = id
	pms.s = make([]*Player, elements.ConstPlayerMax)
	pms.mps = make(map[int32]uint16, 64)
}

// Occupy Assign and jion plays a player object
func (pms *PlayManager) Occupy(sock int32, addr net.IP, port int) (*Player, util.NetHandle, error) {
	pe := playerPool.Get().(*Player)
	pe.addr = addr
	pe.port = port

	var i uint16
	pms.sync.Lock()
	for i = 0; i < elements.ConstPlayerMax; i++ {
		key := ((i + pms.seqID) & elements.ConstPlayerIDMask)
		hash := key & (elements.ConstPlayerMax - 1)
		if pms.s[hash] == nil {
			handle := util.NetHandle{}
			handle.Generate(pms.d, 0, int32(key), sock)
			pe.Handle = handle
			pms.seqID = key + 1
			pms.s[hash] = pe
			pms.s[hash].ref = 2
			pms.mps[sock] = key
			pms.sz++
			pms.sync.Unlock()
			return pe, handle, nil
		}
	}

	pms.sync.Unlock()
	playerPool.Put(pe)

	return nil, util.NetHandle{}, errPlayerFull
}

// Grap Pick up a player object and take it
func (pms *PlayManager) Grap(h *util.NetHandle) *Player {
	pms.sync.Lock()
	hash := uint32(h.HandleID()) & uint32(elements.ConstPlayerMax-1)
	if pms.s[hash] != nil && pms.s[hash].Handle.HandleID() == h.HandleID() {
		pe := pms.s[hash]
		pe.ref++
		pms.sync.Unlock()
		return pe
	}
	pms.sync.Unlock()
	return nil
}

// Erase removes the Player from PlayManager
func (pms *PlayManager) Erase(h *util.NetHandle) {
	pms.sync.Lock()
	hash := uint32(h.HandleID()) & uint32(elements.ConstPlayerMax-1)
	if pms.s[hash] != nil && pms.s[hash].Handle.HandleID() == h.HandleID() {
		pe := pms.s[hash]
		pms.s[hash] = nil
		if _, ok := pms.mps[pe.Handle.SocketID()]; ok {
			delete(pms.mps, pe.Handle.SocketID())
		}
		pms.sz--
		pe.ref--
		if pe.ref <= 0 {
			playerPool.Put(pe)
		}

	}
	pms.sync.Unlock()
}

// Release  Release control
func (pms *PlayManager) Release(pe *Player) {
	pms.sync.Lock()
	pe.ref--
	if pe.ref <= 0 {
		playerPool.Put(pe)
	}
	pms.sync.Unlock()
}

// ToHandleID Socket conversion to the corresponding handle id
func (pms *PlayManager) ToHandleID(sock int32) uint16 {
	pms.sync.Lock()
	if v, ok := pms.mps[sock]; ok {
		pms.sync.Unlock()
		return v
	}
	pms.sync.Unlock()
	return 0
}

// Size the PlayManager of number
func (pms *PlayManager) Size() int {
	return pms.sz
}

// GetHandls Get all the player object handles
func (pms *PlayManager) GetHandls() []util.NetHandle {

	pms.sync.Lock()
	i := 0
	hs := make([]util.NetHandle, pms.sz)
	for _, pe := range pms.s {
		if pe == nil {
			continue
		}

		hs[i] = pe.Handle
		i++
	}
	pms.sync.Unlock()
	return hs
}
