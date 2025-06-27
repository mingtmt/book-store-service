package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mingtmt/book-store/internal/books/domain"
	"github.com/mingtmt/book-store/internal/books/infrastructure/persistence/booksdb"
	"github.com/mingtmt/book-store/pkg/logger"
)

type BookRepository struct {
	db *booksdb.Queries
}

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{
		db: booksdb.New(pool),
	}
}

func (r *BookRepository) Create(ctx context.Context, book domain.Book) (*domain.Book, error) {
	var price pgtype.Numeric
	if err := price.Scan(book.Price); err != nil {
		logger.Error("failed to convert price", err, map[string]interface{}{
			"book_id": book.ID,
			"price":   book.Price,
		})
		return nil, err
	}

	id := pgtype.UUID{
		Bytes: uuid.MustParse(book.ID),
		Valid: true,
	}

	created, err := r.db.CreateBook(ctx, booksdb.CreateBookParams{
		ID:     id,
		Title:  book.Title,
		Author: book.Author,
		Price:  price,
	})
	if err != nil {
		logger.Error("failed to create book in database", err, map[string]interface{}{
			"book_id": book.ID,
			"title":   book.Title,
			"author":  book.Author,
			"price":   book.Price,
		})
		return nil, err
	}

	logger.Info("book created successfully", map[string]interface{}{
		"book_id": created.ID.String(),
		"title":   created.Title,
		"author":  created.Author,
		"price":   created.Price,
	})

	return &domain.Book{
		ID:        created.ID.String(),
		Title:     created.Title,
		Author:    created.Author,
		Price:     book.Price,
		CreatedAt: created.CreatedAt.Time,
	}, nil
}
