package tcp

import (
	"bnet/com"
	"bnet/inet"
	"bnet/log"
	"bnet/util"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// Socket会话
type tcpSession struct {
	*SessionOption

	pInterface inet.INet

	// Socket原始连接
	conn      net.Conn
	connGuard sync.RWMutex

	// 退出同步器
	exitSync sync.WaitGroup

	id int64
	closing *inet.ClosePack

	readBuf		*inet.SocketBuffer
	readDecodeBuf	*inet.SocketBuffer
	writeBuf	*inet.SocketBuffer
	handler  inet.IMessageHandler

	writeCh chan *inet.MessagePack
	quitChan  chan bool
	innerEventChan chan interface{}
	closeOnce sync.Once

	cacheBufferList chan []byte
	bufferPool *com.BufferPool
	identify string
}

func (self *tcpSession) setConn(conn net.Conn) {
	self.connGuard.Lock()
	self.conn = conn
	self.connGuard.Unlock()
}

func (self *tcpSession) Conn() net.Conn {
	self.connGuard.RLock()
	defer self.connGuard.RUnlock()
	return self.conn
}

func (self *tcpSession) Peer() inet.INet {
	return self.pInterface
}

// 取原始连接
func (self *tcpSession) Raw() interface{} {
	return self.Conn()
}

func (self *tcpSession) SetID(id int64) {
	self.id = id
	if self.pInterface.NetType() == inet.TcpClient {
		self.identify = self.pInterface.NetName() + "-session"
	} else {
		self.identify = fmt.Sprintf("%v-session-%v", self.pInterface.NetName(), id)
	}
}

func (self *tcpSession) ID() int64 {
	return self.id
}
func (self *tcpSession) Identify() string {
	return self.identify
}
func (self *tcpSession) Close(closeInfo string) {
	self.innerClose(inet.NewClosePack(inet.CloseReasonCustom, closeInfo))
}
func (self *tcpSession) innerClose(reason *inet.ClosePack) {
	//var err error
	self.closeOnce.Do(func() {
		atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&self.closing)), unsafe.Pointer(reason))
		conn := self.Conn()
		if conn != nil {
			// 关闭读
			tcpConn := conn.(*net.TCPConn)
			// 关闭读
			tcpConn.CloseRead()
			// 手动读超时
			tcpConn.SetReadDeadline(time.Now())
		}
		close(self.quitChan)
	})
}

func (self *tcpSession) GetCacheBufferMessage(size int) *inet.MessagePack {
	//if size <= cacheWriteBufferSize {
	//		select {
	//		case buf := <- self.cacheBufferList:
	//			buf = buf[:0]
	//			return &inet.MessagePack{
	//				Data:buf,
	//				CacheBuffer:buf,
	//			}
	//		default:
	//			break
	//		}
	//}
	buf, _ := self.bufferPool.Get(size)
	return &inet.MessagePack{
		Data:buf,
		CacheBuffer:buf,
	}
}

func (self *tcpSession) writeMessage(msgPack *inet.MessagePack) {
	select {
	case <- self.quitChan:
		log.Error("tcpSession:writeMessage:failed:session has closed:msgId=%v", msgPack.MsgId)
		return
	case self.writeCh <- msgPack:
		return
	default:
		log.Trace("tcpSession:writeMessage:msgId=%v", msgPack.MsgId)
		break
	}

	timer := com.NewTimer(time.Second)
	defer timer.Recycle()

	var uid int64
	select {
	case <- self.quitChan:
		log.Error("tcpSession:writeMessage:failed:session has closed:msgId=%v", msgPack.MsgId)
		return
	case self.writeCh <- msgPack:
		return
	case <- timer.Get():
		uid = util.GenerateUniqueId()
		log.Error("tcpSession:writeMessage:Timeout:Identify=%v,timeout:timeoutSeq=%v,msgId=%v", self.Identify(),
			uid, msgPack.MsgId)
		break
	}

	select {
	case <-self.quitChan:
		log.Error("tcpSession:writeMessage:failed:session has closed:msgId=%v", msgPack.MsgId)
		return
	case self.writeCh <- msgPack:
		log.Error("tcpSession:writeMessage:Finish After Timeout:Identify=%v,timeout:timeoutSeq=%v,msgId=%v", self.Identify(),
			uid, msgPack.MsgId)
		return
	}
}

