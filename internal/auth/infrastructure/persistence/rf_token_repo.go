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
	"github.com/mingtmt/book-store/pkg/logger"
)

type RefreshTokenRepository struct {
	db *authsdb.Queries
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db: authsdb.New(pool),
	}
}

func (r *RefreshTokenRepository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
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

func (r *RefreshTokenRepository) GetRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
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

func (r *RefreshTokenRepository) RevokeRefreshToken(ctx context.Context, token string) error {
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

func (r *RefreshTokenRepository) DeleteExpiredRefreshTokens(ctx context.Context) error {
	err := r.db.DeleteExpiredRefreshTokens(ctx)
	if err != nil {
		logger.Error("failed to delete expired refresh tokens", err, nil)
		return err
	}

	return nil
}
