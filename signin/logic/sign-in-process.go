package logic

<<<<<<< .mineimport (
	"fmt"
=======import (
	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/common/module"
	"github.com/yamakiller/game/pactum"
	"github.com/yamakiller/magicNet/engine/actor"
	"github.com/yamakiller/magicNet/network"
	"log"
	"reflect"
	"time"
)
>>>>>>> .theirs
<<<<<<< .mine	"github.com/yamakiller/magicNet/service/implement"
)

=======func Decode(msgType reflect.Type, data []byte) (interface{}, error) {

	wrap := reflect.Indirect(reflect.New(msgType.Elem())).Addr().Interface().(proto.Message)
	err := proto.Unmarshal(data, wrap)
	if err != nil {
		return nil, err
	}
	return wrap, nil
}

>>>>>>> .theirs//SignInProc Log in to the logical processor
type SignInProc struct {
}

//OnProccess Processing logic
<<<<<<< .minefunc (sip *SignInProc) OnProccess(event implement.INetMethodEvent) {
	fmt.Println("Sign-in Proccess")
=======func (sip *SignInProc) OnProccess(context actor.Context, message interface{}) {
	event := message.(*module.InNetClientChunkEvent)
	iwrap, err := Decode(event.Type, event.Wrap)
	if err != nil {
		log.Println("Decoding exception will close the connection")
		network.OperClose(event.Socket)
		return
	}

>>>>>>> .theirs	wrap, success := iwrap.(*pactum.LoginRequest)
	if !success {
		log.Println("Decoding exception, unable to convert the encoding to the" +
			"corresponding protocol object, will close the connection")
		network.OperClose(event.Socket)
		return
	}
	//Reply packet
	//pactum.Response{State:0,Message:"ok"}

	var(
		response *pactum.LoginResponse
	)

	if response,err = sip.handlerLogin(wrap.Account,wrap.Password,wrap.Origin);err != nil {
		log.Println(err)
		return
	}

	replyData, err := proto.Marshal(response)
	if err != nil {
		log.Println("Reply code failed, connection will be closed")
		network.OperClose(event.Socket)
		return
	}

	agreement.AgentParser(agreement.ConstInParser).Assemble(nil,
		agreement.ConstPactumVersion,
		event.Handle,
		"pactum.GatewayRegisterResponse", replyData, int32(len(replyData)))
}

//检查账号密码 如何已经登录剔除旧用户
func (sip *SignInProc) handlerLogin(account, password, origin string) (response *pactum.LoginResponse, err error) {

	var (
		info     = &UserInfo{}
	)
	response = &pactum.LoginResponse{}
	response.Rep.State = 0
	response.Rep.Message = "ok"

	defer func() {
		if err != nil {
			response.Rep.State = 1
			response.Rep.Message = err.Error()
		}
	}()

	if CheckUserSession(account) == true {
		//其他地方登录强制下线  更新登录时间
	}

	if info, err = GetUserByAccount(account); err != nil {
		log.Println(err)
		return
	}
	//用户密码md5
	if info.Password != password {
		err = ErrLoginFailed
		return
	}

	if err = SetUserSession(account, &UserSession{State: UserStateLogin, Time: time.Now().Unix()}); err != nil {
		log.Println(err)
		return
	}

	return
}
