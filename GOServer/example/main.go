package main

import (
	"os"
	"rpc-go/goserver/config"
	"rpc-go/goserver/example/rpc_methods" // 将rpc 的方法进行注册
	"rpc-go/goserver/service"
	"rpc-go/goserver/service/register"
)

func main() {
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
