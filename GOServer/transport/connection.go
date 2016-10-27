package transport

import (
	"net"
	"rpc-go/goserver/codec"
	"sync"
)

type ConnectionsStruct struct {
	Lock sync.RWMutex
	//key是连接。value 为1 是正常连接。为0为已经断开的连接。为0的情况暂时不考虑。断开的都直接删掉了
	ConnMap map[*JumeiConn]int8
}

type JumeiConn struct {
	Conn net.Conn
}

func init() {
	Connections.ConnMap = make(map[*JumeiConn]int8)
}

// Connections rpc的映射哈希表
var Connections ConnectionsStruct

//Send RPC 服务器发送数据到调用端
//参数response
func (jc *JumeiConn) Send(response string) (err error) {

	return jc.sendWithStatu(200, response)
}

//SendError RPC 服务器发送数据到调用端
//参数response
func (jc *JumeiConn) SendError(response string) (err error) {
	return jc.sendWithStatu(500, response)
}

//sendWithStatu RPC 服务器发送数据到调用端
//参数response
func (jc *JumeiConn) sendWithStatu(statu int, response string) (err error) {
	response, err = codec.WrapS2CData(statu, response)
	jc.Conn.Write([]byte(response))
	return
}

//AddConnection 连接进入
func AddConnection(conn *JumeiConn) {
	Connections.Lock.RLock()
	if _, v := Connections.ConnMap[conn]; v == false {
		Connections.Lock.RUnlock()
		Connections.Lock.Lock()
		if _, v = Connections.ConnMap[conn]; v == false {
			Connections.ConnMap[conn] = 1
		}
		Connections.Lock.Unlock()
	} else {
		Connections.Lock.RUnlock()
	}
}

//连接断开
func (jc *JumeiConn) CloseConn() {
	Connections.Lock.RLock()
	if _, v := Connections.ConnMap[jc]; v == true {
		Connections.Lock.RUnlock()
		Connections.Lock.Lock()
		if _, v = Connections.ConnMap[jc]; v == true {
			if jc.Conn != nil {
				jc.Conn.Close()
				jc.Conn = nil
			}
			delete(Connections.ConnMap, jc)
		}
		Connections.Lock.Unlock()
	} else {
		Connections.Lock.RUnlock()
	}
}
