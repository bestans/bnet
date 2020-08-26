package bnet

import (
	"bnet/inet"
	netimpl "bnet/net"
	"bnet/net/tcp"
)

func NewTcpServer() *tcp.TcpServer {
	return netimpl.NewNet(inet.TcpServer, nil).(*tcp.TcpServer)
}

func NewTcpClient() *tcp.TcpClient {
	return netimpl.NewNet(inet.TcpClient, nil).(*tcp.TcpClient)
}
