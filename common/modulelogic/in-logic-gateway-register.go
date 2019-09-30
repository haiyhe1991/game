package modulelogic

import (
	"github.com/gogo/protobuf/proto"
	"github.com/yamakiller/game/common/agreement"
	"github.com/yamakiller/game/common/module"
	"github.com/yamakiller/game/pactum"
	"github.com/yamakiller/magicNet/engine/logger"
	"github.com/yamakiller/magicNet/network"
	"github.com/yamakiller/magicNet/service"
	"github.com/yamakiller/magicNet/service/implement"
)

//GatewayRegisterProc Gateway registration process
type GatewayRegisterProc struct {
}

//OnProccess Gateway registration method
func (grc *GatewayRegisterProc) OnProccess(event implement.INetMethodEvent) {
	request := event.(*module.InNetMethodClientEvent)
	c := module.GetWarehouse().GrapSocket(request.Socket)
	if c == nil {
		logger.Error(request.Context.Self().GetID(), "Exception to find the target connection")
		return
	}
	defer module.GetWarehouse().Release(c)
	grc.defaultProccess(c, request)
}

func (grc *GatewayRegisterProc) defaultProccess(c interface{}, request *module.InNetMethodClientEvent) {
	client := c.(module.IInNetClient)
	srvClient := c.(service.IService)

	iwrap := client.UnPacket(request.Name, request.Wrap)
	if iwrap == nil {
		srvClient.LogError("Decoding exception will close the connection")
		network.OperClose(request.Socket)
		return
	}

	wrap, success := iwrap.(*pactum.GatewayRegisterRequest)
	if !success {
		srvClient.LogError("Decoding exception, unable to convert the encoding to the" +
			"corresponding protocol object, will close the connection")
		network.OperClose(request.Socket)
		return
	}

	client.SetID(uint64(wrap.GetId()))
	oldClient := module.GetWarehouse().Register(request.Socket, wrap.GetId())
	if oldClient != nil {
		network.OperClose(oldClient.GetSocket())
		module.GetWarehouse().Release(oldClient)
		srvClient.LogDebug("Old Gateway connection exists, need to kick off, old connection")
	}

	//Reply packet
	replyData, err := proto.Marshal(&pactum.GatewayRegisterResponse{Code: 0, Message: "Success"})
	if err != nil {
		srvClient.LogError("Reply code failed, connection will be closed")
		network.OperClose(request.Socket)
		return
	}

	replyData = agreement.AgentParser(agreement.ConstInParser).Assemble(client.GetKeyPair(),
		agreement.ConstPactumVersion,
		request.H,
		"pactum.GatewayRegisterResponse",
		replyData,
		int32(len(replyData)))

	if err := network.OperWrite(request.Socket, replyData, len(replyData)); err != nil {
		srvClient.LogError("Reply code failed, connection will be closed:%+v", err)
		network.OperClose(request.Socket)
		return
	}

	srvClient.LogDebug("Gateway %d Register Complete", wrap.GetId())
}
