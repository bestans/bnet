package log

type ILogger interface {
	Trace(format string, a ...interface{})
	Error(format string, a ...interface{})
}
