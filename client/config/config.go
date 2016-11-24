package config

import (
	"fmt"
	"os"
	"time"

	"github.com/burntsushi/toml"
)

//RPC rpc的调用信息
type EndPoint struct {
	URI        string `toml:"uri"`
	User       string `toml:"user"`
	Secret     string `toml:"secret"`
	Compressor string `toml:"compressor,omitempty"`
	//从uri中分离出来的 :// 的前半部分
	NetType string
	//从uri中分离出来的 :// 的后半部分
	NetURI string
}

//RPCMap rpc的调用信息hash表
type EndPointMap struct {
	RpcSecretKey string               `toml:"rpc_secret_key"`
	TimeOut      time.Duration        `toml:"time_out"`
	Maps         map[string]*EndPoint `toml:"endpoint"`
	Threadnum    int                  `toml:"threadnum"`
	TestTime     int                  `toml:"test_time"`
}

var RPCEndPointMap EndPointMap

func LoadConfig(configpath ...string) {
	path := "conf/config.toml"
	if len(configpath) > 0 {
		path = configpath[0]
	}
	if _, err := toml.DecodeFile(path, &RPCEndPointMap); err != nil {
		fmt.Println("failed to load " + path + ", error:" + err.Error())
		os.Exit(1)
	}
}
