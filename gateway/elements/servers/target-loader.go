package servers

import (
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/st/hash"
)

//NewLoadSet Assign a load set
func NewLoadSet() *TargetLoadSet {
	return &TargetLoadSet{set: make(map[string]*TargetLoader)}
}

//NewLoader Assign an equalizer
func NewLoader(replicas int) *TargetLoader {
	return &TargetLoader{Map: *hash.New(replicas)}
}

//TargeObject target
type TargeObject struct {
	ID     uint32
	Target *actor.PID
}

//TargetLoadSet Target load set
type TargetLoadSet struct {
	set map[string]*TargetLoader
}

//Add Add a cluster type
func (tlset *TargetLoadSet) Add(name string, tl *TargetLoader) {
	tlset.set[name] = tl
}

//Get Return to a cluster
func (tlset *TargetLoadSet) Get(name string) *TargetLoader {
	if v, ok := tlset.set[name]; ok {
		return v
	}
	return nil
}

//TargetLoader Provide load balancing management for servers
type TargetLoader struct {
	hash.Map
}

//AddTarget Join a target service
func (t *TargetLoader) AddTarget(key string, v *TargeObject) {
	t.Lock()
	defer t.Unlock()
	t.UnAdd(key, v)
}

//RemoveTarget Delete a target service
func (t *TargetLoader) RemoveTarget(key string) {
	t.Lock()
	defer t.Unlock()
	t.UnRemove(key)
}

//GetTarget Return Return a service target
func (t *TargetLoader) GetTarget(key string) *TargeObject {
	t.RLock()
	defer t.RUnlock()
	r, err := t.UnGet(key)
	if err != nil {
		return nil
	}
	return r.(*TargeObject)
}
