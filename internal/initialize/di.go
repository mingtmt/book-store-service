package initialize

import (
	"github.com/jackc/pgx/v5/pgxpool"
	bookService "github.com/mingtmt/book-store/internal/books/application"
	bookController "github.com/mingtmt/book-store/internal/books/controller"
	bookRepo "github.com/mingtmt/book-store/internal/books/infrastructure/persistence"
	userService "github.com/mingtmt/book-store/internal/users/application"
	userController "github.com/mingtmt/book-store/internal/users/controller"
	userRepo "github.com/mingtmt/book-store/internal/users/infrastructure/persistence"
)

type Container struct {
	BookHandler *bookController.BookHandler
	UserHandler *userController.UserHandler
}

func NewContainer(dbPool *pgxpool.Pool) *Container {
	// Users
	userRepo := userRepo.NewUserRepository(dbPool)
	authService := userService.NewAuthService(userRepo)
	userHandler := userController.NewUserHandler(authService)

	// Books
	bookRepo := bookRepo.NewBookRepository(dbPool)
	bookService := bookService.NewBookService(bookRepo)
	bookHandler := bookController.NewBookHandler(bookService)

	return &Container{
		BookHandler: bookHandler,
		UserHandler: userHandler,
	}
}
