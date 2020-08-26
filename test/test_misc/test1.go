package main

import (
	"bnet/protoc"
	"fmt"
	"github.com/golang/protobuf/proto"
	"reflect"
)

type TestStru struct {
	value int
}

func Change(v interface{}) *TestStru  {
	return v.(*TestStru)
}
func TestMiscTest1a()  {
	opt := &protoc.BagOpt{
		Ret:100,
	}
	newOpt := reflect.New(reflect.TypeOf(protoc.BagOpt{})).Interface().(proto.Message).(*protoc.BagOpt)
	newOpt.Ret = 101

	newOpt2 := reflect.New(reflect.TypeOf((*protoc.BagOpt)(nil)).Elem()).Elem().Interface().(protoc.BagOpt)
	newOpt2.Ret = 102
	fmt.Println(opt.Ret, newOpt.Ret, newOpt2.Ret)

	v := Change(newOpt)
	fmt.Println(v)
}

func TestMiscTest1b()  {
}

func main()  {
	TestMiscTest1a()
}
