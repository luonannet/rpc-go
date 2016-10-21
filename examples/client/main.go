package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"rpc-go/codec"
	"rpc-go/service"
	"strings"
)

func main() {
	rpcServerAddress := "127.0.0.1:9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp", rpcServerAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	receiveChan := make(chan string, 0)
	go receiveData(conn, receiveChan)
	sender(conn)
	receivedata := <-receiveChan
	fmt.Println(receivedata)

}
func sender(conn net.Conn) {
	words := "hello world! 你好 中文"
	dataString, err := codec.InitRpcData(service.RPC_Client_Prefix+"Example", "RpcTest1Handler", words)
	if err != nil {
		fmt.Println("send err", err.Error())
		return
	}
	conn.Write([]byte(codec.WrapC2SData("RPC", dataString)))

}

//receiveData client 一直监听并读取连接中的数据
func receiveData(conn net.Conn, receivechan chan string) (err error) {
	//defer conn.Close()
	var dataBox []byte
	var size int
	for {
		readData := make([]byte, 1024)
		size, err = conn.Read(readData)
		if err != nil {
			if err != io.EOF {
				log.Println("not eof error:", err.Error())
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
				log.Println("解码出错:", dataString, unWrapErr.Error())
				dataBox = dataBox[0:0]
				conn.Close()
				continue
			}

			dataBox = dataBox[0:0]
			// 如果还有剩余数据，那么继续处理
			if len(leftstring) > 0 {
				dataBox = []byte(leftstring)
				leftstring = ""
				goto dealdata
			} else {
				conn.Close()
				receivechan <- data
				return
			}
		}

	}

}
