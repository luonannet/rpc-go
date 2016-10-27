package service

import (
	"bytes"
	"compress/zlib"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"rpc-go/codec"
	"rpc-go/config"
	"rpc-go/service/register"
	"rpc-go/transport"
	"strings"
)

//JumeiTCPService 聚美 的tcp rpc 服务
type JumeiTCPService struct {
	conf       *config.RpcServiceConfig
	reg        register.Register
	Registered bool
	listener   net.Listener
}

//NewService rpc服务
func NewService(conf *config.RpcServiceConfig) *JumeiTCPService {
	var newService JumeiTCPService
	newService.conf = conf

	listener, linstenErr := net.Listen(conf.NetType, conf.IP+conf.Port)
	if linstenErr != nil {
		config.Logger.Error("server err:", linstenErr.Error())
		return nil
	}

	newService.listener = listener

	return &newService
}

//Run rpc 服务开始运行
func (srvs *JumeiTCPService) Run() {
	defer srvs.listener.Close()
	go listenExitSignal()
	out := fmt.Sprintf("start ‘%s’ and  listen at address: %v", srvs.conf.ServiceName, srvs.listener.Addr())
	config.Logger.Critical(out)
	for {
		conn, connErr := srvs.listener.Accept()
		// 检测是否超过最大连接数限制。
		if len(transport.Connections.ConnMap) >= srvs.conf.MaxConnection {
			response, _ := codec.WrapS2CData(503, fmt.Sprintf("Maximum number of connections is %d ", srvs.conf.MaxConnection))
			conn.Write([]byte(response))
			continue
		}
		if connErr == nil {
			//正常的连接，则开启服务接收
			go srvs.ServerHandleConn(conn)
		} else {
			config.Logger.Error(connErr.Error())
			os.Exit(1)
		}
	}
}

//监听退出信号。
func listenExitSignal() {
	sign := make(chan os.Signal, 1)
	go goExit(sign)
	signal.Notify(sign, os.Kill, os.Interrupt)
}
func goExit(sign chan os.Signal) {
	select {
	case _ = <-sign:
		{
			config.Logger.Critical(" system exit ")
			config.Logger.Flush()
			config.Logger.Close()
			os.Exit(0)
		}
	}
}

//ErrorHandler 错误输出 以及断开连接
func (srvs *JumeiTCPService) ErrorHandler(conn *transport.JumeiConn) {
	if errstr := recover(); errstr != nil {
		outErr := fmt.Sprintf("%v", errstr)
		conn.S2CSendError(outErr)
		config.Logger.Error(outErr)
	}
	//无论是否有错误都关闭。
	conn.CloseConn()
}

//ServerHandleConn 正常的连接，则开启服务接收
func (srvs *JumeiTCPService) ServerHandleConn(conn net.Conn) {
	config.Logger.Infof("client %s come in \n", conn.RemoteAddr())
	jumeiConn := new(transport.JumeiConn)
	jumeiConn.Conn = conn
	transport.AddConnection(jumeiConn)
	defer srvs.ErrorHandler(jumeiConn)
	datachan := make(chan transport.JumeiTextRPC, 10)

	go transport.ReceiveClientData(jumeiConn, datachan)
	for {
		jumeiTextRPC := <-datachan
		switch jumeiTextRPC.Command {
		case "RPC:GZ":
			b := bytes.NewReader([]byte(jumeiTextRPC.Data))
			var out bytes.Buffer
			r, _ := zlib.NewReader(b)
			io.Copy(&out, r)
			jumeiTextRPC.Data = string(out.Bytes())
			fallthrough
		case "RPC":
			rpcData, err := codec.UnwrapRPC(jumeiTextRPC.Data)
			if err == nil {
				var rpcParam codec.RPCParam
				json.Unmarshal([]byte(rpcData.Data), &rpcParam)
				// 有些php调用的class 类似Example\a\b 的。先替换成
				rpcParam.Class = strings.Replace(rpcParam.Class, "\\", "/", -1)
				callFunc := register.GetHandler(rpcParam.Class + "." + rpcParam.Method)
				if callFunc == nil {
					notExist := fmt.Sprintf("RPC Method:%s not exist", rpcParam.Class+"."+rpcParam.Method)
					jumeiConn.S2CSendError(notExist)
					config.Logger.Errorf(notExist)
					continue
				}
				response, apiErr := callFunc(jumeiConn, rpcParam.Params)
				//将rpc执行结果返回到客户端
				if apiErr != nil {
					config.Logger.Errorf("rpc execute err: %s", apiErr.Error())
					err = jumeiConn.S2CSendError(apiErr.Error())
					if err != nil {
						config.Logger.Errorf("response err: %s", err.Error())
					}
				} else {
					err = jumeiConn.S2CSend(response)
					if err != nil {
						config.Logger.Errorf("response err: %s", err.Error())
					}
				}
				//因为是短连接，发送完response 后即刻关闭连接
				return
			}
			err = jumeiConn.S2CSendError(err.Error())
			config.Logger.Errorf("decode in jumeiTextRPC err: %s", err.Error())
			return

		default:
			config.Logger.Errorf("command must be 'RPC' ,now is %s .", jumeiTextRPC.Command)
			return
		}
	}
}
