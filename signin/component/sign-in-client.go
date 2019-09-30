package component

import (
	"bytes"
	"strconv"
	"sync/atomic"

	"github.com/yamakiller/game/common"
	"github.com/yamakiller/game/common/module"
	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/service/implement"
)

const (
	constSignInClientName = "SignIn/Client"
)

var (
	signInClientSerial = uint32(1)
)

//NewSignInManager xxx
func NewSignInManager() *module.InNetClientManage {
	return &module.InNetClientManage{NetClientManager: implement.NetClientManager{Malloc: &SignInAllocer{}}}
}

//SignInAllocer Sign In client allocator
type SignInAllocer struct {
	module.InNetClientAllocer
}

//New Assign connection client service
func (sia *SignInAllocer) New() implement.INetClient {
	newSerial := atomic.AddUint32(&signInClientSerial, 1)
	r := service.Make(constSignInClientName+strconv.Itoa(int(newSerial)), func() service.IService {
		h := &SignInClient{InNetClient: module.SpawnInNetClient()}
		if h.GetRecvBuffer() == nil {
			h.SetRecvBuffer(bytes.NewBuffer([]byte{}))
			h.GetRecvBuffer().Grow(common.ConstInClientBufferLimit)
		} else {
			h.GetRecvBuffer().Reset()
		}

		h.SetAuth(0)
		h.SetRef(0)

		h.Init()
		return h
	}).(*SignInClient)
	return r
}

//Delete Delete and stop the service
func (sia *SignInAllocer) Delete(p implement.INetClient) {
	sia.InNetClientAllocer.Delete(p)
}

//SignInClient Login client service object
type SignInClient struct {
	module.InNetClient
}

//Init Initialize the login client service object
func (sic *SignInClient) Init() {
	sic.InNetClient.Init()
	sic.Dispatch = sic.onDispatch
}

func (sic *SignInClient) onDispatch(event *module.InNetMethodClientEvent) {
	f := module.GetLogicMethod(event.Name, module.LogicInstance())
	if f == nil {
		sic.LogError("Agreement %s is not registered", event.Name)
		return
	}

	f(event)
}
