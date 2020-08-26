package netimpl

import (
	"bnet/inet"
	"fmt"
)

var netCreateFuncMap = map[inet.NetType]inet.NetCreateFunc{}

func RegisterNetCreator(f inet.NetCreateFunc) {
	// 临时实例化一个，获取类型
	dummyNet := f(nil)
	if _, ok := netCreateFuncMap[dummyNet.NetType()]; ok {
		panic(fmt.Sprintf("duplicate netType type: %v", dummyNet.NetType()))
	}
	netCreateFuncMap[dummyNet.NetType()] = f
}

func NewNet(netType inet.NetType, config interface{}) inet.INet {
	creatoFunc := netCreateFuncMap[netType]
	if creatoFunc == nil {
		panic(fmt.Sprintf("invalid net:netType=%v", netType))
	}
	return creatoFunc(config)
}
