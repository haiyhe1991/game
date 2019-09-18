package visitors

import (
	"errors"
	"net"
	"sync"
)

const (
	ConstVisitorMax = 65535
)

var (
	ErrVisitorFull = errors.New("visitor is full")
)

type VisitorAllocer func() interface{}
type VisitorFree func(p interface{})
type VisitorComparison func(a interface{}, b interface{}) int

//VisitorManager
type VisitorManager struct {
	Allocer    VisitorAllocer
	Free       VisitorFree
	Comparison VisitorComparison

	s     []*Visitor
	mps   map[int32]uint16
	sz    int
	seqID uint16
	sync.Mutex
}

// Spawned Initialize Visitor management module
func (vms *VisitorManager) Spawned() {
	vms.s = make([]*Visitor, ConstVisitorMax)
	vms.mps = make(map[int32]uint16, 64)
	vms.seqID = 1
}

func (vms *VisitorManager) Occupy(sock int32, addr net.IP, port int) (interface{}, int32, error) {
	v := vms.Allocer().(*Visitor) //clientPool.Get().(*Client)
	v.Addr = addr
	v.Port = port

	var i uint16
	vms.Lock()

	for i = 0; i < ConstVisitorMax; i++ {
		key := ((i + vms.seqID) & ConstVisitorMax)
		hash := key & (ConstVisitorMax - 1)
		if vms.s[hash] == nil {
			//handle := util.NetHandle{}
			//handle.Generate(cms.d, 0, int32(key), sock)
			//client.Handle = handle
			vms.seqID = key + 1
			vms.s[hash] = v
			vms.s[hash].ref = 2
			vms.mps[sock] = key
			vms.sz++
			vms.Unlock()
			return v, int32(key), nil
		}
	}
	vms.Unlock()
	vms.Free(v)

	return nil, 0, ErrVisitorFull
}

func (vms *VisitorManager) Grap(key uint32) interface{} {
	vms.Lock()
	defer vms.Unlock()
	hash := key & uint32(ConstVisitorMax-1)
	if vms.s[hash] != nil && vms.Comparison(vms.s[hash], key) == 0 {
		pe := vms.s[hash]
		pe.ref++
		return pe
	}
	return nil
}

// Erase removes the Visitor from VisitorManager
func (vms *VisitorManager) Erase(key uint32) {
	vms.Lock()
	defer vms.Unlock()
	hash := key & uint32(ConstVisitorMax-1)
	if vms.s[hash] != nil && vms.Comparison(vms.s[hash], key) == 0 {
		pe := vms.s[hash]
		vms.s[hash] = nil
		if _, ok := vms.mps[pe.Sock]; ok {
			delete(vms.mps, pe.Sock)
		}
		vms.sz--
		pe.ref--
		if pe.ref <= 0 {
			vms.Free(pe)
		}

	}
}

// Release  Release control
func (vms *VisitorManager) Release(v *Visitor) {
	vms.Lock()
	defer vms.Unlock()
	v.ref--
	if v.ref <= 0 {
		vms.Free(v)
	}
}

// ToHandleID Socket conversion to the corresponding handle id
/*func (cms *ClientManager) ToHandleID(sock int32) (uint16, error) {
	cms.Lock()
	defer cms.Unlock()
	if v, ok := cms.mps[sock]; ok {
		return v, nil
	}
	return 0, errors.New("unknown id")
}*/

// Size the PlayManager of number
func (vms *VisitorManager) Size() int {
	return vms.sz
}

// GetHandls Get all the player object handles
/*func (vms *VisitorManager) GetHandls() []util.NetHandle {
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
}*/
