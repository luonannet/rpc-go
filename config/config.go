package config

import (
	"log"

	"github.com/burntsushi/toml"
)

const (
	//RPCSecretKey rpc的secret key
	RPCSecretKey string = "769af463a39f077a0340a189e9c1ec28"
)

var (
	ServiceConf *RpcServiceConfig
)

type RpcServiceConfig struct {
	ServiceName   string `toml:"service_name"`
	IP            string `toml:"ip"`
	Port          string `toml:"port"`
	MaxConnection int    `toml:"max_connection"` // 最大连接数
	SecretKey     string `toml:"secret_key"`
}

func LoadConfig(configpath ...string) (err error) {
	path := "conf/config.toml"
	if len(configpath) > 0 {
		path = configpath[0]
	}
	ServiceConf = new(RpcServiceConfig)
	if _, err = toml.DecodeFile(path, ServiceConf); err != nil {
		log.Panicln("failed to load " + path + ", error:" + err.Error())
		return
	}
	return nil
}
