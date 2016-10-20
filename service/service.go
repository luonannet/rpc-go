package service

import (
	"encoding/json"
	"log"
	"net"
	"rpc-go/codec"
	"rpc-go/transport"
	"strings"
	"time"
)

const (
	RPC_Client_Prefix string = "RpcClient_"
)

type JumeiTCPService struct {
	Name       string
	Handlers   map[string]HandlerFunction
	Registered bool
}

type HandlerFunction func(conn *transport.JumeiConn, request interface{})

func NewService() *JumeiTCPService {
	var newService JumeiTCPService
	newService.Handlers = make(map[string]HandlerFunction)
	return &newService
}

func (service *JumeiTCPService) Init(servicename string) {
	service.Name = servicename
}

func (service *JumeiTCPService) RegisterHandler(handname string, handFunc HandlerFunction) {

	handname = RPC_Client_Prefix + handname
	service.Handlers[handname] = handFunc
}

func (service *JumeiTCPService) ServerHandleConn(conn net.Conn) {
	log.Println("client come in:", conn.RemoteAddr())
	jumeiConn := new(transport.JumeiConn)
	jumeiConn.Conn = conn
	//短连接
	jumeiConn.Conn.SetDeadline(time.Now().Add(time.Second * 1))
	datachan := make(chan transport.JumeiTextRPC, 0)
	go transport.ReceiveData(jumeiConn.Conn, datachan)
	for {
		jumeiTextRpc := <-datachan
		switch jumeiTextRpc.Command {
		case "RPC":
			rpcData, err := codec.UnwrapRPC(jumeiTextRpc.Data)
			var rpcParam codec.RPCParam
			json.Unmarshal([]byte(rpcData.Data), &rpcParam)
			if err == nil {
				rpcParam.Class = strings.Replace(rpcParam.Class, "\\", "/", -1)
				callFunc := service.Handlers[rpcParam.Class+"."+rpcParam.Method]
				callFunc(jumeiConn, rpcParam.Params)
			} else {
				log.Println("ServerHandleConn err:", err.Error())
			}
		default:
		}

	}
}