func (self *tcpSession) rawSend(msg interface{}, flush bool) {
	msgPack, err := self.handler.EncodeMessage(self, msg)
	if err != nil {
		log.Error("tcpSession:rawSend:err=%v,msg=%v", err, msg)
		return
	}
	msgPack.Flush = flush
	self.writeMessage(msgPack)
}

// 发送封包（不立即发送）
func (self *tcpSession) Send(msg interface{}) {
	self.rawSend(msg, false)
}

//发送封包（立即发送）
func (self *tcpSession) SendAndFlush(msg interface{}) {
	self.rawSend(msg, true)
}

var flushMessagePack = &inet.MessagePack{Flush:true}
//通知缓冲区发送消息
func (self *tcpSession) Flush() {
	self.writeMessage(flushMessagePack)
}

func (self *tcpSession) IsManualClosed() bool {
	return atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&self.closing))) != nil
}

func (self *tcpSession) handleEvent(event *inet.EventPack) {
	if event.Data != nil {
		log.Trace("%v:handleEvent=%v,data=%v", self.Identify(), event.Event, event.Data)
	} else {
		log.Trace("%v:handleEvent=%v", self.Identify(), event.Event)
	}
	self.handler.HandleEvent(self, event)
}

func (self *tcpSession) onInnerEvent(data interface{}) {
	select {
	case self.innerEventChan <- data:
	case <- self.quitChan:
	}
}
func (self *tcpSession) processInnerEvent(innerEvent interface{}) {
	switch innerEvent.(type) {
	case *inet.ClosePack:
		self.innerClose(innerEvent.(*inet.ClosePack))
	case *inet.EventPack:
		self.handleEvent(innerEvent.(*inet.EventPack))
	default:
		log.Error("processInnerEvent:unhandled event,%v", innerEvent)
	}
}
func (self *tcpSession) serveMain() {
	for {
		select {
		case innerEvent := <- self.innerEventChan:
			self.processInnerEvent(innerEvent)
		case <- self.quitChan:
			// 等待2个任务结束
			self.exitSync.Wait()
			//self.handler.SessionUnregistered(self)
			// 将会话从管理器移除
			self.Peer().(inet.SessionManager).Remove(self)
			self.handleEvent(inet.NewEventPack(inet.EventSessionClose, self.closing))
			return
		}
	}
}

func (self *tcpSession) serveWrite()  {
	defer func() {
		// 完整关闭
		conn := self.Conn()
		if conn != nil {
			conn.Close()
		}

		// 通知完成
		self.exitSync.Done()
	}()

	for {
		select {
		case <- self.quitChan:
			return

		case data := <- self.writeCh:
			if data.IsEnd {	//用if判断替换quitChan，吞吐量提高1%
				return
			}
			writeFunc := func(writeData []byte) error {
				for {
					if len(writeData) <= 0 {
						break
					}
					n, err := self.conn.Write(writeData)
					if err != nil {
						self.onInnerEvent(inet.NewClosePack(inet.CloseReasonCustom, "ErrConnWrite"))
						return err
					}
					if n >= len(writeData){
						break
					}
					writeData = writeData[n:]
				}
				return nil
			}
			var err error
			for {
				if data.Flush { //当消息flush且writeBuf为空时
					if data.Data == nil {
						err = writeFunc(self.writeBuf.ReadByteData())
						break
					} else if self.writeBuf.ReadableBytes() <= 0 {
						err = writeFunc(data.Data)
						break
					}
				}
				if self.writeBuf.WritableBytes() < len(data.Data) {
					if err = writeFunc(self.writeBuf.ReadByteData()); err != nil {
						break
					}
					self.writeBuf.Reset()
					self.writeBuf.WriteBytes(data.Data)
				} else {
					self.writeBuf.WriteBytes(data.Data)
				}
				if data.Flush {
					err = writeFunc(self.writeBuf.ReadByteData())
					if err != nil {
						break
					}
					self.writeBuf.Reset()
				}
				break
			}
			//返还cacheBuffer
			self.bufferPool.Put(data.CacheBuffer)
			if err != nil {
				return
			}
		}
	}
}

