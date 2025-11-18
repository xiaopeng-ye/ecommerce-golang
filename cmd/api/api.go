package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/xiaopeng-ye/ecommerce-golang/service/cart"
	"github.com/xiaopeng-ye/ecommerce-golang/service/order"
	"github.com/xiaopeng-ye/ecommerce-golang/service/product"
	"github.com/xiaopeng-ye/ecommerce-golang/service/user"
	"github.com/xiaopeng-ye/ecommerce-golang/utils"
)

type APIServer struct {
	addr   string
	db     *sql.DB
	server *http.Server
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	// Health check endpoints (not behind /api/v1 prefix)
	router.HandleFunc("/health", s.healthCheckHandler).Methods("GET")
	router.HandleFunc("/ready", s.readinessCheckHandler).Methods("GET")

	// API v1 routes
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)

	cartHandler := cart.NewHandler(productStore, orderStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	// Create HTTP server with timeouts
	s.server = &http.Server{
		Addr:         s.addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *APIServer) Shutdown(ctx context.Context) error {
	utils.Info("Shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}

// Close immediately closes the server
func (s *APIServer) Close() error {
	utils.Info("Closing HTTP server...")
	return s.server.Close()
}

// healthCheckHandler returns 200 if the service is alive
func (s *APIServer) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

// readinessCheckHandler returns 200 if the service is ready to accept traffic
func (s *APIServer) readinessCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check if database connection exists
	if s.db == nil {
		utils.Error("Database connection is nil during readiness check")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "not ready",
			"reason": "database connection not initialized",
		})
		return
	}

	// Check database connection
	if err := s.db.Ping(); err != nil {
		utils.WithField("error", err.Error()).Error("Database ping failed during readiness check")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "not ready",
			"reason": "database connection failed",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ready",
	})
}
