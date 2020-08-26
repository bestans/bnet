package main

import (
	"bnet"
	"bnet/com"
	"bnet/com/session/protocoder"
	"bnet/inet"
	"bnet/net/tcp"
	_ "bnet/protoc"
	"fmt"
)

type ClientHandler struct {
	com.CoreSessionHandler
}

func (self *ClientHandler) ReadMessage(ses inet.ISession, decodeMsg interface{}) (err error) {
	//fmt.Println("server:", decodeMsg)
	ses.Send(decodeMsg)
	return nil
}
func TestPBServer() {
	server := bnet.NewTcpServer()
	//server.SetSessionOption(tcp.SetMessageHandler(&ClientHandler{}))
	server.SetSessionOption(tcp.SetProtocolCoder(protocoder.DefaultPBCoder()))
	server.Start()
	for {
		input := ""
		fmt.Scanf("%s\n", &input)
		fmt.Println(input)
	}
}
var totalCount int
func recvStrMessage(ses inet.ISession, msg interface{}) {
	ses.SendAndFlush(msg)
}
func TestStrServer() {

	server := bnet.NewTcpServer()
	//server.SetSessionOption(tcp.SetMessageHandler(&ClientHandler{}))
	server.SetSessionOption(tcp.SetProtocolCoder(protocoder.DefaultStringCoder()),
		tcp.SetMessageHandlerRecvFunc(recvStrMessage))
	server.Start()
	input := ""
	fmt.Scanf("%s\n", &input)
	fmt.Println(input)
}
func main()  {
	TestStrServer()
}
