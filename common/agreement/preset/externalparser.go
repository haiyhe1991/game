package preset

import (
	"bytes"

	"github.com/yamakiller/magicNet/network/parser"
)

//ExternalParser External analysis
type ExternalParser struct {
}

// Analysis Play protocol data analysis
//------------------------------------------------------------------------------------------
//  7 Bit Version | 12 Bit data length | 5 Bit AgreementName | AgreementName | data packet |
//------------------------------------------------------------------------------------------
func (epr *ExternalParser) Analysis(keyPair interface{},
	data *bytes.Buffer) (string, uint64, []byte, error) {

	if data.Len() < constAgreeHeader {
		return "", 0, nil, nil
	}

	head := getAgreementHeader(data)

	dl := getAgreementDataLength(head)
	nl := getAgreementNameLength(head)

	if (dl + nl) > uint32(data.Len()) {
		data.UnreadByte()
		data.UnreadByte()
		data.UnreadByte()
		return "", 0, nil, nil
	}

	if int((dl + nl)) > AgreeSingleLimit {
		data.Reset()
		return "", 0, nil, ErrProtoIllegal
	}

	pname := string(data.Next(int(nl)))
	pdata := data.Next(int(dl))

	return pname, 0, pdata, nil
}

// Assemble Assembly Protocol Data
func (epr *ExternalParser) Assemble(keyPair interface{},
	version int32,
	handle uint64,
	agreementName string,
	data []byte,
	length int32) []byte {

	nameLength := int32(len([]rune(agreementName)))
	buffer := bytes.NewBuffer([]byte{})

	buffer.Write(assembleHeader(version, length, nameLength))
	buffer.WriteString(agreementName)
	buffer.Write(data)

	return buffer.Bytes()
}

var _ parser.IParser = &ExternalParser{}
