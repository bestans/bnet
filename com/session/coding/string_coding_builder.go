package coding

import (
	"bnet/inet"
)

type StringCodingBuilder struct {
}

func (s *StringCodingBuilder) EncodeMessage(encodeMsg interface{}) ([]byte, error) {
	msg := encodeMsg.(string)
	return []byte(msg), nil
}

func (s *StringCodingBuilder) DecodeMessage(ses inet.ISession, data []byte) (decodeMsg interface{}, err error) {
	msg := string(data)
	if ses.Peer().NetType() == inet.TcpServer {
		ses.Send("1")
	}
	return msg, nil
}
