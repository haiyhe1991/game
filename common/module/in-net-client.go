package module

import (
	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/util"
)

//
type InNetClient struct {
	implement.NetClient
	service.Service

	handle util.NetHandle
	sock   int32
}

//SetID Setting the client ID
func (inc *InNetClient) SetID(h uint64) {
	inc.handle.SetValue(h)
}

//GetID Returns the client ID
func (inc *InNetClient) GetID() uint64 {
	return inc.handle.GetValue()
}

//GetSocket Returns the client socket
func (inc *InNetClient) GetSocket() int32 {
	return inc.sock
}

//SetSocket Setting the client socket
func (inc *InNetClient) SetSocket(sock int32) {
	inc.sock = sock
}

//GetAuth return to certification time
func (inc *InNetClient) GetAuth() uint64 {
	return 0
}

//SetAuth Setting the time for authentication
func (inc *InNetClient) SetAuth(v uint64) {
}

//GetKeyPair Return key object
func (inc *InNetClient) GetKeyPair() interface{} {
	return nil
}

//BuildKeyPair Build key pair
func (inc *InNetClient) BuildKeyPair() {

}

//GetKeyPublic Return key publicly available information
func (inc *InNetClient) GetKeyPublic() string {
	return ""
}

//Shutdown Terminate this client service
func (inc *InNetClient) Shutdown() {
	inc.Service.Shutdown()
}
