package tcp

type ServerOptionSetFunc func(*ServerOption)

type ServerOption struct {
	SessionOption
	ListenIp string
	ListenPort int
	ServerName string
}

var defaultServerOption = &ServerOption{
	ListenIp:"127.0.0.1",
	ListenPort:7777,
	ServerName:"defaultServer",
}
type CoreServerOption struct {
	ServerOption *ServerOption
	CoreTcpSessionOption
	ops []ServerOptionSetFunc
}

func (self *CoreServerOption) Startup() error {
	self.ServerOption = &ServerOption{}
	*self.ServerOption = *defaultServerOption

	self.CoreTcpSessionOption.SessionOption = &self.ServerOption.SessionOption
	err := self.CoreTcpSessionOption.Startup()
	if err != nil {
		return err
	}
	for _, op := range self.ops {
		op(self.ServerOption)
	}
	return nil
}

func (self *CoreServerOption) SetServerOption(ops...ServerOptionSetFunc)  {
	self.ops = ops
}

func SetServerAddress(ip string, port int) ServerOptionSetFunc {
	return func(option *ServerOption) {
		option.ListenIp = ip
		option.ListenPort = port
	}
}
