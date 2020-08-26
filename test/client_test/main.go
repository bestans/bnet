package main

import (
	"bnet"
	"bnet/com"
	"bnet/com/session/protocoder"
	"bnet/inet"
	"bnet/net/tcp"
	"bnet/protoc"
	"fmt"
	"time"
)

type ClientHandler struct {
	com.CoreSessionHandler
}

func (self *ClientHandler) ReadMessage(ses inet.ISession, decodeMsg interface{}) (err error) {
	//fmt.Println("client:", decodeMsg)
	totalReply++
	if totalReply >= 100000 {
		exit <- true
	}
	ses.Send("1")
	return nil
}

var totalReply int
var exit = make(chan bool, 1)
func TestPBClient() {

	client := bnet.NewTcpClient()
	//client.SetSessionOption(tcp.SetMessageHandler(&ClientHandler{}))
	client.SetSessionOption(tcp.SetProtocolCoder(protocoder.DefaultPBCoder()))
	client.Start()

	input := ""
	fmt.Scanf("%s\n", &input)
	msg := &protoc.BagOpt{
		Ret:111,
		ItemId:1000,
	}
	client.Session().Send(msg)
	client.Session().Send(msg)
	client.Session().Send(msg)
	msg.ItemId = 1022
	client.Session().SendAndFlush(msg)
	//client.Session().Flush()

	start := time.Now()
	<- exit
	fmt.Println("cost_time:", time.Since(start))

	fmt.Scanf("%s\n", &input)
}
var totalLen = 0
var msgCount = 0
func recvStrMessage(ses inet.ISession, msg interface{}) {
	msgCount++
	totalLen += len(msg.(string))
	if msgCount < 10 {
		ses.SendAndFlush(msg)
	} else {
		exit <- true
	}
}
type XXStr struct {

}

func (xs *XXStr) String() string {
	return "XXStr"
}
func TestStringClient() {
	client := bnet.NewTcpClient()
	//client.SetSessionOption(tcp.SetMessageHandler(&ClientHandler{}))
	client.SetSessionOption(
		tcp.SetProtocolCoder(protocoder.DefaultStringCoder()),
		tcp.SetMessageHandlerRecvFunc(recvStrMessage))
	client.Start()
	start := time.Now()
	b := make([]byte, 65532)
	client.Session().SendAndFlush(string(b))
	<- exit
	client.CloseAllSession("1111")
	fmt.Printf("times:%v,costTime=%v\n", totalLen, time.Since(start))
	time.Sleep(time.Second)
}

func recvStrMessageTPS(ses inet.ISession, msg interface{}, maxCount int, count *int, exitCh chan bool) {
	*count++
	if *count >= maxCount {
		exitCh <- true
	}
	ses.SendAndFlush(msg)
}
func TestStringClientTPS() {
	client := bnet.NewTcpClient()
	//client.SetSessionOption(tcp.SetMessageHandler(&ClientHandler{}))
	var count int
	exitCh := make(chan bool)
	client.SetSessionOption(
		tcp.SetProtocolCoder(protocoder.DefaultStringCoder()),
		tcp.SetMessageHandlerRecvFunc(func(session inet.ISession, i interface{}) {
			recvStrMessageTPS(session, i, 100000, &count, exitCh)
		}))
	client.Start()
	start := time.Now()
	client.Session().SendAndFlush("1")
	<- exitCh
	client.CloseAllSession("1111")
	fmt.Printf("times:%v,costTime=%v\n", count, time.Since(start))
	time.Sleep(time.Second)
}
func main()  {
	TestStringClientTPS()
}
