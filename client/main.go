package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"rpc-jumei/codec"
	"strconv"
	"strings"
	"time"
)

func main() {
	server := "127.0.0.1:9999"
	tcpAddr, err := net.ResolveTCPAddr("tcp", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connect success")
	go ReceiveData(conn)
	sender(conn)

	time.Sleep(10 * time.Second)

}
func sender(conn net.Conn) {
	words := "hello world!"
	dataString, err := codec.InitRpcData("TestObj", "RpcTestHandler", words)
	if err != nil {
		fmt.Println("send err", err.Error())
		return
	}

	conn.Write([]byte(codec.WrapC2SData("RPC", dataString)))

}

// client 一直监听并读取连接中的数据
func ReceiveData(conn net.Conn) (err error) {
	defer conn.Close()
	var dataBox []byte
	var size int
	for {
		readData := make([]byte, 1024)
		size, err = conn.Read(readData)
		if err != nil {
			if err != io.EOF {
				log.Println("not eof error:", err.Error())
				return
			} else {
				conn.Close()
				log.Println("client closed :", conn.RemoteAddr())
				return
			}
		}
		dataBox = append(dataBox, readData[0:size]...)

		//读完后，进行解包
		//dealdata:
		dataString := string(dataBox)
		sepNumber := strings.Count(dataString, "\n")
		if sepNumber < 2 {
			//如果数据不够
			commandLengthIndex := strings.Index(dataString, "\n")
			if commandLengthIndex < 1 {
				log.Println("脏数据1:", commandLengthIndex)
				dataBox = dataBox[0:0]
				conn.Close()
				continue
			}
			commandLength, comandLengthErr := strconv.Atoi(dataString[:commandLengthIndex])
			if comandLengthErr != nil || commandLength > 8000 {
				//如果收到的数据前面几个字符不是规定格式的，那么说明是脏数据。丢弃
				log.Println("脏数据2:", comandLengthErr, commandLength)
				dataBox = dataBox[0:0]
				conn.Close()
				continue
			}
			//如果格式是对的，只是长度不够，那么继续等待
			continue
		} else if sepNumber >= 2 {
			data, unWrapErr := codec.UnWrapS2CData(dataString)
			if unWrapErr != nil {
				// 如果解包出现问题，说明数据已经乱了。则丢掉之前的数据
				log.Println("数据已经乱:", dataString, unWrapErr.Error())
				dataBox = dataBox[0:0]
				conn.Close()
				continue
			}
			fmt.Println(data)
		}

	}

}
