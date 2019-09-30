package preset

import (
	"bytes"
	"encoding/binary"
	"errors"
)

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
	constAgreeHandleBit       = 64
	constAgreeHandle          = 8
)

var (
	//ErrProtoIllegal An illegal agreement
	ErrProtoIllegal = errors.New("An illegal agreement")
)

var (
	//AgreeSingleLimit External communication single packet size limit
	AgreeSingleLimit = 2048
	//InsideAgreeSingleLimit Internal communication packet size limit
	InsideAgreeSingleLimit = 2056
)

// SetSingleLimit Set the communication ticket size limit
func SetSingleLimit(singleLimit int) {
	AgreeSingleLimit = singleLimit - constAgreeHeader
	InsideAgreeSingleLimit = AgreeSingleLimit + constAgreeHandle
}

func getAgreementHeader(data *bytes.Buffer) uint32 {
	headByte := make([]byte, 4)
	headByte[0], _ = data.ReadByte()
	headByte[1], _ = data.ReadByte()
	headByte[2], _ = data.ReadByte()
	headByte[3] = 0

	head := binary.BigEndian.Uint32(headByte)
	return head
}

func getAgreementDataLength(header uint32) uint32 {
	return (header >> constAgreeDataLengthShift) & constAgreeDataLengthMask
}

func getAgreementNameLength(header uint32) uint32 {
	return (header >> constAgreeNameLengthShift) & constAgreeNameLengthMask
}

func assembleHeader(version int32, dataLength int32, nameLength int32) []byte {
	var headInt uint32
	headInt = uint32((version&constAgreeVersionMask)<<constAgreeVersionShift) |
		uint32((dataLength&constAgreeDataLengthMask)<<constAgreeDataLengthShift) |
		uint32((nameLength&constAgreeNameLengthMask)<<constAgreeNameLengthShift)

	head := make([]byte, 4)
	binary.BigEndian.PutUint32(head, headInt)

	return head[:3]
}
