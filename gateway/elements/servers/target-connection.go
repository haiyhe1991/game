package servers

//TargetConnection Connection to the target server configuration information
type TargetConnection struct {
	ID        int32
	VirtaulID uint32
	Socket    int32
	Name      string
	Addr      string
	Desc      string
}
