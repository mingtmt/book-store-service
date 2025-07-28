package app

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mingtmt/book-store/configs"
	"github.com/mingtmt/book-store/internal/routes"
	"github.com/mingtmt/book-store/internal/utils/validation"
)

type Module interface {
	Routes() routes.Route
}

type Application struct {
	config  *configs.Config
	router  *gin.Engine
	modules []Module
}

func NewApplication(cfg *configs.Config) *Application {
	validation.InitValidator()
	loadEnv()

	r := gin.Default()
	modules := []Module{
		NewUserModule(),
	}
	routes.RegisterRoutes(r, getModuleRoutes(modules)...)

	return &Application{
		config:  cfg,
		router:  r,
		modules: modules,
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

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

}
