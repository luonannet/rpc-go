package transport

import (
	"errors"
	"fmt"
	"io"
	"net"
	"rpc-go/client/codec"
	"rpc-go/client/config"
	"strings"
	"time"
)

type JumeiConn struct {
	Conn net.Conn
	*JumeiEndPoint
}
type JumeiEndPoint struct {
	*net.TCPAddr
	*config.EndPoint
}

//ParseEndPoint 解析uri数据
func ParseEndPoint(endpoint string) (*JumeiEndPoint, error) {
	jta := new(JumeiEndPoint)
	rpc := config.RPCEndPointMap.Maps[endpoint]
	if rpc == nil {
		return nil, errors.New("this rpc is not exist .please check config toml ")
	}
	jta.EndPoint = rpc
	uriItems := strings.Split(rpc.URI, "://")
	if len(uriItems) != 2 {
		return nil, errors.New("invalidate uri value")
	}
	jta.NetType = uriItems[0]
	jta.NetURI = uriItems[1]
	tcpAddr, err := net.ResolveTCPAddr(jta.NetType, jta.NetURI)
	if err != nil {
		return nil, err
	}
	jta.TCPAddr = tcpAddr
	return jta, nil
}

//Call JumeiEndPoint 每次call创建一个连接，用完后关闭
//class 是rpc的类名 或者struct名
//method是rpc的方法名
//param 是调用rpc的参数
func (rpcAddr *JumeiEndPoint) Call(class, method, param string, compress bool) (response string, err error) {
	var conn net.Conn
	conn, err = net.DialTCP(rpcAddr.NetType, nil, rpcAddr.TCPAddr)
	if err != nil {
		return "", err
	}
	conn.SetDeadline(time.Now().Add(time.Millisecond * config.RPCEndPointMap.TimeOut))
	dataString, err := codec.InitCallRPC(codec.RPC_Client_Prefix+class, method, param)
	if err != nil {
		return "", err
	}
	command := "RPC"
	// 如果需要压缩数据。
	if compress == true {
		command = "RPC:GZ"
		dataString = string(codec.DoZlibCompress([]byte(dataString)))
	}
	receiveChan := make(chan string, 0)
	go receiveData(conn, receiveChan, compress)
	conn.Write([]byte(codec.WrapC2SData(command, dataString)))
	response = <-receiveChan
	closeConn(conn)
	return response, nil
}

//关闭连接
func closeConn(conn net.Conn) {
	if conn != nil {
		conn.Close()
	}
	conn = nil
}

//receiveData client 一直监听并读取连接中的数据
func receiveData(conn net.Conn, receivechan chan string, compress bool) (err error) {
	defer closeConn(conn)
	var dataBox []byte
	var size int
	for {
		readData := make([]byte, 1024)
		size, err = conn.Read(readData)
		if err != nil {
			if err != io.EOF {
				fmt.Println("not eof error:", err.Error())
				receivechan <- err.Error()
				return
			}
			// 此时 err == EOF 。
			continue
		}
		dataBox = append(dataBox, readData[0:size]...)
		//读完后，进行解包
	dealdata:
		var dataString string
		if compress == true {
			//如果是压缩的request ，那么收到的是压缩的response .
			//解压后再返回
			dataString = string(codec.DoZlibUnCompress(dataBox))
		} else {
			dataString = string(dataBox)
		}
		sepNumber := strings.Count(dataString, "\n")
		if sepNumber >= 2 {
			data, leftstring, unWrapErr := codec.UnWrapS2CData(dataString)
			if unWrapErr != nil {
				// 如果解包出现问题，说明数据已经乱了。则丢掉之前的数据
				unWrapErrStr := fmt.Sprintf("decode %s ,and error %s", dataString, unWrapErr.Error())
				dataBox = dataBox[0:0]
				receivechan <- unWrapErrStr
				return
			}

			dataBox = dataBox[0:0]
			// 如果还有剩余数据，那么继续处理
			if len(leftstring) > 0 {
				dataBox = []byte(leftstring)
				leftstring = ""
				goto dealdata
			} else {
				receivechan <- data
				return
			}
		}

	}

}

// func (jc *JumeiConn) SendData(data string) {
// 	words := "hello world! 你好 中文"
// 	dataString, err := codec.InitCallRPC(codec.RPC_Client_Prefix+"Example", "RpcTest1Handler", words)
// 	if err != nil {
// 		fmt.Println("send err", err.Error())
// 		return
// 	}
// 	jc.Conn.Write([]byte(codec.WrapC2SData("RPC", dataString)))
// }
