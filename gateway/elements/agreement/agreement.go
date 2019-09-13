package agreement

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/yamakiller/game/gateway/elements"
)

type protoRegister struct {
	proto interface{}
	route string
	auth  bool
}

const (
	constAgreeHeader          = 3
	constAgreeHeaderBit       = 24
	constAgreeMakeUPBit       = 8
	constAgreeVersionBit      = 7
	constAgreeVersionShift    = (constAgreeHeaderBit + constAgreeMakeUPBit) - constAgreeVersionBit
	constAgreeVersionMask     = 0x7F
	constAgreeDataLengthBit   = 12
	constAgreeDataLengthShift = constAgreeVersionShift - constAgreeDataLengthBit
	constAgreeDataLengthMask  = 0xfff
	constAgreeNameLengthBit   = 5
	constAgreeNameLengthShift = constAgreeDataLengthShift - constAgreeNameLengthBit
	constAgreeNameLengthMask  = 0x1f

	constAgreeSingleLimit = (elements.ConstPlayerBufferLimit >> 1) - constAgreeHeader
)

var (
	//ErrProtoIllegal An illegal agreement
	ErrProtoIllegal = errors.New("An illegal agreement")
)

//ForwardMessage Forward data event
type ForwardMessage struct {
	Handle        uint64
	AgreementName string
	ServerName    string
	Data          []byte
}

//CheckConnectMessage Check the connection status event to achieve automatic reconnection after disconnection
type CheckConnectMessage struct {
}

// ExtAnalysis Play protocol data analysis
func ExtAnalysis(data *bytes.Buffer) (string, []byte, error) {

	if data.Len() < constAgreeHeader {
		return "", nil, nil
	}

	headByte := make([]byte, 4)
	headByte[0] = 0
	headByte[1], _ = data.ReadByte()
	headByte[2], _ = data.ReadByte()
	headByte[3], _ = data.ReadByte()

	head := binary.BigEndian.Uint32(headByte)

	dl := (head >> constAgreeDataLengthShift) & constAgreeDataLengthMask
	nl := (head >> constAgreeNameLengthShift) & constAgreeNameLengthMask

	if (dl + nl) > uint32(data.Len()) {
		data.UnreadByte()
		data.UnreadByte()
		data.UnreadByte()
		return "", nil, nil
	}

	if (dl + nl) > constAgreeSingleLimit {
		data.Reset()
		return "", nil, ErrProtoIllegal
	}

	pname := string(data.Next(int(nl)))
	pdata := data.Next(int(dl))

	return pname, pdata, nil
}
