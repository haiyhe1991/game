package gateway

type protoRegister struct {
	proto interface{}
	route string
	auth  bool
}
