package register

import (
	"rpc-go/transport"
	"sync"
)

const (
	RPC_Client_Prefix string = "RpcClient_"
)

//rpc 的执行体
type HandlerFunction func(conn *transport.JumeiConn, request interface{})
type Register struct {
	Handlers map[string]HandlerFunction
	Lock     sync.Mutex
}

var Reg Register

func (reg *Register) RegisterHandler(handname string, handFunc HandlerFunction) {
	handname = RPC_Client_Prefix + handname
	reg.Lock.Lock()
	if reg.Handlers == nil {
		reg.Handlers = make(map[string]HandlerFunction, 0)
	}
	reg.Handlers[handname] = handFunc
	reg.Lock.Unlock()
}
