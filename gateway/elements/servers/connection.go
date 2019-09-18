package servers

import (
	"bytes"
	"sync"

	"github.com/yamakiller/game/common/agreement"
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

//DataAnalysis Analytic data protocol
func (cn *Connection) DataAnalysis() (string, uint64, []byte, error) {
	return agreement.AgentParser(agreement.ConstInParser).Analysis(cn.data)
}

//DataWrite Write data to buffer
func (cn *Connection) DataWrite(d []byte) (int, error) {
	return cn.data.Write(d)
}

// DataLen The size of the data that has been received
func (cn *Connection) DataLen() int {
	return cn.data.Len()
}

// ClearData The size of the data that has been received
func (cn *Connection) ClearData() {
	n := cn.data.Len()
	if n > 0 {
		cn.data.Next(n)
	}
}
