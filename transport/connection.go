package transport

import (
	"net"
	"rpc-go/codec"
)

type JumeiConn struct {
	Conn net.Conn
}

//S2CSend RPC 服务器发送数据到调用端
//参数response
func (jc *JumeiConn) S2CSend(response string) (err error) {
	response, err = codec.WrapS2CData(response)
	jc.Conn.Write([]byte(response))
	return
}
