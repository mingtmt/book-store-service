package persistence

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	id := pgtype.UUID{
		Bytes: userID,
		Valid: true,
	}
	expiresAtPg := pgtype.Timestamp{
		Time:  expiresAt,
		Valid: true,
	}

	err := r.db.CreateRefreshToken(ctx, authsdb.CreateRefreshTokenParams{
		UserID:    id,
		Token:     token,
		ExpiresAt: expiresAtPg,
	})
	if err != nil {
		logger.Error("failed to create refresh token", err, map[string]interface{}{
			"user_id": userID.String(),
			"token":   token,
		})
		return err
	}

	return nil
}

func (r *AuthRepository) GetRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	rfToken, err := r.db.GetRefreshToken(ctx, token)
	if err != nil {
		if errs.Is(err, sql.ErrNoRows) {
			logger.Error("refresh token not found", err, map[string]interface{}{
				"token": token,
			})
			return nil, err
		}

		logger.Error("failed to get refresh token", err, map[string]interface{}{
			"token": token,
		})
		return nil, err
	}

	return &domain.RefreshToken{
		ID:        rfToken.ID.String(),
		UserID:    rfToken.UserID.String(),
		Token:     rfToken.Token,
		ExpiresAt: rfToken.ExpiresAt.Time,
		Revoked:   rfToken.Revoked.Bool,
		CreatedAt: rfToken.CreatedAt.Time,
	}, nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	err := r.db.RevokeRefreshToken(ctx, token)
	if err != nil {
		logger.Error("failed to revoke refresh token", err, map[string]interface{}{
			"token": token,
		})
		return err
	}

	logger.Info("refresh token revoked successfully", map[string]interface{}{
		"token": token,
	})

	return nil
}

func (r *AuthRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	err := r.db.DeleteExpiredRefreshTokens(ctx)
	if err != nil {
		logger.Error("failed to delete expired refresh tokens", err, nil)
		return err
	}

	return nil
}
