package servers

import (
	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/timer"
)

//TargetConnection Connection to the target server configuration information
type TargetConnection struct {
	ID         int32
	Name       string
	Addr       string
	Desc       string
	TimeOut    uint64
	OutChanMax int

	virtualID uint32
	socket    int32
	startTime uint64
	stat      implement.NetConnectEtat
}

//GetVirtualID Returns the virtual ID mapped to the Hash storage module
func (trc *TargetConnection) GetVirtualID() uint32 {
	return trc.virtualID
}

//SetVirtualID Setting the virtual ID mapped to the Hash storage module
func (trc *TargetConnection) SetVirtualID(vID uint32) {
	trc.virtualID = vID
}

//GetName Retruns the target name
func (trc *TargetConnection) GetName() string {
	return trc.Name
}

//GetAddr Return target address
func (trc *TargetConnection) GetAddr() string {
	return trc.Addr
}

//GetEtat Returns the target in connection status
func (trc *TargetConnection) GetEtat() implement.NetConnectEtat {
	return trc.stat
}

//SetEtat Setting the target in connection status
func (trc *TargetConnection) SetEtat(stat implement.NetConnectEtat) {
	trc.stat = stat
}

//GetOutSize Returns the target out chan buffer size
func (trc *TargetConnection) GetOutSize() int {
	return trc.OutChanMax
}

//GetSocket Returns the target socket value
func (trc *TargetConnection) GetSocket() int32 {
	return trc.socket
}

//SetSocket Setting the target socket value
func (trc *TargetConnection) SetSocket(s int32) {
	trc.socket = s
}

//IsTimeout Calculate whether it times out and return the timeout in milliseconds
func (trc *TargetConnection) IsTimeout() uint64 {
	if trc.TimeOut == 0 {
		return 0
	}

	n := timer.Now() - trc.startTime
	if n > trc.TimeOut {
		return n
	}

	return 0
}

//RestTimeout Reset connection  start time
func (trc *TargetConnection) RestTimeout() {
	trc.startTime = timer.Now()
}
