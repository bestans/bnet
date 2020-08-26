package com

import (
	"bnet/inet"
	"sync"
	"sync/atomic"
)

type CoreSessionManager struct {
	sesById sync.Map // 使用Id关联会话

	sesIDGen int64 // 记录已经生成的会话ID流水号

	count int64 // 记录当前在使用的会话数量
}

func (self *CoreSessionManager) SetIDBase(base int64) {
	atomic.StoreInt64(&self.sesIDGen, base)
}

func (self *CoreSessionManager) Count() int {
	return int(atomic.LoadInt64(&self.count))
}

func (self *CoreSessionManager) Add(ses inet.ISession) {

	id := atomic.AddInt64(&self.sesIDGen, 1)

	atomic.AddInt64(&self.count, 1)

	ses.(interface {
		SetID(int64)
	}).SetID(id)

	self.sesById.Store(id, ses)
}

func (self *CoreSessionManager) Remove(ses inet.ISession) {

	self.sesById.Delete(ses.ID())

	atomic.AddInt64(&self.count, -1)
}

// 获得一个连接
func (self *CoreSessionManager) GetSession(id int64) inet.ISession {
	if v, ok := self.sesById.Load(id); ok {
		return v.(inet.ISession)
	}

	return nil
}

func (self *CoreSessionManager) VisitSession(callback func(inet.ISession) bool) {

	self.sesById.Range(func(key, value interface{}) bool {

		return callback(value.(inet.ISession))

	})
}

func (self *CoreSessionManager) CloseAllSession(closeInfo string) {

	self.VisitSession(func(ses inet.ISession) bool {

		ses.Close(closeInfo)

		return true
	})
}

// 活跃的会话数量
func (self *CoreSessionManager) SessionCount() int {

	v := atomic.LoadInt64(&self.count)

	return int(v)
}

