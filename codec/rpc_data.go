/**
*
*
**/
package codec

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"rpc-go/config"
	"time"
)

//RPCData rpc的data json数据
type RPCData struct {
	Data      string `json:"data"`
	Signature string `json:"signature"`
}

//RPCData rpc的data json数据
type RPCParam struct {
	Version   string      `json:"version"`
	User      string      `json:"user"`
	Password  string      `json:"password"`
	Timestamp int64       `json:"timestamp"`
	Class     string      `json:"class"`
	Method    string      `json:"method"`
	Params    interface{} `json:"params"`
}

//RPC2ClientData
type S2CData struct {
	Data  string `json:"data"`
	Statu int    `json:"statu"`
}

func InitRpcData(class, method, params string) (result string, err error) {
	var rpcData RPCData
	var rpcParam RPCParam
	rpcParam.Version = "2.0"
	rpcParam.Method = method
	rpcParam.Class = class
	rpcParam.User = ""
	rpcParam.Password = ""
	rpcParam.Timestamp = time.Now().Unix()
	rpcParam.Params = params

	var rpcdataBytes []byte
	rpcdataBytes, err = json.Marshal(rpcParam)
	if err != nil {
		return
	}
	mdtcry := md5.New()
	mdtcry.Write([]byte(rpcdataBytes))
	mdtcry.Write([]byte("&"))
	mdtcry.Write([]byte(config.RPCSecretKey))
	rpcData.Signature = hex.EncodeToString(mdtcry.Sum(nil))
	rpcData.Data = string(rpcdataBytes)
	rpcdataBytes, err = json.Marshal(rpcData)
	if err != nil {
		return
	}
	result = string(rpcdataBytes)
	return
}
