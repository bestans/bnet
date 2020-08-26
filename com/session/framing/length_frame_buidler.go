package framing

import (
	"bnet/com/session/sharecoder"
)

//动态长度
type LengthFrameBuilder struct {
}

var defaultLengthFrameBuilder = &LengthFrameBuilder{}
func DefaultLengthFrameBuilder() *LengthFrameBuilder {
	return defaultLengthFrameBuilder
}

func (self *LengthFrameBuilder) TryDecodeFrame(readData []byte) (messageData []byte, readLen int, success bool, err error) {
	lenLen, x, err := sharecoder.DecodeVarint(readData)
	if err != nil {
		return
	}
	msgLen := int(x)
	readLen = lenLen + msgLen
	if len(readData) < readLen {
		return
	}

	return readData[lenLen:msgLen], readLen, true, nil
}

func (self *LengthFrameBuilder) EncodeFrameLenSize(msgLen int) int {
	return sharecoder.EncodeVarintSize(uint64(msgLen))
}

func (self *LengthFrameBuilder) EncodeFrame(msgLen int, buf []byte) ([]byte, error) {
	return sharecoder.EncodeVarint(buf, uint64(msgLen))
}
