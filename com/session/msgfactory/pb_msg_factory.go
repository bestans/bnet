package msgfactory

import (
	"bnet/inet"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
)

var pbMessageMap = make(map[int32]reflect.Type)
var pbMessageReflectType2Type = make(map[reflect.Type]int32)

var pbBaseMessage inet.IProtocolCoder
func RegisterPBMessage(msgType int32, msgReflectType reflect.Type)  {
	if reflect.New(msgReflectType.Elem()).Interface().(proto.Message) == nil {
		panic(fmt.Sprintf("RegisterMessage:cannot new message, msgType(%v)", msgType))
	}

	if _, ok := pbMessageMap[msgType]; ok {
		panic(fmt.Sprintf("RegisterMessage:duplicate msgType(%v)", msgType))
	}
	pbMessageMap[msgType] = msgReflectType.Elem()
	pbMessageReflectType2Type[msgReflectType] = msgType
}

func RegisterPBBaseMessage(baseMsg inet.IProtocolCoder)  {
	pbBaseMessage = baseMsg
}

var pbNewMessageFailed = errors.New("cannot new message")

func GetPBMessageType(msgReflectType reflect.Type) int32 {
	if msgType, ok := pbMessageReflectType2Type[msgReflectType]; ok {
		return msgType
	}
	return 0
}
func NewPBMessage(msgType int32) (proto.Message, error) {
	if msgReflectType, ok := pbMessageMap[msgType]; ok {
		return reflect.New(msgReflectType).Interface().(proto.Message), nil
	}
	return nil, pbNewMessageFailed
}

type PBMessageFactory struct {}

func (self *PBMessageFactory) NewMessage(msgType int32) (proto.Message, error) {
	return NewPBMessage(msgType)
}

var pbMsgFactory PBMessageFactory
func GetPBMessageFactory() *PBMessageFactory {
	return &pbMsgFactory
}
