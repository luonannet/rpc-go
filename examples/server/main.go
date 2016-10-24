package main

import (
	"fmt"
	"os"
	"rpc-go/config"
	_ "rpc-go/examples/server/rpc_methods"
	"rpc-go/service"
	_ "rpc-go/service/register"
	"runtime"
	"runtime/pprof"
	"time"
)

func main() {
	defer saveHeapProfile()
	config.LoadConfig()
	srvs := service.NewService()
	srvs.Init("test.rpc.jumei.com", "tcp", "172.20.4.19:9999")
	srvs.Run()

}

func saveHeapProfile() {
	runtime.GC()
	f, err := os.Create(fmt.Sprintf("server_%s.prof", time.Now().Format("2006_01_02_03_04_05")))

	if err != nil {
		return
	}
	defer f.Close()
	//pprof.WriteHeapProfile(f)
	pprof.Lookup("heap").WriteTo(f, 1)
}
