package common

const (
	//ConstClientBufferLimit External network client buffer cap
	ConstClientBufferLimit = 4096
	//ConstInClientBufferLimit Intranet client read buffer cap
	ConstInClientBufferLimit = ConstClientBufferLimit + 16
)
