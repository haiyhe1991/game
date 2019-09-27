package logic

import "github.com/yamakiller/game/common/module"

func init() {
	module.FactoryInstance().Register("logic.SignOutProc", &SignOutProc{})
	module.FactoryInstance().Register("logic.SignInProc", &SignInProc{})
}
