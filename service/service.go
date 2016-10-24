package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"rpc-go/codec"
	"rpc-go/config"
	"rpc-go/service/register"
	"rpc-go/transport"
	"strings"
)

//JumeiTCPService 聚美 的tcp rpc 服务
type JumeiTCPService struct {
	Name       string
	reg        register.Register
	Registered bool
	listener   net.Listener
}

//New
func NewService() *JumeiTCPService {
	var newService JumeiTCPService
	//	newService.Handlers = make(map[string]HandlerFunction)
	return &newService
}

var logger *log.Logger

//Init servicename 是为下一步分布式rpc留的服务名称
//port 是此rpc的服务端口
func (srvs *JumeiTCPService) Init(servicename, netType, address string) {
	logPath := fmt.Sprintf("%s.log", servicename)
	file, err := os.OpenFile(logPath, os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err.Error())
	}
	logger = log.New(file, "", log.LstdFlags|log.Llongfile)
	// logger = log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)
	listener, linstenErr := net.Listen(netType, address)
	if linstenErr != nil {
		logger.Println("server err:", linstenErr.Error())
		return
	}

	srvs.listener = listener
	out := fmt.Sprintf("server listen at port:%v", listener.Addr())
	logger.Println(out)
	fmt.Printf("start %s and ", servicename)
}

//Run rpc 服务开始运行
func (srvs *JumeiTCPService) Run() {
	defer srvs.listener.Close()
	fmt.Println("run it  at ", srvs.listener.Addr())
	for {
		conn, connErr := srvs.listener.Accept()
		fmt.Println(len(transport.Connections.ConnMap), config.ServiceConf.MaxConnection)
		if len(transport.Connections.ConnMap) >= config.ServiceConf.MaxConnection {
			response, _ := codec.WrapS2CData(503, fmt.Sprintf("Maximum number of connections is %d ", config.ServiceConf.MaxConnection))
			conn.Write([]byte(response))
			continue
		}
		if connErr == nil {
			go srvs.ServerHandleConn(conn)
		} else {
			logger.Println(connErr.Error())
			os.Exit(1)
		}

	}
}

//
func (srvs *JumeiTCPService) ErrorHandler(conn *transport.JumeiConn) {
	if errstr := recover(); errstr != nil {
		outErr := fmt.Sprintf("%v", errstr)
		conn.S2CSendError(outErr)
		logger.Println(outErr)
	}
	conn.CloseConn()
}

func (srvs *JumeiTCPService) ServerHandleConn(conn net.Conn) {
	logger.Printf("client %s come in:", conn.RemoteAddr())
	jumeiConn := new(transport.JumeiConn)
	jumeiConn.Conn = conn
	transport.AddConnection(jumeiConn)
	defer srvs.ErrorHandler(jumeiConn)
	datachan := make(chan transport.JumeiTextRPC, 0)

	go transport.ReceiveClientData(jumeiConn, datachan)
	for {
		jumeiTextRpc := <-datachan
		switch jumeiTextRpc.Command {
		case "RPC":
			rpcData, err := codec.UnwrapRPC(jumeiTextRpc.Data)
			if err == nil {
				var rpcParam codec.RPCParam
				json.Unmarshal([]byte(rpcData.Data), &rpcParam)
				rpcParam.Class = strings.Replace(rpcParam.Class, "\\", "/", -1)
				callFunc := register.Reg.Handlers[rpcParam.Class+"."+rpcParam.Method]
				if callFunc == nil {
					notExist := fmt.Sprintf("RPC Method:%s not exist", rpcParam.Class+"."+rpcParam.Method)
					jumeiConn.S2CSendError(notExist)
					logger.Println(notExist)
					continue
				}

				callFunc(jumeiConn, rpcParam.Params)
				//短连接，发送完response 后即刻关闭连接
				return
			} else {
				err = jumeiConn.S2CSendError(err.Error())
				logger.Printf("decode in jumeiTextRPC err: %s", err.Error())
			}

		default:
			logger.Printf("command must be 'RPC' ,now is %s .", jumeiTextRpc.Command)
		}
	}
}
