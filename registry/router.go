package registry

import (
	"log"
	"reflect"
	"rpc-go/service"
)

// 通过将要发布的rpc 使用反射的方式注册出去的办法
// 但感觉不方便 .暂时不用。
func Register(srv *service.JumeiTCPService, class, method interface{}) {
	classType := reflect.TypeOf(class)
	log.Println(classType.Name())

}
