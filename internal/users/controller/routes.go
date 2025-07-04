package controller

import "github.com/gin-gonic/gin"

func RegisterUserRoutes(r *gin.RouterGroup, handler *UserHandler) {
	r.POST("/register", handler.RegisterUser)
}
