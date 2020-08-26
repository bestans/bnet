package main

import (
	"bnet"
	"fmt"
)

func main()  {
	client := bnet.NewTcpClient()
	client.Start()
	for {
		input := ""
		fmt.Scanf("%s\n", &input)
		client.SendAndFlush(input)
	}
}
