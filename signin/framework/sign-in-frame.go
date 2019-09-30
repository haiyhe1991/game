package framework

import (
	"errors"

	"github.com/yamakiller/game/common/module"
	"github.com/yamakiller/game/signin/component"
	"github.com/yamakiller/game/signin/global"
	_ "github.com/yamakiller/game/signin/logic"
	"github.com/yamakiller/magicNet/core"
	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/util"
)

//SignInFrame Sign in server framework
type SignInFrame struct {
	core.DefaultStart
	core.DefaultEnv
	core.DefaultLoop

	core.DefaultService
	core.DefaultCMDLineOption

	id    int32
	addr  string
	max   int
	ccmax int
	keep  uint64

	startupScript string

	snetListen *component.SignInListen
}

//InitService Initialize the login service
func (sif *SignInFrame) InitService() error {
	if err := sif.DefaultService.InitService(); err != nil {
		return err
	}

	signInEnv := util.GetEnvMap(util.GetEnvRoot(), "sign-in")
	if signInEnv == nil {
		return errors.New("Sign-In configuration information does not exist ")
	}

	sif.id = int32(util.GetEnvInt(signInEnv, "id", 1))
	sif.addr = util.GetEnvString(signInEnv, "addr", "0.0.0.0:7851")
	sif.max = util.GetEnvInt(signInEnv, "max", 1024)
	sif.ccmax = util.GetEnvInt(signInEnv, "chan-max", 1024)
	sif.keep = uint64(util.GetEnvInt64(signInEnv, "client-keep", 1000))
	sif.startupScript = util.GetEnvString(signInEnv,
		"client-startup-script-file",
		"./script/sign-in-client-register.lua")

	//Initialize the client connection repository
	module.NewWarehouse(component.NewSignInManager())

	//register net event
	scirpt := module.InNetScript{}
	scirpt.Execution(sif.startupScript, module.LogicInstance())

	/*if err := module.ReadisEnvAnalysis(signInEnv); err != nil {
		return err
	}*/

	sif.snetListen = func() *component.SignInListen {
		return service.Make(global.ConstNetworkServiceName, func() service.IService {
			h := &component.SignInListen{InNetListen: module.SpawnInNetListen(
				module.GetWarehouse(),
				&module.InNetListenDeleate{},
				sif.addr,
				sif.ccmax,
				sif.max,
				sif.keep)}
			h.Init()
			return h
		}).(*component.SignInListen)
	}()

	return nil
}

//CloseService Close service
func (sif *SignInFrame) CloseService() {
	if sif.snetListen != nil {
		sif.snetListen.Shutdown()
		sif.snetListen = nil
	}

	module.RedisClose()
}
