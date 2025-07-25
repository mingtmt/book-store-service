package initialize

import (
	"github.com/jackc/pgx/v5/pgxpool"
	authService "github.com/mingtmt/book-store/internal/auth/application"
	authController "github.com/mingtmt/book-store/internal/auth/controller"
	authPersitence "github.com/mingtmt/book-store/internal/auth/infrastructure/persistence"
	bookService "github.com/mingtmt/book-store/internal/books/application"
	bookController "github.com/mingtmt/book-store/internal/books/controller"
	bookRepo "github.com/mingtmt/book-store/internal/books/infrastructure/persistence"
)

type Container struct {
	BookHandler *bookController.BookHandler
	AuthHandler *authController.AuthHandler
}

func NewContainer(dbPool *pgxpool.Pool) *Container {
	// Authentication
	authRepo := authPersitence.NewAuthRepository(dbPool)
	rfTokenRepo := authPersitence.NewRefreshTokenRepository(dbPool)
	authService := authService.NewAuthService(authRepo, rfTokenRepo)
	authHandler := authController.NewAuthHandler(authService)

	// Books
	bookRepo := bookRepo.NewBookRepository(dbPool)
	bookService := bookService.NewBookService(bookRepo)
	bookHandler := bookController.NewBookHandler(bookService)

	return &Container{
		BookHandler: bookHandler,
		AuthHandler: authHandler,
	}
}
