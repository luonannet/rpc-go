package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"rpc-go/codec"
	"rpc-go/service/register"
	"strings"
	"sync/atomic"
	"time"
)

func main() {
	rpcServerAddress := flag.String("add", "172.20.4.19:9999", "the species we are studying")

	tcpAddr, err := net.ResolveTCPAddr("tcp", *rpcServerAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	starttime = time.Now()
	//启动多线程.每个线程循环测试
	for i := 0; i < 50; i++ {
		go startAnet(tcpAddr)
	}
	// sign := make(chan os.Signal, 0)
	// signal.Notify(sign, os.Kill)
	// <-sign
	time.Sleep(time.Second * 60)

	fmt.Println(time.Now().Sub(starttime), "---", i)
}
func callTest() {

}

//用来计算服务器压力测试的开始时间
var starttime time.Time

//在指定时间内的 request 数
var i int32

func startAnet(tcpAddr *net.TCPAddr) {

	for {
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Println(time.Now().Sub(starttime), "---", i)
			//	fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			continue
			//	os.Exit(1)
		}

		receiveChan := make(chan string, 0)
		go receiveData(conn, receiveChan)

		sender(conn)
		_ = <-receiveChan
		closeConn(conn)
		atomic.AddInt32(&i, 1)
		//	time.Sleep(time.Microsecond * 1)
		//fmt.Println(receivedata)
		//	receivedata := <-receiveChan

	}
}

func sender(conn net.Conn) {
	words := "hello world! 你好 中文"
	dataString, err := codec.InitCallRPC(register.RPC_Client_Prefix+"Example", "RpcTest1Handler", words)
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
