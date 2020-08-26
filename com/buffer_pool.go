package com

import (
	"fmt"
	"sync"
)

type BufferPool struct {
	pools []*sync.Pool
	startByteCount int
}

var defaultStartByteBit = 5
var defaultPoolSize = 6
//按字节大小分段：32 64 128 256 512 1024
func NewDefaultBufferPool() *BufferPool {
	return NewBufferPool(defaultStartByteBit, defaultPoolSize)
}
func NewBufferPool(startByte, size int) *BufferPool {
	if startByte <= 1 || startByte + size - 1  >= 15 {
		fmt.Printf("invalid param:size,max bufer size=%v", 1 <<14)
		return nil
	}
	pools := &BufferPool{
		pools:make([]*sync.Pool, size),
		startByteCount:startByte,
	}
	for i := 0; i < len(pools.pools); i++ {
		count := 1 << (i + startByte)
		pools.pools[i] = &sync.Pool{
			New: func() interface{} {
				return make([]byte, 0, count)
			},
		}
	}
	return pools
}
func (bp *BufferPool) Get(bufSize int) ([]byte, int) {
	srcSize := bufSize
	if bufSize > 1 {
		bufSize = bufSize - 1
	}
	bufSize = bufSize >> bp.startByteCount
	if bufSize == 0 {
		return bp.pools[0].Get().([]byte), 0
	}
	for i := 1; i < len(bp.pools); i++ {
		bufSize = bufSize & ^(1 << (i-1))
		if bufSize == 0 {
			return bp.pools[i].Get().([]byte), i
		}
	}
	return make([]byte, 0, srcSize), -1
}
func (bp *BufferPool) Put(buf []byte) int {
	bufSize := cap(buf)
	for i := 0; i < len(bp.pools); i++ {
		if bufSize & ^(1 << (i + bp.startByteCount)) == 0 {
			bp.pools[i].Put(buf[:0])
			return i
		}
	}
	return -1
}
