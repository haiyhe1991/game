package gateway

import (
	"bytes"
	"net"
)

const (
	clientTypePlay   = 0
	clientTypeServer = 1
)

type stat struct {
	lastRecvTm uint64
	lastSendTm uint64

	recvBytes uint64
	sendBytes uint64
}

type client struct {
	id     uint64
	playID uint64
	tpe    int
	data   *bytes.Buffer
	addr   net.IP
	port   int
	sta    stat
}

func (ct *client) registerSocket(gatewayID int32, sock int32) {
	ct.id = ((uint64(gatewayID) << 32) & uint64(sock))
}
