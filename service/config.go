package service

import (
	"log"

	"github.com/burntsushi/toml"
)

var (
	Conf *Config
)

type Config struct {
	RPC map[string]RpcConfigStruct `toml:"jumei_rpc"`
}
type RpcConfigStruct struct {
	Service    string `toml:"service"`
	User       string `toml:"user"`
	Secret     string `toml:"secret"`
	Ver        string `toml:"ver"`
	MethodName string `toml:"method_name"`
	Compressor string `toml:"compressor"`
}

func LoadConfig(configpath string) (err error) {
	if len(configpath) == 0 {
		configpath = "conf/config.toml"
	}
	Conf = new(Config)
	if _, err = toml.DecodeFile(configpath, Conf); err != nil {
		log.Panicln("failed to load " + configpath + ", error:" + err.Error())
		return
	}

	return nil
}
