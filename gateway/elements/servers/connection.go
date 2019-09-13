package servers

import (
	"bytes"
	"sync"
)

// ConnStat Connection status information
type ConnStat struct {
	online        uint64
	lastRecvTime  uint64
	lastWriteTime uint64
	recvCount     uint64
	writeCount    uint64
}

// Connection Service Connection Object
type Connection struct {
	ID   int32 //Unique in the cluster
	Sock int32
	Addr string
	data *bytes.Buffer
	stat ConnStat

	sync sync.Mutex
}

func (cn *Connection) ClearData() {
	n := cn.data.Len()
	if n > 0 {
		cn.data.Next(n)
	}
}
