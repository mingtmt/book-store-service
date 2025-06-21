package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mingtmt/book-store/internal/initialize"
	"github.com/mingtmt/book-store/internal/middleware"
	"github.com/mingtmt/book-store/pkg/logger"
)

func main() {
	logger.InitLogger()
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		logger.Error("Error loading .env file", err, nil)
		os.Exit(1)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database connection
	initialize.InitPostgres()

	// Set up Gin router
	r := gin.Default()

	// Register middleware
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.ErrorHandler())

	// Register routes
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	logger.Info("🚀 Server running", map[string]interface{}{"url": "http://localhost:" + port})
	if err := r.Run(":" + port); err != nil {
		logger.Error("Failed to run server", err, nil)
		os.Exit(1)
	}
}
