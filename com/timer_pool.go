package com

import (
	"sync"
	"time"
)
var _timerPool = sync.Pool{}
type TimerInPool struct {
	timer *time.Timer
}

func (t *TimerInPool) Get() <-chan time.Time {
	return t.timer.C
}
func (t *TimerInPool) Recycle() {
	putTimer(t)
}
func NewTimer(d time.Duration) *TimerInPool {
	if v := _timerPool.Get(); v != nil {
		t := v.(*TimerInPool)
		t.timer.Reset(d)
		return t
	}
	return &TimerInPool{
		timer:time.NewTimer(d),
	}
}
func putTimer(t *TimerInPool) {
	if !t.timer.Stop() {
		select {
		case <- t.timer.C:
		default:
		}
	}

	_timerPool.Put(t)
}

var _timerWithAfterFuncPool = sync.Pool{}

type timerWithAfterFunc struct {
	t *time.Timer
	f func()
}
// GetTimerWithAfterFunc get a timer from pool or create from time.AfterFunc
func GetTimerWithAfterFunc(d time.Duration, f func()) *timerWithAfterFunc {
	if v := _timerWithAfterFuncPool.Get(); v != nil {
		t := v.(*timerWithAfterFunc)
		t.f = f
		t.t.Reset(d)
		return t
	}
	tf := &timerWithAfterFunc{f: f}
	tf.t = time.AfterFunc(d, func() {
		tf.f()
	})
	return tf
}

// PutTimerWithAfterFunc stop a timer and return into pool
func PutTimerWithAfterFunc(t *timerWithAfterFunc) {
	t.t.Stop()
	// time.AfterFunc使用的timer.C是nil, 不需要clean
	_timerWithAfterFuncPool.Put(t)
}
