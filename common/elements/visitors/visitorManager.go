package visitors

/*const (
	//ConstVisitorMax xxx
	ConstVisitorMax = 65535
)

var (
	//ErrVisitorFull Too many visitors
	ErrVisitorFull = errors.New("visitor is full")
)

//VisitorAllocer Visitor allocator
type VisitorAllocer func() IVisitor

//VisitorFree Visitor releaser
type VisitorFree func(p IVisitor)

type IVisitorManager interface {
	Spawned()
	Size() int
	SetAllocer(f VisitorAllocer)
	SetFree(f VisitorFree)
	GetKeys() []int32
	ToKey(sock int32) (uint16, error)
	Occupy(sock int32, addr net.IP, port int) (IVisitor, int32, error)
	Grap(key int32) IVisitor
	Erase(key int32)
	Release(v IVisitor)
}

//VisitorManager Visitor manager
type VisitorManager struct {
	allocer VisitorAllocer
	free    VisitorFree

	s     []IVisitor
	mps   map[int32]uint16
	sz    int
	seqID uint32
	sync.Mutex
}

// Spawned Initialize Visitor management module
func (vms *VisitorManager) Spawned() {
	vms.s = make([]IVisitor, ConstVisitorMax)
	vms.mps = make(map[int32]uint16, 64)
	vms.seqID = 1
}

// Size the PlayManager of number
func (vms *VisitorManager) Size() int {
	return vms.sz
}

//SetAllocer setting Distributor
func (vms *VisitorManager) SetAllocer(f VisitorAllocer) {
	vms.allocer = f
}

//SetFree Set releaser
func (vms *VisitorManager) SetFree(f VisitorFree) {
	vms.free = f
}

// GetKeys returns all the player object keys
func (vms *VisitorManager) GetKeys() []int32 {
	vms.Lock()
	defer vms.Unlock()
	i := 0
	hs := make([]int32, vms.sz)
	for _, client := range vms.s {
		if client == nil {
			continue
		}

		hs[i] = client.GetKey()
		i++
	}

	return hs
}

//ToKey Socket conversion to the corresponding handle id
func (vms *VisitorManager) ToKey(sock int32) (uint16, error) {
	vms.Lock()
	defer vms.Unlock()
	if v, ok := vms.mps[sock]; ok {
		return v, nil
	}
	return 0, errors.New("unknown id")
}

//Occupy Occupy a visitor resource
func (vms *VisitorManager) Occupy(sock int32, addr net.IP, port int) (IVisitor, int32, error) {
	v := vms.allocer().(IVisitor)
	v.SetAddr(addr)
	v.SetPort(port)

	var i uint32
	vms.Lock()

	for i = 0; i < ConstVisitorMax; i++ {
		key := ((i + vms.seqID) & ConstVisitorMax)
		hash := key & (ConstVisitorMax - 1)
		if vms.s[hash] == nil {
			v.SetKey(int32(key))
			v.SetSocket(sock)
			vms.seqID = key + 1
			vms.s[hash] = v
			vms.s[hash].IncRef()
			vms.s[hash].IncRef()
			vms.mps[sock] = uint16(key)
			vms.sz++
			vms.Unlock()
			return v, int32(key), nil
		}
	}
	vms.Unlock()
	vms.free(v)

	return nil, 0, ErrVisitorFull
}

// Grap Return a visitor object and add a reference counter
func (vms *VisitorManager) Grap(key int32) IVisitor {
	vms.Lock()
	defer vms.Unlock()
	hash := uint32(key) & uint32(ConstVisitorMax-1)
	if vms.s[hash] != nil && vms.s[hash].GetKey() == key {
		pe := vms.s[hash]
		pe.IncRef()
		return pe
	}
	return nil
}

// Erase removes the Visitor from VisitorManager
func (vms *VisitorManager) Erase(key int32) {
	vms.Lock()
	defer vms.Unlock()
	hash := uint32(key) & uint32(ConstVisitorMax-1)
	if vms.s[hash] != nil && vms.s[hash].GetKey() == key {
		pe := vms.s[hash]
		vms.s[hash] = nil
		if _, ok := vms.mps[pe.GetSocket()]; ok {
			delete(vms.mps, pe.GetSocket())
		}
		vms.sz--

		if pe.DecRef() <= 0 {
			vms.free(pe)
		}

	}
}

// Release  Release control
func (vms *VisitorManager) Release(v IVisitor) {
	vms.Lock()
	defer vms.Unlock()
	if v.DecRef() <= 0 {
		vms.free(v)
	}
}*/
