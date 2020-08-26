package inet

import (
	"errors"
)

type Buffer struct {
	buf   []byte // encode/decode byte stream
	index int    // read point
	size int
}

// NewBuffer allocates a new Buffer and initializes its internal data to
// the contents of the argument slice.
func NewBuffer(e []byte) *Buffer {
	return &Buffer{buf: e, size: len(e)}
}

// Reset resets the Buffer, ready for marshaling a new protocol buffer.
func (p *Buffer) Reset() {
	p.index = 0        // for reading
}

// SetBuf replaces the internal buffer with the slice,
// ready for unmarshaling the contents of the slice.
func (p *Buffer) SetBuf(s []byte) {
	p.buf = s
	p.index = 0
	p.size = len(s)
}

// Bytes returns the contents of the Buffer.
func (p *Buffer) Bytes() []byte { return p.buf[0:p.index] }
func (p *Buffer) DataLen() int {
	return p.index
}

func (p *Buffer) AddIndex(add int) {
	p.index += add
	if p.index > p.size {
		p.index = p.size
	}
}

func (p *Buffer) LeftBytes() []byte {
	return p.buf[p.index:]
}


// EncodeVarint writes a varint-encoded integer to the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
var encodeVarintCompare = uint64(0x7F)
func (p *Buffer) EncodeVarint(x uint64) error {
	for {
		if p.index >= p.size {
			return errors.New("EncodeVarint failed: buffer not enough")
		}
		if x > encodeVarintCompare {
			p.buf[p.index] = byte(x&0x7f|0x80)
			p.index++
			x >>= 7
		} else {
			p.buf[p.index] = byte(x&0x7f)
			p.index++
			return nil
		}
	}
}

func EncodeVarintSize(x uint64) (size int) {
	for {
		if x > encodeVarintCompare {
			size++
			x >>= 7
		} else {
			size++
			return
		}
	}
}

func (p *Buffer) DecodeVarint() (x uint64, err error) {
	i := p.index
	l := len(p.buf)

	for shift := uint(0); shift < 64; shift += 7 {
		if i >= l {
			err = errors.New("decodeVarint failed")
			return
		}
		b := p.buf[i]
		i++
		x |= (uint64(b) & 0x7F) << shift
		if b < 0x80 {
			p.index = i
			return
		}
	}

	// The number is too large to represent in a 64-bit value.
	err = errors.New("decodeVarint failed")
	return
}

func (p *Buffer) WriteByes(data []byte) {
	p.index += copy(p.buf[p.index:], data)
}
