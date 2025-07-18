package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/internal/middleware"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	r.Use(middleware.AuthMiddleware())
	api := r.Group("/api/v1")
	for _, route := range routes {
		route.Register(api)
	}
}
