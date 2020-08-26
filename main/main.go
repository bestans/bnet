package main

import (
	"bnet"
	"fmt"
)

func main()  {
	server := bnet.NewTcpServer()
	client := bnet.NewTcpClient()
	server.Start()
	client.Start()
	for {
		input := ""
		fmt.Scanf("%s\n", &input)
		client.SendAndFlush(input)
	}
}
