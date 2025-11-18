package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
)

// Logger is a structured logger
type Logger struct {
	level      LogLevel
	format     string // "json" or "text"
	output     io.Writer
	fields     map[string]interface{}
	serviceName string
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp   string                 `json:"timestamp"`
	Level       string                 `json:"level"`
	Message     string                 `json:"message"`
	Service     string                 `json:"service,omitempty"`
	File        string                 `json:"file,omitempty"`
	Line        int                    `json:"line,omitempty"`
	Fields      map[string]interface{} `json:"fields,omitempty"`
}

var defaultLogger *Logger

// InitLogger initializes the default logger
func InitLogger(level, format, serviceName string) {
	defaultLogger = NewLogger(level, format, serviceName)
}

// NewLogger creates a new logger instance
func NewLogger(level, format, serviceName string) *Logger {
	return &Logger{
		level:       LogLevel(strings.ToLower(level)),
		format:      strings.ToLower(format),
		output:      os.Stdout,
		fields:      make(map[string]interface{}),
		serviceName: serviceName,
	}
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := &Logger{
		level:       l.level,
		format:      l.format,
		output:      l.output,
		fields:      make(map[string]interface{}),
		serviceName: l.serviceName,
	}
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	newLogger.fields[key] = value
	return newLogger
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	newLogger := &Logger{
		level:       l.level,
		format:      l.format,
		output:      l.output,
		fields:      make(map[string]interface{}),
		serviceName: l.serviceName,
	}
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return newLogger
}

func (l *Logger) log(level LogLevel, message string) {
	if !l.shouldLog(level) {
		return
	}

	_, file, line, _ := runtime.Caller(2)
	// Extract just the filename
	parts := strings.Split(file, "/")
	if len(parts) > 0 {
		file = parts[len(parts)-1]
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     string(level),
		Message:   message,
		Service:   l.serviceName,
		File:      file,
		Line:      line,
		Fields:    l.fields,
	}

	var output string
	if l.format == "json" {
		jsonBytes, err := json.Marshal(entry)
		if err != nil {
			log.Printf("Error marshaling log entry: %v", err)
			return
		}
		output = string(jsonBytes)
	} else {
		// Text format
		output = fmt.Sprintf("[%s] %s: %s", entry.Timestamp, strings.ToUpper(entry.Level), entry.Message)
		if len(l.fields) > 0 {
			fieldsJSON, _ := json.Marshal(l.fields)
			output += fmt.Sprintf(" fields=%s", string(fieldsJSON))
		}
		output += fmt.Sprintf(" (%s:%d)", entry.File, entry.Line)
	}

	fmt.Fprintln(l.output, output)
}

func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		DebugLevel: 0,
		InfoLevel:  1,
		WarnLevel:  2,
		ErrorLevel: 3,
		FatalLevel: 4,
	}
	return levels[level] >= levels[l.level]
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	l.log(DebugLevel, message)
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.log(InfoLevel, message)
}

// Warn logs a warning message
func (l *Logger) Warn(message string) {
	l.log(WarnLevel, message)
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.log(ErrorLevel, message)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string) {
	l.log(FatalLevel, message)
	os.Exit(1)
}

// Package-level logging functions using the default logger
func Debug(message string) {
	if defaultLogger != nil {
		defaultLogger.Debug(message)
	}
}

func Info(message string) {
	if defaultLogger != nil {
		defaultLogger.Info(message)
	}
}

func Warn(message string) {
	if defaultLogger != nil {
		defaultLogger.Warn(message)
	}
}

func Error(message string) {
	if defaultLogger != nil {
		defaultLogger.Error(message)
	}
}

func Fatal(message string) {
	if defaultLogger != nil {
		defaultLogger.Fatal(message)
	}
}

func WithField(key string, value interface{}) *Logger {
	if defaultLogger != nil {
		return defaultLogger.WithField(key, value)
	}
	return NewLogger("info", "json", "app")
}

func WithFields(fields map[string]interface{}) *Logger {
	if defaultLogger != nil {
		return defaultLogger.WithFields(fields)
	}
	return NewLogger("info", "json", "app")
}
