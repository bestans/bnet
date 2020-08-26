package tcp

import (
	"bnet/com"
	"bnet/inet"
	netimpl "bnet/net"
	"fmt"
	"net"
	"sync"
	"time"
)

func init() {
	netimpl.RegisterNetCreator(func(interface{}) inet.INet {
		p := &TcpClient{}
		return p
	})
}

type TcpClient struct {
	com.CoreRunningTag
	com.CoreSessionManager
	CoreClientOption

	defaultSes *tcpSession

	tryConnTimes int // 尝试连接次数

	sesEndSignal sync.WaitGroup

	reconDur time.Duration
}

func (self *TcpClient) Start() inet.INet {
	if err := self.Startup(); err != nil {
		fmt.Println(err)
		return nil
	}
	self.defaultSes = newSession(nil, self, &self.ClientOption.SessionOption)

	self.WaitStopFinished()

	if self.IsRunning() {
		return self
	}

	go self.connect(self.GetConnectAddress())

	return self
}

func (self *TcpClient) Stop(closeInfo string) {
	if !self.IsRunning() {
		return
	}

	if self.IsStopping() {
		return
	}

	self.StartStopping()

	// 通知发送关闭
	self.defaultSes.Close(closeInfo)

	// 等待线程结束
	self.WaitStopFinished()
}

func (self *TcpClient) Session() inet.ISession {
	return self.defaultSes
}

func (self *TcpClient) SendAndFlush(msg interface{})  {
	//if !self.IsRunning() {
	//	fmt.Println("tcpclient is not run")
	//	return
	//}
	self.defaultSes.SendAndFlush(msg)
}

func (self *TcpClient) ReconnectDuration() time.Duration {

	return self.reconDur
}

func (self *TcpClient) SetReconnectDuration(v time.Duration) {
	self.reconDur = v
}

func (self *TcpClient) Port() int {
	return self.ClientOption.ConnectPort
}

const reportConnectFailedLimitTimes = 3

// 连接器，传入连接地址和发送封包次数
func (self *TcpClient) connect(address string) {

	self.SetRunning(true)

	for {
		self.tryConnTimes++

		// 尝试用Socket连接地址
		conn, err := net.Dial("tcp", address)

		self.defaultSes.setConn(conn)

		// 发生错误时退出
		if err != nil {

			if self.tryConnTimes <= reportConnectFailedLimitTimes {
				//log.Errorf("#tcp.connect failed(%s) %v", self.Name(), err.Error())

				if self.tryConnTimes == reportConnectFailedLimitTimes {
					//log.Errorf("(%s) continue reconnecting, but mute log", self.Name())
				}
			}

			// 没重连就退出
			if self.ReconnectDuration() == 0 || self.IsStopping() {
				break
			}

			// 有重连就等待
			time.Sleep(self.ReconnectDuration())

			// 继续连接
			continue
		}

		self.sesEndSignal.Add(1)

		self.ApplyConnOption(conn)

		self.defaultSes.Start()

		self.tryConnTimes = 0

		self.sesEndSignal.Wait()

		self.defaultSes.setConn(nil)

		// 没重连就退出/主动退出
		if self.IsStopping() || self.ReconnectDuration() == 0 {
			break
		}

		// 有重连就等待
		time.Sleep(self.ReconnectDuration())

		// 继续连接
		continue

	}

	self.SetRunning(false)

	self.EndStopping()
}

func (self *TcpClient) NetType() inet.NetType {
	return inet.TcpClient
}

func (self *TcpClient) NetName() string  {
	return self.ClientOption.ClientName
}
