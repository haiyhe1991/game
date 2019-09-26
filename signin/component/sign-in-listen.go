package component

import (
	"github.com/yamakiller/game/common/module"
)

//SignInListen Login monitoring service
type SignInListen struct {
	module.InNetListen
}

//Init Login service object initialization
func (sil *SignInListen) Init() {
	sil.NetListenService.Init()
}
