package tcp

import (
	"bnet/com/session/msghandler"
	"bnet/inet"
	"net"
	"time"
)

type SessionOptionSetFunc func(*SessionOption)

type SessionOption struct {
	ReadBufferSize  int
	WriteBufferSize int
	NoDelay         bool
	maxPacketSize   int

	readTimeout  time.Duration
	writeTimeout time.Duration

	Handler  inet.IMessageHandler
}

//默认sessionOption
var defaultSessionOption = &SessionOption{
	ReadBufferSize:  65535,
	WriteBufferSize: 65535,
	NoDelay:         true,
	maxPacketSize:   65535,
	readTimeout:     0,
	writeTimeout:    0,
}

type CoreTcpSessionOption struct {
	SessionOption *SessionOption
	ops []SessionOptionSetFunc
}

func (self *CoreTcpSessionOption) Startup() error {
	if self.SessionOption == nil {
		self.SessionOption = &SessionOption{}
	}
	*self.SessionOption = *defaultSessionOption
	if self.SessionOption.Handler == nil {
		self.SessionOption.Handler = msghandler.NewMessageHandler()
	}

	for _, op := range self.ops {
		op(self.SessionOption)
	}
	self.ops = nil
	return nil
}

func (self *CoreTcpSessionOption) SetSessionOption(ops...SessionOptionSetFunc) {
	self.ops = ops
}

func (self *CoreTcpSessionOption) ApplyConnOption(conn net.Conn) {
	if cc, ok := conn.(*net.TCPConn); ok {
		if self.SessionOption.ReadBufferSize > 0 {	//预留一次性读不完的情况
			cc.SetReadBuffer(self.SessionOption.ReadBufferSize + 1024)
		}
		if self.SessionOption.WriteBufferSize > 0 {
			cc.SetWriteBuffer(self.SessionOption.WriteBufferSize + 1024)
		}
		cc.SetNoDelay(self.SessionOption.NoDelay)
	}
}

func SetMessageHandler(handler inet.IMessageHandler) SessionOptionSetFunc {
	return func(option *SessionOption) {
		option.Handler = handler
	}
}

func SetMessageHandlerRecvFunc(recvFunc inet.ReceiveMsgFunc) SessionOptionSetFunc {
	return func(option *SessionOption) {
		option.Handler.SetReceiveMessageFunc(recvFunc)
	}
}
func SetMessageHandlerEventFunc(f inet.HandleEventFunc) SessionOptionSetFunc {
	return func(option *SessionOption) {
		option.Handler.SetHandleEventFunc(f)
	}
}

func SetProtocolCoder(coder inet.IProtocolCoder) SessionOptionSetFunc {
	return func(option *SessionOption) {
		option.Handler.SetProtocolCoder(coder)
	}
}

func SetSessionBufferSize(readBufferSize, writeBufferSize int, noDelay bool) SessionOptionSetFunc {
	return func(option *SessionOption) {
		option.ReadBufferSize = readBufferSize
		option.WriteBufferSize = writeBufferSize
		option.NoDelay = noDelay
	}
}

func SetSessionSocketDeadline(read, write time.Duration) SessionOptionSetFunc {
	return func(option *SessionOption) {
		option.readTimeout = read
		option.writeTimeout = write
	}
}

func SetSessionMaxPacketSize(maxSize int) SessionOptionSetFunc {
	return func(option *SessionOption) {
		option.maxPacketSize = maxSize
	}
}
