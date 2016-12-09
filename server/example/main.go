package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"rpc-go/server/config"
	"rpc-go/server/example/rpc_methods"
	"rpc-go/server/service"
	"rpc-go/server/service/register"

	_ "net/http/pprof"
)

//the build -X args
var (
	GitHash   = "####"
	BuildTime = "0000-00-00 00:00:00"
	Version   = "dev"
)

func main() {
	//载入配置文件。默认地址在conf/config.toml
	var configFile string
	var showVer bool
	flag.StringVar(&configFile, "c", "conf/config.toml", " set config file path")
	flag.BoolVar(&showVer, "version", false, "show the version and exit")
	flag.Parse()
	if showVer {
		fmt.Printf("Version:%s, GitHash:%s, BuildTime:%s\n\n", Version, GitHash, BuildTime)
		os.Exit(0)
	}

	conf, err := config.LoadConfig(configFile)
	if err != nil {
		config.Logger.Error("load Config Error:", err.Error())
		os.Exit(1)
	}
	defer config.Logger.Flush()

	//pprof
	go func() {
		defer func() {
			if x := recover(); x != nil {
				log.Println(x)
			}
		}()
		addr := ":6060"
		config.Logger.Critical("API Service Running on :", addr)
		config.Logger.Critical(http.ListenAndServe(addr, nil))
	}()

	srvs := service.NewService(conf)
	register.RegisterHandler("Example.RpcTest1Handler", rpc_methods.TestExample.RpcTest1Handler)
	register.RegisterHandler("Example.SayHello", rpc_methods.TestExample.SayHello)
	srvs.Run()
}
