package tcp

import (
	"bnet/util"
)

type ClientOptionSetFunc func(*ClientOption)

type ClientOption struct {
	SessionOption
	ConnectIp string
	ConnectPort int
	connectAddress string
	ClientName string
}

var defaultClientOption = &ClientOption{
	ConnectIp:"127.0.0.1",
	ConnectPort:7777,
	ClientName:"defaultClient",
}
type CoreClientOption struct {
	ClientOption *ClientOption
	CoreTcpSessionOption
	ops []ClientOptionSetFunc
}

func (self *CoreClientOption) Startup() error {
	if self.ClientOption == nil {
		self.ClientOption = &ClientOption{}
	}
	*self.ClientOption = *defaultClientOption
	self.CoreTcpSessionOption.SessionOption = &self.ClientOption.SessionOption
	if err := self.CoreTcpSessionOption.Startup(); err != nil {
		return err
	}

	for _, op := range self.ops {
		op(self.ClientOption)
	}

	self.ClientOption.connectAddress = util.JoinAddress(self.ClientOption.ConnectIp, self.ClientOption.ConnectPort)
	return nil
}

func (self *CoreClientOption) GetConnectAddress() string {
	return self.ClientOption.connectAddress
}

func (self *CoreClientOption) SetClientOption(ops...ClientOptionSetFunc)  {
	self.ops = ops
}

func SetConnectAddress(ip string, port int) ClientOptionSetFunc {
	return func(option *ClientOption) {
		option.ConnectIp = ip
		option.ConnectPort = port
	}
}
