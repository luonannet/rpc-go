package rpc_methods

import (
	"fmt"
	"rpc-go/service/register"
	"rpc-go/transport"
)

type Example struct {
}

func init() {
	register.Reg.RegisterHandler("Example.sayHello", TestExample.SayHello)
	register.Reg.RegisterHandler("Example.RpcTest1Handler", TestExample.RpcTest1Handler)

}

var TestExample Example

func (this *Example) SayHello(conn *transport.JumeiConn, request interface{}) {
	response := fmt.Sprintf("sayHello response :%s", request)
	//	fmt.Println(response)
	conn.S2CSend(response)
}
func (this *Example) RpcTest1Handler(conn *transport.JumeiConn, request interface{}) {
	// var ff []string
	// _ = ff[2:4]
	//	panic(errors.New("panic from RpcTest1Handler "))
	response := fmt.Sprintf("sayHello response :%s", request)
	//	fmt.Println(response)
	conn.S2CSend(response)
}
