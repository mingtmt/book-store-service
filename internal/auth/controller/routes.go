package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/internal/middleware"
)

func RegisterUserRoutes(r *gin.RouterGroup, handler *AuthHandler) {
	// Public routes
	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)
	r.POST("/refresh", handler.RefreshToken)

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(middleware.JWTAuth())
	{
		authorized.POST("logout", handler.LogoutUser)
	}
}
