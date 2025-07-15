package initialize

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	authController "github.com/mingtmt/book-store/internal/auth/controller"
	bookController "github.com/mingtmt/book-store/internal/books/controller"
	"github.com/mingtmt/book-store/internal/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(container *Container) *gin.Engine {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(middleware.RequestID())
	r.Use(middleware.ErrorHandler())

	apiPrefix := "/v1/api"

	authGroup := r.Group(apiPrefix + "/auth")
	authController.RegisterUserRoutes(authGroup, container.AuthHandler)
	bookGroup := r.Group(apiPrefix + "/books")
	bookController.RegisterBookRoutes(bookGroup, container.BookHandler)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
