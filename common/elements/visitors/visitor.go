package visitors

import (
	"bytes"
	"net"
)

//VisitorKeyPair Visitor the key pair
type VisitorKeyPair struct {
}

// VisitorStat Visitor status data
type VisitorStat struct {
	online        uint64
	lastRecvTime  uint64
	lastWriteTime uint64
	recvCount     uint64
	writeCount    uint64
}

// UpdateWrite Update write status data
func (stat *VisitorStat) UpdateWrite(tts uint64, bytes uint64) {
	stat.lastWriteTime = tts
	stat.writeCount += bytes
}

// UpdateRecv Update read status data
func (stat *VisitorStat) UpdateRecv(tts uint64, bytes uint64) {
	stat.lastRecvTime = tts
	stat.recvCount += bytes
}

// UpdateOnline Update time online
func (stat *VisitorStat) UpdateOnline(tts uint64) {
	stat.online = tts
}

//IVisitor Visitor interface
type IVisitor interface {
	GetAuth() uint64
	SetAuth(v uint64)
	GetKey() int32
	SetKey(v int32)
	GetSocket() int32
	SetSocket(v int32)
	GetAddr() net.IP
	SetAddr(v net.IP)
	GetPort() int
	SetPort(v int)
	GetData() *bytes.Buffer
	SetData(v *bytes.Buffer)
	GetKeyPair() *VisitorKeyPair
	GetStat() *VisitorStat
	RestRef()
	IncRef()
	DecRef() int
}

// VisitorStat xxx
/*type VisitorStat struct {
	online        uint64
	lastRecvTime  uint64
	lastWriteTime uint64
	recvCount     uint64
	writeCount    uint64
}

// VisitorKey xxx
type VisitorKey struct {
}

// UpdateWrite Update write status data
func (vss *VisitorStat) UpdateWrite(tts uint64, bytes uint64) {
	vss.lastWriteTime = tts
	vss.writeCount += bytes
}

// UpdateRecv Update read status data
func (vss *VisitorStat) UpdateRecv(tts uint64, bytes uint64) {
	vss.lastRecvTime = tts
	vss.recvCount += bytes
}

// UpdateOnline Update time online
func (vss *VisitorStat) UpdateOnline(tts uint64) {
	vss.online = tts
}

// Visitor Gateway to service connection object
type Visitor struct {
	Sock int32
	Addr net.IP
	Port int
	Key  VisitorKey
	Stat VisitorStat
	Auth uint64

	data *bytes.Buffer
	ref  int
}

//DRLen Get the length of the read buffer data
func (v *Visitor) DRLen() int {
	return v.data.Len()
}

//DRWrite Write data to the read buffer
func (v *Visitor) DRWrite(p []byte) (int, error) {
	return v.data.Write(p)
}

// GetData xxx
func (v *Visitor) GetData() *bytes.Buffer {
	return v.data
}

// SetData xxxx
func (v *Visitor) SetData(d *bytes.Buffer) {
	v.data = d
}

//RefRest xxxx
func (v *Visitor) RefRest() {
	v.ref = 0
}*/
