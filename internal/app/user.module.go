package app

import (
	"github.com/mingtmt/book-store/internal/handler"
	"github.com/mingtmt/book-store/internal/repository"
	"github.com/mingtmt/book-store/internal/routes"
	"github.com/mingtmt/book-store/internal/service"
)

type UserModule struct {
	routes routes.Route
}

func NewUserModule() *UserModule {
	userRepo := repository.NewInMemUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	userRoutes := routes.NewUserRoutes(userHandler)

	return &UserModule{
		routes: userRoutes,
	}
}

func (um *UserModule) Routes() routes.Route {
	return um.routes
}
