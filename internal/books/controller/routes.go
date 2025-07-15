package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/internal/middleware"
)

func RegisterBookRoutes(r *gin.RouterGroup, handler *BookHandler) {
	// Public routes
	r.GET("/", handler.GetAllBooks)
	r.GET(":id", handler.GetBook)

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(middleware.JWTAuth())
	{
		authorized.POST("", handler.CreateBook)
		authorized.PUT(":id", handler.UpdateBook)
		authorized.DELETE(":id", handler.DeleteBook)
	}
}
