package transport

import (
	"io"
	"log"
	"rpc-go/codec"
	"strings"
)

// 聚美定义的rpc text协议，格式是
// commandLent
type JumeiTextRPC struct {
	Command string
	Data    string
}

// 一直监听并读取连接中的数据
func ReceiveClientData(jmConn *JumeiConn, dataChan chan JumeiTextRPC) (err error) {
	var dataBox []byte

	var size int
	for {
		readData := make([]byte, 1024)
		size, err = jmConn.Conn.Read(readData)
		if err != nil {
			if err != io.EOF {
				log.Println("not eof error:", err.Error())

				dataBox = dataBox[0:0]

				return
			} else {
				log.Println("client closed :")
				dataBox = dataBox[0:0]
				return
			}
		}
		dataBox = append(dataBox, readData[0:size]...)
		readData = nil

		//读完后，进行解包
	dealdata:
		dataString := string(dataBox)
		sepNumber := strings.Count(dataString, "\n")
		if sepNumber >= 4 {

			command, data, leftstring, unWrapErr := codec.UnWrapC2SData(dataString)
			if unWrapErr != nil {
				// 如果解包出现问题，说明数据已经乱了。则丢掉之前的数据
				log.Println("解码出错:", dataString, unWrapErr.Error())
				dataBox = dataBox[0:0]
				continue
			}
			var jumeiTextRpc JumeiTextRPC
			jumeiTextRpc.Command = command
			jumeiTextRpc.Data = data
			dataChan <- jumeiTextRpc

			dataBox = dataBox[0:0]
			// 如果还有剩余数据，那么继续处理
			if len(leftstring) > 0 {
				dataBox = []byte(leftstring)
				leftstring = ""
				goto dealdata
			}
			return
		}

	}

}
