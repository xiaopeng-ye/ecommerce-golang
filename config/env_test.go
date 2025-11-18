package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	value := getEnv("TEST_VAR", "default")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %s", value)
	}

	// Test with non-existing env var
	value = getEnv("NON_EXISTING_VAR", "default")
	if value != "default" {
		t.Errorf("Expected 'default', got %s", value)
	}
}

func TestGetEnvAsInt(t *testing.T) {
	// Test with valid integer
	os.Setenv("TEST_INT", "123")
	defer os.Unsetenv("TEST_INT")

	value := getEnvAsInt("TEST_INT", 456)
	if value != 123 {
		t.Errorf("Expected 123, got %d", value)
	}

	// Test with invalid integer
	os.Setenv("TEST_INVALID_INT", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_INT")

	value = getEnvAsInt("TEST_INVALID_INT", 456)
	if value != 456 {
		t.Errorf("Expected fallback 456, got %d", value)
	}

	// Test with non-existing env var
	value = getEnvAsInt("NON_EXISTING_INT", 789)
	if value != 789 {
		t.Errorf("Expected fallback 789, got %d", value)
	}
}

func TestGetEnvRequiredInTestEnvironment(t *testing.T) {
	// Set ENVIRONMENT to test
	os.Setenv("ENVIRONMENT", "test")
	defer os.Unsetenv("ENVIRONMENT")

	// Test DB_PASSWORD with test default
	os.Unsetenv("DB_PASSWORD")
	value := getEnvRequired("DB_PASSWORD")
	if value != "test-password" {
		t.Errorf("Expected 'test-password', got %s", value)
	}

	// Test JWT_SECRET with test default
	os.Unsetenv("JWT_SECRET")
	value = getEnvRequired("JWT_SECRET")
	if value != "test-secret-key-for-testing-only-do-not-use-in-production" {
		t.Errorf("Expected test JWT secret, got %s", value)
	}
}

func TestGetEnvRequiredWithValue(t *testing.T) {
	// Set environment to test
	os.Setenv("ENVIRONMENT", "test")
	defer os.Unsetenv("ENVIRONMENT")

	// Test with actual value set
	os.Setenv("DB_PASSWORD", "actual_password")
	defer os.Unsetenv("DB_PASSWORD")

	value := getEnvRequired("DB_PASSWORD")
	if value != "actual_password" {
		t.Errorf("Expected 'actual_password', got %s", value)
	}
}

func TestInitConfig(t *testing.T) {
	// Save current environment
	originalEnv := os.Getenv("ENVIRONMENT")
	defer func() {
		if originalEnv != "" {
			os.Setenv("ENVIRONMENT", originalEnv)
		} else {
			os.Unsetenv("ENVIRONMENT")
		}
	}()

	// Set test environment
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("PORT", "9000")
	os.Setenv("LOG_LEVEL", "debug")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("LOG_LEVEL")

	config := initConfig()

	if config.Environment != "test" {
		t.Errorf("Expected environment 'test', got %s", config.Environment)
	}

	if config.Port != "9000" {
		t.Errorf("Expected port '9000', got %s", config.Port)
	}

	if config.LogLevel != "debug" {
		t.Errorf("Expected log level 'debug', got %s", config.LogLevel)
	}

	// Verify defaults
	if config.DBUser != "root" {
		t.Errorf("Expected default DB user 'root', got %s", config.DBUser)
	}

	if config.LogFormat != "json" {
		t.Errorf("Expected default log format 'json', got %s", config.LogFormat)
	}
}

func TestConfigStructure(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "3307")
	defer os.Unsetenv("ENVIRONMENT")
	defer os.Unsetenv("DB_HOST")
	defer os.Unsetenv("DB_PORT")

	config := initConfig()

	// Test DBAddress is formatted correctly
	expectedAddress := "testhost:3307"
	if config.DBAddress != expectedAddress {
		t.Errorf("Expected DB address '%s', got %s", expectedAddress, config.DBAddress)
	}

	if config.DBHost != "testhost" {
		t.Errorf("Expected DB host 'testhost', got %s", config.DBHost)
	}

	if config.DBPort != "3307" {
		t.Errorf("Expected DB port '3307', got %s", config.DBPort)
	}
}
