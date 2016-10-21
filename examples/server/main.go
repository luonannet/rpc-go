package main

import (
	"rpc-go/server/rpc_methods"
	"rpc-go/service"
)

func main() {
	srvs := service.NewService()
	srvs.Init("test.rpc.jumei.com", "tcp", "127.0.0.1:9999")
	srvs.RegisterHandler("Example/a/b.sayHello", rpc_methods.TestExample.SayHello)
	srvs.RegisterHandler("Example.sayHello", rpc_methods.TestExample.SayHello)
	srvs.RegisterHandler("Example.RpcTest1Handler", rpc_methods.TestExample.RpcTest1Handler)
	srvs.RegisterHandler("Example.RpcTest2Handler", rpc_methods.RpcTest2Handler)
	srvs.Run()
}
