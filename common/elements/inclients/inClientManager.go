package inclients

import (
	"bytes"
	"sync"

	"github.com/yamakiller/game/common/elements/visitors"
	"github.com/yamakiller/game/gateway/constant"
)

var clientPool = sync.Pool{
	New: func() interface{} {
		b := new(InClient)
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

//InClientManager Local client management
type InClientManager struct {
	visitors.VisitorManager
}

// Spawned Initialize Client management module
func (cms *InClientManager) Spawned() {
	cms.SetAllocer(func() visitors.IVisitor {
		return clientPool.Get().(*InClient)
	})

	cms.SetFree(func(p visitors.IVisitor) {
		c := p.(*InClient)
		clientPool.Put(c)
	})

	cms.Spawned()
}
