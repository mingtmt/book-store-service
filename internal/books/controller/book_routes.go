package controller

import "github.com/gin-gonic/gin"

func RegisterBookRoutes(r *gin.RouterGroup, handler *BookHandler) {
	r.POST("/", handler.CreateBook)
}
