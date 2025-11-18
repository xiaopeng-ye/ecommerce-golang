package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/xiaopeng-ye/ecommerce-golang/cmd/api"
	"github.com/xiaopeng-ye/ecommerce-golang/config"
	"github.com/xiaopeng-ye/ecommerce-golang/db"
	"github.com/xiaopeng-ye/ecommerce-golang/utils"
)

var (
	// These variables can be set during build time using ldflags
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Initialize structured logger
	utils.InitLogger(config.Envs.LogLevel, config.Envs.LogFormat, "ecommerce-api")

	// Log application startup information
	utils.WithFields(map[string]interface{}{
		"version":     Version,
		"build_time":  BuildTime,
		"git_commit":  GitCommit,
		"environment": config.Envs.Environment,
	}).Info("Starting ecommerce API server")

	// Initialize database connection
	database, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})

	if err != nil {
		utils.WithField("error", err.Error()).Fatal("Failed to connect to database")
	}

	initStorage(database)

	// Create API server with port from config
	addr := fmt.Sprintf(":%s", config.Envs.Port)
	server := api.NewAPIServer(addr, database)

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		utils.WithField("address", addr).Info("Starting HTTP server")
		serverErrors <- server.Run()
	}()

	// Channel to listen for interrupt or terminate signal from the OS
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown
	select {
	case err := <-serverErrors:
		utils.WithField("error", err.Error()).Fatal("Server error")

	case sig := <-shutdown:
		utils.WithField("signal", sig.String()).Info("Starting shutdown")

		// Give outstanding requests a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Asking listener to shutdown and shed load
		if err := server.Shutdown(ctx); err != nil {
			utils.WithField("error", err.Error()).Error("Graceful shutdown did not complete")
			if err := server.Close(); err != nil {
				utils.WithField("error", err.Error()).Fatal("Could not stop server gracefully")
			}
		}

		// Close database connection
		if err := database.Close(); err != nil {
			utils.WithField("error", err.Error()).Error("Failed to close database connection")
		}

		utils.Info("Server stopped gracefully")
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		utils.WithField("error", err.Error()).Fatal("Failed to ping database")
	}

	utils.Info("Database connection established successfully")
}
