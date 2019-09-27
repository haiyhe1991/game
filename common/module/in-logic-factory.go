package module

import (
	"reflect"
	"sync"
)

var (
	once sync.Once
	//logicWarehouse Logical object warehouse
	logicWarehouse *InLogicWarehouse
)

//FactoryInstance Logical warehouse
func FactoryInstance() *InLogicWarehouse {
	once.Do(func() {
		logicWarehouse = &InLogicWarehouse{factory: inLogicFactory{factorys: make(map[string]reflect.Type)},
			warehouses: make(map[string]interface{})}
	})
	return logicWarehouse
}

//inLogicFactory Factory that generates logical objects
type inLogicFactory struct {
	factorys map[string]reflect.Type
}

//Spawn Generate object
func (ilf *inLogicFactory) spawn(name string) interface{} {
	f, success := ilf.factorys[name]
	if !success {
		return nil
	}

	return reflect.Indirect(reflect.New(f.Elem())).Addr().Interface()
}

func (ilf *inLogicFactory) register(name string, k interface{}) {
	ilf.factorys[name] = reflect.TypeOf(k)
}

//InLogicWarehouse xxxx
type InLogicWarehouse struct {
	factory    inLogicFactory
	warehouses map[string]interface{}
	sync       sync.Mutex
}

//Get Back to a generated object
func (ilw *InLogicWarehouse) Get(name string) interface{} {
	ilw.sync.Lock()
	defer ilw.sync.Unlock()

	v, success := ilw.warehouses[name]
	if success {
		return v
	}
	v = ilw.factory.spawn(name)
	if v != nil {
		ilw.warehouses[name] = v
	}

	return v
}

//Register Register a logical object
func (ilw *InLogicWarehouse) Register(name string, inst interface{}) {
	ilw.sync.Lock()
	defer ilw.sync.Unlock()

	ilw.factory.register(name, inst)
}

//Destory destroy
func (ilw *InLogicWarehouse) Destory() {
	ilw.sync.Lock()
	defer ilw.sync.Unlock()
	ilw.warehouses = make(map[string]interface{})
	ilw.factory.factorys = make(map[string]reflect.Type)
}
