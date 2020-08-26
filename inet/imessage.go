package inet

import "github.com/golang/protobuf/proto"

//消息工厂
type IMessageFactory interface {
	NewMessage(msgTyp int32) (proto.Message, error)
}

type IProtocolCoder interface {
	ProtocolSize(msg interface{}) (totalSize int, param interface{})
	ProtocolEncode(msg interface{}, buf []byte, param interface{}) ([]byte, error)
	ProtocolDecode(buf []byte) (msg interface{}, err error, msgType int32)
}

type ReceiveMsgFunc func(ISession, interface{})
type HandleEventFunc func(session ISession, event *EventPack)
