package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"rpc-go/codec"
	"rpc-go/service/register"
	"strings"
	"time"
)

func main() {
	rpcServerAddress := "172.20.4.19:9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp", rpcServerAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	receiveChan := make(chan string, 0)
	var totalNum int
	starttime := time.Now()

	for totalNum < 50000 {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			os.Exit(1)
		}

		go receiveData(conn, receiveChan)

		sender(conn)
		receivedata := <-receiveChan
		closeConn(conn)
		time.Sleep(time.Microsecond * 1)
		totalNum++
		fmt.Println(receivedata)
		//	receivedata := <-receiveChan

	}
	fmt.Println(time.Now().Sub(starttime))
}
func sender(conn net.Conn) {
	words := "hello world! 你好 中文"
	dataString, err := codec.InitRpcData(register.RPC_Client_Prefix+"Example", "RpcTest1Handler", words)
	if err != nil {
		fmt.Println("send err", err.Error())
		return
	}
	conn.Write([]byte(codec.WrapC2SData("RPC", dataString)))
}
func closeConn(conn net.Conn) {
	conn.Close()
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
				log.Println("not eof error:", err.Error())
				receivechan <- err.Error()
				return
			}
			//log.Println("---readData:", string(readData))
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

				receivechan <- "解码出错"
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
