package inet

type INet interface {
	// 开启端，传入地址
	Start() INet

	// 停止通讯端
	Stop(closeReason string)

	// Peer的类型(protocol.type)，例如tcp.Connector/udp.Acceptor
	NetType() NetType
	NetName() string
}

type ITcpNetServer interface {
	// 开启端，传入地址
	Start() INet

	// 停止通讯端
	Stop()

	// Peer的类型(protocol.type)，例如tcp.Connector/udp.Acceptor
	NetType() string
}

type ITcpNetClient interface {
	// 开启端，传入地址
	Start() INet

	// 停止通讯端
	Stop()

	// Peer的类型(protocol.type)，例如tcp.Connector/udp.Acceptor
	NetType() string
}

type IStartup interface {
	Startup() error
}

//net中的handler注册初始化
type IInitializerNet interface {
	Initializer() error
}
