package domain

import (
	"context"
	"time"
)

type Book struct {
	ID        string
	Title     string
	Author    string
	Price     string
	CreatedAt time.Time
}

type BookRepository interface {
	CreateBook(ctx context.Context, book Book) (string, error)
	GetBookByID(ctx context.Context, id string) (*Book, error)
	GetAllBooks(ctx context.Context) ([]Book, error)
	UpdateBook(ctx context.Context, book Book) (*Book, error)
	DeleteBookByID(ctx context.Context, id string) error
}
