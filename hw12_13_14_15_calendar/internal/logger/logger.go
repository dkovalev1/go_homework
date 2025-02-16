package logger

import (
	"fmt"
	"strings"
)

// Log levels.
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
)

// Map string to log level.
func getLogLevel(level string) int {
	switch strings.ToLower(level) {
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warn":
		return WARN
	case "error":
		return ERROR
	default:
		return INFO
	}
}

type Logger struct { // TODO
	logLevel int // Default level
}

func New(level string) *Logger {
	return &Logger{logLevel: getLogLevel(level)}
}

func (l Logger) Debug(msg string) {
	if l.logLevel <= DEBUG {
		fmt.Printf("[DEBUG] %s\n", msg)
	}
}

func (l Logger) Info(msg string) {
	if l.logLevel <= INFO {
		fmt.Printf("[INFO] %s\n", msg)
	}
}

func (l Logger) Warn(msg string) {
	if l.logLevel <= WARN {
		fmt.Printf("[WARN] %s\n", msg)
	}
}

func (l Logger) Error(msg string) {
	if l.logLevel <= ERROR {
		fmt.Printf("[ERROR] %s\n", msg)
	}
}
