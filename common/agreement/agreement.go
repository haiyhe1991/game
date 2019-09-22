package agreement

import (
	"github.com/yamakiller/game/common/agreement/preset"
	"github.com/yamakiller/magicNet/network/parser"
)

const (
	//ConstArgeeVersion Protocol version number
	ConstArgeeVersion = 1
)

//ForwardEvent Forward data event
type ForwardEvent struct {
	Handle     uint64
	PactunName string
	ServoName  string
	Data       []byte
}

//CheckConnectMessage Check the connection status event to achieve automatic reconnection after disconnection
type CheckConnectMessage struct {
}

//CertificationConfirmation Confirm that the login has been successful and change the verification status of the connection.
type CertificationConfirmation struct {
	Handle uint64
}

const (
	//ConstInParser xxx
	ConstInParser = 0
	//ConstExParser xxx
	ConstExParser = 1
)

var (
	inParser = &preset.InsideParser{}
	exParser parser.IParser
)

// SetParser Setting up parsers
func SetParser(pser parser.IParser) {
	exParser = pser
}

//AgentParser Provide Protocol Resolution Agent
func AgentParser(mode int) parser.IParser {
	switch mode {
	case 0:
		return inParser
	case 1:
		if exParser == nil {
			exParser = &preset.ExternalParser{}
		}
		return exParser
	}
	return nil
}
