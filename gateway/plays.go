package gateway

type plays struct {
	csocks map[uint64]*client
	cplays map[uint64]*client
}
