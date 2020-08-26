package sharecoder

import (
	"fmt"
	"sync"
	"testing"
)

//cost time 32s
func TestDecodeEncodeSize(t *testing.T)  {
	w := &sync.WaitGroup{}
	cutLen := uint64(10)
	w.Add(int(cutLen))
	for fi := uint64(0); fi < cutLen; fi++ {
		index := fi
		go func() {
			for i := index; i < 0x00000000FFFFFFFF; i += cutLen {
				data, _ := EncodeVarint(nil, i)
				if len(data) != EncodeVarintSize(i) {
					panic("1111")
				}
			}
			w.Done()
		}()
	}
	w.Wait()
}

//cost time 72s
func TestDecodeEncodeSize2(t *testing.T)  {
	w := &sync.WaitGroup{}
	cutLen := uint64(10)
	w.Add(int(cutLen))
	for fi := uint64(0); fi < cutLen; fi++ {
		index := fi
		go func() {
			for i := index; i < 0x00000000FFFFFFFF; i += cutLen {
				value := i << 32
				data, _ := EncodeVarint(nil, value)
				if len(data) != EncodeVarintSize(value) {
					panic("1111")
				}
			}
			w.Done()
		}()
	}
	w.Wait()
}

//cost time 42s
func TestDecodeEncode(t *testing.T)  {
	w := &sync.WaitGroup{}
	cutLen := uint64(10)
	w.Add(int(cutLen))
	for fi := uint64(0); fi < cutLen; fi++ {
		index := fi
		go func() {
			for i := index; i < 0x00000000FFFFFFFF; i += cutLen {
				data, _ := EncodeVarint(nil, i)
				readLen, x, _ := DecodeVarint(data)
				if readLen != len(data) || x != i {
					panic("1111")
				}
			}
			w.Done()
		}()
	}
	w.Wait()
}
//cost time 90s
func TestDecodeEncode2(t *testing.T)  {
	w := &sync.WaitGroup{}
	cutLen := uint64(10)
	w.Add(int(cutLen))
	for fi := uint64(0); fi < cutLen; fi++ {
		index := fi
		go func() {
			for i := index; i < 0x00000000FFFFFFFF; i += cutLen {
				value := i << 32
				data, _ := EncodeVarint(nil, value)
				readLen, x, _ := DecodeVarint(data)
				if readLen != len(data) || x != value {
					panic("1111")
				}
			}
			w.Done()
		}()
	}
	w.Wait()
}

func TestEncodeString(t *testing.T) {
	str := "111"
	str2 := "222"
	buf := make([]byte, 0, len(str) +len(str2))
	buf, _ = EncodeString(buf, str)
	buf, _ = EncodeString(buf, "")
	buf, _ = EncodeString(buf, str2)
	buf, _ = EncodeString(buf, "afda")
	fmt.Println("buf", string(buf))
	if string(buf) != str + str2 {
		t.Errorf("failed")
	}
	dStr, buf, _ := DecodeString(buf, len(str))
	if dStr != str {
		t.Errorf("failed")
	}
	fmt.Println(dStr)
	dStr, buf, _ = DecodeString(buf, len(str2))
	if dStr != str2 {
		t.Errorf("failed")
	}
	fmt.Println(dStr)
}
