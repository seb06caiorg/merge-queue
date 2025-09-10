package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents different log levels.
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// Logger provides structured logging functionality.
type Logger struct {
	level  LogLevel
	logger *log.Logger
}

// NewLogger creates a new Logger instance.
func NewLogger(level LogLevel) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0), // We'll format ourselves.
	}
}

// NewDefaultLogger creates a logger with info level.
func NewDefaultLogger() *Logger {
	return NewLogger(InfoLevel)
}

// Debug logs a debug message.
func (l *Logger) Debug(message string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.log("DEBUG", message, args...)
	}
}

// Info logs an info message.
func (l *Logger) Info(message string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.log("INFO", message, args...)
	}
}

// Warn logs a warning message.
func (l *Logger) Warn(message string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.log("WARN", message, args...)
	}
}

// Error logs an error message.
func (l *Logger) Error(message string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.log("ERROR", message, args...)
	}
}

// log formats and logs a message.
func (l *Logger) log(level, message string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formattedMessage := fmt.Sprintf(message, args...)
	logLine := fmt.Sprintf("[%s] %s: %s", timestamp, level, formattedMessage)
	l.logger.Println(logLine)
}

// SetLevel sets the minimum log level.
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

// GetLevel returns the current log level.
func (l *Logger) GetLevel() LogLevel {
	return l.level
}

// LogLevelFromString converts a string to LogLevel.
func LogLevelFromString(level string) LogLevel {
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}
