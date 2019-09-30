package module

import (
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/service/implement"
	"github.com/yamakiller/magicNet/timer"
	"github.com/yamakiller/magicNet/util"
)

//Dispatcher Event distributor
type Dispatcher func(event *InNetMethodClientEvent)

//SpawnInNetClient  Produce an intranet client object
func SpawnInNetClient() InNetClient {
	return InNetClient{}
}

//IInNetClient Intranet client interface
type IInNetClient interface {
	implement.INetClient
	UnPacket(name string, data []byte) interface{}
}

//InNetClient Intranet client base class
type InNetClient struct {
	implement.NetClient
	service.Service
	Dispatch Dispatcher

	handle util.NetHandle
	sock   int32
}

//Init Client service initialization
func (inc *InNetClient) Init() {
	inc.Service.Init()
	inc.RegisterMethod(&network.NetChunk{}, inc.OnRecv)
	inc.RegisterMethod(&actor.Stopped{}, inc.Stoped)
}

//SetID Setting the client ID
func (inc *InNetClient) SetID(h uint64) {
	inc.handle.SetValue(h)
}

//GetID Returns the client ID
func (inc *InNetClient) GetID() uint64 {
	return inc.handle.GetValue()
}

//GetSocket Returns the client socket
func (inc *InNetClient) GetSocket() int32 {
	return inc.sock
}

//SetSocket Setting the client socket
func (inc *InNetClient) SetSocket(sock int32) {
	inc.sock = sock
}

//GetAuth return to certification time
func (inc *InNetClient) GetAuth() uint64 {
	return 0
}

//SetAuth Setting the time for authentication
func (inc *InNetClient) SetAuth(v uint64) {
}

//GetKeyPair Return key object
func (inc *InNetClient) GetKeyPair() interface{} {
	return nil
}

//BuildKeyPair Build key pair
func (inc *InNetClient) BuildKeyPair() {

}

//GetKeyPublic Return key publicly available information
func (inc *InNetClient) GetKeyPublic() string {
	return ""
}

//OnRecv Overloaded data reception
func (inc *InNetClient) OnRecv(context actor.Context, message interface{}) {
	wrap := message.(*network.NetChunk)
	if inc.GetSocket() != wrap.Handle {
		inc.LogError("OnRecv: Illegal socket, data drop:%d-%d", inc.GetSocket(), wrap.Handle)
		return
	}

	var (
		space  int
		writed int
		wby    int
		pos    int

		err error
	)

	for {
		space = inc.GetRecvBuffer().Cap() - inc.GetRecvBuffer().Len()
		wby = len(wrap.Data) - writed
		if space > 0 && wby > 0 {
			if space > wby {
				space = wby
			}

			_, err = inc.GetRecvBuffer().Write(wrap.Data[pos : pos+space])
			if err != nil {
				inc.LogError("OnRecv: error %+v socket %d", err, wrap.Handle)
				network.OperClose(wrap.Handle)
				break
			}

			pos += space
			writed += space

			inc.GetStat().UpdateRead(timer.Now(), uint64(space))
		}

		for {
			// Decomposition of Packets
			err = inc.analysis(context)
			if err != nil {
				if err == implement.ErrAnalysisSuccess {
					continue
				} else if err != implement.ErrAnalysisProceed {
					inc.LogError("OnRecv: error %+v socket %d closing client", err, wrap.Handle)
					network.OperClose(wrap.Handle)
					return
				}
			}

			if writed >= len(wrap.Data) {
				return
			}
			break
		}
	}
}

//Stoped xxx
func (inc *InNetClient) Stoped(context actor.Context, message interface{}) {
	inc.LogDebug("Stoped: Socket-%d", inc.GetSocket())
	inc.SetSocket(0)
	inc.Service.Stoped(context, message)
}

func (inc *InNetClient) analysis(context actor.Context) error {
	name, h, wrap, err := agreement.AgentParser(agreement.ConstInParser).Analysis(inc.GetKeyPair(),
		inc.GetRecvBuffer())
	if err != nil {
		return err
	}

	if wrap == nil {
		return implement.ErrAnalysisProceed
	}

	inc.Dispatch(&InNetMethodClientEvent{
		Context: context,
		InNetMethodEvent: InNetMethodEvent{
			H: h,
			NetMethodEvent: implement.NetMethodEvent{
				Name:   name,
				Socket: inc.GetSocket(),
				Wrap:   wrap,
			},
		},
	})
	return implement.ErrAnalysisSuccess
}

//Decode Decoding method
func (inc *InNetClient) decode(msgType reflect.Type, data []byte) (interface{}, error) {
	wrap := reflect.Indirect(reflect.New(msgType.Elem())).Addr().Interface().(proto.Message)
	err := proto.Unmarshal(data, wrap)
	if err != nil {
		return nil, err
	}

	return wrap, nil
}

//UnPacket 解包
func (inc *InNetClient) UnPacket(name string, data []byte) interface{} {
	msgType := proto.MessageType(name)
	if msgType == nil {
		return nil
	}
	wrap, err := inc.decode(msgType, data)
	if err != nil {
		return nil
	}
	return wrap
}

//Shutdown Terminate this client service
func (inc *InNetClient) Shutdown() {
	inc.Service.Shutdown()
}
