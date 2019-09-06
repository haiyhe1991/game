package gateway

import "reflect"

type routeTable struct {
	tb map[interface{}]protoRegister
}

func (rt *routeTable) register(proto interface{}, route string) {
	rt.tb[reflect.TypeOf(proto)] = protoRegister{proto: proto, route: route}
}

func (rt *routeTable) isAllowable(msg interface{}, route string) bool {
	if p, ok := rt.tb[reflect.TypeOf(msg)]; ok {
		if p.route == route {
			return true
		}
	}
	return false
}
