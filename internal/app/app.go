package app

import (
	"github.com/gin-gonic/gin"
	"github.com/mingtmt/book-store/configs"
	"github.com/mingtmt/book-store/internal/routes"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config *configs.Config
	router *gin.Engine
}

func NewApplication(cfg *configs.Config) *Application {
	r := gin.Default()
	modules := []Module{
		NewUserModule(),
	}
	routes.RegisterRoutes(r, getModuleRoutes(modules)...)

	return &Application{
		config: cfg,
		router: r,
	}
}

func (a *Application) Run() error {
	return a.router.Run(a.config.ServerAddress)
}

func getModuleRoutes(modules []Module) []routes.Route {
	routeList := make([]routes.Route, len(modules))
	for i, module := range modules {
		routeList[i] = module.Routes()
	}

	return routeList
}
