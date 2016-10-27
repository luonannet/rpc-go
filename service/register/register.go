package register

import (
	"rpc-go/transport"
	"sync"
)

const (
	RPC_Client_Prefix string = "RpcClient_"
)

//rpc 的执行体
type HandlerFunction func(conn *transport.JumeiConn, request interface{}) (response string, err error)
type Register struct {
	Handlers map[string]HandlerFunction
	Lock     sync.Mutex
}

var reg Register

func RegisterHandler(handname string, handFunc HandlerFunction) {
	handname = RPC_Client_Prefix + handname
	reg.Lock.Lock()
	if reg.Handlers == nil {
		reg.Handlers = make(map[string]HandlerFunction, 0)
	}
	reg.Handlers[handname] = handFunc
	reg.Lock.Unlock()
}

//根据rpc的注册名获取对应的处理接口
func GetHandler(handname string) HandlerFunction {
	return reg.Handlers[handname]
}
