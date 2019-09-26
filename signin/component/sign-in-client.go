package component

import (
	"reflect"
	"unsafe"

	"github.com/yamakiller/game/common/module"
)

type SignInClient struct {
	module.InNetClient
}

func (sic *SignInClient) Init() {
	sic.InNetClient.Init()
	script := module.InNetScript{}
	script.Execution("./script/sign_in_client_register.lua", unsafe.Pointer(sic), reflect.TypeOf(sic))
}
