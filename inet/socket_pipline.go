package inet

type ISocketDecode interface {
	Decode(in SocketBuffer) error
}

type ISocketEncode interface {
	Encode(ses ISession, msg interface{}, out SocketBuffer) error
}
