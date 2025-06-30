package controller

import "github.com/gin-gonic/gin"

func RegisterBookRoutes(r *gin.RouterGroup, handler *BookHandler) {
	r.POST("/", handler.CreateBook)
	r.GET("/:id", handler.GetBook)
	r.GET("/", handler.GetAllBooks)
	r.PUT("/:id", handler.UpdateBook)
}
