package com

import (
	"bnet/com/session/coding"
	"bnet/com/session/framing"
	"bnet/inet"
)

type SessionProcessNull struct {

}

func (s SessionProcessNull) SessionRegistered(ses inet.ISession) error {
	return nil
}

func (s SessionProcessNull) SessionUnregistered(ses inet.ISession) error {
	return nil
}

func (s SessionProcessNull) SessionOnError(ses inet.ISession) error {
	return nil
}

type CoreSessionHandler struct {
	coding.StringCodingBuilder
	framing.LengthFixedFrameBuilder
	SessionProcessNull
}

func (c CoreSessionHandler) ReadMessage(ses inet.ISession, decodeMsg interface{}) (err error) {
	//fmt.Println(decodeMsg)
	return nil
}
