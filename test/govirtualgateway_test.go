package test

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/pactum"
	"github.com/yamakiller/magicNet/util"
)

type SignClientTest struct {
	ID     int32
	Auth   bool
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
	evt := messsage.(*NetEventHandle)
	//msg := evt.Data
	msg := pactum.GatewayRegisterResponse{}
	if err := proto.Unmarshal(evt.Data, &msg); err != nil {
		fmt.Printf("注册回复数据失败:%+v\n", err)
	}

	sc.Auth = true

	fmt.Printf("注册成功:%+v\n", msg)

}

func (sc *SignClientTest) onSendSignInRequest() {
	reqWrap := pactum.LoginRequest{}
	reqWrap.Account = "test_" + strconv.Itoa(int(sc.ID))
	reqWrap.Password = reqWrap.Account
	reqWrap.Origin = "google"

	h := util.NetHandle{}
	h.Generate(sc.ID, 123, 456)

	reqWrapData, _ := proto.Marshal(&reqWrap)

	reqWrapData = agreement.AgentParser(agreement.ConstInParser).Assemble(nil,
		1,
		h.GetValue(),
		"pactum.LoginRequest",
		reqWrapData,
		int32(len(reqWrapData)))

	sc.client.Write(reqWrapData)
}

func runClient(w *sync.WaitGroup, serID int32) {
	defer w.Done()
	c := &SignClientTest{client: NewClientTest()}
	c.ID = serID
	c.client.Analysis = c.Analy
	c.client.RegisterMethod(&pactum.HandshakeResponse{}, c.onHandshake)
	c.client.RegisterMethod(&pactum.GatewayRegisterResponse{}, c.onRegisterResponse)
	if err := c.client.Connect("127.0.0.1:7852"); err != nil {
		fmt.Printf("连接服务器失败:%+v\n", err)
		return
	}

	ick := 0
	for {
		if c.Auth {
			c.onSendSignInRequest()
			ick++
			if ick > 100 {
				break
			}
		}
		time.Sleep(time.Duration(100) * time.Millisecond)
	}

	c.client.Shutdown()
	c.client.Wait()
}

func Test_VirtaulGateway(t *testing.T) {
	var clientNum = 1
	var w sync.WaitGroup
	w.Add(clientNum)
	for icon := 0; icon < clientNum; icon++ {
		go runClient(&w, int32(clientNum+100))
	}

	w.Wait()

	fmt.Printf("测试结束\n")
}
