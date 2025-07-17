package persistence

import (
	"context"
	"database/sql"
	errs "errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mingtmt/book-store/internal/auth/domain"
	"github.com/mingtmt/book-store/internal/auth/infrastructure/persistence/authsdb"
	"github.com/mingtmt/book-store/pkg/errors"
	"github.com/mingtmt/book-store/pkg/logger"
)

type AuthRepository struct {
	db *authsdb.Queries
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: authsdb.New(pool),
	}
}

func (r *AuthRepository) RegisterUser(ctx context.Context, user *domain.Auth) (string, error) {
	created, err := r.db.RegisterUser(context.Background(), authsdb.RegisterUserParams{
		Username: user.Username,
		Password: user.Password,
	})
	if err != nil {
		logger.Error("failed to create user in database", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return "", err
	}

	logger.Info("user created successfully", map[string]interface{}{
		"user_id": created.ID.String(),
	})
	return created.ID.String(), nil
}

func (r *AuthRepository) FindByUsername(ctx context.Context, username string) (*domain.Auth, error) {
	user, err := r.db.FindByUsername(ctx, username)
	if err != nil {
		if errs.Is(err, sql.ErrNoRows) {
			logger.Error("failed to find user by username", err, map[string]interface{}{
				"username": username,
			})
			return nil, errors.ErrUserNotFound
		}

		return nil, err
	}

	if !user.ID.Valid {
		logger.Warn("user not found", map[string]interface{}{
			"username": username,
		})
		return nil, nil
	}

	return &domain.Auth{
		ID:       user.ID.String(),
		Username: user.Username,
		Password: user.Password,
	}, nil
}
