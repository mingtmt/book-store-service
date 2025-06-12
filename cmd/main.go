package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mingtmt/book-store/internal/initialize"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	initialize.InitPostgres()
}
