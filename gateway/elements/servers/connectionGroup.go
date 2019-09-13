package servers

//ConnectionGroup Providing packet management for service connection
type ConnectionGroup struct {
	group []Connection
}

//Register Register a Service Connection Object
func (cgs *ConnectionGroup) Register(id int32, addr string) {
	for i := 0; i < len(cgs.group); i++ {
		if cgs.group[i].ID == id {
			cgs.group[i].Addr = addr
			return
		}
	}

	cgs.group = append(cgs.group, Connection{ID: id, Addr: addr})
}

//FindSocket Find connection service based on SOCKET
func (cgs *ConnectionGroup) FindSocket(sock int32) *Connection {
	for i := 0; i < len(cgs.group); i++ {
		if cgs.group[i].Sock == sock {
			return &cgs.group[i]
		}
	}
	return nil
}

//HashConnection Get service based on hash ? I hope to be optimized later
func (cgs *ConnectionGroup) HashConnection(id int32) *Connection {
	hash := (id % int32(len(cgs.group)))
	return &cgs.group[hash]
}
