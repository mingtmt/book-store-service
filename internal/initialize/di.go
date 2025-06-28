package initialize

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mingtmt/book-store/internal/books/application"
	"github.com/mingtmt/book-store/internal/books/controller"
	"github.com/mingtmt/book-store/internal/books/infrastructure/persistence"
)

type Container struct {
	BookHandler *controller.BookHandler
}

func NewContainer(dbPool *pgxpool.Pool) *Container {
	// Books
	bookRepo := persistence.NewBookRepository(dbPool)
	bookService := application.NewBookService(bookRepo)
	bookHandler := controller.NewBookHandler(bookService)

	return &Container{
		BookHandler: bookHandler,
	}
}
