/**
*
*
**/
package codec

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"
)

//RPCData rpc的data json数据
type RPCData struct {
	Data      string `json:"data"`
	Signature string `json:"signature"`
}

//RPCParam rpc的data json数据
type RPCParam struct {
	Version   string      `json:"version"`
	User      string      `json:"user"`
	Password  string      `json:"password"`
	Timestamp int64       `json:"timestamp"`
	Class     string      `json:"class"`
	Method    string      `json:"method"`
	Params    interface{} `json:"params"`
}

//S2CData 服务器端发往客户端的数据结构
type S2CData struct {
	Data  string `json:"data"`
	Statu int    `json:"statu"`
}

// InitCallRPC 客户端调用jmTextRPC 初始化数据
func InitCallRPC(class, method, params string) (result string, err error) {
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
	mdtcry.Write([]byte("769af463a39f077a0340a189e9c1ec28"))
	rpcData.Signature = hex.EncodeToString(mdtcry.Sum(nil))
	rpcData.Data = string(rpcdataBytes)
	rpcdataBytes, err = json.Marshal(rpcData)
	if err != nil {
		return
	}
	result = string(rpcdataBytes)
	return
}
