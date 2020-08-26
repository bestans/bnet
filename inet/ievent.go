package inet

import "fmt"

// 事件
type Event interface {

	// 事件对应的会话
	Session() ISession

	// 事件携带的消息
	Message() interface{}
}

// 消息收发器
type MessageTransmitter interface {

	// 接收消息
	OnRecvMessage(ses ISession) (msg interface{}, err error)

	// 发送消息
	OnSendMessage(ses ISession, msg interface{}) error
}

// 处理钩子(参数输入, 返回输出, 不给MessageProccessor处理时，可以将Event设置为nil)
type EventHooker interface {

	// 入站(接收)的事件处理
	OnInboundEvent(input Event) (output Event)

	// 出站(发送)的事件处理
	OnOutboundEvent(input Event) (output Event)
}

// 用户端处理
type EventCallback func(ev Event)


type CloseReason int32

const (
	CloseReasonCustom     CloseReason = iota // 普通IO断开
	CloseReasonError			//session内部错误，导致session断开
	CloseReasonRemoteClose		//远程断开
	CloseReasonErrConnRead		//读错误
)

func (cr CloseReason) String() string {
	return [...]string{
		"CloseReasonCustom",
		"CloseReasonError",
		"CloseReasonRemoteClose",
		"CloseReasonErrConnRead",
	}[cr]
}

type EventType int
const (
	EventSessionOpen EventType = iota
	EventSessionClose
)
var etString = [...]string{
	"EventSessionOpen",
	"EventSessionClose",
}
func (et EventType) String() string {
	return etString[et]
}
type ClosePack struct {
	CloseReason CloseReason
	Data interface{}
}

func NewClosePack(reason CloseReason, data interface{}) *ClosePack {
	return &ClosePack{
		CloseReason:reason,
		Data: data,
	}
}

func (cp *ClosePack) String() string {
	if cp == nil {
		return "nil"
	}
	return fmt.Sprintf("[%v,info=%v]", cp.CloseReason, cp.Data)
}

type EventPack struct {
	Event EventType
	Data interface{}
}

func NewEventPack(event EventType, data interface{}) *EventPack {
	return &EventPack{event, data}
}
