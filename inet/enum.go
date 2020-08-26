package inet

type NetType int16
const (
	TcpServer NetType	= iota
	TcpClient
)

func (nt NetType) String() string {
	return [...]string{
		"TcpServer",
		"TcpClient",
	}[nt]
}
