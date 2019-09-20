package route

import "errors"

var (
	//ErrRouteAgreeUnDefined Error informat Agreement not defined
	ErrRouteAgreeUnDefined = errors.New("Agreement not defined")
	//ErrRouteAgreeUnRegister Error informat Agreement not register
	ErrRouteAgreeUnRegister = errors.New("Agreement not register")
	//ErrRoutePlayerUnverified Error informat Player Unverified
	ErrRoutePlayerUnverified = errors.New("Player Unverified")
	//ErrRouteServiceNotStarted Error informat Service has not started yet
	ErrRouteServiceNotStarted = errors.New("Service has not started yet")
)

// NewTable Generate a routing table
func NewTable() *Table {
	return &Table{tb: make(map[interface{}]Data, 64)}
}

//Table Routing address table
type Table struct {
	tb map[interface{}]Data
}

//Register Register routing data
func (t *Table) Register(agreement interface{}, agreementName, serviceName string, auth bool) {
	t.tb[agreement] = Data{Agreement: agreement, AgreementName: agreementName, ServiceName: serviceName, Auth: auth}
}

// Sreach Query a routing information
func (t *Table) Sreach(proto interface{}) *Data {
	if p, ok := t.tb[proto]; ok {
		return &p
	}
	return nil
}
