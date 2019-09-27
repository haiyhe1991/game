package test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"

	"github.com/gogo/protobuf/proto"

	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/pactum"
	"github.com/yamakiller/magicNet/util"
)

type SignClientTest struct {
	ID     int32
	client *Client
}

func (sc *SignClientTest) Analy(b *bytes.Buffer) (string, uint64, []byte, error) {
	return agreement.AgentParser(agreement.ConstInParser).Analysis(nil, b)
}

func (sc *SignClientTest) onHandshake(message interface{}) {
	event := message.(*NetEventHandle)

	wrap := pactum.HandshakeResponse{}
	if err := proto.Unmarshal(event.Data, &wrap); err != nil {
		panic(err)
	}

	reqWrap := pactum.GatewayRegisterRequest{}
	reqWrap.Id = sc.ID

	reqWrapData, _ := proto.Marshal(&reqWrap)

	h := util.NetHandle{}

	h.Generate(sc.ID, 123, 456)

	reqWrapData = agreement.AgentParser(agreement.ConstInParser).Assemble(nil,
		1,
		h.GetValue(),
		"pactum.GatewayRegisterRequest",
		reqWrapData,
		int32(len(reqWrapData)))

	sc.client.Write(reqWrapData)
}

func (sc *SignClientTest) onRegisterResponse(messsage interface{}) {
	//event := message.(*NetEventHandle)
}

func runClient(w *sync.WaitGroup) {
	defer w.Done()
	c := &SignClientTest{client: NewClientTest()}
	c.client.Analysis = c.Analy
	c.client.RegisterMethod(&pactum.HandshakeResponse{}, c.onHandshake)
	c.client.RegisterMethod(&pactum.GatewayRegisterResponse{}, c.onRegisterResponse)
	if err := c.client.Connect("127.0.0.1:7852"); err != nil {
		fmt.Printf("连接服务器失败：%+v\n", err)
		return
	}

	fmt.Printf("连接服务器成功\n")
	w.Wait()
	fmt.Printf("退出连接\n")
}

func Test_VirtaulGateway(t *testing.T) {
	var clientNum = 1
	var w sync.WaitGroup
	w.Add(clientNum)
	for icon := 0; icon < clientNum; icon++ {
		go runClient(&w)
	}

	w.Wait()

	fmt.Printf("测试结束\n")
}
