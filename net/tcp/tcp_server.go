package tcp

import (
	"bnet/com"
	"bnet/inet"
	"bnet/log"
	netimpl "bnet/net"
	"bnet/util"
	"net"
	"time"
)

func init() {
	netimpl.RegisterNetCreator(func(config interface{}) inet.INet {
		p := &TcpServer{
			connChan:make(chan net.Conn, 10),
		}
		return p
	})
}

type TcpServer struct {
	com.CoreStartUp
	CoreServerOption
	com.CoreSessionManager
	com.CoreRunningTag

	// 保存侦听器
	listener net.Listener
	config *ServerOption

	connChan chan net.Conn
}

func (self *TcpServer) Port() int {
	if self.listener == nil {
		return 0
	}

	return self.listener.Addr().(*net.TCPAddr).Port
}

func (self *TcpServer) IsReady() bool {
	return self.IsRunning()
}

// 异步开始侦听
func (self *TcpServer) Start() inet.INet {
	err := self.StartUp(&self.CoreServerOption)
	if err != nil {
		return nil
	}

	self.WaitStopFinished()

	if self.IsRunning() {
		return self
	}

	ln, err := net.Listen("tcp", self.ListenAddress())

	if err != nil {
		self.SetRunning(false)
		return self
	}

	self.listener = ln.(net.Listener)
	go self.serveAccept()
	go self.serveMain()
	return self
}

func (self *TcpServer) ListenAddress() string {
	return util.JoinAddress(self.ServerOption.ListenIp, self.ServerOption.ListenPort)
}

func (self *TcpServer) serveMain() {
	for {
		select {
		case conn := <- self.connChan:
			self.onNewSession(conn)
		}
	}
}

func (self *TcpServer) serveAccept() {
	self.SetRunning(true)

	for {
		conn, err := self.listener.Accept()

		if self.IsStopping() {
			break
		}

		if err == nil {
			self.connChan <- conn
		}else{
			if nerr, ok := err.(net.Error); ok && nerr.Temporary(){
				time.Sleep(time.Millisecond)
				continue
			}

			log.Error("serveAccept:%v accept failed(%v)", self.NetName(), err)
			break
		}
	}

	self.SetRunning(false)
	self.EndStopping()
}

func (self *TcpServer) onNewSession(conn net.Conn) {
	self.ApplyConnOption(conn)

	ses := newSession(conn, self, &self.ServerOption.SessionOption)

	ses.Start()
}

// 停止侦听器
func (self *TcpServer) Stop(closeInfo string) {
	if !self.IsRunning() {
		return
	}

	if self.IsStopping() {
		return
	}

	self.StartStopping()

	self.listener.Close()

	close(self.connChan)

	// 断开所有连接
	self.CloseAllSession(closeInfo)

	// 等待线程结束
	self.WaitStopFinished()
}

func (self *TcpServer) NetType() inet.NetType {
	return inet.TcpServer
}

func (self *TcpServer) NetName() string {
	return self.ServerOption.ServerName
}
