package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println(err.Error())
	}
	for {
		conn, connErr := listener.Accept()
		if connErr != nil {
			fmt.Println(connErr.Error())
			return
		}
		go dealconn(conn)
	}

}

var n int

func dealconn(conn net.Conn) {
	n++
	fmt.Println(n)

}
