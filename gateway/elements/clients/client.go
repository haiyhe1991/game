package clients

import (
	"bytes"
	"errors"
	"net"

	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/magicNet/util"
)

var (
	errPlayerFull = errors.New("player is full")
)

// ClientStat 连接者状态信息
type ClientStat struct {
	online        uint64
	lastRecvTime  uint64
	lastWriteTime uint64
	recvCount     uint64
	writeCount    uint64
}

// UpdateWrite Update write status data
func (cst *ClientStat) UpdateWrite(tts uint64, bytes uint64) {
	cst.lastWriteTime = tts
	cst.writeCount += bytes
}

// UpdateRecv Update read status data
func (cst *ClientStat) UpdateRecv(tts uint64, bytes uint64) {
	cst.lastRecvTime = tts
	cst.recvCount += bytes
}

// UpdateOnline Update time online
func (cst *ClientStat) UpdateOnline(tts uint64) {
	cst.online = tts
}

//Client External client connection object
type Client struct {
	Handle util.NetHandle
	Auth   uint64
	data   *bytes.Buffer
	addr   net.IP
	port   int
	Stat   ClientStat
	ref    int
}

//DataLen Get the length of the read buffer data
func (ct *Client) DataLen() int {
	return ct.data.Len()
}

//DataWrite Write data to the read buffer
func (ct *Client) DataWrite(p []byte) (int, error) {
	return ct.data.Write(p)
}

//DataAnalysis Play protocol data analysis
func (ct *Client) DataAnalysis() (string, []byte, error) {
	name, _, data, err := agreement.AgentParser(agreement.ConstExParser).Analysis(ct.data)
	return name, data, err
}
