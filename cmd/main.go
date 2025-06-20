package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mingtmt/book-store/internal/initialize"
	"github.com/mingtmt/book-store/internal/middleware"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
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
	r.Use(middleware.ErrorHandler())

	// Register routes
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})

	log.Printf("🚀 Server running at http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
