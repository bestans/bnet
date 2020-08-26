package com

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestBufferPool(t *testing.T)  {
	startByte := 5
	poolSize := 5
	p := NewBufferPool(startByte, poolSize)
	for i := 1; i < (1<<15); i++ {
		buf, index := p.Get(i)
		if len(buf) != 0 {
			t.Error("buf len not 0:", len(buf))
		}
		p.Put(buf[0:10])
		if i <= 1<< startByte {
			if cap(buf) != 1 << startByte || index != 0 {
				t.Error(cap(buf), 1 << startByte)
				return
			}
		} else if i > 1<<(startByte + poolSize-1) {
			if cap(buf) != i || index != -1{
				t.Error(cap(buf), 1<<(startByte + poolSize-1))
				return
			}
		} else {
			if cap(buf) != 1 << (startByte + index) {
				t.Error(cap(buf), i, index, 1 << (startByte+index))
				return
			}
		}
	}
}

func TestBufferPool2(t *testing.T)  {
	startByte := 5
	poolSize := 6
	p := NewBufferPool(startByte, poolSize)
	testSize := []int{0, 10,31,32,33,63,64,127,128,256,512,1024,2048,4096}
	for _, value := range testSize {
		buf, index := p.Get(value)
		putIndex := p.Put(buf[:10])
		fmt.Printf("reqSize=%v,resSize=%v,bufLen=%v,index=%v,putIndex=%v\n", value, cap(buf), len(buf), index, putIndex)
		if putIndex != index {
			t.Error("failed")
		}
	}
}

type NetTimer struct {
	timer *time.Timer
	flag int32
}

func NewNetTimer() *NetTimer {
	return &NetTimer{
		timer: time.NewTimer(time.Second),
		flag:  0,
	}
}

func (t *NetTimer) Stop()  {
	if !atomic.CompareAndSwapInt32(&t.flag, 0, 1) {
		return
	}

	if !t.timer.Stop() {
		<- t.timer.C

		fmt.Println("stop wait finish")
	} else {
		fmt.Println("stop direct")
	}
}
