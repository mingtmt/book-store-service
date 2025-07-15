// @title Book Store API
// @version 1.0
// @description This is a sample server for a bookstore.
// @BasePath /v1/api
package main

import (
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mingtmt/book-store/docs"
	"github.com/mingtmt/book-store/internal/initialize"
	"github.com/mingtmt/book-store/pkg/logger"
	"github.com/mingtmt/book-store/pkg/token"
)

func main() {
	// Initialize logger
	logger.InitLogger()
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file", err, nil)
		os.Exit(1)
	}

	// Initialize JWT keys
	token.InitKeys()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database connection
	dbPool, err := initialize.InitPostgres()
	if err != nil {
		logger.Error("failed to connect to database:", err, nil)
	}
	defer dbPool.Close()

	// Initialize DI container
	container := initialize.NewContainer(dbPool)

	r := initialize.SetupRouter(container)

	logger.Info("🚀 Server running", map[string]interface{}{"url": "http://localhost:" + port})
	if err := r.Run(":" + port); err != nil {
		logger.Error("Failed to run server", err, nil)
		os.Exit(1)
	}
}
