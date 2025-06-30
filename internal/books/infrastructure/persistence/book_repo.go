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
		})
		return nil, err
	}

	logger.Info("book created successfully", map[string]interface{}{
		"book_id": created.ID.String(),
	})

	return &domain.Book{
		ID:        created.ID.String(),
		Title:     created.Title,
		Author:    created.Author,
		Price:     book.Price,
		CreatedAt: created.CreatedAt.Time,
	}, nil
}

func (r *BookRepository) GetByID(ctx context.Context, id string) (*domain.Book, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	bookID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}
	b, err := r.db.GetBook(ctx, bookID)
	if err != nil {
		return nil, err
	}
	priceStr := ""
	if b.Price.Valid {
		priceStr = b.Price.Int.String()
	}
	return &domain.Book{
		ID:        b.ID.String(),
		Title:     b.Title,
		Author:    b.Author,
		Price:     priceStr,
		CreatedAt: b.CreatedAt.Time,
	}, nil
}

func (r *BookRepository) GetAll(ctx context.Context) ([]domain.Book, error) {
	rows, err := r.db.ListBooks(ctx)
	if err != nil {
		return nil, err
	}
	var books []domain.Book
	for _, row := range rows {
		priceStr := ""
		if row.Price.Valid {
			priceStr = row.Price.Int.String()
		}
		books = append(books, domain.Book{
			ID:        row.ID.String(),
			Title:     row.Title,
			Author:    row.Author,
			Price:     priceStr,
			CreatedAt: row.CreatedAt.Time,
		})
	}
	return books, nil
}

func (r *BookRepository) UpdateBook(ctx context.Context, book domain.Book) (*domain.Book, error) {
	// This method is not implemented in the original code snippet.
	// If you need to implement it, you can follow a similar pattern as Create and GetByID.
	return nil, nil
}
