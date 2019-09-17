package servers

import (
	"strconv"

	"github.com/yamakiller/magicNet/util"
)

//ConnectionGroup Providing packet management for service connection
type ConnectionGroup struct {
	g *util.Consistent
}

//Register Register a Service Connection Object
func (cgs *ConnectionGroup) Register(name string, id int32, addr string) {
	key := name + strconv.Itoa(int(id))
	c, err := cgs.g.Get(key)
	if err != nil {
		c.(*Connection).Addr = addr
		return
	}


	cgs.g.Push(key, &Connection{ID: id, Addr: addr})
}

//FindSocket Find connection service based on SOCKET
func (cgs *ConnectionGroup) FindSocket(sock int32) *Connection {

	f := func(key interface{}, val interface{}) int {
		if val.(*Connection).Sock == val.(int32) {
			return 0
		}
		return -1
	}

	reulst := cgs.g.Sreach(sock, f)
	if reulst == nil {
		return nil
	}
	return reulst.(*Connection)
}

//HashConnection Get service based on hash ? I hope to be optimized later
func (cgs *ConnectionGroup) HashConnection(name string, id int32) (*Connection, error) {
	v, err := cgs.g.Get(name + strconv.Itoa(int(id)))
	if err != nil {
		return nil, err
	}

	return v.(*Connection), nil
}
