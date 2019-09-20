package inclients

import (
	"bytes"
	"net"

	"github.com/yamakiller/game/common/elements/visitors"
	"github.com/yamakiller/game/common/agreement"
)

//InClient Internal network client
type InClient struct {
	id      int32
	key     int32
	sock    int32
	data    *bytes.Buffer
	addr    net.IP
	port    int
	keyPair visitors.VisitorKeyPair
	stat    visitors.VisitorStat
	ref     int
}

//GetAuth Get verification time tts
func (client *InClient) GetAuth() uint64 {
	return 0
}

//SetAuth Set verification time tts
func (client *InClient) SetAuth(v uint64) {

}

//GetID returns the client id
func (client *InClient) GetID() int32 {
	return client.id
}

//SetID setting the client id
func (client *InClient) SetID(v int32) {
	client.id = v
}

/*//SetGateway Set Gateway ID
func (client *InClient) SetGateway(id int32) {
	client.handle.Generate(id, client.handle.WorldID(), client.handle.HandleID(), client.handle.SocketID())
}

// GetGateway Get Gateway ID
func (client *InClient) GetGateway() int32 {
	return client.handle.GatewayID()
}*/

//GetKey Get key value
func (client *InClient) GetKey() int32 {
	return client.key
}

//SetKey Set key value
func (client *InClient) SetKey(v int32) {
	client.key = v
}

//GetSocket Get Socket ID
func (client *InClient) GetSocket() int32 {
	return client.sock
}

//SetSocket Set Socket ID
func (client *InClient) SetSocket(v int32) {
	client.sock = v
}

//GetAddr Get Address[IP]
func (client *InClient) GetAddr() net.IP {
	return client.addr
}

//SetAddr Set Addrees[IP]
func (client *InClient) SetAddr(v net.IP) {
	client.addr = v
}

//GetPort Get network Port
func (client *InClient) GetPort() int {
	return client.port
}

//SetPort Set newwork port
func (client *InClient) SetPort(v int) {
	client.port = v
}

//GetData xxx
func (client *InClient) GetData() *bytes.Buffer {
	return client.data
}

//SetData xxx
func (client *InClient) SetData(v *bytes.Buffer) {
	client.data = v
}

//GetKeyPair xxx
func (client *InClient) GetKeyPair() *visitors.VisitorKeyPair {
	return &client.keyPair
}

//GetStat xxx
func (client *InClient) GetStat() *visitors.VisitorStat {
	return &client.stat
}

//IncRef xxx
func (client *InClient) IncRef() {
	client.ref++
}

//DecRef xxx
func (client *InClient) DecRef() int {
	client.ref--
	return client.ref
}

//RestRef xxx
func (client *InClient) RestRef() {
	client.ref = 0
}

//Analysis InClient protocol data analysis
func (client *InClient) Analysis() (string, uint64, []byte, error) {
	return agreement.AgentParser(agreement.ConstInParser).Analysis(client.GetData())
}
