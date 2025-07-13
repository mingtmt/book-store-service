// @title Book Store API
// @version 1.0
// @description This is a sample server for a bookstore.
// @BasePath /v1/api
package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mingtmt/book-store/docs"
	authController "github.com/mingtmt/book-store/internal/auth/controller"
	bookController "github.com/mingtmt/book-store/internal/books/controller"
	"github.com/mingtmt/book-store/internal/initialize"
	"github.com/mingtmt/book-store/internal/middleware"
	"github.com/mingtmt/book-store/pkg/logger"
	"github.com/mingtmt/book-store/pkg/token"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Set up Gin router
	r := gin.Default()

	// Register middleware
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(middleware.RequestID())
	r.Use(middleware.ErrorHandler())

	// Initialize DI container
	container := initialize.NewContainer(dbPool)

	// Public endpoints
	authGroup := r.Group("v1/api/auth")
	authController.RegisterUserRoutes(authGroup, container.AuthHandler)

	// Protected endpoints
	api := r.Group("v1/api")
	api.Use(middleware.JWTAuth())

	bookGroup := api.Group("/books")
	bookController.RegisterBookRoutes(bookGroup, container.BookHandler)

	// Serve Swagger UI at /swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("🚀 Server running", map[string]interface{}{"url": "http://localhost:" + port})
	if err := r.Run(":" + port); err != nil {
		logger.Error("Failed to run server", err, nil)
		os.Exit(1)
	}
}
