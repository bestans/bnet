package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

func MgErr(msgs ...interface{}) error {
	return errors.New(StringMerge(msgs...))
}

func StringMerge(strs ...interface{}) string {
	var buffer bytes.Buffer
	for _, str := range strs {
		buffer.WriteString(fmt.Sprintf("%v", str))
	}
	return buffer.String()
}

// 判断网络错误
func IsEOFOrNetReadError(err error) bool {
	if err == io.EOF {
		return true
	}
	ne, ok := err.(*net.OpError)
	return ok && ne.Op == "read"
}

// 将host和端口合并为(host:port)格式的地址
func JoinAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

var globalUniqueIdSeq int64
func GenerateUniqueId() int64 {
	return atomic.AddInt64(&globalUniqueIdSeq, 1)
}
