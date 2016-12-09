package main

import (
	"fmt"
	"rpc-go/client/config"
	"rpc-go/client/transport"
	"sync/atomic"
	"time"
)

func main() {
	//载入配置文件。默认地址在conf/config.toml
	config.LoadConfig("conf/config.toml")
	endPointAddr, err := transport.ParseEndPoint("prod")
	if err != nil {
		fmt.Println(err.Error())
	}

	starttime = time.Now()
	//启动多线程测试.每个线程循环测试

	for i := 0; i < config.RPCEndPointMap.Threadnum; i++ {
		go testCallRPC(endPointAddr)
	}
	time.Sleep(time.Second * time.Duration(config.RPCEndPointMap.TestTime))
	fmt.Printf("线程数:%d,调用%d次,共耗时:%f秒", config.RPCEndPointMap.Threadnum, successNum, time.Now().Sub(starttime).Seconds())
}

//用来计算服务器压力测试的开始时间
var starttime time.Time

//在指定时间内完成的rpc 请求数
var successNum int32

//testCallRPC 循环压力测试 一个rpc
func testCallRPC(endPointAddr *transport.JumeiEndPoint) {
	senddata := "hello world! 你好 中文"
	for {
		response, err := endPointAddr.Call("Example", "SayHello", senddata, false)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			_ = response

			atomic.AddInt32(&successNum, 1)
		}
	}
}
