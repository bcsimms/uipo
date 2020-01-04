package util

import (
	"log"
)

// The ahtentication type, determined by the arugments provided
type logLevelType int

var LogLevel logLevelType

const (
	LegLevelNone logLevelType = iota
	LogLevelInfo
	LogLevelDebug
	LogLevelTrace
)

func LogInfo(logMsg string) {

	if LogLevel >= LogLevelInfo {
		log.SetPrefix("Info: ")
		log.Println(logMsg)
	}

}

func LogDebug(logMsg string) {
	if LogLevel >= LogLevelDebug {
		log.SetPrefix("Debug: ")
		log.Println(logMsg)
	}
}

func LogTrace(logMsg string) {
	if LogLevel >= LogLevelTrace {
		log.SetPrefix("Trace: ")
		log.Println(logMsg)
	}
}
