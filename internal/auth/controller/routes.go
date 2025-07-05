package controller

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r *gin.RouterGroup, handler *AuthHandler) {
	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)
}
