package service

const (
	RPC_Client_Prefix string = "RpcClient_"
)

func (srvs *JumeiTCPService) RegisterHandler(handname string, handFunc HandlerFunction) {
	handname = RPC_Client_Prefix + handname
	srvs.Handlers[handname] = handFunc
}
