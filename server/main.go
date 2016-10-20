package main

import (
	"fmt"
	"log"
	"net"
	"rpc-go/service"
	"rpc-go/transport"
)

func main() {
	service := service.NewService()
	service.Init("test.rpc.jumei.com")
	var testobj Example
	service.RegisterHandler("Example/a/b.sayHello", testobj.sayHello)
	service.RegisterHandler("Example.sayHello", testobj.sayHello)
	listener, linstenErr := net.Listen("tcp", ":9999")
	if linstenErr == nil {
		log.Println("server listen at port:", listener.Addr())
		for {
			conn, connErr := listener.Accept()
			if connErr == nil {
				go service.ServerHandleConn(conn)
			} else {
				log.Println(connErr.Error())
			}
		}
	} else {
		log.Println(linstenErr.Error())
	}
}

type Example struct {
	//service.HandlerFunction
}

func (this *Example) sayHello(conn *transport.JumeiConn, request interface{}) {
	response := fmt.Sprintf("RpcTest1Handler response :%s", request)
	conn.S2CSend(response)
}

func (this *Example) RpcTest2Handler(conn *transport.JumeiConn, request interface{}) {
	response := fmt.Sprintf("RpcTest2Handler response :%s", request)
	conn.S2CSend(response)

}
