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

//封装client->server 成聚美rpc text格式
func WrapC2SData(command, data string) string {
	return fmt.Sprintf("%d\n%s\n%d\n%s\n", len(command), command, len(data), data)
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

	dataLength, confErr := strconv.Atoi(dataLengthStr)
	if confErr != nil {
		err = confErr
		return
	}
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
