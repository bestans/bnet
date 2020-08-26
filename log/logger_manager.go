package log

import (
	"fmt"
	"time"
)

var extendLogger ILogger

func SetLogger(logger ILogger)  {
	extendLogger = logger
}

func Trace(format string, a ...interface{}) {
	if len(a) > 0 {
		fmt.Printf("%v:TRACE:bNet:%v\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, a...))
	} else {
		fmt.Printf("%v:TRACE:bNet:%v\n", time.Now().Format("2006-01-02 15:04:05"), format)
	}

	if extendLogger != nil {
		extendLogger.Trace(format, a...)
	}
}

func Error(format string, a ...interface{}) {
	if len(a) > 0 {
		fmt.Printf("%v:ERROR:bNet:%v\n", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, a))
	} else {
		fmt.Printf("%v:ERROR:bNet:%v\n", time.Now().Format("2006-01-02 15:04:05"), format)
	}
	if extendLogger != nil {
		extendLogger.Error(format, a)
	}
}
