package protocoder

import (
	"bnet/inet"
	"github.com/golang/protobuf/proto"
	"reflect"
)

type PBProtocolCoder struct {
	baseMessageReflectType reflect.Type
}

func NewPBProtocolCoder(baseMessageReflectType reflect.Type) *PBProtocolCoder {
	return &PBProtocolCoder{baseMessageReflectType:baseMessageReflectType}
}

type newMarshaler interface {
	XXX_Size() int
	XXX_Marshal(b []byte, deterministic bool) ([]byte, error)
}

type newUnmarshaler interface {
	XXX_Unmarshal([]byte) error
}

func (self *PBProtocolCoder) ProtocolSize(msg interface{}) int32 {
	return int32(msg.(newMarshaler).XXX_Size())
}

func (self *PBProtocolCoder) ProtocolEncode(msg interface{}, buffer *inet.Buffer) error {
	_, err := msg.(newMarshaler).XXX_Marshal(buffer.LeftBytes(), false)
	return err
}

func (self *PBProtocolCoder) ProtocolDecode(buffer *inet.Buffer) (interface{}, error) {
	newMsg := reflect.New(self.baseMessageReflectType).Interface().(proto.Message)
	return newMsg.(newMarshaler).XXX_Marshal(buffer.Bytes(), false)
}
