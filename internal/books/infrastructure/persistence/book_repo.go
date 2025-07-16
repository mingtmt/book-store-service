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

func (r *BookRepository) CreateBook(ctx context.Context, book domain.Book) (string, error) {
	var price pgtype.Numeric
	if err := price.Scan(book.Price); err != nil {
		logger.Error("failed to convert price", err, map[string]interface{}{
			"book_id": book.ID,
			"price":   book.Price,
		})
		return "", err
	}

	created, err := r.db.CreateBook(ctx, booksdb.CreateBookParams{
		Title:  book.Title,
		Author: book.Author,
		Price:  price,
	})
	if err != nil {
		logger.Error("failed to create book in database", err, map[string]interface{}{
			"book_id": book.ID,
		})
		return "", err
	}

	logger.Info("book created successfully", map[string]interface{}{
		"book_id": created.ID.String(),
	})

	return created.ID.String(), nil
}

func (r *BookRepository) GetBookByID(ctx context.Context, id string) (*domain.Book, error) {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	bookID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}
	book, err := r.db.GetBook(ctx, bookID)
	if err != nil {
		return nil, err
	}

	return &domain.Book{
		ID:        book.ID.String(),
		Title:     book.Title,
		Author:    book.Author,
		Price:     book.Price.Int.String(),
		CreatedAt: book.CreatedAt.Time,
	}, nil
}

func (r *BookRepository) GetAllBooks(ctx context.Context) ([]domain.Book, error) {
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
	parsedUUID, err := uuid.Parse(book.ID)
	if err != nil {
		return nil, err
	}

	bookID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	var price pgtype.Numeric
	if err := price.Scan(book.Price); err != nil {
		logger.Error("failed to convert price", err, map[string]interface{}{
			"book_id": book.ID,
			"price":   book.Price,
		})
		return nil, err
	}

	updated, err := r.db.UpdateBook(ctx, booksdb.UpdateBookParams{
		ID:     bookID,
		Title:  book.Title,
		Author: book.Author,
		Price:  price,
	})
	if err != nil {
		logger.Error("failed to update book in database", err, map[string]interface{}{
			"book_id": book.ID,
		})
		return nil, err
	}

	logger.Info("book updated successfully", map[string]interface{}{
		"book_id": updated.ID.String(),
	})

	return &domain.Book{
		ID:        updated.ID.String(),
		Title:     updated.Title,
		Author:    updated.Author,
		Price:     updated.Price.Int.String(),
		CreatedAt: updated.CreatedAt.Time,
	}, nil
}

func (r *BookRepository) DeleteBookByID(ctx context.Context, id string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	bookID := pgtype.UUID{
		Bytes: parsedUUID,
		Valid: true,
	}

	err = r.db.DeleteBook(ctx, bookID)
	if err != nil {
		logger.Error("failed to delete book from database", err, map[string]interface{}{
			"book_id": id,
		})
		return err
	}

	logger.Info("book deleted successfully", map[string]interface{}{
		"book_id": id,
	})
	return nil
}
