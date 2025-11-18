package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xiaopeng-ye/ecommerce-golang/utils"
)

// Mock database connection for testing
func setupMockDB() *sql.DB {
	// For health check tests, we'll use nil and handle it appropriately
	return nil
}

func TestHealthCheckHandler(t *testing.T) {
	// Initialize logger to avoid nil pointer
	utils.InitLogger("info", "json", "test")

	server := &APIServer{
		addr: ":8080",
		db:   nil,
	}

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	server.healthCheckHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %s", response["status"])
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got %s", contentType)
	}
}

func TestReadinessCheckHandlerWithoutDB(t *testing.T) {
	// Initialize logger to avoid nil pointer
	utils.InitLogger("info", "json", "test")

	server := &APIServer{
		addr: ":8080",
		db:   nil, // No database connection
	}

	req := httptest.NewRequest("GET", "/ready", nil)
	w := httptest.NewRecorder()

	server.readinessCheckHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status 503, got %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "not ready" {
		t.Errorf("Expected status 'not ready', got %s", response["status"])
	}

	if response["reason"] != "database connection not initialized" {
		t.Errorf("Expected reason 'database connection not initialized', got %s", response["reason"])
	}
}

func TestNewAPIServer(t *testing.T) {
	db := setupMockDB()
	server := NewAPIServer(":8080", db)

	if server == nil {
		t.Fatal("Expected server to be created")
	}

	if server.addr != ":8080" {
		t.Errorf("Expected addr to be ':8080', got %s", server.addr)
	}

	if server.db != db {
		t.Error("Expected db to be set correctly")
	}
}
