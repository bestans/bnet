package inet

import (
	"errors"
	"fmt"
)

//readIndex:标记目前可读数据起始索引
//writeIndex:标记目前已写数据索引
type SocketBuffer struct {
	data       []byte
	size       int
	writeIndex int
	readIndex  int
}

//触发update(readDecodeBuf剩余数据写入readBuf)的剩余可写字节数
//一个tcp包最大1448个字节
var triggerUpdateLeftWritableBytes = 2000

func NewSocketBuffer(size int) *SocketBuffer {
	if size <= 0 {
		panic(errors.New(fmt.Sprintf("Size must be positive, size=%v", size)))
	}

	b := &SocketBuffer{
		size: size,
		data: make([]byte, size),
	}
	return b
}

//写入buffer
func (b *SocketBuffer) WriteBytes(buf []byte) int {
	n := copy(b.data[b.writeIndex:], buf)
	b.writeIndex += n
	return n
}

// Size returns the size of the buffer
func (b *SocketBuffer) Size() int {
	return b.size
}

func (b *SocketBuffer) TriggerUpdate() bool {
	return b.WritableBytes() < triggerUpdateLeftWritableBytes
}

func (b *SocketBuffer) AddWriteIndex(writeSize int)  {
	b.writeIndex += writeSize
	if b.writeIndex > b.size {
		b.writeIndex = b.size
	}
}

func (b *SocketBuffer) AddReadIndex(readSize int)  {
	b.readIndex += readSize
	if b.readIndex > b.writeIndex {
		b.readIndex = b.writeIndex
	}
}

//可写字节数
func (b *SocketBuffer) WritableBytes() int {
	return b.size - b.writeIndex
}

//可读字节数
func (b *SocketBuffer) ReadableBytes() int {
	return b.writeIndex - b.readIndex
}

//可读字节内容
func (b *SocketBuffer) ReadByteData() []byte {
	return b.data[b.readIndex:b.writeIndex]
}

//可写字节内容
func (b *SocketBuffer) WriteByteData() []byte {
	return b.data[b.writeIndex:]
}

//buffer满了
func (b *SocketBuffer) Full() bool {
	return b.ReadableBytes() >= b.size
}

// Reset resets the buffer so it has no content.
func (b *SocketBuffer) Reset() {
	b.writeIndex = 0
	b.readIndex = 0
}
