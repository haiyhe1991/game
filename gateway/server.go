package gateway

import (
	"bytes"
	"sync"
)

const (
	constConnectMax        = 128
	constConnectChanMax    = 256
	constConnectPushErrMax = 16
)

type connector struct {
	id   int32 //Unique in the cluster
	sock int32
	data *bytes.Buffer
	stat remoteStat
	addr string
	sync sync.Mutex
}

type singleService struct {
	cs []connector
}

func (sse *singleService) get(id int32) *connector {
	// if id == 0 时间 随机分布，暂时不支持此功能
	hash := id & (constConnectMax - 1)
	return &sse.cs[hash]
}

type servers struct {
	ms map[string]singleService
}

func (srv *servers) registerService(name string) {
	srv.ms[name] = singleService{cs: make([]connector, constConnectMax)}
}

func (srv *servers) reginsterConnector(name string, id int32, addr string) {
	if _, ok := srv.ms[name]; !ok {
		return
	}
	hash := id & (constConnectMax - 1)
	pc := &srv.ms[name].cs[hash]
	pc.id = id
	pc.data = bytes.NewBuffer([]byte{})
	pc.addr = addr
}

func (srv *servers) get(name string) *singleService {
	if s, ok := srv.ms[name]; ok {
		return &s
	}
	return nil
}

//? 需要优化
func (srv *servers) getConnector(sock int32) *connector {
	for _, v := range srv.ms {
		for i := 0; i < constConnectMax; i++ {
			if v.cs[i].sock == sock {
				return &v.cs[i]
			}
		}
	}
	return nil
}
