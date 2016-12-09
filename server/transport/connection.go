package transport

import (
	"rpc-go/server/codec"
	"net"
	"sync"
)

//ConnectionsStruct 带读写锁的链接集合
type ConnectionsStruct struct {
	Lock sync.RWMutex
	//key是连接。value 为1 是正常连接。为0为已经断开的连接。为0的情况暂时不考虑。断开的都直接删掉了
	// 	ConnMap map[*JumeiConn]int8
}

type JumeiConn struct {
	Conn net.Conn
}

func init() {
	//	Connections.ConnMap = make(map[*JumeiConn]int8)
}

// Connections rpc的映射哈希表
var Connections ConnectionsStruct

//Send RPC 服务器发送数据到调用端
//response 发送的内容
//compress 是否压缩内容
func (jc *JumeiConn) Send(response string, compress bool) (err error) {

	return jc.sendWithStatu(200, response, compress)
}

//SendError RPC 服务器发送错误信息到调用端
//response 返回的错误内容
//compress 是否压缩内容
func (jc *JumeiConn) SendError(response string, compress bool) (err error) {

	return jc.sendWithStatu(500, response, compress)
}

//sendWithStatu RPC 服务器发送状态码和数据到调用端
//参数response
//compress 是否压缩内容
func (jc *JumeiConn) sendWithStatu(statu int, response string, compress bool) (err error) {
	response, err = codec.WrapS2CData(statu, response)
	var out []byte
	if compress == true {
		out = codec.DoZlibCompress([]byte(response))
	} else {
		out = []byte(response)
	}
	jc.Conn.Write(out)
	return
}

//AddConnection 连接进入
func AddConnection(conn *JumeiConn) {
	// Connections.Lock.RLock()
	// if _, v := Connections.ConnMap[conn]; v == false {
	// 	Connections.Lock.RUnlock()
	// 	Connections.Lock.Lock()
	// 	if _, v = Connections.ConnMap[conn]; v == false {
	// 		Connections.ConnMap[conn] = 1
	// 	}
	// 	Connections.Lock.Unlock()
	// } else {
	// 	Connections.Lock.RUnlock()
	// }
}

//CloseConn 连接断开
func (jc *JumeiConn) CloseConn() {
	// Connections.Lock.RLock()
	// if _, v := Connections.ConnMap[jc]; v == true {
	// 	Connections.Lock.RUnlock()
	// 	Connections.Lock.Lock()
	// 	if _, v = Connections.ConnMap[jc]; v == true {
	if jc.Conn != nil {
		jc.Conn.Close()
		jc.Conn = nil
	}
	// 		delete(Connections.ConnMap, jc)
	// 	}
	// 	Connections.Lock.Unlock()
	// } else {
	// 	Connections.Lock.RUnlock()
	// }
}

// 获取此时连接数
func GetConnNumber() (result int) {
	// Connections.Lock.RLock()
	// result = len(Connections.ConnMap)
	// Connections.Lock.RUnlock()
	return result
}
