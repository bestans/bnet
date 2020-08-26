package protocoder

import (
	"bnet/com/session/msgfactory"
	"bnet/inet"
	"errors"
)

type PBBaseMessage struct {
	MsgType 	int32
	Message		interface{}
	MessageHead interface{}
}

var (
	pbBaseMessageEncodeFailed = errors.New("pb BaseMessageEncode failed")
	pbBaseMessageDecodeFailed = errors.New("pb BaseMessageDecode failed")
)

func (self *PBBaseMessage) BaseMessageSize() int32 {
	return int32(inet.EncodeVarintSize(uint64(self.MsgType)) + self.Message.(newMarshaler).XXX_Size())
}

func (self *PBBaseMessage) BaseMessageEncode(b *inet.Buffer) error {
	err := b.EncodeVarint(uint64(self.MsgType))
	if err != nil {
		return pbBaseMessageEncodeFailed
	}
	_, err = self.Message.(newMarshaler).XXX_Marshal(b.LeftBytes(), false)
	if err != nil {
		return pbBaseMessageEncodeFailed
	}
	return nil
}

func (self *PBBaseMessage) BaseMessageDecode(b *inet.Buffer) error {
	msgType, err := b.DecodeVarint()
	if err != nil {
		return err
	}
	self.MsgType = int32(msgType)
	newMsg, err := msgfactory.GetPBMessageFactory().NewMessage(self.MsgType)
	if err != nil {
		return err
	}
	err = newMsg.(newUnmarshaler).XXX_Unmarshal(b.LeftBytes())
	if err != nil {
		return pbBaseMessageDecodeFailed
	}
	self.Message = newMsg
	return nil
}
