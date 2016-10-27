package main

import (
	"fmt"
	"rpc-go/goclient/codec"
	"rpc-go/goclient/config"
	"rpc-go/goclient/transport"
	"sync/atomic"
	"time"
)

func main() {
	//载入配置文件。默认地址在conf/config.toml
	config.LoadConfig("conf/config.toml")
	tcpAddr, err := transport.ParseURI("user")
	if err != nil {
		fmt.Println(err.Error())
	}
	starttime = time.Now()
	//启动多线程.每个线程循环测试
	for i := 0; i < 100; i++ {
		go startAnet(tcpAddr)
	}
	time.Sleep(time.Second * 10)
	fmt.Printf("耗时:%f秒,调用%d次", time.Now().Sub(starttime).Seconds(), i)
}

//用来计算服务器压力测试的开始时间
var starttime time.Time

//在指定时间内的 request 数
var i int32

func startAnet(tcpAddr *transport.JumeiTcpAddr) {
	for {
		tcpAddr.Call(codec.RPC_Client_Prefix+"Example", "RpcTest1Handler", "hello world! 你好 中文")

		atomic.AddInt32(&i, 1)
	}
}
