package clients

import (
	"bytes"
	"net"

	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/common/elements/visitors"

	"github.com/yamakiller/magicNet/util"
)

// Client 客户端对象
type Client struct {
	handle  util.NetHandle
	auth    uint64
	data    *bytes.Buffer
	addr    net.IP
	port    int
	keyPair visitors.VisitorKeyPair
	stat    visitors.VisitorStat
	ref     int
}

//GetAuth Get verification time tts
func (client *Client) GetAuth() uint64 {
	return client.auth
}

//SetAuth Set verification time tts
func (client *Client) SetAuth(v uint64) {
	client.auth = v
}

//GetKeyValue xxx
func (client *Client) GetKeyValue() uint64 {
	return client.handle.GetValue()
}

//SetGateway Set Gateway ID
func (client *Client) SetGateway(id int32) {
	client.handle.Generate(id, client.handle.WorldID(), client.handle.HandleID(), client.handle.SocketID())
}

// GetGateway Get Gateway ID
func (client *Client) GetGateway() int32 {
	return client.handle.GatewayID()
}

//GetKey Get key value
func (client *Client) GetKey() int32 {
	return client.handle.HandleID()
}

//SetKey Set key value
func (client *Client) SetKey(v int32) {
	client.handle.Generate(client.handle.GatewayID(), client.handle.WorldID(), v, client.handle.SocketID())
}

//GetSocket Get Socket ID
func (client *Client) GetSocket() int32 {
	return client.handle.SocketID()
}

//SetSocket Set Socket ID
func (client *Client) SetSocket(v int32) {
	client.handle.Generate(client.handle.GatewayID(), client.handle.WorldID(), client.handle.HandleID(), v)
}

//GetAddr Get Address[IP]
func (client *Client) GetAddr() net.IP {
	return client.addr
}

//SetAddr Set Addrees[IP]
func (client *Client) SetAddr(v net.IP) {
	client.addr = v
}

//GetPort Get network Port
func (client *Client) GetPort() int {
	return client.port
}

//SetPort Set newwork port
func (client *Client) SetPort(v int) {
	client.port = v
}

//GetData xxx
func (client *Client) GetData() *bytes.Buffer {
	return client.data
}

//SetData xxx
func (client *Client) SetData(v *bytes.Buffer) {
	client.data = v
}

//GetKeyPair xxx
func (client *Client) GetKeyPair() *visitors.VisitorKeyPair {
	return &client.keyPair
}

//GetStat xxx
func (client *Client) GetStat() *visitors.VisitorStat {
	return &client.stat
}

//IncRef xxx
func (client *Client) IncRef() {
	client.ref++
}

//DecRef xxx
func (client *Client) DecRef() int {
	client.ref--
	return client.ref
}

//RestRef xxx
func (client *Client) RestRef() {
	client.ref = 0
}

//Analysis Play protocol data analysis
func (client *Client) Analysis() (string, []byte, error) {
	name, _, data, err := agreement.AgentParser(agreement.ConstExParser).Analysis(client.GetData())
	return name, data, err
}

//var (
//	errPlayerFull = errors.New("player is full")
//)

/*// ClientStat 连接者状态信息
type ClientStat struct {
	online        uint64
	lastRecvTime  uint64
	lastWriteTime uint64
	recvCount     uint64
	writeCount    uint64
}

// UpdateWrite Update write status data
func (cst *ClientStat) UpdateWrite(tts uint64, bytes uint64) {
	cst.lastWriteTime = tts
	cst.writeCount += bytes
}

// UpdateRecv Update read status data
func (cst *ClientStat) UpdateRecv(tts uint64, bytes uint64) {
	cst.lastRecvTime = tts
	cst.recvCount += bytes
}

// UpdateOnline Update time online
func (cst *ClientStat) UpdateOnline(tts uint64) {
	cst.online = tts
}

//Client External client connection object
type Client struct {
	Handle util.NetHandle
	Auth   uint64
	data   *bytes.Buffer
	addr   net.IP
	port   int
	Stat   ClientStat
	ref    int
}

//DataLen Get the length of the read buffer data
func (ct *Client) DataLen() int {
	return ct.data.Len()
}

//DataWrite Write data to the read buffer
func (ct *Client) DataWrite(p []byte) (int, error) {
	return ct.data.Write(p)
}

//DataAnalysis Play protocol data analysis
func (ct *Client) DataAnalysis() (string, []byte, error) {
	name, _, data, err := agreement.AgentParser(agreement.ConstExParser).Analysis(ct.data)
	return name, data, err
}*/
