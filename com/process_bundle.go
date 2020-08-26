package com

import (
	"bnet/inet"
	"errors"
)

type MessagePoster interface {
	// 投递一个消息到Hooker之前
	ProcEvent(ev inet.Event)
}

type CoreProcBundle struct {
	transmit inet.MessageTransmitter
	hooker   inet.EventHooker
	callback inet.EventCallback
}

func (self *CoreProcBundle) GetBundle() *CoreProcBundle {
	return self
}

func (self *CoreProcBundle) SetTransmitter(v inet.MessageTransmitter) {
	self.transmit = v
}

func (self *CoreProcBundle) SetHooker(v inet.EventHooker) {
	self.hooker = v
}

func (self *CoreProcBundle) SetCallback(v inet.EventCallback) {
	self.callback = v
}

var notHandled = errors.New("Processor: Transimitter nil")

func (self *CoreProcBundle) ReadMessage(ses inet.ISession) (msg interface{}, err error) {

	if self.transmit != nil {
		return self.transmit.OnRecvMessage(ses)
	}

	return nil, notHandled
}

func (self *CoreProcBundle) SendMessage(ev inet.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnOutboundEvent(ev)
	}

	if self.transmit != nil && ev != nil {
		self.transmit.OnSendMessage(ev.Session(), ev.Message())
	}
}

func (self *CoreProcBundle) ProcEvent(ev inet.Event) {

	if self.hooker != nil {
		ev = self.hooker.OnInboundEvent(ev)
	}

	if self.callback != nil && ev != nil {
		self.callback(ev)
	}
}
