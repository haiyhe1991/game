package gateway

import (
	"bytes"
	"errors"
	"net"
	"sync"

	"github.com/yamakiller/magicNet/engine/util"
)

const (
	//palyerMax play number of max
	constPalyerMax = 65535
	//playerIDMask play ID of mask
	constPlayerIDMask = 0xFF
)

var (
	errPlayerFull = errors.New("player is full")
)

type remoteInfo struct {
	addr net.IP
	port int
}

type remoteStat struct {
	olt uint64

	lrt uint64
	lwt uint64

	rb uint64
	wb uint64
}

func (st *remoteStat) updateWrite(tts uint64, bytes uint64) {
	st.lwt = tts
	st.wb += bytes
}

func (st *remoteStat) updateRecv(tts uint64, bytes uint64) {
	st.lrt = tts
	st.rb += bytes
}

func (st *remoteStat) update(tts uint64) {
	st.olt = tts
}

type player struct {
	handle util.NetHandle
	data   *bytes.Buffer
	stat   remoteStat
	addrs  remoteInfo
	ref    int
}

var playerPool = sync.Pool{
	New: func() interface{} {
		b := new(player)
		if b.data == nil {
			b.data = bytes.NewBuffer([]byte{})
		} else {
			b.data.Reset()
		}
		b.ref = 0
		return b
	},
}

type gatePlays struct {
	d     int32
	s     []*player
	seqID uint16
	sync  sync.Mutex
}

func (gpls *gatePlays) init(gatewayID int32) {
	gpls.d = gatewayID
	gpls.s = make([]*player, constPalyerMax)
}

func (gpls *gatePlays) alloc(addr net.IP, port int) *player {
	pe := playerPool.Get().(*player)
	pe.addrs.addr = addr
	pe.addrs.port = port
	return pe
}

func (gpls *gatePlays) register(sock int32, pe *player) (util.NetHandle, error) {
	var i uint16
	gpls.sync.Lock()

	for i = 0; i < constPalyerMax; i++ {
		key := ((i + gpls.seqID) & constPlayerIDMask)
		hash := key & (constPalyerMax - 1)
		if gpls.s[hash] == nil {
			handle := util.NetHandle{}
			handle.Generate(gpls.d, 0, int32(key), sock)
			pe.handle = handle
			gpls.seqID = key + 1
			gpls.s[hash] = pe
			gpls.s[hash].ref = 2
			gpls.sync.Unlock()
			return handle, nil
		}
	}
	gpls.sync.Unlock()
	return util.NetHandle{}, errPlayerFull
}

func (gpls *gatePlays) remove(handle *util.NetHandle) {
	gpls.sync.Lock()
	hash := uint32(handle.HandleID()) & uint32(constPalyerMax-1)
	if gpls.s[hash] != nil && gpls.s[hash].handle.HandleID() == handle.HandleID() {
		pe := gpls.s[hash]
		gpls.s[hash] = nil
		pe.ref--
		if pe.ref <= 0 {
			playerPool.Put(pe)
		}

	}
	gpls.sync.Unlock()
}

func (gpls *gatePlays) removeAll() {
	gpls.sync.Lock()
	for k, pe := range gpls.s {
		if pe == nil {
			continue
		}

		pe.ref--
		if pe.ref <= 0 {
			playerPool.Put(pe)
		}
		gpls.s[k] = nil
	}
	gpls.sync.Unlock()
}

func (gpls *gatePlays) grap(handle *util.NetHandle) *player {
	gpls.sync.Lock()
	hash := uint32(handle.HandleID()) & uint32(constPalyerMax-1)
	if gpls.s[hash] != nil && gpls.s[hash].handle.HandleID() == handle.HandleID() {
		gpls.s[hash].ref++
		gpls.sync.Unlock()
		return gpls.s[hash]
	}
	gpls.sync.Unlock()
	return nil
}

func (gpls *gatePlays) free(pe *player) {
	gpls.sync.Lock()
	pe.ref--
	if pe.ref <= 0 {
		playerPool.Put(pe)
	}
	gpls.sync.Unlock()
}
