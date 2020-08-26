package inet


type NetCreateFunc func(interface{}) INet

type EncodeFunc func(interface{})(*MessagePack, error)
type DecodeFunc func(data []byte)(interface{}, int, error)
type ProcessErrorFunc func(err error)
