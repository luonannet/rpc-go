package main

import (
	"flag"
	"os"
	"rpc-go/server/config"
	"rpc-go/server/example/rpc_methods"
	"rpc-go/server/service"
	"rpc-go/server/service/register"
)

func main() {
	//载入配置文件。默认地址在conf/config.toml
	var configFile string
	flag.StringVar(&configFile, "c", "conf/config.toml", " set config file path")
	flag.Parse()
	conf, err := config.LoadConfig(configFile)
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
