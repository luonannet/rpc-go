/**
*
*
**/
package codec

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	RPC_Client_Prefix string = "RpcClient_"
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

//封装client->server 成聚美rpc text格式
func WrapC2SData(command, data string) string {
	return fmt.Sprintf("%d\n%s\n%d\n%s\n", len(command), command, len(data), data)
}

//解封server->client
func UnWrapS2CData(originData string) (data, leftString string, err error) {
	dataLengthIndex := strings.Index(originData, "\n")
	if dataLengthIndex <= 0 || dataLengthIndex > len(originData) {
		err = errors.New("dataLengthIndex length is invalidate")
		return
	}
	dataLengthStr := originData[:dataLengthIndex]
	dataLength, confErr := strconv.Atoi(dataLengthStr)
	if confErr != nil {
		err = confErr
		return
	}
	if dataLength <= 0 || dataLength > len(originData) {
		err = errors.New("data length is invalidate")
		return
	}
	data = originData[dataLengthIndex+1 : dataLengthIndex+1+dataLength]
	leftString = originData[dataLengthIndex+1+dataLength+1:]
	return
}

func UnwrapRPC(rpcData string) (receivedata RPCData, err error) {

	err = json.Unmarshal([]byte(rpcData), &receivedata)

	return
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
