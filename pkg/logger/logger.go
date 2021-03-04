package logger

import (
	"fmt"
	"log"
)

const (
	// Trace log level
	Trace = 1000 * iota
	// Debug log level
	Debug
	// Info log level
	Info
	// Warn log level
	Warn
	// Error log level
	Error
)

const (
	defaultLogLevel = Info
)

// Logger provides logging functionality with various log levels
type Logger struct {
	LogLevel int
}

// New creates a new logger with the specified log level
func New(level int) Logger {
	return Logger{LogLevel: level}
}

// ParseLogLevel determines the desired log level from an input string
func ParseLogLevel(level string) int {
	switch level {
	case "Trace":
		return Trace
	case "Debug":
		return Debug
	case "Info":
		return Info
	case "Warn":
		return Warn
	case "Error":
		return Error
	}

	return defaultLogLevel
}

// Trace logs when the log level satisfies the severity
func (logger *Logger) Trace(input string) {
	if logger.LogLevel <= Trace {
		fmt.Println(input)
	}
}

// Debug logs when the log level satisfies the severity
func (logger *Logger) Debug(input string) {
	if logger.LogLevel <= Debug {
		fmt.Println(input)
	}
}

// Info logs when the log level satisfies the severity
func (logger *Logger) Info(input string) {
	if logger.LogLevel <= Info {
		fmt.Println(input)
	}
}

// Warn logs when the log level satisfies the severity
func (logger *Logger) Warn(input string) {
	if logger.LogLevel <= Warn {
		fmt.Println(input)
	}
}

// Error logs when the log level satisfies the severity
func (logger *Logger) Error(input string) {
	if logger.LogLevel <= Error {
		fmt.Println(input)
	}
}

// Fatal proxies the log.Fatal to consolidate access to a single logger
func (logger *Logger) Fatal(input interface{}) {
	log.Fatal(input)
}
