package msghandler

import (
	"bnet/inet"
	"fmt"
)

type ZlMessageHandler struct {
	MessageHandler
}

func (self *ZlMessageHandler) EncodeMessage(ses inet.ISession, msg interface{}) (msgPack *inet.MessagePack, err error) {
	msgSize, param := self.protoCoder.ProtocolSize(msg)
	lenSize := self.frame.EncodeFrameLenSize(msgSize)
	totalSize := msgSize + lenSize
	msgPack = ses.GetCacheBufferMessage(totalSize)
	msgPack.Data, err = self.frame.EncodeFrame(msgSize, msgPack.Data)
	if err != nil {
		return
	}
	msgPack.Data, err = self.protoCoder.ProtocolEncode(msg, msgPack.Data, param)
	if err != nil {
		return
	}
	return
}
func (self *ZlMessageHandler) DecodePartMessage(ses inet.ISession, data []byte) (msgInfo string, err error) {
	msgData, readLen, _, err := self.frame.TryDecodeFrame(data)
	if err != nil {
		ses.Close("ZlMessageHandler:TryDecodeFrame failed")
		return
	}
	_, _, msgType := self.protoCoder.ProtocolDecode(msgData)
	msgInfo = fmt.Sprintf("[msgType=%v,len=%v]", msgType, readLen)
	return
}
func (self *ZlMessageHandler) DecodeMessage(ses inet.ISession, data []byte) (readLen int, err error) {
	msgData, readLen, success, err := self.frame.TryDecodeFrame(data)
	if err != nil {
		ses.Close("ZlMessageHandler:TryDecodeFrame failed")
		return
	}
	if !success {
		return
	}
	decodeMsg, err, _ := self.protoCoder.ProtocolDecode(msgData)
	if err == nil {
		self.recvFunc(ses, decodeMsg)
	}
	return
}
