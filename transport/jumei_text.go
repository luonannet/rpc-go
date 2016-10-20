package transport

import (
	"io"
	"log"
	"net"
	"rpc-go/codec"
	"strconv"
	"strings"
)

// 聚美定义的rpc text协议，格式是
// commandLent
type JumeiTextRPC struct {
	Command string
	Data    string
}

// 一直监听并读取连接中的数据
func ReceiveData(conn net.Conn, dataChan chan JumeiTextRPC) (err error) {
	defer conn.Close()
	var dataBox []byte
	var size int
	for {
		readData := make([]byte, 1024)
		size, err = conn.Read(readData)
		if err != nil {
			if err != io.EOF {
				log.Println("not eof error:", err.Error())
				return
			} else {
				conn.Close()
				log.Println("client closed :", conn.RemoteAddr())
				return
			}
		}
		dataBox = append(dataBox, readData[0:size]...)

		//读完后，进行解包
	dealdata:
		dataString := string(dataBox)
		sepNumber := strings.Count(dataString, "\n")
		if sepNumber < 4 {
			//如果数据不够
			commandLengthIndex := strings.Index(dataString, "\n")
			if commandLengthIndex < 1 {
				log.Println("脏数据1:", commandLengthIndex)
				dataBox = dataBox[0:0]
				conn.Close()
				continue
			}
			commandLength, comandLengthErr := strconv.Atoi(dataString[:commandLengthIndex])
			if comandLengthErr != nil || commandLength > 8000 {
				//如果收到的数据前面几个字符不是规定格式的，那么说明是脏数据。丢弃
				log.Println("脏数据2:", comandLengthErr, commandLength)
				dataBox = dataBox[0:0]
				conn.Close()
				continue
			}
			//如果格式是对的，只是长度不够，那么继续等待
			continue
		} else if sepNumber >= 4 {
			command, data, leftstring, unWrapErr := codec.UnWrapC2SData(dataString)
			if unWrapErr != nil {
				// 如果解包出现问题，说明数据已经乱了。则丢掉之前的数据
				log.Println("数据已经乱:", dataString, unWrapErr.Error())
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
		}

	}

}
