package clients

import (
	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/util"
)

//GClient Gateway Connection from the Internet
type GClient struct {
	implement.NetClient

	handle util.NetHandle
	auth   uint64
}

//SetID Setting the client ID
func (gct *GClient) SetID(h *util.NetHandle) {
	gct.handle = *h
}

//GetID Returns the client ID
func (gct *GClient) GetID() *util.NetHandle {
	return &gct.handle
}

//GetSocket Returns the client socket
func (gct *GClient) GetSocket() int32 {
	return gct.handle.GetSocket()
}

//SetSocket Setting the client socket
func (gct *GClient) SetSocket(sock int32) {
	gct.handle.Generate(gct.handle.GetServiceID(), gct.handle.GetHandle(), sock)
}

//GetAuth return to certification time
func (gct *GClient) GetAuth() uint64 {
	return gct.auth
}

//SetAuth Setting the time for authentication
func (gct *GClient) SetAuth(v uint64) {
	gct.auth = v
}

//GetKeyPair Return key object
func (gct *GClient) GetKeyPair() interface{} {
	return nil
}

//BuildKeyPair Build key pair
func (gct *GClient) BuildKeyPair() {

}

//GetKeyPublic Return key publicly available information
func (gct *GClient) GetKeyPublic() string {
	return ""
}
