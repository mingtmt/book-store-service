package main

import (
	"github.com/mingtmt/book-store/configs"
	"github.com/mingtmt/book-store/internal/app"
)

func main() {
	// Initialize the configuration
	cfg := configs.NewConfig()

	// Initialize the application
	application := app.NewApplication(cfg)

	if err := application.Run(); err != nil {
		panic(err)
	}
}
