package gateway

import (
	"errors"
)

var (
	errRouteAgreeUnDefined   = errors.New("Agreement not defined")
	errRouteAgreeUnRegister  = errors.New("Agreement not register")
	errRoutePlayerUnverified = errors.New("Player Unverified")
)

type routeTable struct {
	tb map[interface{}]protoRegister
}

func (rt *routeTable) register(proto interface{}, route string, auth bool) {
	rt.tb[proto] = protoRegister{proto: proto, route: route, auth: auth}
}

func (rt *routeTable) get(proto interface{}) *protoRegister {
	if p, ok := rt.tb[proto]; ok {
		return &p
	}
	return nil
}

/*func (rt *routeTable) isAllowable(msg interface{}, route string) bool {
	if p, ok := rt.tb[reflect.TypeOf(msg)]; ok {
		if p.route == route {
			return true
		}
	}
	return false
}*/
