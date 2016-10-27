/**
*
*
**/
package codec

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
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

//解封client->server
func UnWrapC2SData(originData string) (command, data, leftString string, err error) {
	commandLengthIndex := strings.Index(originData, "\n")
	if commandLengthIndex <= 0 || commandLengthIndex > len(originData) {
		err = errors.New("command length is invalidate")

		return
	}
	commandLengthStr := originData[:commandLengthIndex]
	commandLength, confErr := strconv.Atoi(commandLengthStr)
	if confErr != nil {
		err = confErr
		return
	}
	command = originData[commandLengthIndex+1 : commandLengthIndex+1+commandLength]
	dataLengthIndex := strings.Index(originData[commandLengthIndex+1+commandLength+3:], "\n")
	if dataLengthIndex <= 0 || dataLengthIndex > len(originData) {
		err = errors.New("data length is invalidate")
		return
	}
	dataLengthStr := originData[commandLengthIndex+1+commandLength+1 : commandLengthIndex+1+commandLength+1+dataLengthIndex+1+1]
	var dataLength int

	// if dataLengthStr == '?' {
	// 	//数据长度, 为 '?' 时读取数据直到出现 "\n".
	// 	dataLength = strings.Index(originData[commandLengthIndex+1+commandLength+1+dataLengthIndex+1+1+1:], "\n")
	// 	fmt.Println(originData[commandLengthIndex+1+commandLength+1+dataLengthIndex+1+1+1:], "----", dataLength)
	// } else {
	//数据长度为正常的数字
	dataLength, err = strconv.Atoi(dataLengthStr)
	if err != nil {
		return
	}
	//	}
	data = originData[commandLengthIndex+1+commandLength+1+dataLengthIndex+1+2 : commandLengthIndex+1+commandLength+1+dataLengthIndex+1+dataLength+2]
	// 剩余部分可能是粘包数据，提出来和下一个包进行组合
	leftString = originData[commandLengthIndex+1+commandLength+1+dataLengthIndex+1+dataLength+3:]
	return
}

//封装server->client 成聚美rpc text格式
func WrapS2CData(statu int, data string) (string, error) {

	var s2cData S2CData
	s2cData.Data = data
	s2cData.Statu = statu
	dataBytes, dataErr := json.Marshal(&s2cData)
	if dataErr != nil {
		log.Fatalln(dataErr.Error())
		return "", dataErr
	}
	return fmt.Sprintf("%d\n%s\n", len(dataBytes), string(dataBytes)), nil
}

func UnwrapRPC(rpcData string) (receivedata RPCData, err error) {

	err = json.Unmarshal([]byte(rpcData), &receivedata)

	return
}
