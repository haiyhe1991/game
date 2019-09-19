package clients

import (
	"bytes"
	"sync"

	"github.com/yamakiller/game/common/elements/visitors"

	"github.com/yamakiller/game/gateway/constant"
)

var clientPool = sync.Pool{
	New: func() interface{} {
		b := new(Client)
		if b.GetData() == nil {
			b.SetData(bytes.NewBuffer([]byte{}))
			b.GetData().Grow(constant.ConstPlayerBufferLimit)
		} else {
			b.GetData().Reset()
		}
		b.SetAuth(0)
		b.RestRef()
		return b
	},
}

//ClientManager client Manager
type ClientManager struct {
	visitors.VisitorManager
}

// Spawned Initialize Client management module
func (cms *ClientManager) Spawned() {
	cms.SetAllocer(func() visitors.IVisitor {
		return clientPool.Get().(*Client)
	})

	cms.SetFree(func(p visitors.IVisitor) {
		c := p.(*Client)
		clientPool.Put(c)
	})
	cms.Spawned()
}
