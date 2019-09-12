package gateway

import (
	"bytes"
	"encoding/binary"
	"errors"
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

	constAgreeSingleLimit = (constPlayerBufferLimit >> 1) - constAgreeHeader
)

var (
	errProtoIllegal = errors.New("An illegal agreement")
)

func extAgreeAnalysis(data *bytes.Buffer) (string, []byte, error) {

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
		return "", nil, errProtoIllegal
	}

	pname := string(data.Next(int(nl)))
	pdata := data.Next(int(dl))

	return pname, pdata, nil
}
