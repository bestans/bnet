package inet

type ErrorCode int16

const (
	ErrSuccess ErrorCode	= iota
	ErrEncode
	ErrConnWrite
	ErrSendWhenSessionClose
	ErrConnRead
	ErrConnReadRemoteClose
	ErrDecode
)

func (ec ErrorCode) String() string {
	return [...]string{
		"ErrSuccess",
		"ErrEncode",
		"ErrConnWrite",
		"ErrSendWhenSessionClose",
		"ErrConnRead",
		"ErrConnReadRemoteClose",
		"ErrDecode",
	}[ec]
}