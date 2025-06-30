package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/mingtmt/book-store/internal/books/domain"
)

type BookRepository interface {
	Create(ctx context.Context, book domain.Book) (*domain.Book, error)
	GetByID(ctx context.Context, id string) (*domain.Book, error)
	GetAll(ctx context.Context) ([]domain.Book, error)
}

type BookService struct {
	repo BookRepository
}

func NewBookService(repo BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(ctx context.Context, title, author, price string) (*domain.Book, error) {
	book := domain.Book{
		ID:     uuid.New().String(),
		Title:  title,
		Author: author,
		Price:  price,
	}
	return s.repo.Create(ctx, book)
}

func (s *BookService) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BookService) GetAll(ctx context.Context) ([]domain.Book, error) {
	return s.repo.GetAll(ctx)
}