func (self *tcpSession) serveRead()  {
	defer func() {
		// 通知完成
		self.exitSync.Done()
	}()

	for {
		n, err := self.conn.Read(self.readBuf.WriteByteData())
		if err != nil {
			if !util.IsEOFOrNetReadError(err) {
				self.onInnerEvent(inet.NewClosePack(inet.CloseReasonErrConnRead, err))
			} else {
				self.onInnerEvent(inet.NewClosePack(inet.CloseReasonRemoteClose, err))
			}
			return
		}
		self.readBuf.AddWriteIndex(n)
		self.readBuf, self.readDecodeBuf = self.readDecodeBuf, self.readBuf

		for {
			readLen, err := self.handler.DecodeMessage(self, self.readDecodeBuf.ReadByteData())
			if err != nil {
				self.onInnerEvent(inet.NewClosePack(inet.CloseReasonCustom, "DecodeMessage failed"))
				return
			}
			//if readLen > self.SessionOption.ReadBufferSize {
			//	msgType, _ := self.handler.DecodePartMessage(self, self.readDecodeBuf.ReadByteData())
			//	self.onInnerEvent(inet.NewClosePack(inet.CloseReasonCustom, fmt.Sprintf("part message:msgType=%v", msgType)))
			//	return
			//}
			if readLen <= 0 {
				break
			}
			self.readDecodeBuf.AddReadIndex(readLen)
		}
		if self.readDecodeBuf.ReadableBytes() <= 0 {
			//readDecodeBuf都解析完了
			self.readDecodeBuf.Reset()
		}
		if self.readDecodeBuf.ReadableBytes() >= 0 && self.readDecodeBuf.TriggerUpdate() {
			//还有没解析完的readBuffer，并且触发了更新
			if self.readDecodeBuf.Full() {
				msgInfo, err := self.handler.DecodePartMessage(self, self.readDecodeBuf.ReadByteData())
				self.onInnerEvent(inet.NewClosePack(inet.CloseReasonCustom, fmt.Sprintf("readDecodeBuf full:part message:msgType=%v,err=%v,readBufferSize=%v", msgInfo, err, self.SessionOption.ReadBufferSize)))
				return
			}
			self.readBuf.WriteBytes(self.readDecodeBuf.ReadByteData())
			self.readDecodeBuf.Reset()
		} else {
			//没解析完，并且没有触发数据更新
			self.readBuf, self.readDecodeBuf = self.readDecodeBuf, self.readBuf
		}
	}
}

// 启动会话的各种资源
func (self *tcpSession) Start() {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&self.closing)), nil)

	// 需要接收和发送线程同时完成时才算真正的完成
	self.exitSync.Add(2)

	// 将会话添加到管理器, 在线程处理前添加到管理器(分配id), 避免ID还未分配,就开始使用id的竞态问题
	if sesMana, ok := self.Peer().(inet.SessionManager); ok {
		sesMana.Add(self)
	}
	//self.handler.SessionRegistered(self)

	go self.serveMain()

	// 启动并发接收goroutine
	go self.serveRead()

	// 启动并发发送goroutine
	go self.serveWrite()

	self.onInnerEvent(inet.NewEventPack(inet.EventSessionOpen, nil))
}

func newSession(conn net.Conn, p inet.INet, sessionOption *SessionOption) *tcpSession {
	self := &tcpSession{
		conn:          conn,
		SessionOption: sessionOption,
		pInterface:    p,
		readBuf:       inet.NewSocketBuffer(sessionOption.ReadBufferSize),
		readDecodeBuf: inet.NewSocketBuffer(sessionOption.ReadBufferSize),
		writeBuf:      inet.NewSocketBuffer(sessionOption.WriteBufferSize),
		handler:       sessionOption.Handler,
		writeCh:       make(chan *inet.MessagePack, 10),
		quitChan:      make(chan bool, 1),
		innerEventChan: make(chan interface{}, 5),
		cacheBufferList: make(chan []byte, cacheWriteBufferNum),
		bufferPool: com.NewDefaultBufferPool(),
	}
	//初始化cacheBufferList
	for i := 0; i < cacheWriteBufferNum; i++ {
		self.cacheBufferList <- make([]byte, 0, cacheWriteBufferSize)
	}
	return self
}

