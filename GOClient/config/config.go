package config

import (
	"fmt"
	"os"

	"github.com/burntsushi/toml"
)

//RPC rpc的调用信息
type RPC struct {
	URI        string `toml:"uri"`
	User       string `toml:"user"`
	Secret     string `toml:"secret"`
	Compressor string `toml:"compressor"`
	//从uri中分离出来的 :// 的前半部分
	NetType string
	//从uri中分离出来的 :// 的后半部分
	NetURI string
}

//RPCMap rpc的调用信息hash表
type RPCMap struct {
	Maps map[string]RPC `toml:"rpc"`
}

var RPCConfig RPCMap

func LoadConfig(configpath ...string) {
	path := "conf/config.toml"
	if len(configpath) > 0 {
		path = configpath[0]
	}
	if _, err := toml.DecodeFile(path, &RPCConfig); err != nil {
		fmt.Println("failed to load " + path + ", error:" + err.Error())
		os.Exit(1)
	}
}
