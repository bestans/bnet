package sharecoder

import (
	"errors"
	"unsafe"
)

var encodeVarintCompare = uint64(0x7F)
func EncodeVarint(b []byte, x uint64) ([]byte, error) {
	for {
		if x > encodeVarintCompare {
			b = append(b, byte(x&0x7f|0x80))
			x >>= 7
		} else {
			b = append(b, byte(x&0x7f))
			return b, nil
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

func DecodeVarint(buf []byte) (readLen int, x uint64, err error) {
	l := len(buf)
	for shift := uint(0); shift < 64; shift += 7 {
		if readLen >= l {
			err = errors.New("decodeVarint failed")
			return
		}
		b := buf[readLen]
		readLen++
		x |= (uint64(b) & 0x7F) << shift
		if b < 0x80 {
			return
		}
	}

	// The number is too large to represent in a 64-bit value.
	err = errors.New("decodeVarint failed")
	return
}

func EncodeString(buf []byte, str string) ([]byte, error) {
	strLen := len(str)
	bufLen := len(buf)
	if cap(buf) - bufLen < strLen {
		return buf, errors.New("buf not enough")
	}
	buf = buf[:bufLen+strLen]
	copy(buf[bufLen:], *((*[]byte)(unsafe.Pointer(&str))))
	return buf, nil
}

func DecodeString(buf []byte, strLen int) (string, []byte, error) {
	if len(buf) < strLen {
		return "", buf, errors.New("buf not enough")
	}
	return string(buf[:strLen]), buf[strLen:], nil
}
