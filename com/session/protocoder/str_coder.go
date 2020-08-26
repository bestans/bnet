package protocoder

import (
	"bnet/com/session/sharecoder"
)

type StringCoder struct {

}
var defaultStringCoder = &StringCoder{}

func DefaultStringCoder() *StringCoder {
	return defaultStringCoder
}

func (self *StringCoder) ProtocolSize(msg interface{}) (totalSize int, param interface{}) {
	return len(msg.(string)), nil
}

func (self *StringCoder) ProtocolEncode(msg interface{}, buf []byte, param interface{}) ([]byte, error) {
	return sharecoder.EncodeString(buf, msg.(string))
}

func (self *StringCoder) ProtocolDecode(buf []byte) (str interface{}, err error, msgType int32) {
	str, _, err = sharecoder.DecodeString(buf, len(buf))
	return
}
