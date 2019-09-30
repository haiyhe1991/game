package logic

import (
	"github.com/yamakiller/game/common/module"
	"github.com/yamakiller/game/common/modulelogic"
)

func init() {
	module.FactoryInstance().Register("logic.GatewayRegisterProc", &modulelogic.GatewayRegisterProc{})
	module.FactoryInstance().Register("logic.SignOutProc", &SignOutProc{})
	module.FactoryInstance().Register("logic.SignInProc", &SignInProc{})
}
