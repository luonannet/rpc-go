package service

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"rpc-go/goserver/codec"
	"rpc-go/goserver/config"
	"rpc-go/goserver/service/register"
	"rpc-go/goserver/transport"
	"strings"
	"time"
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

	listener, linstenErr := net.Listen(conf.NetType, conf.Port)
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
		connectionNum := transport.GetConnNumber()
		if connectionNum >= srvs.conf.MaxConnection {
			warnning := fmt.Sprintf("Maximum number of connections is %d ,but now is  %d ", srvs.conf.MaxConnection, connectionNum)
			config.Logger.Warnf(warnning)
			response, _ := codec.WrapS2CData(503, warnning)
			conn.Write([]byte(response))
			conn.Close()
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
			config.Logger.Critical(" system exit success")
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
		conn.SendError(outErr, false)
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
	datachan := make(chan transport.JumeiTextRPC, 100)
	jumeiConn.Conn.SetDeadline(time.Now().Add(time.Second * 10))
	go transport.ReceiveClientData(jumeiConn, datachan)
	for {
		jumeiTextRPC := <-datachan
		compress := false
		switch jumeiTextRPC.Command {
		case "RPC:GZ":
			//数据采用了压缩.
			jumeiTextRPC.Data = string(codec.DoZlibUnCompress([]byte(jumeiTextRPC.Data)))
			compress = true
			fallthrough
		case "RPC":
			rpcData, err := codec.UnwrapRPC(jumeiTextRPC.Data)
			if err != nil {
				//解码客户端数据的时候就发生错误
				err = jumeiConn.SendError(err.Error(), compress)
				config.Logger.Errorf("decode in jumeiTextRPC err: %s", err.Error())
				return
			}
			var rpcParam codec.RPCParam
			json.Unmarshal([]byte(rpcData.Data), &rpcParam)
			// 有些php调用的class 类似Example\a\b 的。先替换成
			rpcParam.Class = strings.Replace(rpcParam.Class, "\\", "/", -1)
			callFunc := register.GetHandler(rpcParam.Class + "." + rpcParam.Method)
			if callFunc == nil {
				notExist := fmt.Sprintf("RPC Method:%s not exist", rpcParam.Class+"."+rpcParam.Method)
				jumeiConn.SendError(notExist, compress)
				config.Logger.Errorf(notExist)
				return
			}

			response, apiErr := callFunc(jumeiConn, rpcParam.Params)
			//将rpc执行结果返回到客户端.如果有错误 那么把错误返回。
			if apiErr != nil {
				config.Logger.Errorf("rpc execute err: %s", apiErr.Error())
				err = jumeiConn.SendError(apiErr.Error(), compress)
				if err != nil {
					config.Logger.Errorf("response err: %s", err.Error())
				}
			} else {
				//返回的结果是否压缩是根据request是否压缩来决定，request压缩则response也压缩，
				err = jumeiConn.Send(response, compress)
				if err != nil {
					config.Logger.Errorf("response err: %s", err.Error())
				}
			}
			//因为是短连接，发送完response 后即刻返回并关闭连接
			return

		default:
			config.Logger.Errorf("command must be 'RPC' or 'RPC:GZ' ,now is %s .", jumeiTextRPC.Command)
			return
		}
	}
}
