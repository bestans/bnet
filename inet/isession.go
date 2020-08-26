package inet

import (
	"bytes"
	"fmt"
)

// 长连接
type ISession interface {

	// 获得原始的Socket连接
	Raw() interface{}

	// 获得Session归属的Peer
	Peer() INet

	// 发送消息，消息需要以指针格式传入
	Send(msg interface{})

	SendAndFlush(msg interface{})

	Flush()

	// 断开
	Close(closeInfo string)

	// 标示ID
	ID() int64
	GetCacheBufferMessage(size int) *MessagePack
}


// 会话访问
type SessionAccessor interface {

	// 获取一个连接
	GetSession(int64) ISession

	// 遍历连接
	VisitSession(func(ISession) bool)

	// 连接数量
	SessionCount() int

	// 关闭所有连接
	CloseAllSession(closeInfo string)
}

// 完整功能的会话管理
type SessionManager interface {
	SessionAccessor

	Add(ISession)
	Remove(ISession)
	Count() int

	// 设置ID开始的号
	SetIDBase(base int64)
}

type MessagePack struct {
	MsgId 	int32
	Data 	[]byte
	Flush	bool
	CacheBuffer   []byte
	IsEnd bool
}

func (mp *MessagePack) GetCacheBuffer() []byte {
	return mp.CacheBuffer
}

type ErrorPack struct {
	Err error
	ErrCode ErrorCode
	Info interface{}
}

func (ep *ErrorPack) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(ep.ErrCode.String())
	buffer.WriteString("")
	if ep.Info != nil {
		buffer.WriteString(fmt.Sprintf(",%v", ep.Info))
	}
	if ep.Err != nil {
		buffer.WriteString(fmt.Sprintf(",%v", ep.Err))
	}
	buffer.WriteString("]")
	return buffer.String()
}
func NewErrorPack(code ErrorCode, err error, info interface{}) *ErrorPack  {
	return &ErrorPack{
		Err: err,
		ErrCode:code,
		Info: info,
	}
}
