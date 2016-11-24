package rpc_methods

import (
	"fmt"
	"rpc-go/server/config"
	"rpc-go/server/transport"
)

type Example struct {
}

func init() {
	//	register.Reg.RegisterHandler("Example.sayHello", TestExample.SayHello)
	//register.RegisterHandler("Example.RpcTest1Handler", TestExample.RpcTest1Handler)

}

var TestExample Example

func (this *Example) SayHello(conn *transport.JumeiConn, request interface{}) (response string, err error) {
	response = fmt.Sprintf("sayHello response :%s", request)
	config.Logger.Info("SayHello response:", response)
	return
}
func (this *Example) RpcTest1Handler(conn *transport.JumeiConn, request interface{}) (response string, err error) {
	// var ff []string
	// _ = ff[2:4]
	//	panic(errors.New("panic from RpcTest1Handler "))
	response = fmt.Sprintf("RpcTest1Handler response :%s", request)
	config.Logger.Info("RpcTest1Handler response:", response)

	return
}
