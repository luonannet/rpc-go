package rpc_methods

import (
	"fmt"
	"rpc-go/service/register"
	"rpc-go/transport"
)

func init() {
	register.Reg.RegisterHandler("Example/a/b.sayHello", RpcTest2Handler)
}

func RpcTest2Handler(conn *transport.JumeiConn, request interface{}) {
	response := fmt.Sprintf("RpcTest2Handler response :%s", request)
	conn.S2CSend(response)
}
