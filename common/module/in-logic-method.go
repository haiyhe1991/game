package module

import (
	"sync"

	"github.com/yamakiller/magicNet/service/implement"
)

var (
	onceMethod sync.Once

	logicMethodMap *implement.NetMethodDispatch
)

//LogicInstance 获取逻辑映射器
func LogicInstance() *implement.NetMethodDispatch {
	onceMethod.Do(func() {
		logicMethodMap = implement.NewMethodDispatch()
	})

	return logicMethodMap
}
