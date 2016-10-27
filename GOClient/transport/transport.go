package transport

import (
	"errors"
	"fmt"
	"io"
	"net"
	"rpc-go/goclient/codec"
	"rpc-go/goclient/config"
	"strings"
)

type JumeiConn struct {
	Conn net.Conn
	JumeiTcpAddr
}
type JumeiTcpAddr struct {
	*net.TCPAddr
	config.RPC
}

//ParseURI 解析uri数据
func ParseURI(rpcApi string) (*JumeiTcpAddr, error) {
	jta := new(JumeiTcpAddr)
	rpc := config.RPCConfig.Maps[rpcApi]
	uriItems := strings.Split(rpc.URI, "://")
	if len(uriItems) != 2 {
		return nil, errors.New("uri格式不对")
	}
	jta.NetType = uriItems[0]
	jta.NetURI = uriItems[1]
	tcpAddr, err := net.ResolveTCPAddr(jta.NetType, jta.NetURI)
	if err != nil {
		//		fmt.Fprintf(os.Stderr, " error: %s", err.Error())
		return nil, err
	}
	jta.TCPAddr = tcpAddr
	return jta, nil
}

//Call 每次call创建一个连接。
func (rpcAddr *JumeiTcpAddr) Call(class, method, param string) (response string, err error) {
	var conn net.Conn

	conn, err = net.DialTCP(rpcAddr.NetType, nil, rpcAddr.TCPAddr)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	dataString, err := codec.InitCallRPC(class, method, param)
	if err != nil {
		fmt.Println("send err", err.Error())
		return "", err
	}
	receiveChan := make(chan string, 0)
	go receiveData(conn, receiveChan)
	conn.Write([]byte(codec.WrapC2SData("RPC", dataString)))
	response = <-receiveChan
	closeConn(conn)
	return response, nil
}
func closeConn(conn net.Conn) {
	if conn != nil {
		conn.Close()
	}
	conn = nil
}

//receiveData client 一直监听并读取连接中的数据
func receiveData(conn net.Conn, receivechan chan string) (err error) {
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
		dataString := string(dataBox)
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
