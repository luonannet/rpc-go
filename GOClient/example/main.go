package main

import (
	"fmt"
	"rpc-go/goclient/config"
	"rpc-go/goclient/transport"
	"sync/atomic"
	"time"
)

func main() {
	//载入配置文件。默认地址在conf/config.toml
	config.LoadConfig("conf/config.toml")
	endPointAddr, err := transport.ParseEndPoint("user")
	if err != nil {
		fmt.Println(err.Error())
	}

	starttime = time.Now()
	//启动多线程.每个线程循环测试
	for i := 0; i < 100; i++ {
		go testCallRPC(endPointAddr)
	}
	time.Sleep(time.Second * 10)
	fmt.Printf("耗时:%f秒,调用%d次", time.Now().Sub(starttime).Seconds(), num)
}

//用来计算服务器压力测试的开始时间
var starttime time.Time

//在指定时间内完成的rpc 请求数
var num int32

//testCallRPC 循环压力测试 一个rpc
func testCallRPC(endPointAddr *transport.JumeiEndPoint) {
	for {
		endPointAddr.Call("Example", "RpcTest1Handler", "hello world! 你好 中文")

		atomic.AddInt32(&num, 1)
	}
}
