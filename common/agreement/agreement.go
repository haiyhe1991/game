package agreement

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

//CertificationConfirmation Confirm that the login has been successful and change the verification status of the connection.
type CertificationConfirmation struct {
	Handle uint64
}

func getAgreementHeader(data *bytes.Buffer) uint32 {
	headByte := make([]byte, 4)
	headByte[0] = 0
	headByte[1], _ = data.ReadByte()
	headByte[2], _ = data.ReadByte()
	headByte[3], _ = data.ReadByte()

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
	return head[1:]
}

//InsideAnalysis Inside protocol data analysis
//-----------------------------------------------------------------------------------------------------------
//  7 Bit Version | 12 Bit data length | 5 Bit AgreementName | 64 Bit Handle | AgreementName | data packet |
//-----------------------------------------------------------------------------------------------------------
func InsideAnalysis(data *bytes.Buffer) (string, uint64, []byte, error) {
	if data.Len() < (constAgreeHeader + constAgreeHandle) {
		return "", 0, nil, nil
	}

	head := getAgreementHeader(data)

	dl := getAgreementDataLength(head)
	nl := getAgreementNameLength(head)

	if (dl + nl + constAgreeHandle) > uint32(data.Len()) {
		data.UnreadByte()
		data.UnreadByte()
		data.UnreadByte()
		return "", 0, nil, nil
	}

	if int((dl + nl + constAgreeHandle)) > InsideAgreeSingleLimit {
		data.Reset()
		return "", 0, nil, ErrProtoIllegal
	}

	pname := string(data.Next(int(nl)))
	pdata := data.Next(int(dl))
	phandle := data.Next(int(constAgreeHandle))

	return pname, binary.BigEndian.Uint64(phandle), pdata, nil
}

//InsideAssemble Assembly Inside protocol
func InsideAssemble(version int32, handle uint64, agreementName string, data []byte, length int32) []byte {
	nameLength := int32(len([]rune(agreementName)))
	buffer := bytes.NewBuffer([]byte{})
	handleBuffer := make([]byte, 8)

	buffer.Write(assembleHeader(version, length, nameLength))
	buffer.WriteString(agreementName)
	buffer.Write(data)

	binary.BigEndian.PutUint64(handleBuffer, handle)
	buffer.Write(handleBuffer)

	return buffer.Bytes()
}

// ExtAnalysis Play protocol data analysis
//------------------------------------------------------------------------------------------
//  7 Bit Version | 12 Bit data length | 5 Bit AgreementName | AgreementName | data packet |
//------------------------------------------------------------------------------------------
func ExtAnalysis(data *bytes.Buffer) (string, []byte, error) {

	if data.Len() < constAgreeHeader {
		return "", nil, nil
	}

	head := getAgreementHeader(data)

	dl := getAgreementDataLength(head)
	nl := getAgreementNameLength(head)

	if (dl + nl) > uint32(data.Len()) {
		data.UnreadByte()
		data.UnreadByte()
		data.UnreadByte()
		return "", nil, nil
	}

	if int((dl + nl)) > AgreeSingleLimit {
		data.Reset()
		return "", nil, ErrProtoIllegal
	}

	pname := string(data.Next(int(nl)))
	pdata := data.Next(int(dl))

	return pname, pdata, nil
}

// ExtAssemble Assembly Protocol Data
func ExtAssemble(version int32, agreementName string, data []byte, length int32) []byte {
	nameLength := int32(len([]rune(agreementName)))
	//count := constAgreeHeader + nameLength + length
	buffer := bytes.NewBuffer([]byte{})

	buffer.Write(assembleHeader(version, length, nameLength))
	buffer.WriteString(agreementName)
	buffer.Write(data)

	return buffer.Bytes()
}
