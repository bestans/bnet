package framing

import (
	"encoding/binary"
	"errors"
)

//固定4字节int32长度
type LengthFixedFrameBuilder struct {
}

var defaultLengthFixedFrameBuilder = &LengthFixedFrameBuilder{}
func DefaultLengthFixedFrameBuilder() *LengthFixedFrameBuilder {
	return defaultLengthFixedFrameBuilder
}

//长度需要的字节数
var lengthFixedNeedBytes = 4

func (self *LengthFixedFrameBuilder) EncodeFrameLenSize(msgLen int) int {
	return lengthFixedNeedBytes
}

func (self *LengthFixedFrameBuilder) TryDecodeFrame(readData []byte) (messageData []byte, readLen int, success bool, err error) {
	if len(readData) < lengthFixedNeedBytes {
		return
	}
	msgLen := int(binary.BigEndian.Uint32(readData))
	readLen = msgLen + lengthFixedNeedBytes
	if len(readData) < readLen {
		return
	}
	return readData[lengthFixedNeedBytes:readLen], readLen, true, nil
}

var errBytesNotEnough = errors.New("LengthFixedFrameBuilder bytes not enough")
func (self *LengthFixedFrameBuilder) EncodeFrame(msgLen int, buf []byte) ([]byte, error) {
	bufLen := len(buf)
	bufCap := cap(buf)
	if bufCap - bufLen < lengthFixedNeedBytes {
		return buf, errBytesNotEnough
	}
	buf = buf[:bufLen+lengthFixedNeedBytes]
	binary.BigEndian.PutUint32(buf[bufLen:bufLen+lengthFixedNeedBytes], uint32(msgLen))
	return buf, nil
}
