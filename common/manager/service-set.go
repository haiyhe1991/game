package manager

import (
	"sync"

	"github.com/yamakiller/magicNet/engine/actor"
)

//NewSSets xxx
func NewSSets() *SSets {
	return &SSets{s: make(map[string]actor.PID)}
}

//SSets Service set
type SSets struct {
	s map[string]actor.PID
	sync.Mutex
}

//Sreach Search for the PID of the corresponding service
func (sset *SSets) Sreach(name string) *actor.PID {
	sset.Lock()
	defer sset.Unlock()

	if d, ok := sset.s[name]; ok {
		return &d
	}

	return nil
}

//Push Add a service
func (sset *SSets) Push(name string, pid *actor.PID) {
	sset.Lock()
	defer sset.Unlock()
	sset.s[name] = *pid
}

//Erase xxx
func (sset *SSets) Erase(name string) {
	sset.Lock()
	defer sset.Unlock()
	if _, ok := sset.s[name]; ok {
		delete(sset.s, name)
		return
	}
}
