package utils

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewLogger(t *testing.T) {
	logger := NewLogger("info", "json", "test-service")

	if logger == nil {
		t.Fatal("Expected logger to be created")
	}

	if logger.level != InfoLevel {
		t.Errorf("Expected level to be info, got %s", logger.level)
	}

	if logger.format != "json" {
		t.Errorf("Expected format to be json, got %s", logger.format)
	}

	if logger.serviceName != "test-service" {
		t.Errorf("Expected service name to be test-service, got %s", logger.serviceName)
	}
}

func TestLoggerWithField(t *testing.T) {
	logger := NewLogger("info", "json", "test-service")
	loggerWithField := logger.WithField("key", "value")

	if loggerWithField.fields["key"] != "value" {
		t.Errorf("Expected field 'key' to be 'value', got %v", loggerWithField.fields["key"])
	}

	// Original logger should not be modified
	if _, exists := logger.fields["key"]; exists {
		t.Error("Original logger should not be modified")
	}
}

func TestLoggerWithFields(t *testing.T) {
	logger := NewLogger("info", "json", "test-service")
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}
	loggerWithFields := logger.WithFields(fields)

	if loggerWithFields.fields["key1"] != "value1" {
		t.Errorf("Expected field 'key1' to be 'value1', got %v", loggerWithFields.fields["key1"])
	}

	if loggerWithFields.fields["key2"] != 123 {
		t.Errorf("Expected field 'key2' to be 123, got %v", loggerWithFields.fields["key2"])
	}
}

func TestLoggerJSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("info", "json", "test-service")
	logger.output = &buf

	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("Expected log output to contain 'test message'")
	}

	// Parse JSON to verify it's valid
	var logEntry LogEntry
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Errorf("Expected valid JSON, got error: %v", err)
	}

	if logEntry.Message != "test message" {
		t.Errorf("Expected message to be 'test message', got %s", logEntry.Message)
	}

	if logEntry.Level != "info" {
		t.Errorf("Expected level to be 'info', got %s", logEntry.Level)
	}

	if logEntry.Service != "test-service" {
		t.Errorf("Expected service to be 'test-service', got %s", logEntry.Service)
	}
}

func TestLoggerTextFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("info", "text", "test-service")
	logger.output = &buf

	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Error("Expected log output to contain 'INFO'")
	}

	if !strings.Contains(output, "test message") {
		t.Error("Expected log output to contain 'test message'")
	}
}

func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name        string
		loggerLevel LogLevel
		logLevel    LogLevel
		shouldLog   bool
	}{
		{"Debug logger logs debug", DebugLevel, DebugLevel, true},
		{"Debug logger logs info", DebugLevel, InfoLevel, true},
		{"Info logger doesn't log debug", InfoLevel, DebugLevel, false},
		{"Info logger logs info", InfoLevel, InfoLevel, true},
		{"Info logger logs warn", InfoLevel, WarnLevel, true},
		{"Warn logger doesn't log info", WarnLevel, InfoLevel, false},
		{"Error logger doesn't log warn", ErrorLevel, WarnLevel, false},
		{"Error logger logs error", ErrorLevel, ErrorLevel, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewLogger(string(tt.loggerLevel), "json", "test")
			logger.output = &buf

			logger.log(tt.logLevel, "test message")

			output := buf.String()
			hasOutput := len(output) > 0

			if hasOutput != tt.shouldLog {
				t.Errorf("Expected shouldLog=%v, got output=%v", tt.shouldLog, hasOutput)
			}
		})
	}
}

func TestLoggerWithFieldsInOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := NewLogger("info", "json", "test-service")
	logger.output = &buf

	logger.WithFields(map[string]interface{}{
		"user_id": 123,
		"action":  "login",
	}).Info("User logged in")

	output := buf.String()

	var logEntry LogEntry
	if err := json.Unmarshal([]byte(output), &logEntry); err != nil {
		t.Fatalf("Expected valid JSON, got error: %v", err)
	}

	if logEntry.Fields["user_id"] != float64(123) {
		t.Errorf("Expected user_id to be 123, got %v", logEntry.Fields["user_id"])
	}

	if logEntry.Fields["action"] != "login" {
		t.Errorf("Expected action to be 'login', got %v", logEntry.Fields["action"])
	}
}

func TestInitLogger(t *testing.T) {
	InitLogger("debug", "json", "test-app")

	if defaultLogger == nil {
		t.Fatal("Expected default logger to be initialized")
	}

	if defaultLogger.level != DebugLevel {
		t.Errorf("Expected default logger level to be debug, got %s", defaultLogger.level)
	}

	// Reset default logger for other tests
	defaultLogger = NewLogger("info", "json", "app")
}
