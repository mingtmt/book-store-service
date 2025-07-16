package application

import (
	"context"

	"github.com/mingtmt/book-store/internal/books/domain"
)

type BookRepository interface {
	CreateBook(ctx context.Context, book domain.Book) (string, error)
	GetBookByID(ctx context.Context, id string) (*domain.Book, error)
	GetAllBooks(ctx context.Context) ([]domain.Book, error)
	UpdateBook(ctx context.Context, book domain.Book) (*domain.Book, error)
	DeleteBookByID(ctx context.Context, id string) error
}

type BookService struct {
	repo BookRepository
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(ctx context.Context, title, author, price string) (string, error) {
	book := domain.Book{
		Title:  title,
		Author: author,
		Price:  price,
	}
	return s.repo.CreateBook(ctx, book)
}

func (s *BookService) GetBookByID(ctx context.Context, id string) (*domain.Book, error) {
	return s.repo.GetBookByID(ctx, id)
}

func (s *BookService) GetAllBooks(ctx context.Context) ([]domain.Book, error) {
	return s.repo.GetAllBooks(ctx)
}

func (s *BookService) UpdateBook(ctx context.Context, id, title, author, price string) (*domain.Book, error) {
	book := domain.Book{
		ID:     id,
		Title:  title,
		Author: author,
		Price:  price,
	}
	return s.repo.UpdateBook(ctx, book)
}

func (s *BookService) DeleteBookByID(ctx context.Context, id string) error {
	return s.repo.DeleteBookByID(ctx, id)
}
