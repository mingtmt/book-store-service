package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mingtmt/book-store/internal/users/domain"
	"github.com/mingtmt/book-store/internal/users/infrastructure/persistence/usersdb"
	"github.com/mingtmt/book-store/pkg/logger"
)

type UserRepository struct {
	db *usersdb.Queries
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: usersdb.New(pool),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	id := pgtype.UUID{
		Bytes: uuid.MustParse(user.ID),
		Valid: true,
	}

	created, err := r.db.CreateUser(context.Background(), usersdb.CreateUserParams{
		ID:       id,
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		logger.Error("failed to create user in database", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, err
	}

	logger.Info("user created successfully", map[string]interface{}{
		"user_id": created.ID.String(),
	})
	return &domain.User{
		ID:       created.ID.String(),
		Username: created.Username,
		Password: created.Password,
	}, nil
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := r.db.FindByUsername(ctx, username)
	if err != nil {
		logger.Error("failed to find user by username", err, map[string]interface{}{
			"username": username,
		})
		return nil, err
	}

	if !user.ID.Valid {
		logger.Warn("user not found", map[string]interface{}{
			"username": username,
		})
		return nil, nil
	}

	return &domain.User{
		ID:       user.ID.String(),
		Username: user.Username,
		Password: user.Password,
	}, nil
}
