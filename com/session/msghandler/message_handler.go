package msghandler

import (
	"bnet/com/session/framing"
	"bnet/com/session/protocoder"
	"bnet/inet"
	"bnet/log"
	"fmt"
)

type MessageHandler struct {
	frame      inet.IFrameProcess
	protoCoder inet.IProtocolCoder
	recvFunc   inet.ReceiveMsgFunc
	eventFunc  inet.HandleEventFunc
}

func defaultRecvFunc(ses inet.ISession, msg interface{})  {
	fmt.Printf("receive message:%v\n", msg)
}

func defaultHandleEvent(session inet.ISession, event *inet.EventPack) {

}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		framing.DefaultLengthFixedFrameBuilder(),
		protocoder.DefaultPBCoder(),
		defaultRecvFunc,
		defaultHandleEvent,
	}
}

func (self *MessageHandler) EncodeMessage(ses inet.ISession, msg interface{}) (msgPack *inet.MessagePack, err error) {
	msgSize, param := self.protoCoder.ProtocolSize(msg)
	lenSize := self.frame.EncodeFrameLenSize(msgSize)
	totalSize := msgSize + lenSize
	msgPack = ses.GetCacheBufferMessage(totalSize)
	msgPack.Data, err = self.frame.EncodeFrame(msgSize, msgPack.Data)
	if err != nil {
		log.Error("bNet:MessageHandler:EncodeMessage:EncodeFrame failed, err=%v, msg=%v", err, msg)
		return
	}
	msgPack.Data, err = self.protoCoder.ProtocolEncode(msg, msgPack.Data, param)
	if err != nil {
		log.Error("bNet:MessageHandler:EncodeMessage:ProtocolEncode failed, err=%v, msg=%v", err, msg)
		return
	}
	return
}
func (self *MessageHandler) DecodePartMessage(ses inet.ISession, data []byte) (msgInfo string, err error) {
	msgData, readLen, _, err := self.frame.TryDecodeFrame(data)
	if err != nil {
		ses.Close("MessageHandler:TryDecodeFrame failed")
		return
	}
	_, _, msgType := self.protoCoder.ProtocolDecode(msgData)
	msgInfo = fmt.Sprintf("[msgType=%v,len=%v]", msgType, readLen)
	return
}
func (self *MessageHandler) DecodeMessage(ses inet.ISession, data []byte) (readLen int, err error) {
	msgData, readLen, success, err := self.frame.TryDecodeFrame(data)
	if err != nil {
		ses.Close("MessageHandler:TryDecodeFrame failed")
		return
	}
	if !success {
		readLen = 0
		return
	}
	decodeMsg, err, _ := self.protoCoder.ProtocolDecode(msgData)
	if err == nil {
		self.recvFunc(ses, decodeMsg)
	} else {
		log.Error("bNet:MessageHandler:DecodeMessage:ProtocolDecode failed, err=%v", err)
	}
	return
}

func (self *MessageHandler) SetReceiveMessageFunc(recvFunc inet.ReceiveMsgFunc) {
	self.recvFunc = recvFunc
}

func (self *MessageHandler) SetProtocolCoder(coder inet.IProtocolCoder)  {
	self.protoCoder = coder
}

func (self *MessageHandler) SetHandleEventFunc(f inet.HandleEventFunc) {
	self.eventFunc = f
}

func (self *MessageHandler) HandleEvent(session inet.ISession, event *inet.EventPack)  {
	self.eventFunc(session, event)
}
