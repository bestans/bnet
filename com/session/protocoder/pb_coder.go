package protocoder

import (
	"bnet/com/session/msgfactory"
	"bnet/com/session/sharecoder"
	"bnet/inet"
	"errors"
	"reflect"
)

type PBCoder struct {
}

var defaultPBCoder = &PBCoder{}

func DefaultPBCoder() *PBCoder  {
	return defaultPBCoder
}

func NewPBCoder() *PBCoder {
	return &PBCoder{}
}

type pbCoderParam struct {
	msgType int32
	msgSize int
}

func (self *PBCoder) ProtocolSize(msg interface{}) (int, interface{}) {
	msgType := msgfactory.GetPBMessageType(reflect.TypeOf(msg))
	msgSize := msg.(newMarshaler).XXX_Size()
	totalSize := msgSize + inet.EncodeVarintSize(uint64(msgType))
	return totalSize, &pbCoderParam{msgType, msgSize}
}

func (self *PBCoder) ProtocolEncode(msg interface{}, buffer []byte, param interface{}) ([]byte, error) {
	p := param.(*pbCoderParam)
	msgType := p.msgType
	if msgType == 0 {
		return nil, errors.New("cannot find message type")
	}
	var err error
	if buffer, err = sharecoder.EncodeVarint(buffer, uint64(msgType)); err != nil {
		return nil, err
	}
	buffer, err = msg.(newMarshaler).XXX_Marshal(buffer, false)
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (self *PBCoder) ProtocolDecode(buf []byte) (newMsg interface{}, err error, msgType int32) {
	readLen, messageType, err := sharecoder.DecodeVarint(buf)
	if err != nil {
		return
	}
	msgType = int32(messageType)
	newMsg, err = msgfactory.NewPBMessage(msgType)
	if err != nil {
		return
	}
	err = newMsg.(newUnmarshaler).XXX_Unmarshal(buf[readLen:])
	return
}
