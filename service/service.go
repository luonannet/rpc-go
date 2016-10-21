package service

import (
	"encoding/json"
	"log"
	"net"
	"os"
	"rpc-go/codec"
	"rpc-go/transport"
	"strings"
)

//rpc 的执行体
type HandlerFunction func(conn *transport.JumeiConn, request interface{})

//JumeiTCPService 聚美 的tcp rpc 服务
type JumeiTCPService struct {
	Name       string
	Handlers   map[string]HandlerFunction
	Registered bool
	listener   net.Listener
}

//New
func NewService() *JumeiTCPService {
	var newService JumeiTCPService
	newService.Handlers = make(map[string]HandlerFunction)
	return &newService
}

//Init servicename 是为下一步分布式rpc留的服务名称
//port 是此rpc的服务端口
func (srvs *JumeiTCPService) Init(servicename, netType, address string) {

	listener, linstenErr := net.Listen(netType, address)
	if linstenErr != nil {
		log.Println("server err:", linstenErr.Error())
		return
	}
	srvs.listener = listener
	log.Println("server listen at port:", listener.Addr())

}

//Run rpc 服务开始运行
func (srvs *JumeiTCPService) Run() net.Listener {

	for {
		conn, connErr := srvs.listener.Accept()
		if connErr == nil {
			go srvs.ServerHandleConn(conn)
		} else {
			log.Println(connErr.Error())
			os.Exit(1)
		}
	}

}

func (service *JumeiTCPService) ServerHandleConn(conn net.Conn) {
	log.Println("client come in:", conn.RemoteAddr())
	jumeiConn := new(transport.JumeiConn)
	jumeiConn.Conn = conn
	datachan := make(chan transport.JumeiTextRPC, 10)
	go transport.ReceiveData(jumeiConn.Conn, datachan)
	for {
		jumeiTextRpc := <-datachan
		switch jumeiTextRpc.Command {
		case "RPC":
			rpcData, err := codec.UnwrapRPC(jumeiTextRpc.Data)
			if err == nil {
				var rpcParam codec.RPCParam
				json.Unmarshal([]byte(rpcData.Data), &rpcParam)
				rpcParam.Class = strings.Replace(rpcParam.Class, "\\", "/", -1)
				callFunc := service.Handlers[rpcParam.Class+"."+rpcParam.Method]
				if callFunc == nil {
					log.Printf("RPC Method:%s not exist", rpcParam.Class+"."+rpcParam.Method)
					continue
				}
				callFunc(jumeiConn, rpcParam.Params)
			} else {
				log.Println("ServerHandleConn err:", err.Error())
			}
		default:
			log.Printf("command must be 'RPC' ,now is %s:", jumeiTextRpc.Command)
		}

	}
}
