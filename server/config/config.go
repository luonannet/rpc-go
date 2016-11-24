package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/burntsushi/toml"
	log "github.com/cihub/seelog"
)

type RpcServiceConfig struct {
	ServiceName   string `toml:"service_name"`
	Port          string `toml:"port"`
	MaxConnection int    `toml:"max_connection"` // 最大连接数
	NetType       string `toml:"net_type"`
	LogFilePath   string `toml:"log_file_path"`
	LogLevels     string `toml:"log_levels"`
}

func LoadConfig(configpath ...string) (*RpcServiceConfig, error) {
	path := "conf/config.toml"
	if len(configpath) > 0 {
		path = configpath[0]
	}
	serviceConf := new(RpcServiceConfig)
	if _, err := toml.DecodeFile(path, serviceConf); err != nil {
		fmt.Println("failed to load " + path + ", error:" + err.Error())
		return nil, err
	}
	initLogger(serviceConf)
	return serviceConf, nil
}

var Logger log.LoggerInterface

func initLogger(conf *RpcServiceConfig) {
	var err error
	var logfile *os.File
	if conf.LogFilePath == "" {
		//如果没有配置日志文件地址，则日志文件放在程序同目录
		pwd, errf := filepath.Abs("./")
		if errf != nil {
			fmt.Println(errf.Error())
		}
		conf.LogFilePath = pwd
	}
	conf.LogFilePath = strings.Replace(conf.LogFilePath, "\\", "/", -1)
	if strings.LastIndex(conf.LogFilePath, "/") != (len(conf.LogFilePath) - 1) {
		conf.LogFilePath = conf.LogFilePath + "/"
	}
	logfile, err = os.OpenFile(conf.LogFilePath+conf.ServiceName+"_"+time.Now().Format("2006-01-02-15-04")+".log", os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err.Error())
	}
	Logger, err = log.LoggerFromConfigAsString(`<seelog levels="` + conf.LogLevels + `">
    <outputs formatid="common">
        <console />
		<file path="` + logfile.Name() + `"/>
    </outputs>
    <formats>
    <format id="common" format="[%Level]%Date(2006/01/02 15:04:05.999)[%File:%Line] %Msg%n"/>
</formats>
</seelog>
	`)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = log.ReplaceLogger(Logger)
	if err != nil {
		fmt.Println(err.Error())
	}
}
