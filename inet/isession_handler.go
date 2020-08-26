package inet

//封包/解包接口
type IFrameProcess interface {
	//解包，session read协程中调用
	TryDecodeFrame(readData []byte) (messageData []byte, readLen int, success bool, err error)
	//封包，session write协程中调用
	EncodeFrame(msgLen int, buf []byte) ([]byte, error)
	EncodeFrameLenSize(msgLen int) int
}

//消息处理
type ICodingProcess interface {
	//消息编码，所有协程可能调用（发送消息时）
	EncodeMessage(ses ISession, encodeMsg interface{})(*MessagePack, error)
	//消息解码，session read协程中调用
	DecodeMessage(ses ISession, data []byte)(decodeMsg interface{}, err error)
}

type IMessageProcess interface {
	ReadMessage(ses ISession, decodeMsg interface{})(err error)
}

type ISessionProcess interface {
	SessionRegistered(ses ISession) error
	SessionUnregistered(ses ISession) error
	SessionOnError(ses ISession) error
}

type IMessageHandler interface {
	EncodeMessage(ses ISession, msg interface{}) (msgPack *MessagePack, err error)
	DecodeMessage(ses ISession, data []byte) (readLen int, err error)
	DecodePartMessage(ses ISession, data []byte) (msgInfo string, err error)
	SetReceiveMessageFunc(ReceiveMsgFunc)
	SetProtocolCoder(coder IProtocolCoder)
	SetHandleEventFunc(f HandleEventFunc)
	HandleEvent(session ISession, event *EventPack)
}
