package rpc_methods

import (
	"fmt"
	"rpc-go/transport"
)

type Example struct {
	//service.HandlerFunction
}

var TestExample Example

func (this *Example) SayHello(conn *transport.JumeiConn, request interface{}) {
	response := fmt.Sprintf("sayHello response :%s", request)
	conn.S2CSend(response)
}
func (this *Example) RpcTest1Handler(conn *transport.JumeiConn, request interface{}) {
	response := fmt.Sprintf("sayHello response :%s", request)
	conn.S2CSend(response)
}

func RpcTest2Handler(conn *transport.JumeiConn, request interface{}) {
	response := fmt.Sprintf("RpcTest2Handler response :%s", request)
	conn.S2CSend(response)

}
