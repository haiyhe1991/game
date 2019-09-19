package servers

import (
	"bytes"

	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/common/elements/visitors"
)

// Connection Service Connection Object
type Connection struct {
	id   int32
	sock int32
	auth uint64
	addr string
	data *bytes.Buffer

	keyPair visitors.VisitorKeyPair
	stat    visitors.VisitorStat
}

//GetAuth returns Auth time
func (c *Connection) GetAuth() uint64 {
	return c.auth
}

//SetAuth setting Auth time
func (c *Connection) SetAuth(v uint64) {
	c.auth = v
}

//GetID Get key value
func (c *Connection) GetID() int32 {
	return c.id
}

//SetID Set key value
func (c *Connection) SetID(v int32) {
	c.id = v
}

//GetSocket Get Socket ID
func (c *Connection) GetSocket() int32 {
	return c.sock
}

//SetSocket Set Socket ID
func (c *Connection) SetSocket(v int32) {
	c.id = v
}

//GetAddr Get Address[IP]
func (c *Connection) GetAddr() string {
	return c.addr
}

//SetAddr Set Addrees[IP]
func (c *Connection) SetAddr(v string) {
	c.addr = v
}

//GetData xxx
func (c *Connection) GetData() *bytes.Buffer {
	return c.data
}

//SetData xxx
func (c *Connection) SetData(v *bytes.Buffer) {
	c.data = v
}

//GetKeyPair xxx
func (c *Connection) GetKeyPair() *visitors.VisitorKeyPair {
	return &c.keyPair
}

//GetStat xxx
func (c *Connection) GetStat() *visitors.VisitorStat {
	return &c.stat
}

//Clear The size of the data that has been received
func (c *Connection) Clear() {
	n := c.data.Len()
	if n > 0 {
		c.data.Next(n)
	}
}

//Analysis Analytic data protocol
func (c *Connection) Analysis() (string, uint64, []byte, error) {
	return agreement.AgentParser(agreement.ConstInParser).Analysis(c.data)
}

/*// ConnStat Connection status information
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
}*/
