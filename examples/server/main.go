package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"rpc-go/config"
	"rpc-go/examples/server/rpc_methods" // 将rpc 的方法进行注册
	"rpc-go/service"
	"rpc-go/service/register"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	//载入配置文件。默认地址在conf/config.toml
	conf, err := config.LoadConfig("conf/config.toml")
	if err != nil {
		config.Logger.Error("load Config Error:", err.Error())
		os.Exit(1)
	}
	defer config.Logger.Flush()
	srvs := service.NewService(conf)
	register.RegisterHandler("Example.RpcTest1Handler", rpc_methods.TestExample.RpcTest1Handler)
	register.RegisterHandler("Example.SayHello", rpc_methods.TestExample.SayHello)
	srvs.Run()
}
