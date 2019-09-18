package preset

import (
	"bytes"
	"encoding/binary"

	"github.com/yamakiller/magicNet/network/parser"
)

//InsideParser Internal analysis
type InsideParser struct {
}

//Analysis Inside protocol data analysis
//-----------------------------------------------------------------------------------------------------------
//  7 Bit Version | 12 Bit data length | 5 Bit AgreementName | 64 Bit Handle | AgreementName | data packet |
//-----------------------------------------------------------------------------------------------------------
func (ipr *InsideParser) Analysis(data *bytes.Buffer) (string, uint64, []byte, error) {
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

//Assemble Assembly Inside protocol
func (ipr *InsideParser) Assemble(version int32, handle uint64, agreeName string, data []byte, length int32) []byte {
	nameLength := int32(len([]rune(agreeName)))
	buffer := bytes.NewBuffer([]byte{})
	handleBuffer := make([]byte, 8)

	buffer.Write(assembleHeader(version, length, nameLength))
	buffer.WriteString(agreeName)
	buffer.Write(data)

	binary.BigEndian.PutUint64(handleBuffer, handle)
	buffer.Write(handleBuffer)

	return buffer.Bytes()
}

var _ parser.IParser = &InsideParser{}
