package clients

import (
	"bytes"
	"errors"
	"net"

	"github.com/yamakiller/game/gateway/elements/agreement"
	"github.com/yamakiller/magicNet/engine/util"
)

var (
	errPlayerFull = errors.New("player is full")
)

// PlayStat 连接者状态信息
type PlayStat struct {
	online        uint64
	lastRecvTime  uint64
	lastWriteTime uint64
	recvCount     uint64
	writeCount    uint64
}

// UpdateWrite Update write status data
func (pst *PlayStat) UpdateWrite(tts uint64, bytes uint64) {
	pst.lastWriteTime = tts
	pst.writeCount += bytes
}

// UpdateRecv Update read status data
func (pst *PlayStat) UpdateRecv(tts uint64, bytes uint64) {
	pst.lastRecvTime = tts
	pst.recvCount += bytes
}

// UpdateOnline Update time online
func (pst *PlayStat) UpdateOnline(tts uint64) {
	pst.online = tts
}

//Player External client connection object
type Player struct {
	Handle util.NetHandle
	Auth   uint64
	data   *bytes.Buffer
	addr   net.IP
	port   int
	Stat   PlayStat
	ref    int
}

//DataLen Get the length of the read buffer data
func (pe *Player) DataLen() int {
	return pe.data.Len()
}

//DataWrite Write data to the read buffer
func (pe *Player) DataWrite(p []byte) (int, error) {
	return pe.data.Write(p)
}

//DataAnalysis Play protocol data analysis
func (pe *Player) DataAnalysis() (string, []byte, error) {
	return agreement.ExtAnalysis(pe.data)
}
